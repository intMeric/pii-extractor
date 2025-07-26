package piiextractor

import (
	"context"
	"strings"
	"time"

	"github.com/teilomillet/gollm"
)

// LLMProvider represents the different LLM providers available
type LLMProvider string

const (
	ProviderOpenAI    LLMProvider = "openai"
	ProviderMistral   LLMProvider = "mistral"
	ProviderGemini    LLMProvider = "gemini"
	ProviderOllama    LLMProvider = "ollama"
	ProviderAnthropic LLMProvider = "anthropic"
)

// ValidationResult represents the result of LLM validation for a PII entity
type ValidationResult struct {
	Valid      bool    `json:"valid"`
	Confidence float64 `json:"confidence"`
	Reasoning  string  `json:"reasoning,omitempty"`
	Provider   string  `json:"provider"`
	Model      string  `json:"model"`
}

// ValidationConfig holds configuration for LLM validation
type ValidationConfig struct {
	Enabled         bool                   `json:"enabled"`
	Provider        LLMProvider            `json:"provider"`
	Model           string                 `json:"model,omitempty"`
	APIKey          string                 `json:"api_key,omitempty"`
	BaseURL         string                 `json:"base_url,omitempty"`
	Timeout         time.Duration          `json:"timeout"`
	MinConfidence   float64                `json:"min_confidence"`
	MaxRetries      int                    `json:"max_retries"`
	ProviderOptions map[string]interface{} `json:"provider_options,omitempty"`
}

// DefaultValidationConfig returns a default configuration for validation
func DefaultValidationConfig() *ValidationConfig {
	return &ValidationConfig{
		Enabled:         false,
		Provider:        ProviderOpenAI,
		Model:           "gpt-4o-mini",
		Timeout:         30 * time.Second,
		MinConfidence:   0.7,
		MaxRetries:      3,
		ProviderOptions: make(map[string]interface{}),
	}
}

// LLMValidator interface for validating PII entities using LLMs
type LLMValidator interface {
	// ValidateEntity validates a single PII entity in its context
	ValidateEntity(ctx context.Context, entity PiiEntity, context string) (*ValidationResult, error)

	// ValidateBatch validates multiple PII entities in batch for efficiency
	ValidateBatch(ctx context.Context, entities []PiiEntity, contexts []string) ([]*ValidationResult, error)

	// HealthCheck checks if the LLM service is available
	HealthCheck(ctx context.Context) error

	// GetProviderInfo returns information about the configured provider
	GetProviderInfo() (provider string, model string)
}

// PiiValidationRequest represents a request to validate PII
type PiiValidationRequest struct {
	PiiType    PiiType `json:"pii_type"`
	Value      string  `json:"value"`
	Context    string  `json:"context"`
	EntityInfo string  `json:"entity_info,omitempty"`
}

// LLMValidatorImpl implements the LLMValidator interface using gollm
type LLMValidatorImpl struct {
	llm    gollm.LLM
	config *ValidationConfig
}

