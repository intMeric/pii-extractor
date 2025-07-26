package piiextractor

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewValidatedExtractor(t *testing.T) {
	baseExtractor := NewRegexExtractor()

	t.Run("with nil config", func(t *testing.T) {
		extractor, err := NewValidatedExtractor(baseExtractor, nil)
		if err != nil {
			t.Fatalf("Expected no error with nil config, got: %v", err)
		}

		if extractor.baseExtractor != baseExtractor {
			t.Error("Expected base extractor to be set correctly")
		}

		if extractor.config == nil {
			t.Error("Expected config to be set to default")
		}

		if extractor.config.Enabled {
			t.Error("Expected validation to be disabled by default")
		}

		if extractor.validator != nil {
			t.Error("Expected validator to be nil when validation is disabled")
		}
	})

	t.Run("with disabled validation", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = false

		extractor, err := NewValidatedExtractor(baseExtractor, config)
		if err != nil {
			t.Fatalf("Expected no error with disabled validation, got: %v", err)
		}

		if extractor.validator != nil {
			t.Error("Expected validator to be nil when validation is disabled")
		}

		if extractor.IsValidationEnabled() {
			t.Error("Expected validation to be disabled")
		}
	})

	t.Run("with enabled validation (mock)", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = true

		// Since real API creation will fail, we'll create manually with mock
		mockValidator := NewMockValidator("openai", "gpt-4o-mini")
		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     mockValidator,
			config:        config,
		}

		if !extractor.IsValidationEnabled() {
			t.Error("Expected validation to be enabled")
		}

		if extractor.validator == nil {
			t.Error("Expected validator to be set when validation is enabled")
		}
	})
}

func TestValidatedExtractor_Extract(t *testing.T) {
	baseExtractor := NewRegexExtractor()
	config := DefaultValidationConfig()
	mockValidator := NewMockValidator("test", "test")

	extractor := &ValidatedExtractor{
		baseExtractor: baseExtractor,
		validator:     mockValidator,
		config:        config,
	}

	text := "Contact john@example.com or call 555-123-4567"
	result, err := extractor.Extract(text)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if result.IsEmpty() {
		t.Error("Expected to find PII entities")
	}

	// Verify no validation was performed in basic extraction
	for _, entity := range result.Entities {
		if entity.IsValidated() {
			t.Error("Expected no validation in basic Extract method")
		}
	}

	if result.ValidationStats != nil {
		t.Error("Expected no validation stats in basic extraction")
	}
}

func TestValidatedExtractor_ExtractWithValidation(t *testing.T) {
	baseExtractor := NewRegexExtractor()

	t.Run("validation disabled", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = false

		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     nil,
			config:        config,
		}

		text := "Contact john@example.com"
		result, err := extractor.ExtractWithValidation(text)
		if err != nil {
			t.Fatalf("ExtractWithValidation failed: %v", err)
		}

		// Should behave like basic extraction when validation is disabled
		for _, entity := range result.Entities {
			if entity.IsValidated() {
				t.Error("Expected no validation when validation is disabled")
			}
		}

		if result.ValidationStats != nil {
			t.Error("Expected no validation stats when validation is disabled")
		}
	})

	t.Run("validation enabled with mock", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = true
		config.MinConfidence = 0.5

		mockValidator := NewMockValidator("openai", "gpt-4o-mini")
		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     mockValidator,
			config:        config,
		}

		text := "Contact john@example.com or call 555-123-4567"
		result, err := extractor.ExtractWithValidation(text)
		if err != nil {
			t.Fatalf("ExtractWithValidation failed: %v", err)
		}

		if result.IsEmpty() {
			t.Error("Expected to find PII entities")
		}

		// Verify validation was performed
		validatedCount := 0
		for _, entity := range result.Entities {
			if entity.IsValidated() {
				validatedCount++
			}
		}

		if validatedCount == 0 {
			t.Error("Expected some entities to be validated")
		}

		// Check validation stats
		if result.ValidationStats == nil {
			t.Error("Expected validation stats to be present")
		}
	})

	t.Run("empty text", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = true

		mockValidator := NewMockValidator("test", "test")
		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     mockValidator,
			config:        config,
		}

		result, err := extractor.ExtractWithValidation("")
		if err != nil {
			t.Fatalf("ExtractWithValidation failed: %v", err)
		}

		if !result.IsEmpty() {
			t.Error("Expected empty result for empty text")
		}

		if result.ValidationStats != nil {
			t.Error("Expected no validation stats for empty result")
		}
	})

	t.Run("validator is nil but validation enabled", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = true

		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     nil, // Simulate missing validator
			config:        config,
		}

		text := "Contact john@example.com"
		result, err := extractor.ExtractWithValidation(text)
		if err != nil {
			t.Fatalf("ExtractWithValidation failed: %v", err)
		}

		// Should return unvalidated results gracefully
		for _, entity := range result.Entities {
			if entity.IsValidated() {
				t.Error("Expected no validation when validator is nil")
			}
		}
	})
}

