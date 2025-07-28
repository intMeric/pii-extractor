package patterns

import (
	"testing"
)

func TestChinaPostalCodes(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Valid Chinese postal codes",
			text:     "北京市的邮编是100000，上海市的邮编是200000。",
			expected: []string{"100000", "200000"},
		},
		{
			name:     "Mixed languages postal codes",
			text:     "Beijing postal code is 100001, Shenzhen is 518000.",
			expected: []string{"100001", "518000"},
		},
		{
			name:     "Various cities",
			text:     "广州510000，深圳518000，成都610000。",
			expected: []string{"510000", "518000", "610000"},
		},
		{
			name:     "No postal codes",
			text:     "这个文本没有邮编。",
			expected: []string{},
		},
		{
			name:     "Invalid format (starts with 0)",
			text:     "012345 is not a valid Chinese postal code.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PostalCodesChina(tc.text)
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

func TestChinaPhones(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Chinese mobile numbers with country code",
			text:     "联系我：+86 138 0013 8000 或者 +86 159-9999-8888。",
			expected: []string{"+86 138 0013 8000", "+86 159-9999-8888"},
		},
		{
			name:     "Chinese mobile numbers without country code",
			text:     "我的手机号是138-0013-8000，办公室电话是159 9999 8888。",
			expected: []string{"138-0013-8000", "159 9999 8888"},
		},
		{
			name:     "Various mobile prefixes",
			text:     "Mobile numbers: 13800138000, 15999998888, 18612345678.",
			expected: []string{"13800138000", "15999998888", "18612345678"},
		},
		{
			name:     "No phone numbers",
			text:     "这个文本没有电话号码。",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PhonesChina(tc.text)
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

func TestChinaStreetAddresses(t *testing.T) {
	testCases := []struct{
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Chinese street addresses",
			text:     "我住在北京市朝阳区建国门外大街1号，公司在上海市浦东新区世纪大道88号。",
			expected: []string{"北京市朝阳区建国门外大街1号", "上海市浦东新区世纪大道88号"},
		},
		{
			name:     "Various address components",
			text:     "地址：广州市天河区珠江新城花城大道123号大厦。",
			expected: []string{"广州市天河区珠江新城花城大道123号"},
		},
		{
			name:     "Different location types",
			text:     "深圳市南山区科技园南区深南大道9988号，东莞市长安镇步步高大道168号。",
			expected: []string{"深圳市南山区科技园南区深南大道9988号", "东莞市长安镇步步高大道168号"},
		},
		{
			name:     "No street addresses", 
			text:     "这个文本没有街道地址。",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StreetAddressesChina(tc.text)
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