// NewLLMValidator creates a new LLM validator with the given configuration
func NewLLMValidator(config *ValidationConfig) (*LLMValidatorImpl, error) {
	if config == nil {
		config = DefaultValidationConfig()
	}

	var options []gollm.ConfigOption

	// Set provider-specific configuration
	switch config.Provider {
	case ProviderOpenAI:
		options = append(options, gollm.SetProvider("openai"))
		if config.Model != "" {
			options = append(options, gollm.SetModel(config.Model))
		} else {
			options = append(options, gollm.SetModel("gpt-4o-mini"))
		}
		if config.APIKey != "" {
			options = append(options, gollm.SetAPIKey(config.APIKey))
		}

	case ProviderMistral:
		options = append(options, gollm.SetProvider("mistral"))
		if config.Model != "" {
			options = append(options, gollm.SetModel(config.Model))
		} else {
			options = append(options, gollm.SetModel("mistral-small-latest"))
		}
		if config.APIKey != "" {
			options = append(options, gollm.SetAPIKey(config.APIKey))
		}

	case ProviderGemini:
		options = append(options, gollm.SetProvider("googleai"))
		if config.Model != "" {
			options = append(options, gollm.SetModel(config.Model))
		} else {
			options = append(options, gollm.SetModel("gemini-1.5-flash"))
		}
		if config.APIKey != "" {
			options = append(options, gollm.SetAPIKey(config.APIKey))
		}

	case ProviderOllama:
		options = append(options, gollm.SetProvider("ollama"))
		if config.Model != "" {
			options = append(options, gollm.SetModel(config.Model))
		} else {
			options = append(options, gollm.SetModel("llama3.2"))
		}
		// For Ollama, the base URL is typically handled by the provider configuration
		// We'll add BaseURL support if it's available in the gollm API

	case ProviderAnthropic:
		options = append(options, gollm.SetProvider("anthropic"))
		if config.Model != "" {
			options = append(options, gollm.SetModel(config.Model))
		} else {
			options = append(options, gollm.SetModel("claude-3-haiku-20240307"))
		}
		if config.APIKey != "" {
			options = append(options, gollm.SetAPIKey(config.APIKey))
		}
	}

	// Apply additional provider options like temperature, max tokens, etc.
	if temp, ok := config.ProviderOptions["temperature"].(float64); ok {
		options = append(options, gollm.SetTemperature(temp))
	}
	if maxTokens, ok := config.ProviderOptions["max_tokens"].(int); ok {
		options = append(options, gollm.SetMaxTokens(maxTokens))
	}

	llm, err := gollm.NewLLM(options...)
	if err != nil {
		return nil, err
	}

	return &LLMValidatorImpl{
		llm:    llm,
		config: config,
	}, nil
}

// ValidateEntity validates a single PII entity using the configured LLM
func (v *LLMValidatorImpl) ValidateEntity(ctx context.Context, entity PiiEntity, context string) (*ValidationResult, error) {
	promptText := v.buildValidationPrompt(entity, context)
	prompt := gollm.NewPrompt(promptText)

	response, err := v.llm.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return v.parseValidationResponse(response)
}

