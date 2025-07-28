package patterns

import "regexp"

// Arabic countries-specific patterns
const (
	PostalCodeArabicPattern    = `\b\d{5}\b`
	PhoneArabicPattern         = `(?:\+(?:966|971|20|962|965|968|973|974|967)[\s\-]?)?(?:0)?[1-9]\d[\s\-]?\d{3}[\s\-]?\d{4}`
	StreetAddressArabicPattern = `(?i)(?:[^\x00-\x7F]+\s*)+(?:شارع|طريق|حي|منطقة|مدينة|قرية|ميدان|كورنيش|جسر|نفق|ساحة|حديقة|مجمع|برج|عمارة|بناية|فيلا|شقة|رقم|ص\.ب)\s*(?:\d+)?`
)

// Arabic countries-specific compiled patterns
var (
	PostalCodeArabicRegex    = regexp.MustCompile(PostalCodeArabicPattern)
	PhoneArabicRegex         = regexp.MustCompile(PhoneArabicPattern)
	StreetAddressArabicRegex = regexp.MustCompile(StreetAddressArabicPattern)
)

// Arabic countries-specific convenience functions
var PostalCodesArabic = func(text string) []string { return Match(text, PostalCodeArabicRegex) }
var PhonesArabic = func(text string) []string { return Match(text, PhoneArabicRegex) }
var StreetAddressesArabic = func(text string) []string { return MatchAddresses(text, StreetAddressArabicRegex) }