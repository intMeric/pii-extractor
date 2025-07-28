package patterns

import (
	"testing"
)

func TestIndiaPostalCodes(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Valid Indian postal codes",
			text:     "New Delhi PIN code is 110001, Mumbai is 400001.",
			expected: []string{"110001", "400001"},
		},
		{
			name:     "Various cities",
			text:     "Bangalore 560001, Chennai 600001, Kolkata 700001.",
			expected: []string{"560001", "600001", "700001"},
		},
		{
			name:     "PIN codes in sentence",
			text:     "Send the package to 110001 New Delhi or 400001 Mumbai.",
			expected: []string{"110001", "400001"},
		},
		{
			name:     "No postal codes",
			text:     "This text has no PIN codes.",
			expected: []string{},
		},
		{
			name:     "Invalid format (starts with 0)",
			text:     "012345 is not a valid Indian PIN code.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PostalCodesIndia(tc.text)
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d postal codes, got %d", len(tc.expected), len(result))
				return
			}
			for i, expected := range tc.expected {
				if result[i] != expected {
					t.Errorf("Expected postal code %s, got %s", expected, result[i])
				}
			}
		})
	}
}

func TestIndiaPhones(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Indian mobile numbers with country code",
			text:     "Call me at +91 98765 43210 or +91 90123 45678.",
			expected: []string{"+91 98765 43210", "+91 90123 45678"},
		},
		{
			name:     "Indian mobile numbers without country code",
			text:     "My number is 9876543210, office is 8012345678.",
			expected: []string{"9876543210", "8012345678"},
		},
		{
			name:     "Landline numbers",
			text:     "Delhi landline: 011-2345-6789, Mumbai: 022 1234 5678.",
			expected: []string{"011-2345-6789", "022 1234 5678"},
		},
		{
			name:     "Mixed formats",
			text:     "Mobile: 98765-43210, Landline: 11 2345 6789.",
			expected: []string{"98765-43210", "11 2345 6789"},
		},
		{
			name:     "No phone numbers",
			text:     "This text has no phone numbers.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PhonesIndia(tc.text)
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d phone numbers, got %d", len(tc.expected), len(result))
				return
			}
			for i, expected := range tc.expected {
				if result[i] != expected {
					t.Errorf("Expected phone %s, got %s", expected, result[i])
				}
			}
		})
	}
}

func TestIndiaStreetAddresses(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Indian street addresses",
			text:     "I live at 123 MG Road and work at 456 Brigade Road.",
			expected: []string{"123 MG Road", "456 Brigade Road"},
		},
		{
			name:     "Addresses with common Indian terms",
			text:     "Visit us at 789 Nehru Nagar, Sector 15 or 321 Gandhi Colony.",
			expected: []string{"789 Nehru Nagar", "321 Gandhi Colony"},
		},
		{
			name:     "Different address types",
			text:     "Residential: 12 Lotus Street, Commercial: 34 Business Park.",
			expected: []string{"12 Lotus Street", "34 Business Park"},
		},
		{
			name:     "Addresses with layout and phase",
			text:     "Plot 56 JP Nagar Phase 2, House 78 Koramangala Layout.",
			expected: []string{"56 JP Nagar Phase 2", "78 Koramangala Layout"},
		},
		{
			name:     "No street addresses",
			text:     "This text contains no street addresses.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StreetAddressesIndia(tc.text)
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d addresses, got %d", len(tc.expected), len(result))
				return
			}
			for i, expected := range tc.expected {
				if result[i] != expected {
					t.Errorf("Expected address %s, got %s", expected, result[i])
				}
			}
		})
	}
}