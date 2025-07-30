package regex

import (
	"strings"
	"testing"

	"github.com/intMeric/pii-extractor/extractors"
	"github.com/intMeric/pii-extractor/extractors/regex/patterns"
)

// Benchmark data - realistic multi-country text with various PII types
const benchmarkText = `
John Doe works at Acme Corp. You can reach him at john.doe@acme.com or call (555) 123-4567.
His office is located at 123 Main Street, New York, NY 10001.
SSN: 123-45-6789, Credit Card: 4111-1111-1111-1111
Server IP: 192.168.1.100, Bitcoin: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa

UK Office: 221B Baker Street, London SW1A 1AA. Contact: +44 20 7946 0958
France Office: 123 rue de la Paix, 75001 Paris. Contact: +33 1 42 96 87 56

Germany Office: Münchner Straße 15, 80331 München. Phone: +49 30 12345678
China Office: 北京市朝阳区建国门外大街1号, 100020. Phone: +86 138 0013 8000

India Office: 123 MG Road, Bangalore 560001. Phone: +91 98765 43210
Russia Office: ул. Тверская, д. 13, 101000 Москва. Phone: +7 495 123-45-67

IBAN: GB82WEST12345698765432, DE89370400440532013000
More emails: marie.dupont@example.fr, hans.mueller@example.de, priya.sharma@example.in
Multiple phones: (555) 987-6543, +44 161 496 0018, +33 6 12 34 56 78
`

func BenchmarkRegexExtractor_Extract(b *testing.B) {
	extractor := NewDefaultExtractor()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, err := extractor.Extract(benchmarkText)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRegexExtractor_ExtractLargeText(b *testing.B) {
	// Create a large text by repeating the benchmark text
	largeText := strings.Repeat(benchmarkText, 100) // ~100KB of text
	extractor := NewDefaultExtractor()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, err := extractor.Extract(largeText)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRegexExtractor_ExtractByType_Email(b *testing.B) {
	extractor := NewDefaultExtractor()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, err := extractor.ExtractByType(benchmarkText, 1) // PiiTypeEmail = 1
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRegexExtractor_ExtractByType_Phone(b *testing.B) {
	extractor := NewDefaultExtractor()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, err := extractor.ExtractByType(benchmarkText, 0) // PiiTypePhone = 0
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRegexExtractor_ExtractSpecificCountries(b *testing.B) {
	config := &extractors.ExtractorConfig{
		Countries: []string{"US", "UK", "France"},
	}
	extractor := NewExtractor(config)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, err := extractor.Extract(benchmarkText)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark individual extraction functions
func BenchmarkExtractEmails(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		ExtractEmails(benchmarkText)
	}
}

func BenchmarkExtractPhonesUS(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		ExtractPhonesUS(benchmarkText)
	}
}

func BenchmarkExtractCreditCards(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		ExtractCreditCards(benchmarkText)
	}
}

// Benchmark context extraction specifically
func BenchmarkExtractContext(b *testing.B) {
	text := benchmarkText
	start := 50
	end := 70
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_ = patterns.ExtractContext(text, start, end)
	}
}

// Memory allocation test
func BenchmarkMemoryAllocations(b *testing.B) {
	extractor := NewDefaultExtractor()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		result, err := extractor.Extract(benchmarkText)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimizations
		_ = result.Total
	}
}