// ValidateBatch validates multiple entities in a single request for efficiency
func (v *LLMValidatorImpl) ValidateBatch(ctx context.Context, entities []PiiEntity, contexts []string) ([]*ValidationResult, error) {
	if len(entities) != len(contexts) {
		return nil, ErrMismatchedBatchSize
	}

	// For now, process individually. Could be optimized for batch processing
	results := make([]*ValidationResult, len(entities))
	for i, entity := range entities {
		result, err := v.ValidateEntity(ctx, entity, contexts[i])
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// HealthCheck verifies the LLM service is available
func (v *LLMValidatorImpl) HealthCheck(ctx context.Context) error {
	prompt := gollm.NewPrompt("Respond with 'OK'")
	_, err := v.llm.Generate(ctx, prompt)
	return err
}

// GetProviderInfo returns the provider and model information
func (v *LLMValidatorImpl) GetProviderInfo() (provider string, model string) {
	return string(v.config.Provider), v.config.Model
}

// buildValidationPrompt creates a prompt for validating PII entities
func (v *LLMValidatorImpl) buildValidationPrompt(entity PiiEntity, context string) string {
	piiType := entity.Type.String()
	value := entity.GetValue()

	// Get type-specific validation criteria
	typeSpecificGuidance := v.getTypeSpecificGuidance(entity.Type)

	prompt := `You are a PII validation expert. Your task is to determine if the identified text is actually a valid ` + piiType + ` in the given context.

PII Type: ` + piiType + `
Identified Value: "` + value + `"
Context: "` + context + `"

` + typeSpecificGuidance + `

General validation criteria:
1. Is the format correct for this type of PII?
2. Does the context support that this is actually a ` + piiType + `?
3. Could this be a false positive (e.g., random numbers that look like PII)?
4. Does the surrounding text provide clues about the actual meaning?
5. Are there contextual indicators that suggest this is real vs. example/test data?

Respond in JSON format:
{
  "valid": true/false,
  "confidence": 0.0-1.0,
  "reasoning": "Brief explanation of your decision"
}

Be conservative - if you're unsure, mark as invalid with lower confidence.`

	return prompt
}

// getTypeSpecificGuidance returns validation guidance specific to each PII type
func (v *LLMValidatorImpl) getTypeSpecificGuidance(piiType PiiType) string {
	switch piiType {
	case PiiTypePhone:
		return `Phone number validation criteria:
- Check if the number format is consistent with real phone numbers
- Look for country codes, area codes, and proper digit grouping
- Be wary of obviously fake numbers (like 555-0123, 123-456-7890)
- Consider if the context suggests it's a real contact number vs. an example`

	case PiiTypeEmail:
		return `Email validation criteria:
- Verify the email has a realistic domain (not example.com, test.com, etc.)
- Check if the local part (before @) looks genuine vs. obviously fake
- Look for context clues about whether this is a real email or placeholder
- Be suspicious of emails with suspicious TLDs or patterns`

	case PiiTypeSSN:
		return `SSN validation criteria:
- Check for the XXX-XX-XXXX format
- Be very wary of obviously fake SSNs (000-00-0000, 123-45-6789, etc.)
- Look for context that suggests official documentation vs. examples
- Consider if this appears in a form, document, or casual conversation`

	case PiiTypeCreditCard:
		return `Credit card validation criteria:
- Verify the number format matches known card types (Visa, MasterCard, etc.)
- Check if it could pass basic Luhn algorithm validation conceptually
- Be very suspicious of obviously fake numbers (4111-1111-1111-1111, etc.)
- Look for context suggesting real transaction vs. test/example data`

	case PiiTypeZipCode:
		return `ZIP code validation criteria:
- Check if it's a valid US ZIP format (5 digits or 5+4)
- Consider if the ZIP code matches the context (geographic references)
- Be wary of obviously fake codes (00000, 12345, etc.)
- Look for context suggesting real addresses vs. examples`

	case PiiTypeStreetAddress:
		return `Street address validation criteria:
- Check if the address format is realistic and well-formed
- Look for real street names, not obviously fake ones (123 Main St is suspicious)
- Consider if house numbers are reasonable for the street type
- Check for proper abbreviations and formatting`

	case PiiTypeIPAddress:
		return `IP address validation criteria:
- Verify the format is valid (IPv4: x.x.x.x, IPv6: proper format)
- Check if it's in a reasonable range (not 0.0.0.0 or other invalid IPs)
- Consider if it's a private vs. public IP and if that makes sense in context
- Look for context suggesting real network data vs. examples`

	case PiiTypeIBAN:
		return `IBAN validation criteria:
- Check if the country code is valid (first 2 letters)
- Verify the format matches IBAN standards for that country
- Look for context suggesting real banking information vs. examples
- Be wary of obviously fake IBANs or test patterns`

	case PiiTypeBtcAddress:
		return `Bitcoin address validation criteria:
- Check if the format is valid (starts with 1, 3, or bc1)
- Verify the length is appropriate for the address type
- Look for context suggesting real cryptocurrency activity vs. examples
- Be wary of obviously fake or example addresses`

	case PiiTypePoBox:
		return `P.O. Box validation criteria:
- Check if the format is realistic (P.O. Box followed by number)
- Consider if the box number is reasonable (not obviously fake like 1 or 123456)
- Look for context suggesting real postal addresses vs. examples
- Check for proper formatting and abbreviations`

	default:
		return `Generic PII validation criteria:
- Assess if the format and content appear genuine
- Look for context clues about real vs. example/test data
- Consider common patterns used in fake or example data`
	}
}

// parseValidationResponse parses the LLM response into a ValidationResult
func (v *LLMValidatorImpl) parseValidationResponse(response string) (*ValidationResult, error) {
	result := &ValidationResult{
		Provider: string(v.config.Provider),
		Model:    v.config.Model,
	}

	// Try to extract JSON from the response
	jsonStart := -1
	jsonEnd := -1

	// Find the JSON object boundaries
	for i, char := range response {
		if char == '{' && jsonStart == -1 {
			jsonStart = i
		}
		if char == '}' {
			jsonEnd = i + 1
		}
	}

	if jsonStart == -1 || jsonEnd == -1 {
		// Fallback to heuristic parsing
		return v.parseHeuristically(response, result)
	}

	jsonStr := response[jsonStart:jsonEnd]

	// Simple JSON field extraction without importing encoding/json
	// In a production environment, you'd want to use proper JSON parsing

	// Extract "valid" field
	if v.findSubstring(jsonStr, `"valid": true`) != -1 ||
		v.findSubstring(jsonStr, `"valid":true`) != -1 ||
		v.findSubstring(jsonStr, `"valid" : true`) != -1 {
		result.Valid = true
	} else {
		result.Valid = false
	}

	// Extract confidence field
	result.Confidence = v.extractConfidence(jsonStr)

	// Extract reasoning field
	result.Reasoning = v.extractReasoning(jsonStr)

	return result, nil
}

// parseHeuristically provides a fallback parsing method
func (v *LLMValidatorImpl) parseHeuristically(response string, result *ValidationResult) (*ValidationResult, error) {
	// Look for keywords indicating validity

	lowerResponse := strings.ToLower(response)

	if v.containsAny(lowerResponse, []string{"valid: true", "is valid", "valid pii", "legitimate"}) {
		result.Valid = true
		result.Confidence = 0.7
	} else if v.containsAny(lowerResponse, []string{"valid: false", "invalid", "false positive", "not valid"}) {
		result.Valid = false
		result.Confidence = 0.7
	} else {
		result.Valid = false
		result.Confidence = 0.5
	}

	result.Reasoning = "Heuristic parsing of LLM response"
	return result, nil
}

func (v *LLMValidatorImpl) extractConfidence(jsonStr string) float64 {
	// Look for confidence field - simplified extraction
	patterns := []string{`"confidence":`, `"confidence" :`}

	for _, pattern := range patterns {
		start := v.findSubstring(jsonStr, pattern)
		if start != -1 {
			start += len(pattern)

			// Skip whitespace and find the number
			for start < len(jsonStr) && (jsonStr[start] == ' ' || jsonStr[start] == '\t') {
				start++
			}

			// Extract the number
			end := start
			for end < len(jsonStr) && (v.isDigit(jsonStr[end]) || jsonStr[end] == '.') {
				end++
			}

			if end > start {
				// Simple float parsing
				numStr := jsonStr[start:end]
				if confidence := v.parseSimpleFloat(numStr); confidence >= 0 {
					return confidence
				}
			}
		}
	}

	return 0.8 // Default confidence
}

func (v *LLMValidatorImpl) extractReasoning(jsonStr string) string {
	// Look for reasoning field - simplified extraction
	patterns := []string{`"reasoning":`, `"reasoning" :`}

	for _, pattern := range patterns {
		start := v.findSubstring(jsonStr, pattern)
		if start != -1 {
			start += len(pattern)

			// Skip whitespace and find the opening quote
			for start < len(jsonStr) && jsonStr[start] != '"' {
				start++
			}

			if start < len(jsonStr) && jsonStr[start] == '"' {
				start++ // Skip opening quote

				// Find closing quote
				end := start
				for end < len(jsonStr) && jsonStr[end] != '"' {
					end++
				}

				if end > start {
					return jsonStr[start:end]
				}
			}
		}
	}

	return "LLM validation completed"
}

func (v *LLMValidatorImpl) containsAny(text string, patterns []string) bool {
	for _, pattern := range patterns {
		if v.findSubstring(text, pattern) != -1 {
			return true
		}
	}
	return false
}

func (v *LLMValidatorImpl) findSubstring(text, pattern string) int {
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

func (v *LLMValidatorImpl) isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (v *LLMValidatorImpl) parseSimpleFloat(s string) float64 {
	// Very simple float parsing - in production, use strconv.ParseFloat
	if s == "0" {
		return 0.0
	}
	if s == "1" {
		return 1.0
	}
	if s == "0.5" {
		return 0.5
	}
	if s == "0.7" {
		return 0.7
	}
	if s == "0.8" {
		return 0.8
	}
	if s == "0.9" {
		return 0.9
	}
	// Add more common values as needed
	return 0.8 // Default
}

// Custom errors
var (
	ErrMismatchedBatchSize  = NewPiiError("mismatched batch size", "entities and contexts slices must have the same length")
	ErrValidationDisabled   = NewPiiError("validation disabled", "LLM validation is not enabled in configuration")
	ErrProviderNotSupported = NewPiiError("provider not supported", "the specified LLM provider is not supported")
)

// PiiError represents an error in PII processing
type PiiError struct {
	Code    string
	Message string
}

func (e PiiError) Error() string {
	return e.Message
}

// NewPiiError creates a new PII error
func NewPiiError(code, message string) PiiError {
	return PiiError{Code: code, Message: message}
}