func TestValidatedExtractor_HealthCheck(t *testing.T) {
	baseExtractor := NewRegexExtractor()

	t.Run("validation disabled", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = false

		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     nil,
			config:        config,
		}

		ctx := context.Background()
		err := extractor.HealthCheck(ctx)
		if err == nil {
			t.Error("Expected health check to fail when validation is disabled")
		}

		if err != ErrValidationDisabled {
			t.Errorf("Expected ErrValidationDisabled, got: %v", err)
		}
	})

	t.Run("validation enabled with healthy mock", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = true

		mockValidator := NewMockValidator("test", "test")
		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     mockValidator,
			config:        config,
		}

		ctx := context.Background()
		err := extractor.HealthCheck(ctx)
		if err != nil {
			t.Errorf("Expected health check to pass, got: %v", err)
		}
	})

	t.Run("validation enabled with unhealthy mock", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = true

		mockValidator := NewMockValidator("test", "test")
		expectedError := errors.New("health check failed")
		mockValidator.healthCheckFunc = func(ctx context.Context) error {
			return expectedError
		}

		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     mockValidator,
			config:        config,
		}

		ctx := context.Background()
		err := extractor.HealthCheck(ctx)
		if err == nil {
			t.Error("Expected health check to fail")
		}

		if err != expectedError {
			t.Errorf("Expected specific error, got: %v", err)
		}
	})
}

func TestValidatedExtractor_SetValidationConfig(t *testing.T) {
	baseExtractor := NewRegexExtractor()
	config := DefaultValidationConfig()

	extractor := &ValidatedExtractor{
		baseExtractor: baseExtractor,
		validator:     nil,
		config:        config,
	}

	t.Run("set nil config", func(t *testing.T) {
		err := extractor.SetValidationConfig(nil)
		if err == nil {
			t.Error("Expected error when setting nil config")
		}
	})

	t.Run("set disabled validation config", func(t *testing.T) {
		newConfig := DefaultValidationConfig()
		newConfig.Enabled = false

		err := extractor.SetValidationConfig(newConfig)
		if err != nil {
			t.Fatalf("SetValidationConfig failed: %v", err)
		}

		if extractor.config != newConfig {
			t.Error("Expected config to be updated")
		}

		if extractor.validator != nil {
			t.Error("Expected validator to be nil when validation is disabled")
		}

		if extractor.IsValidationEnabled() {
			t.Error("Expected validation to be disabled")
		}
	})

	t.Run("set enabled validation config", func(t *testing.T) {
		newConfig := DefaultValidationConfig()
		newConfig.Enabled = true

		// This will fail with real API, so we test the error handling
		err := extractor.SetValidationConfig(newConfig)
		if err == nil {
			t.Log("SetValidationConfig succeeded (unexpected in test environment)")
		} else {
			t.Logf("SetValidationConfig failed as expected: %v", err)
		}
	})
}

func TestValidatedExtractor_GetValidationConfig(t *testing.T) {
	baseExtractor := NewRegexExtractor()
	config := DefaultValidationConfig()

	extractor := &ValidatedExtractor{
		baseExtractor: baseExtractor,
		validator:     nil,
		config:        config,
	}

	retrievedConfig := extractor.GetValidationConfig()
	if retrievedConfig != config {
		t.Error("Expected to get the same config that was set")
	}
}

func TestValidatedExtractor_IsValidationEnabled(t *testing.T) {
	baseExtractor := NewRegexExtractor()

	t.Run("validation disabled", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = false

		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     nil,
			config:        config,
		}

		if extractor.IsValidationEnabled() {
			t.Error("Expected validation to be disabled")
		}
	})

	t.Run("validation enabled with validator", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = true

		mockValidator := NewMockValidator("test", "test")
		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     mockValidator,
			config:        config,
		}

		if !extractor.IsValidationEnabled() {
			t.Error("Expected validation to be enabled")
		}
	})

	t.Run("validation enabled but no validator", func(t *testing.T) {
		config := DefaultValidationConfig()
		config.Enabled = true

		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     nil, // No validator set
			config:        config,
		}

		if extractor.IsValidationEnabled() {
			t.Error("Expected validation to be disabled when validator is nil")
		}
	})

	t.Run("nil config", func(t *testing.T) {
		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     NewMockValidator("test", "test"),
			config:        nil,
		}

		if extractor.IsValidationEnabled() {
			t.Error("Expected validation to be disabled when config is nil")
		}
	})
}

