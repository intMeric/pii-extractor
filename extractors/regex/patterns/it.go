package patterns

import "regexp"

// Italy-specific patterns
const (
	PostalCodeItalyPattern    = `\b(?:0[0-9]|[1-9]\d)\d{3}\b`
	StreetAddressItalyPattern = `(?i)\b\d{1,4}\s+(?:via|viale|piazza|corso|largo|strada|vicolo|piazzale|lungotevere|circonvallazione|passeggiata|salita|discesa|scalinata|rampa)\s+(?:del\s+|della\s+|dei\s+|delle\s+|di\s+)?[a-zàèéìíîòóùú\-']+(?:\s+[a-zàèéìíîòóùú\-']+){0,2}`
)

// Italy-specific compiled patterns
var (
	PostalCodeItalyRegex    = regexp.MustCompile(PostalCodeItalyPattern)
	StreetAddressItalyRegex = regexp.MustCompile(StreetAddressItalyPattern)
)

// Italy-specific convenience functions
var PostalCodesItaly = func(text string) []string { return Match(text, PostalCodeItalyRegex) }
var StreetAddressesItaly = func(text string) []string { return MatchAddresses(text, StreetAddressItalyRegex) }
