package piiextractor

import (
	"context"
	"time"
)

// ValidatedExtractor combines regex-based extraction with LLM validation
type ValidatedExtractor struct {
	baseExtractor PiiExtractor
	validator     LLMValidator
	config        *ValidationConfig
}

// NewValidatedExtractor creates a new validated extractor
func NewValidatedExtractor(baseExtractor PiiExtractor, config *ValidationConfig) (*ValidatedExtractor, error) {
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
		baseExtractor: baseExtractor,
		validator:     validator,
		config:        config,
	}, nil
}

// Extract performs basic extraction without validation (implements PiiExtractor)
func (v *ValidatedExtractor) Extract(text string) (*PiiExtractionResult, error) {
	return v.baseExtractor.Extract(text)
}

// ExtractWithValidation performs extraction with LLM validation (implements ValidatedPiiExtractor)
// Uses the configuration provided during creation
func (v *ValidatedExtractor) ExtractWithValidation(text string) (*PiiExtractionResult, error) {
	config := v.config

	// If validation is disabled, just do regular extraction
	if !config.Enabled {
		return v.baseExtractor.Extract(text)
	}

	// Perform initial regex-based extraction
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

// validateEntities validates all entities in the result
func (v *ValidatedExtractor) validateEntities(ctx context.Context, result *PiiExtractionResult, originalText string, validator LLMValidator, config *ValidationConfig) error {
	for i := range result.Entities {
		entity := &result.Entities[i]

		// Get context for this entity
		context := v.getEntityContext(originalText, entity)

		// Validate with retries
		var validation *ValidationResult
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
func (v *ValidatedExtractor) getEntityContext(text string, entity *PiiEntity) string {
	// For now, return the first context from the entity
	// In a more sophisticated implementation, we could find the entity in the text
	// and extract more precise context
	contexts := entity.GetContexts()
	if len(contexts) > 0 {
		return contexts[0]
	}

	// Fallback: try to find the entity value in the text and extract context
	value := entity.GetValue()
	if value == "" {
		return ""
	}

	// Simple context extraction (could be improved)
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
func (v *ValidatedExtractor) calculateValidationStats(result *PiiExtractionResult, validator LLMValidator) {
	if validator == nil {
		return
	}

	stats := &ValidationStats{}
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

// HealthCheck verifies that the LLM validator is working
func (v *ValidatedExtractor) HealthCheck(ctx context.Context) error {
	if v.validator == nil {
		return ErrValidationDisabled
	}
	return v.validator.HealthCheck(ctx)
}

// SetValidationConfig updates the validation configuration
func (v *ValidatedExtractor) SetValidationConfig(config *ValidationConfig) error {
	if config == nil {
		return NewPiiError("invalid config", "validation config cannot be nil")
	}

	v.config = config

	if config.Enabled {
		validator, err := NewLLMValidator(config)
		if err != nil {
			return err
		}
		v.validator = validator
	} else {
		v.validator = nil
	}

	return nil
}

// GetValidationConfig returns the current validation configuration
func (v *ValidatedExtractor) GetValidationConfig() *ValidationConfig {
	return v.config
}

// IsValidationEnabled returns true if LLM validation is enabled
func (v *ValidatedExtractor) IsValidationEnabled() bool {
	return v.config != nil && v.config.Enabled && v.validator != nil
}
