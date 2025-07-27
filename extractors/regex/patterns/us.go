package patterns

import "regexp"

// US-specific patterns
const (
	PhoneUSPattern          = `\b(?:(?:\+?1[-.\s]?)?(?:\([2-9]\d{2}\)|[2-9]\d{2})[-.\s]?[2-9]\d{2}[-.\s]?\d{4})\b`
	PhonesWithExtsUSPattern = `(?i)(?:(?:\+?1\s*(?:[.-]\s*)?)?(?:\(\s*(?:[2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9])\s*\)|(?:[2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9]))\s*(?:[.-]\s*)?)?(?:[2-9]1[02-9]|[2-9][02-9]1|[2-9][02-9]{2})\s*(?:[.-]\s*)?(?:[0-9]{4})(?:\s*(?:#|x\.?|ext\.?|extension)\s*(?:\d+)?)`
	StreetAddressUSPattern  = `(?i)\d{1,4}\s+[a-z\s]+?\s+(?:street|st|avenue|ave|road|rd|highway|hwy|square|sq|trail|trl|drive|dr|court|ct|park|parkway|pkwy|circle|cir|boulevard|blvd)\b`
	ZipCodeUSPattern        = `\b\d{5}(?:[-\s]\d{4})?\b`
	PoBoxUSPattern          = `(?i)P\.? ?O\.? Box \d+`
	SSNUSPattern            = `(?:\d{3}-\d{2}-\d{4})`
)

// US-specific compiled patterns
var (
	PhoneUSRegex          = regexp.MustCompile(PhoneUSPattern)
	PhonesWithExtsUSRegex = regexp.MustCompile(PhonesWithExtsUSPattern)
	StreetAddressUSRegex  = regexp.MustCompile(StreetAddressUSPattern)
	ZipCodeUSRegex        = regexp.MustCompile(ZipCodeUSPattern)
	PoBoxUSRegex          = regexp.MustCompile(PoBoxUSPattern)
	SSNUSRegex            = regexp.MustCompile(SSNUSPattern)
)

// US-specific convenience functions
var PhonesUS = func(text string) []string { return Match(text, PhoneUSRegex) }
var PhonesWithExtsUS = func(text string) []string { return Match(text, PhonesWithExtsUSRegex) }
var StreetAddressesUS = func(text string) []string { return Match(text, StreetAddressUSRegex) }
var ZipCodesUS = func(text string) []string { return Match(text, ZipCodeUSRegex) }
var PoBoxesUS = func(text string) []string { return Match(text, PoBoxUSRegex) }
var SSNsUS = func(text string) []string { return Match(text, SSNUSRegex) }
