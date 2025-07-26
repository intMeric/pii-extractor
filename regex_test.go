package piiextractor

import (
	"reflect"
	"testing"
)

func TestEmailExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single valid email",
			input:    "Contact me at john.doe@example.com for details",
			expected: []string{"john.doe@example.com"},
		},
		{
			name:     "multiple emails in text",
			input:    "Send to admin@company.org and support@help.co.uk",
			expected: []string{"admin@company.org", "support@help.co.uk"},
		},
		{
			name:     "email with numbers and special chars",
			input:    "User test_user123+tag@domain-name.info needs access",
			expected: []string{"test_user123+tag@domain-name.info"},
		},
		{
			name:     "no emails present",
			input:    "This is just regular text without any email addresses",
			expected: []string{},
		},
		{
			name:     "invalid email formats",
			input:    "Invalid: @domain.com, user@, user@.com, user.domain.com",
			expected: []string{},
		},
		{
			name:     "case insensitive matching",
			input:    "Email: USER@DOMAIN.COM and user@domain.com",
			expected: []string{"USER@DOMAIN.COM", "user@domain.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Emails(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Emails() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

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

func TestCreditCardExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "visa card with spaces",
			input:    "Card: 4111 1111 1111 1111",
			expected: []string{"4111 1111 1111 1111"},
		},
		{
			name:     "visa card with dashes",
			input:    "Payment: 4111-1111-1111-1111",
			expected: []string{"4111-1111-1111-1111"},
		},
		{
			name:     "mastercard format",
			input:    "MC: 5555555555554444",
			expected: []string{"5555555555554444"},
		},
		{
			name:     "multiple cards",
			input:    "Card 1: 4111111111111111, Card 2: 5555 5555 5555 4444",
			expected: []string{"4111111111111111", "5555 5555 5555 4444"},
		},
		{
			name:     "invalid card numbers",
			input:    "Invalid: 1234, 12345678901234567890, abcd-efgh-ijkl-mnop",
			expected: []string{},
		},
		{
			name:     "no credit cards",
			input:    "This text has no payment information",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreditCards(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CreditCards() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIPAddressExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "IPv4 addresses",
			input:    "Server IP: 192.168.1.1, Public: 8.8.8.8",
			expected: []string{"192.168.1.1", "8.8.8.8"},
		},
		{
			name:     "IPv6 address",
			input:    "IPv6: 2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected: []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
		},
		{
			name:     "mixed IP versions",
			input:    "IPv4: 127.0.0.1 and IPv6: ::1",
			expected: []string{"127.0.0.1", "::1"},
		},
		{
			name:     "invalid IP addresses",
			input:    "Invalid: 256.256.256.256, 192.168.1, not.an.ip.address",
			expected: []string{},
		},
		{
			name:     "edge case IPs",
			input:    "Valid: 0.0.0.0, 255.255.255.255",
			expected: []string{"0.0.0.0", "255.255.255.255"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IPs(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("IPs() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIBANExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "valid IBAN",
			input:    "Bank account: GB82WEST12345698765432",
			expected: []string{"GB82WEST12345698765432"},
		},
		{
			name:     "multiple IBANs",
			input:    "Account 1: DE89370400440532013000, Account 2: FR1420041010050500013M02606",
			expected: []string{"DE89370400440532013000", "FR1420041010050500013M02606"},
		},
		{
			name:     "IBAN with different countries",
			input:    "US: US64SVBKUS6S3300958879, IT: IT60X0542811101000000123456",
			expected: []string{"US64SVBKUS6S3300958879", "IT60X0542811101000000123456"},
		},
		{
			name:     "invalid IBAN formats",
			input:    "Invalid: GB82WEST, 1234567890, SHORT",
			expected: []string{},
		},
		{
			name:     "no IBANs present",
			input:    "This text contains no international bank account numbers",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IBANs(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("IBANs() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestBitcoinAddressExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "valid Bitcoin address",
			input:    "Send payment to: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			expected: []string{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"},
		},
		{
			name:     "multiple Bitcoin addresses",
			input:    "Wallet 1: 3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy, Wallet 2: 1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2",
			expected: []string{"3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy", "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2"},
		},
		{
			name:     "Bitcoin address in transaction",
			input:    "Transaction from 1234567890abcdef to 1F1tAaz5x1HUXrCNLbtMDqcw6o5GNn4xqX",
			expected: []string{"1F1tAaz5x1HUXrCNLbtMDqcw6o5GNn4xqX"},
		},
		{
			name:     "invalid Bitcoin addresses",
			input:    "Invalid: 0A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa, short123, toolongaddresshere1234567890",
			expected: []string{},
		},
		{
			name:     "no Bitcoin addresses",
			input:    "This text contains no cryptocurrency addresses",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BtcAddresses(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("BtcAddresses() = %v, expected %v", result, tt.expected)
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

func TestExtractContext(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		start    int
		end      int
		expected string
	}{
		{
			name:     "extract complete sentence",
			text:     "This is the first sentence. Contact me at john@example.com for details. This is the last sentence.",
			start:    39,
			end:      55,
			expected: "Contact me at john@example.com for details.",
		},
		{
			name:     "extract from beginning of text",
			text:     "Email john@example.com today! More text follows.",
			start:    6,
			end:      22,
			expected: "Email john@example.com today!",
		},
		{
			name:     "extract from end of text",
			text:     "Please send the report to admin@company.org",
			start:    26,
			end:      43,
			expected: "Please send the report to admin@company.org",
		},
		{
			name:     "fallback to word context when no sentence boundaries",
			text:     "word1 word2 word3 word4 word5 john@example.com word6 word7 word8 word9 word10 word11 word12 word13",
			start:    30,
			end:      46,
			expected: "word1 word2 word3 word4 word5 john@example.com word6 word7 word8 word9 word10 word11 word12 word13",
		},
		{
			name:     "handle match at text boundaries",
			text:     "john@example.com",
			start:    0,
			end:      16,
			expected: "john@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractContext(tt.text, tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("extractContext() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestExtractSentence(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		start    int
		end      int
		expected string
	}{
		{
			name:     "middle sentence with periods",
			text:     "First sentence. Middle sentence with match. Last sentence.",
			start:    25,
			end:      30,
			expected: "Middle sentence with match.",
		},
		{
			name:     "sentence with exclamation mark",
			text:     "Call me today! My number is 555-1234. Thanks!",
			start:    28,
			end:      37,
			expected: "My number is 555-1234. Thanks!",
		},
		{
			name:     "sentence with question mark",
			text:     "What is your email? Is it john@example.com? Let me know.",
			start:    26,
			end:      42,
			expected: "Is it john@example.com?",
		},
		{
			name:     "start of text without sentence boundary",
			text:     "john@example.com is my email address. More text here.",
			start:    0,
			end:      16,
			expected: "john@example.com is my email address.",
		},
		{
			name:     "end of text without sentence boundary",
			text:     "Please contact me at admin@company.org",
			start:    21,
			end:      38,
			expected: "Please contact me at admin@company.org",
		},
		{
			name:     "no sentence boundaries found",
			text:     "just some words without punctuation",
			start:    5,
			end:      9,
			expected: "just some words without punctuation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSentence(tt.text, tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("extractSentence() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestExtractWordContext(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		start    int
		end      int
		expected string
	}{
		{
			name:     "extract 8 words before and after",
			text:     "w1 w2 w3 w4 w5 w6 w7 w8 w9 match w10 w11 w12 w13 w14 w15 w16 w17 w18",
			start:    30,
			end:      35,
			expected: "w2 w3 w4 w5 w6 w7 w8 w9 match w10 w11 w12 w13 w14 w15 w16 w17 w18",
		},
		{
			name:     "less than 8 words before",
			text:     "w1 w2 w3 match w4 w5 w6 w7 w8 w9 w10 w11 w12",
			start:    9,
			end:      14,
			expected: "w1 w2 w3 match w4 w5 w6 w7 w8 w9 w10 w11",
		},
		{
			name:     "less than 8 words after",
			text:     "w1 w2 w3 w4 w5 w6 w7 w8 w9 match w10 w11 w12",
			start:    30,
			end:      35,
			expected: "w2 w3 w4 w5 w6 w7 w8 w9 match w10 w11 w12",
		},
		{
			name:     "single word match",
			text:     "match",
			start:    0,
			end:      5,
			expected: "match",
		},
		{
			name:     "email address in sentence",
			text:     "Please send the report to john@example.com by tomorrow",
			start:    26,
			end:      42,
			expected: "Please send the report to john@example.com by tomorrow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractWordContext(tt.text, tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("extractWordContext() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestExtractPhonesUS(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
		expectedValue string
		expectedCtx   string
		expectedOccur int
	}{
		{
			name:          "single phone with context",
			input:         "Please call me at (555) 123-4567 for urgent matters.",
			expectedCount: 1,
			expectedValue: "(555) 123-4567",
			expectedCtx:   "Please call me at (555) 123-4567 for urgent matters.",
			expectedOccur: 1,
		},
		{
			name:          "multiple occurrences of same phone",
			input:         "Call 555-123-4567 for support. Emergency line: 555-123-4567. The number 555-123-4567 is available 24/7.",
			expectedCount: 1,
			expectedValue: "555-123-4567",
			expectedCtx:   "Call 555-123-4567 for support.",
			expectedOccur: 3,
		},
		{
			name:          "multiple different phones",
			input:         "Home: (555) 123-4567. Work: 555.987.6543. Mobile: +1 555 111 2222.",
			expectedCount: 3,
			expectedValue: "",
			expectedCtx:   "",
			expectedOccur: 1,
		},
		{
			name:          "phone with sentence context",
			input:         "Our customer service team is available! You can reach us at 555-123-4567 or email support@company.com. We're here to help!",
			expectedCount: 1,
			expectedValue: "555-123-4567",
			expectedCtx:   "You can reach us at 555-123-4567 or email support@company.",
			expectedOccur: 1,
		},
		{
			name:          "no phones in text",
			input:         "This text contains no phone numbers at all.",
			expectedCount: 0,
			expectedValue: "",
			expectedCtx:   "",
			expectedOccur: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractPhonesUS(tt.input)

			if len(result) != tt.expectedCount {
				t.Errorf("ExtractPhonesUS() returned %d phones, expected %d", len(result), tt.expectedCount)
				return
			}

			if tt.expectedCount == 0 {
				return
			}

			if tt.expectedCount == 1 {
				entity := result[0]
				phone, ok := entity.AsPhone()
				if !ok {
					t.Errorf("ExtractPhonesUS() returned non-phone entity")
					return
				}
				if phone.GetValue() != tt.expectedValue {
					t.Errorf("ExtractPhonesUS() phone value = %q, expected %q", phone.GetValue(), tt.expectedValue)
				}
				if phone.GetCount() != tt.expectedOccur {
					t.Errorf("ExtractPhonesUS() phone count = %d, expected %d", phone.GetCount(), tt.expectedOccur)
				}
				if phone.Country != "US" {
					t.Errorf("ExtractPhonesUS() phone country = %q, expected %q", phone.Country, "US")
				}
				contexts := phone.GetContexts()
				if len(contexts) > 0 && tt.expectedCtx != "" && contexts[0] != tt.expectedCtx {
					t.Errorf("ExtractPhonesUS() phone context = %q, expected %q", contexts[0], tt.expectedCtx)
				}
			}

			if tt.expectedCount > 1 {
				for _, entity := range result {
					phone, ok := entity.AsPhone()
					if !ok {
						t.Errorf("ExtractPhonesUS() returned non-phone entity")
						continue
					}
					if phone.Country != "US" {
						t.Errorf("ExtractPhonesUS() phone country = %q, expected %q", phone.Country, "US")
					}
					if phone.GetCount() != 1 {
						t.Errorf("ExtractPhonesUS() phone count = %d, expected %d", phone.GetCount(), 1)
					}
					contexts := phone.GetContexts()
					if len(contexts) == 0 {
						t.Errorf("ExtractPhonesUS() phone contexts should not be empty")
					}
				}
			}
		})
	}
}

func TestExtractEmails(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
		expectedValue string
		expectedCtx   string
		expectedOccur int
	}{
		{
			name:          "single email with context",
			input:         "Please contact me at john.doe@example.com for more information.",
			expectedCount: 1,
			expectedValue: "john.doe@example.com",
			expectedCtx:   "Please contact me at john.doe@example.com for more information.",
			expectedOccur: 1,
		},
		{
			name:          "multiple occurrences of same email",
			input:         "Send reports to admin@company.org. CC admin@company.org on all emails! Important: admin@company.org must be notified.",
			expectedCount: 1,
			expectedValue: "admin@company.org",
			expectedCtx:   "Send reports to admin@company.org.",
			expectedOccur: 3,
		},
		{
			name:          "multiple different emails",
			input:         "Support: help@company.com. Sales: sales@company.com. Admin: admin@company.com.",
			expectedCount: 3,
			expectedValue: "",
			expectedCtx:   "",
			expectedOccur: 1,
		},
		{
			name:          "email with sentence context",
			input:         "Welcome to our service! Please verify your account by clicking the link sent to user123@domain.co.uk. Thank you for joining us!",
			expectedCount: 1,
			expectedValue: "user123@domain.co.uk",
			expectedCtx:   "Please verify your account by clicking the link sent to user123@domain.co.uk.",
			expectedOccur: 1,
		},
		{
			name:          "no emails in text",
			input:         "This text contains no email addresses whatsoever.",
			expectedCount: 0,
			expectedValue: "",
			expectedCtx:   "",
			expectedOccur: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractEmails(tt.input)

			if len(result) != tt.expectedCount {
				t.Errorf("ExtractEmails() returned %d emails, expected %d", len(result), tt.expectedCount)
				return
			}

			if tt.expectedCount == 0 {
				return
			}

			if tt.expectedCount == 1 {
				entity := result[0]
				email, ok := entity.AsEmail()
				if !ok {
					t.Errorf("ExtractEmails() returned non-email entity")
					return
				}
				if email.GetValue() != tt.expectedValue {
					t.Errorf("ExtractEmails() email value = %q, expected %q", email.GetValue(), tt.expectedValue)
				}
				if email.GetCount() != tt.expectedOccur {
					t.Errorf("ExtractEmails() email count = %d, expected %d", email.GetCount(), tt.expectedOccur)
				}
				contexts := email.GetContexts()
				if len(contexts) > 0 && tt.expectedCtx != "" && contexts[0] != tt.expectedCtx {
					t.Errorf("ExtractEmails() email context = %q, expected %q", contexts[0], tt.expectedCtx)
				}
			}

			if tt.expectedCount > 1 {
				for _, entity := range result {
					email, ok := entity.AsEmail()
					if !ok {
						t.Errorf("ExtractEmails() returned non-email entity")
						continue
					}
					if email.GetCount() != 1 {
						t.Errorf("ExtractEmails() email count = %d, expected %d", email.GetCount(), 1)
					}
					contexts := email.GetContexts()
					if len(contexts) == 0 {
						t.Errorf("ExtractEmails() email contexts should not be empty")
					}
				}
			}
		})
	}
}

func TestExtractSSNsUS(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
		expectedValue string
		expectedCtx   string
		expectedOccur int
	}{
		{
			name:          "single SSN with context",
			input:         "Employee Social Security Number: 123-45-6789 for tax records.",
			expectedCount: 1,
			expectedValue: "123-45-6789",
			expectedCtx:   "Employee Social Security Number: 123-45-6789 for tax records.",
			expectedOccur: 1,
		},
		{
			name:          "multiple occurrences of same SSN",
			input:         "SSN 555-44-3333 was entered. Please verify 555-44-3333 is correct. Confirm: 555-44-3333.",
			expectedCount: 1,
			expectedValue: "555-44-3333",
			expectedCtx:   "SSN 555-44-3333 was entered.",
			expectedOccur: 3,
		},
		{
			name:          "multiple different SSNs",
			input:         "Employee 1: 123-45-6789. Employee 2: 987-65-4321. Manager: 555-44-3333.",
			expectedCount: 3,
			expectedValue: "",
			expectedCtx:   "",
			expectedOccur: 1,
		},
		{
			name:          "SSN with sentence context",
			input:         "Please update your records! The new SSN is 999-88-7777 effective immediately. Contact HR for questions.",
			expectedCount: 1,
			expectedValue: "999-88-7777",
			expectedCtx:   "The new SSN is 999-88-7777 effective immediately.",
			expectedOccur: 1,
		},
		{
			name:          "no SSNs in text",
			input:         "This document contains no social security numbers.",
			expectedCount: 0,
			expectedValue: "",
			expectedCtx:   "",
			expectedOccur: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSSNsUS(tt.input)

			if len(result) != tt.expectedCount {
				t.Errorf("ExtractSSNsUS() returned %d SSNs, expected %d", len(result), tt.expectedCount)
				return
			}

			if tt.expectedCount == 0 {
				return
			}

			if tt.expectedCount == 1 {
				entity := result[0]
				ssn, ok := entity.AsSSN()
				if !ok {
					t.Errorf("ExtractSSNsUS() returned non-SSN entity")
					return
				}
				if ssn.GetValue() != tt.expectedValue {
					t.Errorf("ExtractSSNsUS() SSN value = %q, expected %q", ssn.GetValue(), tt.expectedValue)
				}
				if ssn.GetCount() != tt.expectedOccur {
					t.Errorf("ExtractSSNsUS() SSN count = %d, expected %d", ssn.GetCount(), tt.expectedOccur)
				}
				if ssn.Country != "US" {
					t.Errorf("ExtractSSNsUS() SSN country = %q, expected %q", ssn.Country, "US")
				}
				contexts := ssn.GetContexts()
				if len(contexts) > 0 && tt.expectedCtx != "" && contexts[0] != tt.expectedCtx {
					t.Errorf("ExtractSSNsUS() SSN context = %q, expected %q", contexts[0], tt.expectedCtx)
				}
			}

			if tt.expectedCount > 1 {
				for _, entity := range result {
					ssn, ok := entity.AsSSN()
					if !ok {
						t.Errorf("ExtractSSNsUS() returned non-SSN entity")
						continue
					}
					if ssn.Country != "US" {
						t.Errorf("ExtractSSNsUS() SSN country = %q, expected %q", ssn.Country, "US")
					}
					if ssn.GetCount() != 1 {
						t.Errorf("ExtractSSNsUS() SSN count = %d, expected %d", ssn.GetCount(), 1)
					}
					contexts := ssn.GetContexts()
					if len(contexts) == 0 {
						t.Errorf("ExtractSSNsUS() SSN contexts should not be empty")
					}
				}
			}
		})
	}
}

// International PII Tests

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
