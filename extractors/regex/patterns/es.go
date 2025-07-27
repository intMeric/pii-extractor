package patterns

import "regexp"

// Spain-specific patterns
const (
	PostalCodeSpainPattern    = `\b(?:0[1-9]|[1-4]\d|5[0-2])\d{3}\b`
	StreetAddressSpainPattern = `(?i)\b\d{1,4}\s+(?:calle|avenida|plaza|paseo|ronda|travesía|glorieta|carretera|camino|vía|callejón|callejuela|costanilla|corredera|rambla|alameda|boulevard|pasaje)\s+(?:de\s+)?(?:la\s+|el\s+|los\s+|las\s+|del\s+|de\s+los\s+|de\s+las\s+)?[a-zñáéíóúü\-']+(?:\s+[a-zñáéíóúü\-']+){0,2}`
)

// Spain-specific compiled patterns
var (
	PostalCodeSpainRegex    = regexp.MustCompile(PostalCodeSpainPattern)
	StreetAddressSpainRegex = regexp.MustCompile(StreetAddressSpainPattern)
)

// Spain-specific convenience functions
var PostalCodesSpain = func(text string) []string { return Match(text, PostalCodeSpainRegex) }
var StreetAddressesSpain = func(text string) []string { return MatchAddresses(text, StreetAddressSpainRegex) }
