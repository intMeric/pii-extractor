package llm

import (
	"context"
	"fmt"
	"github.com/intMeric/pii-extractor/pii"
	"github.com/intMeric/pii-extractor/extractors"
	"github.com/teilomillet/gollm"
)

// Provider represents an LLM provider
type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderClaude    Provider = "claude"
	ProviderOllama    Provider = "ollama"
	ProviderAzureAI   Provider = "azure"
	ProviderMistral   Provider = "mistral"
	ProviderGemini    Provider = "gemini"
	ProviderAnthropic Provider = "anthropic"
)

// LLMExtractor implements PII extraction using Large Language Models
type LLMExtractor struct {
	name     string
	provider Provider
	model    string
	apiKey   string
	baseURL  string
	config   LLMConfig
	llm      gollm.LLM
}

// LLMConfig contains LLM-specific configuration
type LLMConfig struct {
	Temperature   float32 `json:"temperature"`
	MaxTokens     int     `json:"max_tokens"`
	SystemPrompt  string  `json:"system_prompt"`
	RetryAttempts int     `json:"retry_attempts"`
	Timeout       int     `json:"timeout_seconds"`
}

// NewExtractor creates a new LLM-based PII extractor
func NewExtractor(provider Provider, model string, config *extractors.ExtractorConfig) (*LLMExtractor, error) {
	extractor := &LLMExtractor{
		name:     "llm-extractor",
		provider: provider,
		model:    model,
		config: LLMConfig{
			Temperature:   0.1, // Low temperature for consistent extraction
			MaxTokens:     2048,
			RetryAttempts: 3,
			Timeout:       30,
		},
	}
	
	if config != nil && config.Options != nil {
		if apiKey, ok := config.Options["api_key"].(string); ok {
			extractor.apiKey = apiKey
		}
		if baseURL, ok := config.Options["base_url"].(string); ok {
			extractor.baseURL = baseURL
		}
		if temp, ok := config.Options["temperature"].(float32); ok {
			extractor.config.Temperature = temp
		}
	}
	
	// Initialize gollm LLM
	var options []gollm.ConfigOption
	
	// Set provider-specific configuration
	switch provider {
	case ProviderOpenAI:
		options = append(options, gollm.SetProvider("openai"))
		if model != "" {
			options = append(options, gollm.SetModel(model))
		} else {
			options = append(options, gollm.SetModel("gpt-4o-mini"))
		}
		if extractor.apiKey != "" {
			options = append(options, gollm.SetAPIKey(extractor.apiKey))
		}

	case ProviderMistral:
		options = append(options, gollm.SetProvider("mistral"))
		if model != "" {
			options = append(options, gollm.SetModel(model))
		} else {
			options = append(options, gollm.SetModel("mistral-small-latest"))
		}
		if extractor.apiKey != "" {
			options = append(options, gollm.SetAPIKey(extractor.apiKey))
		}

	case ProviderGemini:
		options = append(options, gollm.SetProvider("googleai"))
		if model != "" {
			options = append(options, gollm.SetModel(model))
		} else {
			options = append(options, gollm.SetModel("gemini-1.5-flash"))
		}
		if extractor.apiKey != "" {
			options = append(options, gollm.SetAPIKey(extractor.apiKey))
		}

	case ProviderOllama:
		options = append(options, gollm.SetProvider("ollama"))
		if model != "" {
			options = append(options, gollm.SetModel(model))
		} else {
			options = append(options, gollm.SetModel("llama3.2"))
		}

	case ProviderAnthropic:
		options = append(options, gollm.SetProvider("anthropic"))
		if model != "" {
			options = append(options, gollm.SetModel(model))
		} else {
			options = append(options, gollm.SetModel("claude-3-haiku-20240307"))
		}
		if extractor.apiKey != "" {
			options = append(options, gollm.SetAPIKey(extractor.apiKey))
		}
	
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	
	// Apply LLM configuration
	options = append(options, gollm.SetTemperature(float64(extractor.config.Temperature)))
	options = append(options, gollm.SetMaxTokens(extractor.config.MaxTokens))
	
	llm, err := gollm.NewLLM(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize LLM: %w", err)
	}
	
	extractor.llm = llm
	return extractor, nil
}

// Extract performs PII extraction using LLM
func (l *LLMExtractor) Extract(text string) (*pii.PiiExtractionResult, error) {
	// Prepare prompt for PII extraction
	prompt := l.buildExtractionPrompt(text)
	
	// Create context for LLM call
	ctx := context.Background()
	
	// Call LLM
	response, err := l.llm.Generate(ctx, gollm.NewPrompt(prompt))
	if err != nil {
		return nil, fmt.Errorf("LLM extraction failed: %w", err)
	}
	
	// Parse response to PiiEntity objects
	entities, err := l.parseExtractionResponse(response, text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}
	
	return pii.NewPiiExtractionResult(entities), nil
}

// ExtractByType extracts specific PII types using LLM
func (l *LLMExtractor) ExtractByType(text string, piiType pii.PiiType) ([]pii.PiiEntity, error) {
	// Prepare type-specific prompt
	prompt := l.buildTypeSpecificPrompt(text, piiType)
	
	// Create context for LLM call
	ctx := context.Background()
	
	// Call LLM
	response, err := l.llm.Generate(ctx, gollm.NewPrompt(prompt))
	if err != nil {
		return nil, fmt.Errorf("LLM type-specific extraction failed: %w", err)
	}
	
	// Parse response to PiiEntity objects
	entities, err := l.parseExtractionResponse(response, text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}
	
	// Filter entities to only include the requested type
	var filtered []pii.PiiEntity
	for _, entity := range entities {
		if entity.Type == piiType {
			filtered = append(filtered, entity)
		}
	}
	
	return filtered, nil
}

// GetSupportedTypes returns PII types this LLM extractor can handle
func (l *LLMExtractor) GetSupportedTypes() []pii.PiiType {
	// LLM can potentially handle all types, but we'll be conservative
	return []pii.PiiType{
		pii.PiiTypePhone,
		pii.PiiTypeEmail,
		pii.PiiTypeSSN,
		pii.PiiTypeZipCode,
		pii.PiiTypeStreetAddress,
		pii.PiiTypeCreditCard,
		pii.PiiTypeIPAddress,
		pii.PiiTypeBtcAddress,
		pii.PiiTypeIBAN,
		pii.PiiTypePoBox,
	}
}

// GetMethod returns the extraction method
func (l *LLMExtractor) GetMethod() extractors.ExtractionMethod {
	return extractors.MethodLLM
}

// GetName returns the extractor name
func (l *LLMExtractor) GetName() string {
	return l.name
}

// GetProvider returns the LLM provider being used
func (l *LLMExtractor) GetProvider() Provider {
	return l.provider
}

// GetModel returns the model being used
func (l *LLMExtractor) GetModel() string {
	return l.model
}

// buildExtractionPrompt creates a prompt for general PII extraction
func (l *LLMExtractor) buildExtractionPrompt(text string) string {
	return fmt.Sprintf(`You are a PII (Personally Identifiable Information) extraction expert. Analyze the following text and extract all PII entities.

Text to analyze:
%s

Extract the following types of PII if present:
- Email addresses
- Phone numbers (US format)
- Social Security Numbers (SSN)
- ZIP codes
- Street addresses
- Credit card numbers
- IP addresses
- Bitcoin addresses
- IBAN numbers
- P.O. Box addresses

Respond in JSON format with an array of objects, each containing:
{
  "type": "email|phone|ssn|zipcode|address|creditcard|ip|bitcoin|iban|pobox",
  "value": "extracted_value",
  "context": "surrounding_text_context"
}

Example response:
[
  {
    "type": "email",
    "value": "john@example.com",
    "context": "Contact me at john@example.com for more info"
  },
  {
    "type": "phone",
    "value": "555-123-4567",
    "context": "Call me at 555-123-4567"
  }
]

If no PII is found, respond with an empty array: []`, text)
}

// buildTypeSpecificPrompt creates a prompt for extracting specific PII types
func (l *LLMExtractor) buildTypeSpecificPrompt(text string, piiType pii.PiiType) string {
	typeStr := piiType.String()
	
	return fmt.Sprintf(`You are a PII extraction expert. Analyze the following text and extract only %s entities.

Text to analyze:
%s

Focus specifically on finding %s in the text. Be precise and only extract genuine %s entities, not false positives.

Respond in JSON format with an array of objects:
[
  {
    "type": "%s",
    "value": "extracted_value",
    "context": "surrounding_text_context"
  }
]

If no %s entities are found, respond with an empty array: []`, typeStr, text, typeStr, typeStr, typeStr, typeStr)
}

// parseExtractionResponse parses the LLM response into PiiEntity objects
func (l *LLMExtractor) parseExtractionResponse(response, originalText string) ([]pii.PiiEntity, error) {
	// Simple JSON parsing without importing encoding/json
	// In production, you'd want to use proper JSON parsing
	
	var entities []pii.PiiEntity
	
	// Find JSON array start and end
	arrayStart := -1
	arrayEnd := -1
	bracketCount := 0
	
	for i, char := range response {
		if char == '[' {
			if arrayStart == -1 {
				arrayStart = i
			}
			bracketCount++
		} else if char == ']' {
			bracketCount--
			if bracketCount == 0 && arrayStart != -1 {
				arrayEnd = i + 1
				break
			}
		}
	}
	
	if arrayStart == -1 || arrayEnd == -1 {
		// No valid JSON array found, return empty result
		return entities, nil
	}
	
	jsonStr := response[arrayStart:arrayEnd]
	
	// Simple object extraction (this is a simplified parser)
	entities = l.extractEntitiesFromJSON(jsonStr, originalText)
	
	return entities, nil
}

// extractEntitiesFromJSON extracts entities from JSON string (simplified parser)
func (l *LLMExtractor) extractEntitiesFromJSON(jsonStr, originalText string) []pii.PiiEntity {
	var entities []pii.PiiEntity
	
	// Look for object patterns in the JSON
	objectStart := -1
	braceCount := 0
	
	for i, char := range jsonStr {
		if char == '{' {
			if objectStart == -1 {
				objectStart = i
			}
			braceCount++
		} else if char == '}' {
			braceCount--
			if braceCount == 0 && objectStart != -1 {
				objectEnd := i + 1
				objectStr := jsonStr[objectStart:objectEnd]
				
				// Parse individual object
				entity := l.parseEntityObject(objectStr, originalText)
				if entity != nil {
					entities = append(entities, *entity)
				}
				
				objectStart = -1
			}
		}
	}
	
	return entities
}

// parseEntityObject parses a single entity object from JSON
func (l *LLMExtractor) parseEntityObject(objectStr, originalText string) *pii.PiiEntity {
	// Extract type, value, and context using simple string parsing
	piiType := l.extractJSONField(objectStr, "type")
	value := l.extractJSONField(objectStr, "value")
	context := l.extractJSONField(objectStr, "context")
	
	if piiType == "" || value == "" {
		return nil
	}
	
	// Convert string type to PiiType
	var entityType pii.PiiType
	switch piiType {
	case "email":
		entityType = pii.PiiTypeEmail
	case "phone":
		entityType = pii.PiiTypePhone
	case "ssn":
		entityType = pii.PiiTypeSSN
	case "zipcode":
		entityType = pii.PiiTypeZipCode
	case "address":
		entityType = pii.PiiTypeStreetAddress
	case "creditcard":
		entityType = pii.PiiTypeCreditCard
	case "ip":
		entityType = pii.PiiTypeIPAddress
	case "bitcoin":
		entityType = pii.PiiTypeBtcAddress
	case "iban":
		entityType = pii.PiiTypeIBAN
	case "pobox":
		entityType = pii.PiiTypePoBox
	default:
		return nil // Unknown type
	}
	
	// Create appropriate PII value object
	var piiValue pii.Pii
	switch entityType {
	case pii.PiiTypeEmail:
		piiValue = pii.NewEmail(value)
	case pii.PiiTypePhone:
		piiValue = pii.NewPhoneUS(value)
	case pii.PiiTypeSSN:
		piiValue = pii.NewSSN(value)
	case pii.PiiTypeZipCode:
		piiValue = pii.NewZipCode(value, "US")
	case pii.PiiTypeStreetAddress:
		piiValue = pii.NewStreetAddress(value, "US")
	case pii.PiiTypeCreditCard:
		piiValue = pii.NewCreditCard(value, "unknown")
	case pii.PiiTypeIPAddress:
		piiValue = pii.NewIPAddress(value, "IPv4")
	case pii.PiiTypeBtcAddress:
		piiValue = pii.NewBtcAddress(value)
	case pii.PiiTypeIBAN:
		piiValue = pii.NewIBAN(value, "unknown")
	case pii.PiiTypePoBox:
		piiValue = pii.NewPoBox(value, "US")
	default:
		return nil
	}
	
	// Add context if provided
	if context != "" {
		// Add context to the PII value (which has BasePii embedded)
		if basePii, ok := piiValue.(interface{ AddContext(string) }); ok {
			basePii.AddContext(context)
		}
	}
	
	entity := pii.PiiEntity{
		Type:  entityType,
		Value: piiValue,
	}
	
	return &entity
}

// extractJSONField extracts a field value from a JSON object string
func (l *LLMExtractor) extractJSONField(objectStr, fieldName string) string {
	// Look for "fieldName": "value" pattern
	patterns := []string{
		fmt.Sprintf(`"%s": "`, fieldName),
		fmt.Sprintf(`"%s":"`, fieldName),
		fmt.Sprintf(`"%s" : "`, fieldName),
	}
	
	for _, pattern := range patterns {
		start := l.findSubstring(objectStr, pattern)
		if start != -1 {
			start += len(pattern)
			
			// Find closing quote
			end := start
			for end < len(objectStr) && objectStr[end] != '"' {
				end++
			}
			
			if end > start {
				return objectStr[start:end]
			}
		}
	}
	
	return ""
}

// findSubstring finds a substring in text
func (l *LLMExtractor) findSubstring(text, pattern string) int {
	if len(pattern) > len(text) {
		return -1
	}

	for i := 0; i <= len(text)-len(pattern); i++ {
		match := true
		for j := 0; j < len(pattern); j++ {
			if text[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}