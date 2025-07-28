package patterns

import (
	"testing"
)

func TestArabicPostalCodes(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Valid Arabic countries postal codes",
			text:     "Saudi Arabia: 12345, UAE: 54321, Egypt: 11111.",
			expected: []string{"12345", "54321", "11111"},
		},
		{
			name:     "Mixed with Arabic text",
			text:     "الرياض 11564، دبي 12345، القاهرة 54321.",
			expected: []string{"11564", "12345", "54321"},
		},
		{
			name:     "Various Gulf countries",
			text:     "Kuwait 13000, Qatar 25000, Bahrain 33000.",
			expected: []string{"13000", "25000", "33000"},
		},
		{
			name:     "No postal codes",
			text:     "هذا النص لا يحتوي على رموز بريدية.",
			expected: []string{},
		},
		{
			name:     "Invalid format (less than 5 digits)",
			text:     "1234 is not a valid postal code.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PostalCodesArabic(tc.text)
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

func TestArabicPhones(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Saudi Arabia phone numbers",
			text:     "اتصل بي على +966 50 123 4567 أو 055-987-6543.",
			expected: []string{"+966 50 123 4567", "055-987-6543"},
		},
		{
			name:     "UAE phone numbers",
			text:     "Dubai: +971 50 123 4567, Abu Dhabi: 02 123 4567.",
			expected: []string{"+971 50 123 4567", "02 123 4567"},
		},
		{
			name:     "Egypt phone numbers",
			text:     "Cairo: +20 10 1234 5678, Alexandria: 03 123 4567.",
			expected: []string{"+20 10 1234 5678", "03 123 4567"},
		},
		{
			name:     "Mixed formats",
			text:     "رقم الهاتف: 966-50-123-4567، الفاكس: +971 4 123 4567.",
			expected: []string{"966-50-123-4567", "+971 4 123 4567"},
		},
		{
			name:     "No phone numbers",
			text:     "هذا النص لا يحتوي على أرقام هواتف.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PhonesArabic(tc.text)
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

func TestArabicStreetAddresses(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Arabic street addresses",
			text:     "أسكن في شارع الملك فهد، الرياض وأعمل في طريق الملك عبدالعزيز، جدة.",
			expected: []string{"شارع الملك فهد", "طريق الملك عبدالعزيز"},
		},
		{
			name:     "Mixed Arabic and English addresses",
			text:     "العنوان: شارع الشيخ زايد، دبي أو Emirates Road, Abu Dhabi.",
			expected: []string{"شارع الشيخ زايد"},
		},
		{
			name:     "Different Arabic address terms",
			text:     "زورونا في حي الملز، مدينة الرياض أو منطقة الزمالك، القاهرة.",
			expected: []string{"حي الملز", "مدينة الرياض", "منطقة الزمالك"},
		},
		{
			name:     "PO Box addresses",
			text:     "صندوق بريد ١٢٣٤، الرياض أو ص.ب 5678، دبي.",
			expected: []string{"صندوق بريد ١٢٣٤", "ص.ب 5678"},
		},
		{
			name:     "No street addresses",
			text:     "هذا النص لا يحتوي على عناوين شوارع.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StreetAddressesArabic(tc.text)
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