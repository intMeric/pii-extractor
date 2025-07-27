package patterns

import (
	"reflect"
	"testing"
)

func TestUKPostalCodeExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard UK postcode format",
			input:    "Please send mail to SW1A 1AA",
			expected: []string{"SW1A 1AA"},
		},
		{
			name:     "UK postcode without space",
			input:    "Location: M11AA Manchester",
			expected: []string{"M11AA"},
		},
		{
			name:     "multiple UK postcodes",
			input:    "From W1A 0AX to EC1A 1BB",
			expected: []string{"W1A 0AX", "EC1A 1BB"},
		},
		{
			name:     "London postcode formats",
			input:    "Office at E1 6AN and home at SW1P 3BU",
			expected: []string{"E1 6AN", "SW1P 3BU"},
		},
		{
			name:     "invalid UK postcode formats",
			input:    "Invalid: 123456, A1B 2C3D, XYZ",
			expected: []string{},
		},
		{
			name:     "no UK postcodes",
			input:    "This text has no British postal codes",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PostalCodesUK(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PostalCodesUK() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestUKAddressExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard UK street address",
			input:    "Address: 123 Oxford Street, London",
			expected: []string{"123 Oxford Street"},
		},
		{
			name:     "UK address with Road",
			input:    "Location: 456 Abbey Road, Westminster",
			expected: []string{"456 Abbey Road"},
		},
		{
			name:     "multiple UK addresses",
			input:    "From 10 Downing Street to 221B Baker Street",
			expected: []string{"10 Downing Street", "221B Baker Street"},
		},
		{
			name:     "UK addresses with various types",
			input:    "Properties: 100 High Street, 200 Church Lane, 300 Market Square",
			expected: []string{"100 High Street", "200 Church Lane", "300 Market Square"},
		},
		{
			name:     "invalid UK addresses",
			input:    "Invalid: Oxford Street, 123, Street without number",
			expected: []string{},
		},
		{
			name:     "no UK addresses",
			input:    "This text has no British street addresses",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StreetAddressesUK(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("StreetAddressesUK() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