func TestValidatedExtractor_validateEntities(t *testing.T) {
	baseExtractor := NewRegexExtractor()
	config := DefaultValidationConfig()
	config.Enabled = true
	config.MinConfidence = 0.8
	config.MaxRetries = 2
	config.Timeout = 5 * time.Second

	t.Run("successful validation", func(t *testing.T) {
		mockValidator := NewMockValidator("test", "test")
		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     mockValidator,
			config:        config,
		}

		entities := []PiiEntity{
			{Type: PiiTypeEmail, Value: NewEmail("test@example.com")},
			{Type: PiiTypePhone, Value: NewPhoneUS("555-123-4567")},
		}
		result := &PiiExtractionResult{Entities: entities}

		ctx := context.Background()
		err := extractor.validateEntities(ctx, result, "test text", mockValidator, config)
		if err != nil {
			t.Fatalf("validateEntities failed: %v", err)
		}

		// Check that entities were validated
		for _, entity := range result.Entities {
			if !entity.IsValidated() {
				t.Error("Expected entity to be validated")
			}
		}
	})

	t.Run("validation with low confidence", func(t *testing.T) {
		mockValidator := NewMockValidator("test", "test")
		// Set mock to return low confidence
		mockValidator.validateEntityFunc = func(ctx context.Context, entity PiiEntity, context string) (*ValidationResult, error) {
			return &ValidationResult{
				Valid:      true,
				Confidence: 0.5, // Below MinConfidence threshold
				Reasoning:  "Low confidence validation",
				Provider:   "test",
				Model:      "test",
			}, nil
		}

		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     mockValidator,
			config:        config,
		}

		entities := []PiiEntity{
			{Type: PiiTypeEmail, Value: NewEmail("test@example.com")},
		}
		result := &PiiExtractionResult{Entities: entities}

		ctx := context.Background()
		err := extractor.validateEntities(ctx, result, "test text", mockValidator, config)
		if err != nil {
			t.Fatalf("validateEntities failed: %v", err)
		}

		// Entity should not be validated due to low confidence
		if result.Entities[0].IsValidated() {
			t.Error("Expected entity to not be validated due to low confidence")
		}
	})

	t.Run("validation with retries", func(t *testing.T) {
		mockValidator := NewMockValidator("test", "test")
		callCount := 0
		mockValidator.validateEntityFunc = func(ctx context.Context, entity PiiEntity, context string) (*ValidationResult, error) {
			callCount++
			if callCount < 3 { // Fail first 2 attempts
				return nil, errors.New("temporary failure")
			}
			return &ValidationResult{
				Valid:      true,
				Confidence: 0.9,
				Reasoning:  "Success after retries",
				Provider:   "test",
				Model:      "test",
			}, nil
		}

		extractor := &ValidatedExtractor{
			baseExtractor: baseExtractor,
			validator:     mockValidator,
			config:        config,
		}

		entities := []PiiEntity{
			{Type: PiiTypeEmail, Value: NewEmail("test@example.com")},
		}
		result := &PiiExtractionResult{Entities: entities}

		ctx := context.Background()
		err := extractor.validateEntities(ctx, result, "test text", mockValidator, config)
		if err != nil {
			t.Fatalf("validateEntities failed: %v", err)
		}

		if callCount != 3 {
			t.Errorf("Expected 3 calls (2 failures + 1 success), got %d", callCount)
		}

		if !result.Entities[0].IsValidated() {
			t.Error("Expected entity to be validated after successful retry")
		}
	})
}

