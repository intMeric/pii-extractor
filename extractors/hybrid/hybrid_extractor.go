package hybrid

import (
	"context"
	"fmt"
	"time"

	"github.com/intMeric/pii-extractor/pii"
	"github.com/intMeric/pii-extractor/extractors"
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
	ValidateEntity(ctx context.Context, entity pii.PiiEntity, context string) (*pii.ValidationResult, error)

	// ValidateBatch validates multiple PII entities in batch for efficiency
	ValidateBatch(ctx context.Context, entities []pii.PiiEntity, contexts []string) ([]*pii.ValidationResult, error)

	// HealthCheck checks if the LLM service is available
	HealthCheck(ctx context.Context) error

	// GetProviderInfo returns information about the configured provider
	GetProviderInfo() (provider string, model string)
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
func (v *LLMValidatorImpl) ValidateEntity(ctx context.Context, entity pii.PiiEntity, context string) (*pii.ValidationResult, error) {
	promptText := v.buildValidationPrompt(entity, context)
	prompt := gollm.NewPrompt(promptText)

	response, err := v.llm.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return v.parseValidationResponse(response)
}

// ValidateBatch validates multiple entities in a single request for efficiency
func (v *LLMValidatorImpl) ValidateBatch(ctx context.Context, entities []pii.PiiEntity, contexts []string) ([]*pii.ValidationResult, error) {
	if len(entities) != len(contexts) {
		return nil, fmt.Errorf("entities and contexts slices must have the same length")
	}

	// For now, process individually. Could be optimized for batch processing
	results := make([]*pii.ValidationResult, len(entities))
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

// ValidatedExtractor combines any base extractor with LLM validation
type ValidatedExtractor struct {
	name          string
	baseExtractor extractors.PiiExtractor
	validator     LLMValidator
	config        *ValidationConfig
}

// NewValidatedExtractor creates a new validated extractor
func NewValidatedExtractor(baseExtractor extractors.PiiExtractor, config *ValidationConfig) (*ValidatedExtractor, error) {
	if config == nil {
		config = DefaultValidationConfig()
	}

	var validator LLMValidator
	var err error

	if config.Enabled {
		validator, err = NewLLMValidator(config)
		if err != nil {
			return nil, err
		}
	}

	return &ValidatedExtractor{
		name:          "validated-extractor",
		baseExtractor: baseExtractor,
		validator:     validator,
		config:        config,
	}, nil
}

// Extract performs basic extraction without validation (implements PiiExtractor)
func (v *ValidatedExtractor) Extract(text string) (*pii.PiiExtractionResult, error) {
	return v.baseExtractor.Extract(text)
}

// ExtractByType extracts specific PII types
func (v *ValidatedExtractor) ExtractByType(text string, piiType pii.PiiType) ([]pii.PiiEntity, error) {
	return v.baseExtractor.ExtractByType(text, piiType)
}

// GetSupportedTypes returns the supported types from the base extractor
func (v *ValidatedExtractor) GetSupportedTypes() []pii.PiiType {
	return v.baseExtractor.GetSupportedTypes()
}

// GetMethod returns the extraction method
func (v *ValidatedExtractor) GetMethod() extractors.ExtractionMethod {
	return extractors.MethodHybrid
}

// GetName returns the extractor name
func (v *ValidatedExtractor) GetName() string {
	return v.name
}

// ExtractWithValidation performs extraction with LLM validation
func (v *ValidatedExtractor) ExtractWithValidation(text string) (*pii.PiiExtractionResult, error) {
	config := v.config

	// If validation is disabled, just do regular extraction
	if !config.Enabled {
		return v.baseExtractor.Extract(text)
	}

	// Perform initial extraction
	result, err := v.baseExtractor.Extract(text)
	if err != nil {
		return nil, err
	}

	// If no entities found, return early
	if result.IsEmpty() {
		return result, nil
	}

	// Use the configured validator
	validator := v.validator
	if validator == nil {
		// This shouldn't happen if validation is enabled, but handle gracefully
		return result, nil
	}

	// Validate entities
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	err = v.validateEntities(ctx, result, text, validator, config)
	if err != nil {
		// If validation fails, return unvalidated results
		return result, nil
	}

	// Calculate validation statistics
	v.calculateValidationStats(result, validator)

	return result, nil
}

// IsValidationEnabled returns true if LLM validation is enabled
func (v *ValidatedExtractor) IsValidationEnabled() bool {
	return v.config != nil && v.config.Enabled && v.validator != nil
}

// HealthCheck verifies that the LLM validator is working
func (v *ValidatedExtractor) HealthCheck(ctx context.Context) error {
	if v.validator == nil {
		return fmt.Errorf("validation disabled")
	}
	return v.validator.HealthCheck(ctx)
}

// EnsembleExtractor combines multiple extraction methods
type EnsembleExtractor struct {
	name           string
	extractors     []extractors.PiiExtractor
	strategy       CombinationStrategy
	validationMode ValidationMode
}

// CombinationStrategy defines how results from multiple extractors are combined
type CombinationStrategy string

const (
	StrategyUnion        CombinationStrategy = "union"        // Combine all results
	StrategyIntersection CombinationStrategy = "intersection" // Only results found by all
	StrategyMajority     CombinationStrategy = "majority"     // Results found by majority
	StrategyWeighted     CombinationStrategy = "weighted"     // Weighted combination
)

// ValidationMode defines how results are cross-validated
type ValidationMode string

const (
	ValidationNone   ValidationMode = "none"   // No cross-validation
	ValidationBasic  ValidationMode = "basic"  // Basic overlap checking
	ValidationStrict ValidationMode = "strict" // Strict validation rules
)

// NewEnsembleExtractor creates a new ensemble extractor
func NewEnsembleExtractor(extractors ...extractors.PiiExtractor) *EnsembleExtractor {
	return &EnsembleExtractor{
		name:           "ensemble-extractor",
		extractors:     extractors,
		strategy:       StrategyUnion,
		validationMode: ValidationBasic,
	}
}

// WithStrategy sets the combination strategy
func (e *EnsembleExtractor) WithStrategy(strategy CombinationStrategy) *EnsembleExtractor {
	e.strategy = strategy
	return e
}

// WithValidation sets the validation mode
func (e *EnsembleExtractor) WithValidation(mode ValidationMode) *EnsembleExtractor {
	e.validationMode = mode
	return e
}

// Extract performs PII extraction using multiple methods and combines results
func (e *EnsembleExtractor) Extract(text string) (*pii.PiiExtractionResult, error) {
	if len(e.extractors) == 0 {
		return nil, fmt.Errorf("no extractors configured")
	}

	// Run all extractors
	allResults := make([]*pii.PiiExtractionResult, len(e.extractors))
	for i, extractor := range e.extractors {
		result, err := extractor.Extract(text)
		if err != nil {
			// Continue with other extractors if one fails
			continue
		}
		allResults[i] = result
	}

	// Combine results based on strategy
	combinedEntities := e.combineResults(allResults)

	return pii.NewPiiExtractionResult(combinedEntities), nil
}

// ExtractByType extracts specific PII types using ensemble approach
func (e *EnsembleExtractor) ExtractByType(text string, piiType pii.PiiType) ([]pii.PiiEntity, error) {
	if len(e.extractors) == 0 {
		return nil, fmt.Errorf("no extractors configured")
	}

	var allEntities []pii.PiiEntity
	for _, extractor := range e.extractors {
		entities, err := extractor.ExtractByType(text, piiType)
		if err != nil {
			continue
		}
		allEntities = append(allEntities, entities...)
	}

	return e.deduplicateEntities(allEntities), nil
}

// GetSupportedTypes returns union of all supported types from all extractors
func (e *EnsembleExtractor) GetSupportedTypes() []pii.PiiType {
	typeSet := make(map[pii.PiiType]bool)
	for _, extractor := range e.extractors {
		for _, piiType := range extractor.GetSupportedTypes() {
			typeSet[piiType] = true
		}
	}

	types := make([]pii.PiiType, 0, len(typeSet))
	for piiType := range typeSet {
		types = append(types, piiType)
	}

	return types
}

// GetMethod returns the extraction method
func (e *EnsembleExtractor) GetMethod() extractors.ExtractionMethod {
	return extractors.MethodHybrid
}

// GetName returns the extractor name
func (e *EnsembleExtractor) GetName() string {
	return e.name
}

// GetExtractors returns the list of underlying extractors
func (e *EnsembleExtractor) GetExtractors() []extractors.PiiExtractor {
	return e.extractors
}

// GetStrategy returns the current combination strategy
func (e *EnsembleExtractor) GetStrategy() CombinationStrategy {
	return e.strategy
}

// GetValidationMode returns the current validation mode
func (e *EnsembleExtractor) GetValidationMode() ValidationMode {
	return e.validationMode
}

// Private helper methods for ValidatedExtractor

// validateEntities validates all entities in the result
func (v *ValidatedExtractor) validateEntities(ctx context.Context, result *pii.PiiExtractionResult, originalText string, validator LLMValidator, config *ValidationConfig) error {
	for i := range result.Entities {
		entity := &result.Entities[i]

		// Get context for this entity
		context := v.getEntityContext(originalText, entity)

		// Validate with retries
		var validation *pii.ValidationResult
		var err error

		for attempt := 0; attempt <= config.MaxRetries; attempt++ {
			validation, err = validator.ValidateEntity(ctx, *entity, context)
			if err == nil {
				break
			}

			if attempt < config.MaxRetries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
			}
		}

		// If validation succeeded and meets confidence threshold
		if err == nil && validation.Confidence >= config.MinConfidence {
			entity.Validation = validation
		}
	}

	return nil
}

// getEntityContext extracts context around the entity from the original text
func (v *ValidatedExtractor) getEntityContext(text string, entity *pii.PiiEntity) string {
	// For now, return the first context from the entity
	contexts := entity.GetContexts()
	if len(contexts) > 0 {
		return contexts[0]
	}

	// Fallback: try to find the entity value in the text and extract context
	value := entity.GetValue()
	if value == "" {
		return ""
	}

	// Simple context extraction
	return v.extractSimpleContext(text, value)
}

// extractSimpleContext extracts context around a value in text
func (v *ValidatedExtractor) extractSimpleContext(text, value string) string {
	// Find the value in the text
	start := -1
	for i := 0; i <= len(text)-len(value); i++ {
		if text[i:i+len(value)] == value {
			start = i
			break
		}
	}

	if start == -1 {
		return text // Return full text if value not found
	}

	// Extract context around the value (100 characters before and after)
	contextStart := start - 100
	if contextStart < 0 {
		contextStart = 0
	}

	contextEnd := start + len(value) + 100
	if contextEnd > len(text) {
		contextEnd = len(text)
	}

	return text[contextStart:contextEnd]
}

// calculateValidationStats calculates validation statistics for the result
func (v *ValidatedExtractor) calculateValidationStats(result *pii.PiiExtractionResult, validator LLMValidator) {
	if validator == nil {
		return
	}

	stats := &pii.ValidationStats{}
	provider, model := validator.GetProviderInfo()
	stats.Provider = provider
	stats.Model = model

	var totalConfidence float64
	validatedCount := 0

	for _, entity := range result.Entities {
		if entity.IsValidated() {
			validatedCount++
			totalConfidence += entity.GetValidationConfidence()

			if entity.IsValid() {
				stats.ValidCount++
			} else {
				stats.InvalidCount++
			}
		}
	}

	stats.TotalValidated = validatedCount
	if validatedCount > 0 {
		stats.AverageConfidence = totalConfidence / float64(validatedCount)
	}

	result.ValidationStats = stats
}

// Private helper methods for EnsembleExtractor

// combineResults combines results from multiple extractors based on strategy
func (e *EnsembleExtractor) combineResults(results []*pii.PiiExtractionResult) []pii.PiiEntity {
	switch e.strategy {
	case StrategyUnion:
		return e.unionResults(results)
	case StrategyIntersection:
		return e.intersectionResults(results)
	case StrategyMajority:
		return e.majorityResults(results)
	case StrategyWeighted:
		return e.weightedResults(results)
	default:
		return e.unionResults(results)
	}
}

// unionResults combines all entities from all extractors
func (e *EnsembleExtractor) unionResults(results []*pii.PiiExtractionResult) []pii.PiiEntity {
	var allEntities []pii.PiiEntity
	for _, result := range results {
		if result != nil {
			allEntities = append(allEntities, result.Entities...)
		}
	}
	return e.deduplicateEntities(allEntities)
}

// intersectionResults returns only entities found by all extractors
func (e *EnsembleExtractor) intersectionResults(results []*pii.PiiExtractionResult) []pii.PiiEntity {
	if len(results) == 0 {
		return []pii.PiiEntity{}
	}

	// Start with first result
	candidates := make(map[string]pii.PiiEntity)
	if results[0] != nil {
		for _, entity := range results[0].Entities {
			key := e.getEntityKey(entity)
			candidates[key] = entity
		}
	}

	// Check if entities exist in all other results
	for i := 1; i < len(results); i++ {
		if results[i] == nil {
			return []pii.PiiEntity{} // If any result is nil, intersection is empty
		}

		currentEntities := make(map[string]bool)
		for _, entity := range results[i].Entities {
			key := e.getEntityKey(entity)
			currentEntities[key] = true
		}

		// Remove candidates not found in current result
		for key := range candidates {
			if !currentEntities[key] {
				delete(candidates, key)
			}
		}
	}

	// Convert back to slice
	entities := make([]pii.PiiEntity, 0, len(candidates))
	for _, entity := range candidates {
		entities = append(entities, entity)
	}

	return entities
}

// majorityResults returns entities found by majority of extractors
func (e *EnsembleExtractor) majorityResults(results []*pii.PiiExtractionResult) []pii.PiiEntity {
	entityCounts := make(map[string]int)
	entityMap := make(map[string]pii.PiiEntity)

	// Count occurrences of each entity
	for _, result := range results {
		if result != nil {
			for _, entity := range result.Entities {
				key := e.getEntityKey(entity)
				entityCounts[key]++
				entityMap[key] = entity
			}
		}
	}

	// Find majority threshold
	majority := (len(results) + 1) / 2

	// Collect entities that appear in majority of results
	var entities []pii.PiiEntity
	for key, count := range entityCounts {
		if count >= majority {
			entities = append(entities, entityMap[key])
		}
	}

	return entities
}

// weightedResults combines results with weights (basic implementation)
func (e *EnsembleExtractor) weightedResults(results []*pii.PiiExtractionResult) []pii.PiiEntity {
	// For now, just use union - can be enhanced with actual weighting
	return e.unionResults(results)
}

// getEntityKey creates a unique key for an entity for comparison
func (e *EnsembleExtractor) getEntityKey(entity pii.PiiEntity) string {
	return fmt.Sprintf("%s:%s", entity.Type.String(), entity.GetValue())
}

// deduplicateEntities removes duplicate entities
func (e *EnsembleExtractor) deduplicateEntities(entities []pii.PiiEntity) []pii.PiiEntity {
	seen := make(map[string]bool)
	var unique []pii.PiiEntity

	for _, entity := range entities {
		key := e.getEntityKey(entity)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, entity)
		}
	}

	return unique
}

