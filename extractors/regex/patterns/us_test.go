package patterns

import (
	"reflect"
	"testing"
)

func TestUSPhoneExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard US phone format",
			input:    "Call me at (555) 123-4567",
			expected: []string{"(555) 123-4567"},
		},
		{
			name:     "phone without parentheses",
			input:    "My number is 555-123-4567",
			expected: []string{"555-123-4567"},
		},
		{
			name:     "phone with dots",
			input:    "Contact: 555.123.4567",
			expected: []string{"555.123.4567"},
		},
		{
			name:     "phone with spaces",
			input:    "Phone: 555 123 4567",
			expected: []string{"555 123 4567"},
		},
		{
			name:     "international format",
			input:    "International: +1 555 123 4567",
			expected: []string{"+1 555 123 4567"},
		},
		{
			name:     "multiple phones",
			input:    "Home: (555) 123-4567, Work: 555.987.6543",
			expected: []string{"(555) 123-4567", "555.987.6543"},
		},
		{
			name:     "no valid phones",
			input:    "Invalid: 123-45, 12345, phone number",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PhonesUS(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PhonesUS() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestUSSSNExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "valid SSN format",
			input:    "My SSN is 123-45-6789",
			expected: []string{"123-45-6789"},
		},
		{
			name:     "multiple SSNs",
			input:    "Employee 1: 123-45-6789, Employee 2: 987-65-4321",
			expected: []string{"123-45-6789", "987-65-4321"},
		},
		{
			name:     "SSN in document",
			input:    "Social Security Number: 555-44-3333 for tax purposes",
			expected: []string{"555-44-3333"},
		},
		{
			name:     "invalid SSN formats",
			input:    "Invalid: 123456789, 123-456-789, 12-34-5678",
			expected: []string{},
		},
		{
			name:     "no SSN present",
			input:    "This document contains no social security numbers",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SSNsUS(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SSNsUS() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestUSZipCodeExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "5 digit zip code",
			input:    "Address: 123 Main St, City, ST 12345",
			expected: []string{"12345"},
		},
		{
			name:     "zip+4 format",
			input:    "Shipping to 90210-1234",
			expected: []string{"90210-1234"},
		},
		{
			name:     "zip+4 with space",
			input:    "Location: 10001 5678",
			expected: []string{"10001 5678"},
		},
		{
			name:     "multiple zip codes",
			input:    "From 90210 to 10001-2345",
			expected: []string{"90210", "10001-2345"},
		},
		{
			name:     "invalid zip formats",
			input:    "Invalid: 1234, 123456, abcde",
			expected: []string{},
		},
		{
			name:     "no zip codes",
			input:    "No postal codes in this text",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZipCodesUS(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ZipCodesUS() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestUSStreetAddressExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard street address",
			input:    "Address: 123 Main Street, Anytown, ST 12345",
			expected: []string{"123 Main Street"},
		},
		{
			name:     "abbreviated street types",
			input:    "Location: 456 Oak Ave and 789 Pine Rd",
			expected: []string{"456 Oak Ave", "789 Pine Rd"},
		},
		{
			name:     "various street types",
			input:    "Properties: 100 First Blvd, 200 Second Dr, 300 Third Ct",
			expected: []string{"100 First Blvd", "200 Second Dr", "300 Third Ct"},
		},
		{
			name:     "address with apartment number",
			input:    "Shipping: 1234 Elm Street Apt 5B",
			expected: []string{"1234 Elm Street"},
		},
		{
			name:     "invalid addresses",
			input:    "Invalid: Main Street, 123, Street without number",
			expected: []string{},
		},
		{
			name:     "no street addresses",
			input:    "This text has no street addresses mentioned",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StreetAddressesUS(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("StreetAddressesUS() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
