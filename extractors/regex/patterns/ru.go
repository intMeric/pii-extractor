package patterns

import "regexp"

// Russia-specific patterns
const (
	PostalCodeRussiaPattern    = `\b[1-6]\d{5}\b`
	PhoneRussiaPattern         = `(?:\+7|8)[\s\-]?\(?(?:3[0-9][0-9]|4[0-9][0-9]|8[0-9][0-9]|9[0-9][0-9])\)?[\s\-]?\d{3}[\s\-]?\d{2}[\s\-]?\d{2}`
	StreetAddressRussiaPattern = `(?i)(?:[^\x00-\x7F]+\s*)+(?:улица|ул\.|проспект|пр\.|переулок|пер\.|площадь|пл\.|набережная|наб\.|бульвар|б-р|шоссе|ш\.|тракт|дорога|линия|аллея|тупик|проезд|спуск|подъем|мост|км|дом|д\.|корпус|корп\.|строение|стр\.|квартира|кв\.)\s*(?:\d+[а-я]?)?`
)

// Russia-specific compiled patterns
var (
	PostalCodeRussiaRegex    = regexp.MustCompile(PostalCodeRussiaPattern)
	PhoneRussiaRegex         = regexp.MustCompile(PhoneRussiaPattern)
	StreetAddressRussiaRegex = regexp.MustCompile(StreetAddressRussiaPattern)
)

// Russia-specific convenience functions
var PostalCodesRussia = func(text string) []string { return Match(text, PostalCodeRussiaRegex) }
var PhonesRussia = func(text string) []string { return Match(text, PhoneRussiaRegex) }
var StreetAddressesRussia = func(text string) []string { return MatchAddresses(text, StreetAddressRussiaRegex) }