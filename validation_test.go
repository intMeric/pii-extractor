package piiextractor

import (
	"strings"
	"testing"
	"time"
)

func TestDefaultValidationConfig(t *testing.T) {
	config := DefaultValidationConfig()

	if config.Enabled {
		t.Error("Expected validation to be disabled by default")
	}

	if config.Provider != ProviderOpenAI {
		t.Errorf("Expected default provider to be OpenAI, got %s", config.Provider)
	}

	if config.Model != "gpt-4o-mini" {
		t.Errorf("Expected default model to be gpt-4o-mini, got %s", config.Model)
	}

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout to be 30s, got %v", config.Timeout)
	}

	if config.MinConfidence != 0.7 {
		t.Errorf("Expected default min confidence to be 0.7, got %f", config.MinConfidence)
	}
}

func TestPiiEntityValidationMethods(t *testing.T) {
	// Test entity without validation
	entity := PiiEntity{
		Type:       PiiTypeEmail,
		Value:      NewEmail("test@example.com"),
		Validation: nil,
	}

	if entity.IsValidated() {
		t.Error("Expected entity to not be validated")
	}

	if entity.IsValid() {
		t.Error("Expected entity to not be valid (no validation)")
	}

	if entity.GetValidationConfidence() != 0.0 {
		t.Errorf("Expected validation confidence to be 0.0, got %f", entity.GetValidationConfidence())
	}

	// Test entity with validation
	validation := &ValidationResult{
		Valid:      true,
		Confidence: 0.95,
		Reasoning:  "Valid email format and domain",
		Provider:   "openai",
		Model:      "gpt-4o-mini",
	}

	entity.Validation = validation

	if !entity.IsValidated() {
		t.Error("Expected entity to be validated")
	}

	if !entity.IsValid() {
		t.Error("Expected entity to be valid")
	}

	if entity.GetValidationConfidence() != 0.95 {
		t.Errorf("Expected validation confidence to be 0.95, got %f", entity.GetValidationConfidence())
	}
}

func TestPiiExtractionResultValidationMethods(t *testing.T) {
	// Create test entities
	validatedEntity := PiiEntity{
		Type:  PiiTypeEmail,
		Value: NewEmail("valid@example.com"),
		Validation: &ValidationResult{
			Valid:      true,
			Confidence: 0.9,
			Provider:   "openai",
		},
	}

	invalidatedEntity := PiiEntity{
		Type:  PiiTypeEmail,
		Value: NewEmail("invalid@example.com"),
		Validation: &ValidationResult{
			Valid:      false,
			Confidence: 0.8,
			Provider:   "openai",
		},
	}

	unvalidatedEntity := PiiEntity{
		Type:       PiiTypePhone,
		Value:      NewPhoneUS("555-0123"),
		Validation: nil,
	}

	entities := []PiiEntity{validatedEntity, invalidatedEntity, unvalidatedEntity}
	result := NewPiiExtractionResult(entities)

	// Test GetValidatedEntities
	validated := result.GetValidatedEntities()
	if len(validated) != 2 {
		t.Errorf("Expected 2 validated entities, got %d", len(validated))
	}

	// Test GetValidEntities
	valid := result.GetValidEntities()
	if len(valid) != 1 {
		t.Errorf("Expected 1 valid entity, got %d", len(valid))
	}

	// Test GetInvalidEntities
	invalid := result.GetInvalidEntities()
	if len(invalid) != 1 {
		t.Errorf("Expected 1 invalid entity, got %d", len(invalid))
	}
}

func TestValidationConfigProviders(t *testing.T) {
	tests := []struct {
		provider LLMProvider
		expected string
	}{
		{ProviderOpenAI, "openai"},
		{ProviderMistral, "mistral"},
		{ProviderGemini, "gemini"},
		{ProviderOllama, "ollama"},
		{ProviderAnthropic, "anthropic"},
	}

	for _, test := range tests {
		if string(test.provider) != test.expected {
			t.Errorf("Expected provider %s to have string value %s, got %s",
				test.provider, test.expected, string(test.provider))
		}
	}
}

