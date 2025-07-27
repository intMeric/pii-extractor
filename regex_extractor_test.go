package piiextractor

import (
	"testing"
)

func TestRegexExtractor_Extract(t *testing.T) {
	extractor := NewDefaultRegexExtractor()

	text := `
		Contact John at john.doe@email.com or call him at (555) 123-4567.
		His address is 123 Main Street and his SSN is 123-45-6789.
		Credit card: 4111-1111-1111-1111
		IP: 192.168.1.1
		Bitcoin: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
		IBAN: GB82WEST12345698765432
		ZIP: 90210
	`

	result, err := extractor.Extract(text)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	// Check that we found some entities
	if result.Total == 0 {
		t.Error("Expected to find PII entities, but got none")
	}

	// Use the built-in stats from PiiExtractionResult
	typeCount := result.Stats

	// Verify we found some core types (not all, due to potential false positives/negatives)
	coreTypes := []PiiType{PiiTypeEmail, PiiTypePhone, PiiTypeSSN}
	for _, coreType := range coreTypes {
		if typeCount[coreType] == 0 {
			t.Errorf("Expected to find at least one %s entity", coreType)
		}
	}

	// Log results for debugging
	t.Logf("Found %d PII entities:", result.Total)
	for _, entity := range result.Entities {
		t.Logf("- Type: %s, Value: %s", entity.Type, entity.GetValue())
	}

	// Additional checks for quality
	if result.Total < 5 {
		t.Errorf("Expected to find at least 5 entities, got %d", result.Total)
	}
}

func TestRegexExtractor_EmptyText(t *testing.T) {
	extractor := NewDefaultRegexExtractor()

	result, err := extractor.Extract("")
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if result.Total != 0 {
		t.Errorf("Expected no entities for empty text, got %d", result.Total)
	}
}

func TestRegexExtractor_TypeAssertions(t *testing.T) {
	extractor := NewDefaultRegexExtractor()

	text := "Contact me at test@example.com or call (555) 123-4567"

	result, err := extractor.Extract(text)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	for _, entity := range result.Entities {
		switch entity.Type {
		case PiiTypeEmail:
			if email, ok := entity.AsEmail(); ok {
				if email.GetValue() == "" {
					t.Error("Email value should not be empty")
				}
			} else {
				t.Error("Failed to cast email entity to Email type")
			}
		case PiiTypePhone:
			if phone, ok := entity.AsPhone(); ok {
				if phone.GetValue() == "" {
					t.Error("Phone value should not be empty")
				}
				if phone.Country != "US" {
					t.Errorf("Expected US phone country, got %s", phone.Country)
				}
			} else {
				t.Error("Failed to cast phone entity to Phone type")
			}
		}
	}
}
