package patterns

import "regexp"

// Germany-specific patterns
const (
	PostalCodeGermanyPattern    = `\b(?:0[1-9]|[1-9]\d)\d{3}\b`
	PhoneGermanyPattern         = `(?:\+49\s?|0)(?:\(\d{2,5}\)|\d{2,5})[\s\-]?\d{6,10}`
	StreetAddressGermanyPattern = `(?i)\b\d{1,4}[a-z]?\s+(?:[a-züäöß\-']+\s+)*(?:straße|str\.|platz|weg|allee|gasse|ring|damm|chaussee|ufer|promenade|avenue|boulevard)\b`
)

// Germany-specific compiled patterns
var (
	PostalCodeGermanyRegex    = regexp.MustCompile(PostalCodeGermanyPattern)
	PhoneGermanyRegex         = regexp.MustCompile(PhoneGermanyPattern)
	StreetAddressGermanyRegex = regexp.MustCompile(StreetAddressGermanyPattern)
)

// Germany-specific convenience functions
var PostalCodesGermany = func(text string) []string { return Match(text, PostalCodeGermanyRegex) }
var PhonesGermany = func(text string) []string { return Match(text, PhoneGermanyRegex) }
var StreetAddressesGermany = func(text string) []string { return MatchAddresses(text, StreetAddressGermanyRegex) }