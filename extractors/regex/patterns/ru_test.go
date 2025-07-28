package patterns

import (
	"testing"
)

func TestRussiaPostalCodes(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Valid Russian postal codes",
			text:     "Москва 101000, Санкт-Петербург 190000, Новосибирск 630000.",
			expected: []string{"101000", "190000", "630000"},
		},
		{
			name:     "Mixed with English text",
			text:     "Moscow postal code is 101000, St. Petersburg is 190000.",
			expected: []string{"101000", "190000"},
		},
		{
			name:     "Various Russian cities",
			text:     "Екатеринбург 620000, Нижний Новгород 603000, Казань 420000.",
			expected: []string{"620000", "603000", "420000"},
		},
		{
			name:     "No postal codes",
			text:     "Этот текст не содержит почтовых индексов.",
			expected: []string{},
		},
		{
			name:     "Invalid format (starts with 0 or 7-9)",
			text:     "012345, 700000, 800000 are not valid Russian postal codes.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PostalCodesRussia(tc.text)
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

func TestRussiaPhones(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Russian mobile numbers with country code",
			text:     "Звоните мне: +7 495 123-45-67 или +7 (812) 987-65-43.",
			expected: []string{"+7 495 123-45-67", "+7 (812) 987-65-43"},
		},
		{
			name:     "Russian numbers with 8 prefix",
			text:     "Мой номер: 8 495 123-45-67, офис: 8(812)987-65-43.",
			expected: []string{"8 495 123-45-67", "8(812)987-65-43"},
		},
		{
			name:     "Mobile numbers",
			text:     "Мобильные: +7 903 123 45 67, 8 916 987 65 43.",
			expected: []string{"+7 903 123 45 67", "8 916 987 65 43"},
		},
		{
			name:     "Various formats",
			text:     "Телефоны: 7-495-123-45-67, 8(812)987-65-43, +7 903 1234567.",
			expected: []string{"7-495-123-45-67", "8(812)987-65-43", "+7 903 1234567"},
		},
		{
			name:     "No phone numbers",
			text:     "Этот текст не содержит номеров телефонов.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PhonesRussia(tc.text)
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

func TestRussiaStreetAddresses(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Russian street addresses",
			text:     "Я живу на улице Тверская, дом 13 и работаю на проспекте Невский, дом 25.",
			expected: []string{"улице Тверская, дом 13", "проспекте Невский, дом 25"},
		},
		{
			name:     "Addresses with abbreviations",
			text:     "Адрес: ул. Пушкина, д. 10, кв. 5 или пр. Ленина, д. 20.",
			expected: []string{"ул. Пушкина, д. 10, кв. 5", "пр. Ленина, д. 20"},
		},
		{
			name:     "Different address types",
			text:     "Переулок Чехова, дом 3, набережная Мойки, дом 12, площадь Победы, дом 1.",
			expected: []string{"Переулок Чехова, дом 3", "набережная Мойки, дом 12", "площадь Победы, дом 1"},
		},
		{
			name:     "Addresses with building parts",
			text:     "Московская улица, дом 15, корпус 2, строение 1, квартира 45.",
			expected: []string{"Московская улица, дом 15, корпус 2, строение 1, квартира 45"},
		},
		{
			name:     "No street addresses",
			text:     "Этот текст не содержит адресов улиц.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StreetAddressesRussia(tc.text)
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