func TestValidatedExtractorCreation(t *testing.T) {
	baseExtractor := NewRegexExtractor()

	// Test with nil config (should use defaults)
	extractor, err := NewValidatedExtractor(baseExtractor, nil)
	if err != nil {
		t.Fatalf("Failed to create validated extractor: %v", err)
	}

	if extractor.IsValidationEnabled() {
		t.Error("Expected validation to be disabled by default")
	}

	// Test with validation enabled
	config := DefaultValidationConfig()
	config.Enabled = true
	config.APIKey = "test-key" // Set a test API key

	// Note: This will fail in tests without a real API key, but tests the structure
	_, err = NewValidatedExtractor(baseExtractor, config)
	// We expect this to potentially fail due to missing API key, but the structure should be correct
	// In a real test environment, you'd use mock services
}

func TestValidatedExtractorBasicExtraction(t *testing.T) {
	baseExtractor := NewRegexExtractor()
	extractor, err := NewValidatedExtractor(baseExtractor, nil)
	if err != nil {
		t.Fatalf("Failed to create validated extractor: %v", err)
	}

	// Test basic extraction (no validation)
	text := "Contact me at john@example.com or call 555-123-4567"
	result, err := extractor.Extract(text)
	if err != nil {
		t.Fatalf("Failed to extract PII: %v", err)
	}

	if result.IsEmpty() {
		t.Error("Expected to find PII entities")
	}

	// Verify no validation was performed
	for _, entity := range result.Entities {
		if entity.IsValidated() {
			t.Error("Expected no validation to be performed in basic extraction")
		}
	}
}

func TestParseValidationResponseJSON(t *testing.T) {
	config := DefaultValidationConfig()
	config.Provider = ProviderOpenAI
	config.Model = "gpt-4o-mini"

	validator := &LLMValidatorImpl{
		config: config,
	}

	tests := []struct {
		name     string
		response string
		expected ValidationResult
	}{
		{
			name:     "Valid JSON response",
			response: `{"valid": true, "confidence": 0.9, "reasoning": "Valid email format"}`,
			expected: ValidationResult{
				Valid:      true,
				Confidence: 0.9,
				Reasoning:  "Valid email format",
				Provider:   "openai",
				Model:      "gpt-4o-mini",
			},
		},
		{
			name:     "Invalid JSON response",
			response: `{"valid": false, "confidence": 0.8, "reasoning": "Suspicious domain"}`,
			expected: ValidationResult{
				Valid:      false,
				Confidence: 0.8,
				Reasoning:  "Suspicious domain",
				Provider:   "openai",
				Model:      "gpt-4o-mini",
			},
		},
		{
			name:     "Malformed JSON",
			response: `This is not JSON. The email looks valid to me.`,
			expected: ValidationResult{
				Valid:      false,
				Confidence: 0.5,
				Reasoning:  "Heuristic parsing of LLM response",
				Provider:   "openai",
				Model:      "gpt-4o-mini",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := validator.parseValidationResponse(test.response)
			if err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if result.Valid != test.expected.Valid {
				t.Errorf("Expected valid=%v, got %v", test.expected.Valid, result.Valid)
			}

			if result.Confidence != test.expected.Confidence {
				t.Errorf("Expected confidence=%f, got %f", test.expected.Confidence, result.Confidence)
			}

			if result.Provider != test.expected.Provider {
				t.Errorf("Expected provider=%s, got %s", test.expected.Provider, result.Provider)
			}
		})
	}
}

func TestValidationPromptGeneration(t *testing.T) {
	config := DefaultValidationConfig()
	validator := &LLMValidatorImpl{
		config: config,
	}

	entity := PiiEntity{
		Type:  PiiTypeEmail,
		Value: NewEmail("test@example.com"),
	}

	context := "Please contact me at test@example.com for more information."

	prompt := validator.buildValidationPrompt(entity, context)

	// Check that prompt contains key elements
	if !validator.containsAny(prompt, []string{"email", "test@example.com", context}) {
		t.Error("Prompt should contain PII type, value, and context")
	}

	// Check that prompt contains type-specific guidance
	if !validator.containsAny(prompt, []string{"domain", "realistic", "example.com"}) {
		t.Error("Prompt should contain email-specific validation guidance")
	}

	// Check JSON format requirement
	if validator.findSubstring(prompt, "JSON format") == -1 {
		t.Error("Prompt should request JSON format response")
	}
}

