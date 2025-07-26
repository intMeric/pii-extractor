package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pii "github.com/intMeric/pii-extractor"
)

func main() {
	// Example text containing PII
	text := `
	Dear John,
	
	Thank you for your application. Please contact us at support@company.com 
	or call us at 555-123-4567. We may also reach you at john.doe@email.com.
	
	For verification, please confirm your SSN: 123-45-6789 and ZIP code: 90210.
	
	Test data: contact test@example.com or call 000-000-0000.
	
	Best regards,
	Customer Service
	`

	// Example 1: Basic extraction without validation
	fmt.Println("=== Basic Extraction (No Validation) ===")
	basicExtractor := pii.NewRegexExtractor()
	basicResult, err := basicExtractor.Extract(text)
	if err != nil {
		log.Fatalf("Basic extraction failed: %v", err)
	}

	fmt.Printf("Found %d PII entities:\n", basicResult.Total)
	for _, entity := range basicResult.Entities {
		fmt.Printf("- %s: %s\n", entity.Type.String(), entity.GetValue())
	}
	fmt.Println()

	// Example 2: Extraction with OpenAI validation (requires API key)
	fmt.Println("=== Extraction with OpenAI Validation ===")
	runValidationExample(text, pii.ProviderOpenAI, "gpt-4o-mini", "your-openai-api-key")

	// Example 3: Extraction with Mistral validation (requires API key)
	fmt.Println("=== Extraction with Mistral Validation ===")
	runValidationExample(text, pii.ProviderMistral, "mistral-small-latest", "your-mistral-api-key")

	// Example 4: Extraction with Ollama (local, no API key needed)
	fmt.Println("=== Extraction with Ollama (Local) ===")
	runOllamaExample(text)

	// Example 5: Custom validation configuration
	fmt.Println("=== Custom Validation Configuration ===")
	runCustomConfigExample(text)
}

func runValidationExample(text string, provider pii.LLMProvider, model, apiKey string) {
	// Create validation configuration
	config := &pii.ValidationConfig{
		Enabled:       true,
		Provider:      provider,
		Model:         model,
		APIKey:        apiKey,
		Timeout:       30 * time.Second,
		MinConfidence: 0.7,
		MaxRetries:    2,
	}

	// Create validated extractor
	baseExtractor := pii.NewRegexExtractor()
	validatedExtractor, err := pii.NewValidatedExtractor(baseExtractor, config)
	if err != nil {
		log.Printf("Failed to create validated extractor: %v", err)
		return
	}

	// Extract with validation using instance configuration
	result, err := validatedExtractor.ExtractWithValidation(text)
	if err != nil {
		log.Printf("Validation extraction failed: %v", err)
		return
	}

	printValidationResults(result)
}

func runOllamaExample(text string) {
	// Ollama configuration (assumes Ollama is running locally)
	config := &pii.ValidationConfig{
		Enabled:       true,
		Provider:      pii.ProviderOllama,
		Model:         "llama3.2",
		BaseURL:       "http://localhost:11434",
		Timeout:       60 * time.Second, // Longer timeout for local models
		MinConfidence: 0.6,              // Lower threshold for local models
		MaxRetries:    1,
	}

	baseExtractor := pii.NewRegexExtractor()
	validatedExtractor, err := pii.NewValidatedExtractor(baseExtractor, config)
	if err != nil {
		log.Printf("Failed to create Ollama extractor: %v", err)
		return
	}

	result, err := validatedExtractor.ExtractWithValidation(text)
	if err != nil {
		log.Printf("Ollama validation failed: %v", err)
		return
	}

	printValidationResults(result)
}

func runCustomConfigExample(text string) {
	// Custom configuration with high confidence threshold
	config := &pii.ValidationConfig{
		Enabled:       true,
		Provider:      pii.ProviderOpenAI,
		Model:         "gpt-4o",
		APIKey:        "your-api-key",
		Timeout:       45 * time.Second,
		MinConfidence: 0.9, // Very high confidence required
		MaxRetries:    3,
		ProviderOptions: map[string]interface{}{
			"temperature": 0.1, // Low temperature for consistent results
		},
	}

	baseExtractor := pii.NewRegexExtractor()
	validatedExtractor, err := pii.NewValidatedExtractor(baseExtractor, config)
	if err != nil {
		log.Printf("Failed to create custom extractor: %v", err)
		return
	}

	// Check if validation is enabled and working
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := validatedExtractor.HealthCheck(ctx); err != nil {
		log.Printf("Validation service health check failed: %v", err)
		return
	}

	result, err := validatedExtractor.ExtractWithValidation(text)
	if err != nil {
		log.Printf("Custom validation failed: %v", err)
		return
	}

	printValidationResults(result)

	// Demonstrate filtering methods
	fmt.Println("\n--- Filtering Results ---")
	validEntities := result.GetValidEntities()
	fmt.Printf("Valid entities (high confidence): %d\n", len(validEntities))

	invalidEntities := result.GetInvalidEntities()
	fmt.Printf("Invalid entities (filtered out): %d\n", len(invalidEntities))

	unvalidatedEntities := []pii.PiiEntity{}
	for _, entity := range result.Entities {
		if !entity.IsValidated() {
			unvalidatedEntities = append(unvalidatedEntities, entity)
		}
	}
	fmt.Printf("Unvalidated entities: %d\n", len(unvalidatedEntities))
}

