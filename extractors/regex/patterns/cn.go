package patterns

import "regexp"

// China-specific patterns
const (
	PostalCodeChinaPattern    = `\b[1-9]\d{5}\b`
	PhoneChinaPattern         = `(?:\+86\s?|0)?1[3-9]\d[\s\-]?\d{4}[\s\-]?\d{4}`
	StreetAddressChinaPattern = `(?i)(?:[^\x00-\x7F]+(?:市|省|区|县|镇|村|街道|路|街|巷|号|弄|里|园|庄|苑|大厦|大楼|中心|广场|公园)+)+`
)

// China-specific compiled patterns
var (
	PostalCodeChinaRegex    = regexp.MustCompile(PostalCodeChinaPattern)
	PhoneChinaRegex         = regexp.MustCompile(PhoneChinaPattern)
	StreetAddressChinaRegex = regexp.MustCompile(StreetAddressChinaPattern)
)

// China-specific convenience functions
var PostalCodesChina = func(text string) []string { return Match(text, PostalCodeChinaRegex) }
var PhonesChina = func(text string) []string { return Match(text, PhoneChinaRegex) }
var StreetAddressesChina = func(text string) []string { return MatchAddresses(text, StreetAddressChinaRegex) }