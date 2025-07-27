package patterns

import (
	"reflect"
	"testing"
)

func TestSpainPostalCodeExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard Spain postal code",
			input:    "Madrid office at 28001",
			expected: []string{"28001"},
		},
		{
			name:     "multiple Spain postal codes",
			input:    "From 08001 Barcelona to 41001 Sevilla",
			expected: []string{"08001", "41001"},
		},
		{
			name:     "Spain postal codes in address",
			input:    "Dirección: Calle Mayor 123, 28013 Madrid",
			expected: []string{"28013"},
		},
		{
			name:     "Canary Islands postal codes",
			input:    "Las Palmas: 35001, Santa Cruz: 38001",
			expected: []string{"35001", "38001"},
		},
		{
			name:     "invalid Spain postal codes",
			input:    "Invalid: 123, 1234, 123456, 00000",
			expected: []string{},
		},
		{
			name:     "no Spain postal codes",
			input:    "This text has no Spanish postal codes",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PostalCodesSpain(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PostalCodesSpain() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSpainAddressExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard Spain street address",
			input:    "Dirección: 123 Calle Mayor, Madrid",
			expected: []string{"123 Calle Mayor"},
		},
		{
			name:     "Spain address with Avenida",
			input:    "Location: 456 Avenida de la Constitución",
			expected: []string{"456 Avenida de la Constitución"},
		},
		{
			name:     "multiple Spain addresses",
			input:    "From 10 Plaza de España to 25 Paseo de Gracia",
			expected: []string{"10 Plaza de España", "25 Paseo de Gracia"},
		},
		{
			name:     "Spain addresses with various types",
			input:    "Properties: 100 Calle Real, 200 Travesía del Carmen, 300 Ronda de Toledo",
			expected: []string{"100 Calle Real", "200 Travesía del Carmen", "300 Ronda de Toledo"},
		},
		{
			name:     "invalid Spain addresses",
			input:    "Invalid: Calle Mayor, 123, Street without number",
			expected: []string{},
		},
		{
			name:     "no Spain addresses",
			input:    "This text has no Spanish street addresses",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StreetAddressesSpain(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("StreetAddressesSpain() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
