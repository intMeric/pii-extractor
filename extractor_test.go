package piiextractor

import (
	"testing"
)

func TestRegexExtractor_Extract(t *testing.T) {
	extractor := NewRegexExtractor()
	
	text := `
		Contact John at john.doe@email.com or call him at (555) 123-4567.
		His address is 123 Main Street and his SSN is 123-45-6789.
		Credit card: 4111111111111111
		IP: 192.168.1.1
		Bitcoin: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
		IBAN: GB82WEST12345698765432
		ZIP: 90210
	`

	entities, err := extractor.Extract(text)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	// Check that we found some entities
	if len(entities) == 0 {
		t.Error("Expected to find PII entities, but got none")
	}

	// Count entities by type
	typeCount := make(map[string]int)
	for _, entity := range entities {
		typeCount[entity.Type]++
	}

	// Verify we found expected types
	expectedTypes := []string{"email", "phone", "street_address", "ssn", "credit_card", "ip_address", "btc_address", "iban", "zip_code"}
	for _, expectedType := range expectedTypes {
		if typeCount[expectedType] == 0 {
			t.Errorf("Expected to find at least one %s entity", expectedType)
		}
	}

	t.Logf("Found %d PII entities:", len(entities))
	for _, entity := range entities {
		t.Logf("- Type: %s, Value: %s", entity.Type, entity.GetValue())
	}
}

func TestRegexExtractor_EmptyText(t *testing.T) {
	extractor := NewRegexExtractor()
	
	entities, err := extractor.Extract("")
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if len(entities) != 0 {
		t.Errorf("Expected no entities for empty text, got %d", len(entities))
	}
}

func TestRegexExtractor_TypeAssertions(t *testing.T) {
	extractor := NewRegexExtractor()
	
	text := "Contact me at test@example.com or call (555) 123-4567"
	
	entities, err := extractor.Extract(text)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	for _, entity := range entities {
		switch entity.Type {
		case "email":
			if email, ok := entity.AsEmail(); ok {
				if email.GetValue() == "" {
					t.Error("Email value should not be empty")
				}
			} else {
				t.Error("Failed to cast email entity to Email type")
			}
		case "phone":
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