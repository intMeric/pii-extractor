package internal

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