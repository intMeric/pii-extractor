package patterns

import "regexp"

// France-specific patterns
const (
	PostalCodeFrancePattern    = `\b(?:0[1-9]|[1-8]\d|9[0-8])\d{3}\b`
	StreetAddressFrancePattern = `(?i)\b\d{1,4}\s+(?:rue|avenue|boulevard|place|impasse|allée|cours|quai|passage|square|villa|cité|résidence|hameau|chemin|route|voie|esplanade|promenade|parvis|mail|galerie|sentier|traverse|venelle)\s+(?:de\s+)?(?:la\s+|le\s+|les\s+|du\s+|des\s+)?[a-zéèàçôöùûîôâêë\-']+(?:\s+[a-zéèàçôöùûîôâêë\-']+){0,2}`
)

// France-specific compiled patterns
var (
	PostalCodeFranceRegex    = regexp.MustCompile(PostalCodeFrancePattern)
	StreetAddressFranceRegex = regexp.MustCompile(StreetAddressFrancePattern)
)

// France-specific convenience functions
var PostalCodesFrance = func(text string) []string { return Match(text, PostalCodeFranceRegex) }
var StreetAddressesFrance = func(text string) []string { return MatchAddresses(text, StreetAddressFranceRegex) }