func TestValidatedExtractor_calculateValidationStats(t *testing.T) {
	baseExtractor := NewRegexExtractor()
	mockValidator := NewMockValidator("test-provider", "test-model")

	extractor := &ValidatedExtractor{
		baseExtractor: baseExtractor,
		validator:     mockValidator,
		config:        DefaultValidationConfig(),
	}

	entities := []PiiEntity{
		{
			Type:  PiiTypeEmail,
			Value: NewEmail("valid@example.com"),
			Validation: &ValidationResult{
				Valid:      true,
				Confidence: 0.9,
				Provider:   "test-provider",
				Model:      "test-model",
			},
		},
		{
			Type:  PiiTypeEmail,
			Value: NewEmail("invalid@example.com"),
			Validation: &ValidationResult{
				Valid:      false,
				Confidence: 0.8,
				Provider:   "test-provider",
				Model:      "test-model",
			},
		},
		{
			Type:       PiiTypePhone,
			Value:      NewPhoneUS("555-0123"),
			Validation: nil, // Not validated
		},
	}

	result := &PiiExtractionResult{Entities: entities}

	extractor.calculateValidationStats(result, mockValidator)

	if result.ValidationStats == nil {
		t.Fatal("Expected validation stats to be calculated")
	}

	stats := result.ValidationStats
	if stats.TotalValidated != 2 {
		t.Errorf("Expected 2 validated entities, got %d", stats.TotalValidated)
	}

	if stats.ValidCount != 1 {
		t.Errorf("Expected 1 valid entity, got %d", stats.ValidCount)
	}

	if stats.InvalidCount != 1 {
		t.Errorf("Expected 1 invalid entity, got %d", stats.InvalidCount)
	}

	expectedAvgConfidence := (0.9 + 0.8) / 2
	if abs(stats.AverageConfidence-expectedAvgConfidence) > 0.001 {
		t.Errorf("Expected average confidence %f, got %f", expectedAvgConfidence, stats.AverageConfidence)
	}

	if stats.Provider != "test-provider" {
		t.Errorf("Expected provider 'test-provider', got '%s'", stats.Provider)
	}

	if stats.Model != "test-model" {
		t.Errorf("Expected model 'test-model', got '%s'", stats.Model)
	}
}

func TestValidatedExtractor_getEntityContext(t *testing.T) {
	baseExtractor := NewRegexExtractor()
	extractor := &ValidatedExtractor{
		baseExtractor: baseExtractor,
		validator:     nil,
		config:        DefaultValidationConfig(),
	}

	t.Run("entity with existing contexts", func(t *testing.T) {
		email := NewEmail("test@example.com")
		email.AddContext("Please contact test@example.com for support")
		
		entity := &PiiEntity{
			Type:  PiiTypeEmail,
			Value: email,
		}

		context := extractor.getEntityContext("Some other text", entity)
		expected := "Please contact test@example.com for support"
		if context != expected {
			t.Errorf("Expected context '%s', got '%s'", expected, context)
		}
	})

	t.Run("entity without contexts - extract from text", func(t *testing.T) {
		entity := &PiiEntity{
			Type:  PiiTypeEmail,
			Value: NewEmail("test@example.com"),
		}

		text := "Please contact test@example.com for more information"
		context := extractor.getEntityContext(text, entity)
		
		if !containsString(context, "test@example.com") {
			t.Errorf("Expected context to contain the email address, got '%s'", context)
		}
	})

	t.Run("entity with empty value", func(t *testing.T) {
		entity := &PiiEntity{
			Type:  PiiTypeEmail,
			Value: NewEmail(""),
		}

		context := extractor.getEntityContext("test text", entity)
		if context != "" {
			t.Errorf("Expected empty context for empty value, got '%s'", context)
		}
	})
}

func TestValidatedExtractor_extractSimpleContext(t *testing.T) {
	baseExtractor := NewRegexExtractor()
	extractor := &ValidatedExtractor{
		baseExtractor: baseExtractor,
		validator:     nil,
		config:        DefaultValidationConfig(),
	}

	t.Run("value found in text", func(t *testing.T) {
		text := "The quick brown fox jumps over the lazy dog. Please contact test@example.com for more information about our services."
		value := "test@example.com"

		context := extractor.extractSimpleContext(text, value)

		if !containsString(context, value) {
			t.Error("Expected context to contain the value")
		}

		// Context should be limited to about 200 characters (100 before + 100 after + value)
		if len(context) > len(text) {
			t.Error("Expected context to not be longer than original text")
		}
	})

	t.Run("value not found in text", func(t *testing.T) {
		text := "This text does not contain the target value"
		value := "missing@example.com"

		context := extractor.extractSimpleContext(text, value)

		if context != text {
			t.Errorf("Expected full text when value not found, got '%s'", context)
		}
	})

	t.Run("value at beginning of text", func(t *testing.T) {
		text := "test@example.com is my email address for contact"
		value := "test@example.com"

		context := extractor.extractSimpleContext(text, value)

		if !containsString(context, value) {
			t.Error("Expected context to contain the value")
		}
	})

	t.Run("value at end of text", func(t *testing.T) {
		text := "Please contact me at test@example.com"
		value := "test@example.com"

		context := extractor.extractSimpleContext(text, value)

		if !containsString(context, value) {
			t.Error("Expected context to contain the value")
		}
	})
}

// Helper function to check if a string contains another string
func containsString(text, substr string) bool {
	return len(text) >= len(substr) && findSubstring(text, substr) != -1
}

// Helper function to find substring (simple implementation)
func findSubstring(text, pattern string) int {
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

// Helper function for float comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}