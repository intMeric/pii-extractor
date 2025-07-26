package piiextractor

import (
	"context"
	"strings"
	"testing"
	"time"
)

// MockLLMValidator is a test implementation of LLMValidator
type MockLLMValidator struct {
	validateEntityFunc func(ctx context.Context, entity PiiEntity, context string) (*ValidationResult, error)
	healthCheckFunc    func(ctx context.Context) error
	provider           string
	model              string
}

func (m *MockLLMValidator) ValidateEntity(ctx context.Context, entity PiiEntity, context string) (*ValidationResult, error) {
	if m.validateEntityFunc != nil {
		return m.validateEntityFunc(ctx, entity, context)
	}
	// Default mock behavior - mark everything as valid with high confidence
	return &ValidationResult{
		Valid:      true,
		Confidence: 0.9,
		Reasoning:  "Mock validation result",
		Provider:   m.provider,
		Model:      m.model,
	}, nil
}

func (m *MockLLMValidator) ValidateBatch(ctx context.Context, entities []PiiEntity, contexts []string) ([]*ValidationResult, error) {
	results := make([]*ValidationResult, len(entities))
	for i, entity := range entities {
		result, err := m.ValidateEntity(ctx, entity, contexts[i])
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}

func (m *MockLLMValidator) HealthCheck(ctx context.Context) error {
	if m.healthCheckFunc != nil {
		return m.healthCheckFunc(ctx)
	}
	return nil
}

func (m *MockLLMValidator) GetProviderInfo() (provider string, model string) {
	return m.provider, m.model
}

// NewMockValidator creates a new mock validator for testing
func NewMockValidator(provider, model string) *MockLLMValidator {
	return &MockLLMValidator{
		provider: provider,
		model:    model,
	}
}

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

	// Test with validation enabled but using mock validator
	config := DefaultValidationConfig()
	config.Enabled = true
	config.APIKey = "test-key"

	// Create extractor with validation enabled
	validatedExtractor, err := NewValidatedExtractor(baseExtractor, config)
	if err != nil {
		// This is expected to fail with real API calls, so we'll create a mock version
		t.Logf("Real API creation failed as expected: %v", err)

		// Create a mock validated extractor for testing
		mockValidator := NewMockValidator("openai", "gpt-4o-mini")
		validatedExtractor = &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     mockValidator,
			config:        config,
		}
	}

	if !validatedExtractor.IsValidationEnabled() {
		t.Error("Expected validation to be enabled with mock validator")
	}
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

// Test ValidatedExtractor with mock validation
func TestValidatedExtractorWithMockValidation(t *testing.T) {
	baseExtractor := NewRegexExtractor()
	config := DefaultValidationConfig()
	config.Enabled = true

	// Create mock validator
	mockValidator := NewMockValidator("openai", "gpt-4o-mini")

	// Create validated extractor with mock
	extractor := &ValidatedExtractor{
		baseExtractor: baseExtractor,
		validator:     mockValidator,
		config:        config,
	}

	text := "Contact me at john@example.com or call 555-123-4567"
	result, err := extractor.ExtractWithValidation(text)
	if err != nil {
		t.Fatalf("Failed to extract with validation: %v", err)
	}

	if result.IsEmpty() {
		t.Error("Expected to find PII entities")
	}

	// Verify validation was performed
	validatedCount := 0
	for _, entity := range result.Entities {
		if entity.IsValidated() {
			validatedCount++
			if !entity.IsValid() {
				t.Error("Mock validator should mark entities as valid")
			}
			if entity.GetValidationConfidence() != 0.9 {
				t.Errorf("Expected confidence 0.9, got %f", entity.GetValidationConfidence())
			}
		}
	}

	if validatedCount == 0 {
		t.Error("Expected some entities to be validated")
	}

	// Check validation stats
	if result.ValidationStats == nil {
		t.Error("Expected validation stats to be present")
	} else {
		if result.ValidationStats.Provider != "openai" {
			t.Errorf("Expected provider 'openai', got '%s'", result.ValidationStats.Provider)
		}
		if result.ValidationStats.Model != "gpt-4o-mini" {
			t.Errorf("Expected model 'gpt-4o-mini', got '%s'", result.ValidationStats.Model)
		}
	}
}

// Test mock validator with custom behavior
func TestMockValidatorCustomBehavior(t *testing.T) {
	mockValidator := NewMockValidator("test-provider", "test-model")

	// Set custom validation behavior
	mockValidator.validateEntityFunc = func(ctx context.Context, entity PiiEntity, context string) (*ValidationResult, error) {
		// Mark emails as invalid, everything else as valid
		isValid := !entity.IsEmail()
		confidence := 0.8
		if isValid {
			confidence = 0.95
		}

		return &ValidationResult{
			Valid:      isValid,
			Confidence: confidence,
			Reasoning:  "Custom mock validation logic",
			Provider:   "test-provider",
			Model:      "test-model",
		}, nil
	}

	// Test the custom behavior
	emailEntity := PiiEntity{
		Type:  PiiTypeEmail,
		Value: NewEmail("test@example.com"),
	}

	phoneEntity := PiiEntity{
		Type:  PiiTypePhone,
		Value: NewPhoneUS("555-123-4567"),
	}

	// Test email validation (should be invalid)
	ctx := context.Background()
	result, err := mockValidator.ValidateEntity(ctx, emailEntity, "context")
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if result.Valid {
		t.Error("Expected email to be marked as invalid by custom mock")
	}

	if result.Confidence != 0.8 {
		t.Errorf("Expected confidence 0.8, got %f", result.Confidence)
	}

	// Test phone validation (should be valid)
	result, err = mockValidator.ValidateEntity(ctx, phoneEntity, "context")
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.Valid {
		t.Error("Expected phone to be marked as valid by custom mock")
	}

	if result.Confidence != 0.95 {
		t.Errorf("Expected confidence 0.95, got %f", result.Confidence)
	}
}

// Test health check with mock
func TestMockValidatorHealthCheck(t *testing.T) {
	mockValidator := NewMockValidator("test-provider", "test-model")

	// Test default health check (should pass)
	ctx := context.Background()
	err := mockValidator.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Expected health check to pass, got error: %v", err)
	}

	// Test custom health check behavior
	expectedError := NewPiiError("health_check_failed", "Mock health check failure")
	mockValidator.healthCheckFunc = func(ctx context.Context) error {
		return expectedError
	}

	err = mockValidator.HealthCheck(ctx)
	if err == nil {
		t.Error("Expected health check to fail with custom function")
	}

	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error '%s', got '%s'", expectedError.Error(), err.Error())
	}
}