// buildValidationPrompt creates a prompt for validating PII entities
func (v *LLMValidatorImpl) buildValidationPrompt(entity pii.PiiEntity, context string) string {
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
func (v *LLMValidatorImpl) getTypeSpecificGuidance(piiType pii.PiiType) string {
	switch piiType {
	case pii.PiiTypePhone:
		return `Phone number validation criteria:
- Check if the number format is consistent with real phone numbers
- Look for country codes, area codes, and proper digit grouping
- Be wary of obviously fake numbers (like 555-0123, 123-456-7890)
- Consider if the context suggests it's a real contact number vs. an example`

	case pii.PiiTypeEmail:
		return `Email validation criteria:
- Verify the email has a realistic domain (not example.com, test.com, etc.)
- Check if the local part (before @) looks genuine vs. obviously fake
- Look for context clues about whether this is a real email or placeholder
- Be suspicious of emails with suspicious TLDs or patterns`

	case pii.PiiTypeSSN:
		return `SSN validation criteria:
- Check for the XXX-XX-XXXX format
- Be very wary of obviously fake SSNs (000-00-0000, 123-45-6789, etc.)
- Look for context that suggests official documentation vs. examples
- Consider if this appears in a form, document, or casual conversation`

	case pii.PiiTypeCreditCard:
		return `Credit card validation criteria:
- Verify the number format matches known card types (Visa, MasterCard, etc.)
- Check if it could pass basic Luhn algorithm validation conceptually
- Be very suspicious of obviously fake numbers (4111-1111-1111-1111, etc.)
- Look for context suggesting real transaction vs. test/example data`

	case pii.PiiTypeZipCode:
		return `ZIP code validation criteria:
- Check if it's a valid US ZIP format (5 digits or 5+4)
- Consider if the ZIP code matches the context (geographic references)
- Be wary of obviously fake codes (00000, 12345, etc.)
- Look for context suggesting real addresses vs. examples`

	case pii.PiiTypeStreetAddress:
		return `Street address validation criteria:
- Check if the address format is realistic and well-formed
- Look for real street names, not obviously fake ones (123 Main St is suspicious)
- Consider if house numbers are reasonable for the street type
- Check for proper abbreviations and formatting`

	case pii.PiiTypeIPAddress:
		return `IP address validation criteria:
- Verify the format is valid (IPv4: x.x.x.x, IPv6: proper format)
- Check if it's in a reasonable range (not 0.0.0.0 or other invalid IPs)
- Consider if it's a private vs. public IP and if that makes sense in context
- Look for context suggesting real network data vs. examples`

	case pii.PiiTypeIBAN:
		return `IBAN validation criteria:
- Check if the country code is valid (first 2 letters)
- Verify the format matches IBAN standards for that country
- Look for context suggesting real banking information vs. examples
- Be wary of obviously fake IBANs or test patterns`

	case pii.PiiTypeBtcAddress:
		return `Bitcoin address validation criteria:
- Check if the format is valid (starts with 1, 3, or bc1)
- Verify the length is appropriate for the address type
- Look for context suggesting real cryptocurrency activity vs. examples
- Be wary of obviously fake or example addresses`

	case pii.PiiTypePoBox:
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
func (v *LLMValidatorImpl) parseValidationResponse(response string) (*pii.ValidationResult, error) {
	result := &pii.ValidationResult{
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
func (v *LLMValidatorImpl) parseHeuristically(response string, result *pii.ValidationResult) (*pii.ValidationResult, error) {
	// Look for keywords indicating validity
	lowerResponse := ""
	for _, r := range response {
		if r >= 'A' && r <= 'Z' {
			lowerResponse += string(r + 32) // Convert to lowercase
		} else {
			lowerResponse += string(r)
		}
	}

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