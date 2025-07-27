package patterns

import (
	"reflect"
	"testing"
)

func TestFrancePostalCodeExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard France postal code",
			input:    "Paris office at 75001",
			expected: []string{"75001"},
		},
		{
			name:     "multiple France postal codes",
			input:    "From 69001 Lyon to 13001 Marseille",
			expected: []string{"69001", "13001"},
		},
		{
			name:     "France postal codes in address",
			input:    "Address: 123 rue de la Paix, 75008 Paris",
			expected: []string{"75008"},
		},
		{
			name:     "overseas territories",
			input:    "Guadeloupe: 97110, Martinique: 97200",
			expected: []string{"97110", "97200"},
		},
		{
			name:     "invalid France postal codes",
			input:    "Invalid: 123, 1234, 123456, 00000",
			expected: []string{},
		},
		{
			name:     "no France postal codes",
			input:    "This text has no French postal codes",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PostalCodesFrance(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PostalCodesFrance() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFranceAddressExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard France street address",
			input:    "Adresse: 123 rue de la Paix, Paris",
			expected: []string{"123 rue de la Paix"},
		},
		{
			name:     "France address with Avenue",
			input:    "Location: 456 avenue des Champs-Élysées",
			expected: []string{"456 avenue des Champs-Élysées"},
		},
		{
			name:     "multiple France addresses",
			input:    "From 10 boulevard Saint-Germain to 25 place de la Bastille",
			expected: []string{"10 boulevard Saint-Germain", "25 place de la Bastille"},
		},
		{
			name:     "France addresses with various types",
			input:    "Properties: 100 rue Victor Hugo, 200 impasse Mozart, 300 allée des Roses",
			expected: []string{"100 rue Victor Hugo", "200 impasse Mozart", "300 allée des Roses"},
		},
		{
			name:     "invalid France addresses",
			input:    "Invalid: rue de la Paix, 123, Street without number",
			expected: []string{},
		},
		{
			name:     "no France addresses",
			input:    "This text has no French street addresses",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StreetAddressesFrance(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("StreetAddressesFrance() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