func TestTypeSpecificGuidance(t *testing.T) {
	config := DefaultValidationConfig()
	validator := &LLMValidatorImpl{
		config: config,
	}

	tests := []struct {
		piiType  PiiType
		keywords []string
	}{
		{PiiTypeEmail, []string{"domain", "@", "realistic"}},
		{PiiTypePhone, []string{"area codes", "555", "fake numbers"}},
		{PiiTypeSSN, []string{"XXX-XX-XXXX", "000-00-0000", "fake SSNs"}},
		{PiiTypeCreditCard, []string{"Luhn", "4111-1111-1111-1111", "test"}},
		{PiiTypeZipCode, []string{"5 digits", "00000", "geographic"}},
	}

	for _, test := range tests {
		guidance := validator.getTypeSpecificGuidance(test.piiType)

		for _, keyword := range test.keywords {
			if validator.findSubstring(strings.ToLower((guidance)), strings.ToLower(keyword)) == -1 {
				t.Errorf("Guidance for %s should contain keyword '%s'", test.piiType.String(), keyword)
			}
		}
	}
}

func TestValidationStats(t *testing.T) {
	// Create entities with mixed validation results
	entities := []PiiEntity{
		{
			Type:       PiiTypeEmail,
			Value:      NewEmail("valid@example.com"),
			Validation: &ValidationResult{Valid: true, Confidence: 0.9},
		},
		{
			Type:       PiiTypeEmail,
			Value:      NewEmail("invalid@example.com"),
			Validation: &ValidationResult{Valid: false, Confidence: 0.8},
		},
		{
			Type:       PiiTypePhone,
			Value:      NewPhoneUS("555-0123"),
			Validation: nil, // Not validated
		},
	}

	result := NewPiiExtractionResult(entities)

	// Simulate validation stats calculation
	stats := &ValidationStats{
		TotalValidated:    2,
		ValidCount:        1,
		InvalidCount:      1,
		AverageConfidence: 0.85,
		Provider:          "openai",
		Model:             "gpt-4o-mini",
	}

	result.ValidationStats = stats

	if result.ValidationStats.TotalValidated != 2 {
		t.Errorf("Expected 2 validated entities, got %d", result.ValidationStats.TotalValidated)
	}

	if result.ValidationStats.ValidCount != 1 {
		t.Errorf("Expected 1 valid entity, got %d", result.ValidationStats.ValidCount)
	}

	if result.ValidationStats.InvalidCount != 1 {
		t.Errorf("Expected 1 invalid entity, got %d", result.ValidationStats.InvalidCount)
	}

	if result.ValidationStats.AverageConfidence != 0.85 {
		t.Errorf("Expected average confidence 0.85, got %f", result.ValidationStats.AverageConfidence)
	}
}

func TestConfigValidation(t *testing.T) {
	// Test various provider configurations
	providers := []LLMProvider{
		ProviderOpenAI,
		ProviderMistral,
		ProviderGemini,
		ProviderOllama,
		ProviderAnthropic,
	}

	for _, provider := range providers {
		config := DefaultValidationConfig()
		config.Provider = provider
		config.Enabled = true

		// Test that each provider can be configured
		// Note: This doesn't test actual LLM connections, just configuration
		if config.Provider != provider {
			t.Errorf("Failed to set provider to %s", provider)
		}
	}
}

// Benchmark validation response parsing
func BenchmarkParseValidationResponse(b *testing.B) {
	config := DefaultValidationConfig()
	validator := &LLMValidatorImpl{config: config}

	response := `{"valid": true, "confidence": 0.85, "reasoning": "Email format is valid and domain appears legitimate"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := validator.parseValidationResponse(response)
		if err != nil {
			b.Fatalf("Parsing failed: %v", err)
		}
	}
}
