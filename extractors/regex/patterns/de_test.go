package patterns

import (
	"testing"
)

func TestGermanyPostalCodes(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Valid German postal codes",
			text:     "My address is 10115 Berlin, and my friend lives in 80331 München.",
			expected: []string{"10115", "80331"},
		},
		{
			name:     "Mixed with other numbers",
			text:     "Call me at 030 12345678 or visit me at 10115 Berlin.",
			expected: []string{"10115"},
		},
		{
			name:     "East German postal codes",
			text:     "Dresden 01067 and Leipzig 04109 are beautiful cities.",
			expected: []string{"01067", "04109"},
		},
		{
			name:     "No postal codes",
			text:     "This text has no postal codes in it.",
			expected: []string{},
		},
		{
			name:     "Invalid format",
			text:     "00123 is not a valid German postal code.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PostalCodesGermany(tc.text)
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

func TestGermanyPhones(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "German phone numbers with country code",
			text:     "Call me at +49 30 12345678 or +49(89)87654321.",
			expected: []string{"+49 30 12345678", "+49(89)87654321"},
		},
		{
			name:     "German phone numbers without country code",
			text:     "My number is 030 12345678 and office is (089) 87654321.",
			expected: []string{"030 12345678", "(089) 87654321"},
		},
		{
			name:     "Mobile numbers",
			text:     "Mobile: 0177 1234567 or 0160-9876543.",
			expected: []string{"0177 1234567", "0160-9876543"},
		},
		{
			name:     "No phone numbers",
			text:     "This text has no phone numbers.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PhonesGermany(tc.text)
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

func TestGermanyStreetAddresses(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "German street addresses",
			text:     "I live at Münchner Straße 15 and work on Unter den Linden 1.",
			expected: []string{"Münchner Straße 15", "Unter den Linden 1"},
		},
		{
			name:     "Street addresses with abbreviations",
			text:     "Visit us at Bahnhofstr. 42 or Königsplatz 8.",
			expected: []string{"Bahnhofstr. 42", "Königsplatz 8"},
		},
		{
			name:     "Different street types",
			text:     "Schillerweg 23, Goethealle 5, and Hauptring 17.",
			expected: []string{"Schillerweg 23", "Goethealle 5", "Hauptring 17"},
		},
		{
			name:     "No street addresses",
			text:     "This text contains no street addresses.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StreetAddressesGermany(tc.text)
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