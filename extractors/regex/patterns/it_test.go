package patterns

import (
	"reflect"
	"testing"
)

func TestItalyPostalCodeExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard Italy postal code",
			input:    "Rome office at 00100",
			expected: []string{"00100"},
		},
		{
			name:     "multiple Italy postal codes",
			input:    "From 20100 Milano to 80100 Napoli",
			expected: []string{"20100", "80100"},
		},
		{
			name:     "Italy postal codes in address",
			input:    "Indirizzo: Via del Corso 123, 00186 Roma",
			expected: []string{"00186"},
		},
		{
			name:     "Sicily and Sardinia postal codes",
			input:    "Palermo: 90100, Cagliari: 09100",
			expected: []string{"90100", "09100"},
		},
		{
			name:     "invalid Italy postal codes",
			input:    "Invalid: 123, 1234, 123456",
			expected: []string{},
		},
		{
			name:     "no Italy postal codes",
			input:    "This text has no Italian postal codes",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PostalCodesItaly(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PostalCodesItaly() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestItalyAddressExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard Italy street address",
			input:    "Indirizzo: 123 Via del Corso, Roma",
			expected: []string{"123 Via del Corso"},
		},
		{
			name:     "Italy address with Piazza",
			input:    "Location: 456 Piazza San Marco",
			expected: []string{"456 Piazza San Marco"},
		},
		{
			name:     "multiple Italy addresses",
			input:    "From 10 Via Veneto to 25 Corso Buenos Aires",
			expected: []string{"10 Via Veneto", "25 Corso Buenos Aires"},
		},
		{
			name:     "Italy addresses with various types",
			input:    "Properties: 100 Via Roma, 200 Viale Europa, 300 Largo Argentina",
			expected: []string{"100 Via Roma", "200 Viale Europa", "300 Largo Argentina"},
		},
		{
			name:     "invalid Italy addresses",
			input:    "Invalid: Via del Corso, 123, Street without number",
			expected: []string{},
		},
		{
			name:     "no Italy addresses",
			input:    "This text has no Italian street addresses",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StreetAddressesItaly(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("StreetAddressesItaly() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
