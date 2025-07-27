package patterns

import "regexp"

// UK-specific patterns
const (
	PostalCodeUKPattern    = `(?i)\b([A-Z]{1,2}\d[A-Z\d]?\s?\d[A-Z]{2})\b`
	StreetAddressUKPattern = `(?i)\b\d{1,4}[a-z]?\s+[a-z\-]+(?:\s+[a-z\-]+)*\s+(?:street|st|road|rd|lane|ln|avenue|ave|place|pl|square|sq|crescent|cres|close|cl|way|drive|dr|court|ct|terrace|ter|gardens|gdns|mews|hill|park|green|common|grove|rise|view|walk|bridge|manor|vale|row|circus|gate|heights|fields|meadow|cottage|house|villa|lodge|chambers|buildings|flats|towers|hall)\b`
)

// UK-specific compiled patterns
var (
	PostalCodeUKRegex    = regexp.MustCompile(PostalCodeUKPattern)
	StreetAddressUKRegex = regexp.MustCompile(StreetAddressUKPattern)
)

// UK-specific convenience functions
var PostalCodesUK = func(text string) []string { return Match(text, PostalCodeUKRegex) }
var StreetAddressesUK = func(text string) []string { return Match(text, StreetAddressUKRegex) }