func printValidationResults(result *pii.PiiExtractionResult) {
	fmt.Printf("Found %d PII entities:\n", result.Total)

	for _, entity := range result.Entities {
		status := "Not validated"
		confidence := ""

		if entity.IsValidated() {
			if entity.IsValid() {
				status = "✓ Valid"
			} else {
				status = "✗ Invalid"
			}
			confidence = fmt.Sprintf(" (confidence: %.2f)", entity.GetValidationConfidence())
		}

		fmt.Printf("- %s: %s [%s%s]\n",
			entity.Type.String(),
			entity.GetValue(),
			status,
			confidence)

		if entity.IsValidated() && entity.Validation.Reasoning != "" {
			fmt.Printf("  Reasoning: %s\n", entity.Validation.Reasoning)
		}
	}

	// Print validation statistics if available
	if result.ValidationStats != nil {
		stats := result.ValidationStats
		fmt.Printf("\nValidation Statistics:\n")
		fmt.Printf("- Provider: %s (%s)\n", stats.Provider, stats.Model)
		fmt.Printf("- Total validated: %d\n", stats.TotalValidated)
		fmt.Printf("- Valid: %d, Invalid: %d\n", stats.ValidCount, stats.InvalidCount)
		fmt.Printf("- Average confidence: %.2f\n", stats.AverageConfidence)
	}
	fmt.Println()
}

// Example of processing specific PII types with validation
func processEmailsWithValidation(text string) {
	config := &pii.ValidationConfig{
		Enabled:       true,
		Provider:      pii.ProviderOpenAI,
		Model:         "gpt-4o-mini",
		APIKey:        "your-api-key",
		MinConfidence: 0.8,
	}

	baseExtractor := pii.NewRegexExtractor()
	validatedExtractor, err := pii.NewValidatedExtractor(baseExtractor, config)
	if err != nil {
		log.Printf("Failed to create extractor: %v", err)
		return
	}

	result, err := validatedExtractor.ExtractWithValidation(text)
	if err != nil {
		log.Printf("Extraction failed: %v", err)
		return
	}

	// Process only validated emails
	emails := result.GetEmails()
	fmt.Printf("Processing %d email addresses:\n", len(emails))

	for _, emailEntity := range emails {
		if email, ok := emailEntity.AsEmail(); ok {
			if emailEntity.IsValid() {
				fmt.Printf("✓ Processing valid email: %s\n", email.GetValue())
				// Process the valid email...
			} else if emailEntity.IsValidated() {
				fmt.Printf("✗ Skipping invalid email: %s (reason: %s)\n",
					email.GetValue(), emailEntity.Validation.Reasoning)
			} else {
				fmt.Printf("? Email not validated: %s\n", email.GetValue())
			}
		}
	}
}

// Example configuration for different use cases
func getConfigForUseCase(useCase string) *pii.ValidationConfig {
	switch useCase {
	case "high-accuracy":
		// For critical applications requiring high accuracy
		return &pii.ValidationConfig{
			Enabled:       true,
			Provider:      pii.ProviderOpenAI,
			Model:         "gpt-4o",
			Timeout:       60 * time.Second,
			MinConfidence: 0.95,
			MaxRetries:    3,
		}

	case "cost-effective":
		// For high-volume processing with cost constraints
		return &pii.ValidationConfig{
			Enabled:       true,
			Provider:      pii.ProviderMistral,
			Model:         "mistral-small-latest",
			Timeout:       30 * time.Second,
			MinConfidence: 0.7,
			MaxRetries:    2,
		}

	case "privacy-focused":
		// For scenarios requiring data privacy (local processing)
		return &pii.ValidationConfig{
			Enabled:       true,
			Provider:      pii.ProviderOllama,
			Model:         "llama3.2",
			BaseURL:       "http://localhost:11434",
			Timeout:       90 * time.Second,
			MinConfidence: 0.6,
			MaxRetries:    1,
		}

	case "fast-processing":
		// For real-time applications requiring speed
		return &pii.ValidationConfig{
			Enabled:       true,
			Provider:      pii.ProviderOpenAI,
			Model:         "gpt-4o-mini",
			Timeout:       15 * time.Second,
			MinConfidence: 0.75,
			MaxRetries:    1,
		}

	default:
		return pii.DefaultValidationConfig()
	}
}
