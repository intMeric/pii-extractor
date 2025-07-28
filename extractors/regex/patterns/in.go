package patterns

import "regexp"

// India-specific patterns
const (
	PostalCodeIndiaPattern    = `\b[1-9]\d{5}\b`
	PhoneIndiaPattern         = `(?:\+91\s?|0)?(?:[6-9]\d{9}|11[\s\-]?\d{4}[\s\-]?\d{4})`
	StreetAddressIndiaPattern = `(?i)\b\d{1,4}[a-z]?\s+(?:[a-z\s\-']+\s+)*(?:road|rd|street|st|lane|ln|nagar|colony|sector|block|phase|plot|house|building|apartment|flat|cross|main|layout|extension|park|garden|circle|square|compound|society|residency|enclave|vihar|kunj|puram|gram|marg|path)\b`
)

// India-specific compiled patterns
var (
	PostalCodeIndiaRegex    = regexp.MustCompile(PostalCodeIndiaPattern)
	PhoneIndiaRegex         = regexp.MustCompile(PhoneIndiaPattern)
	StreetAddressIndiaRegex = regexp.MustCompile(StreetAddressIndiaPattern)
)

// India-specific convenience functions
var PostalCodesIndia = func(text string) []string { return Match(text, PostalCodeIndiaRegex) }
var PhonesIndia = func(text string) []string { return Match(text, PhoneIndiaRegex) }
var StreetAddressesIndia = func(text string) []string { return MatchAddresses(text, StreetAddressIndiaRegex) }