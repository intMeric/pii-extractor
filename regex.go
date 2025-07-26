package piiextractor

import (
	"regexp"
	"strings"
	"unicode"
)

const (
	// US-specific patterns
	PhoneUSPattern          = `(?:(?:\+?\d{1,3}[-.\s*]?)?(?:\(?\d{3}\)?[-.\s*]?)?\d{3}[-.\s*]?\d{4,6})|(?:(?:(?:\(\+?\d{2}\))|(?:\+?\d{2}))\s*\d{2}\s*\d{3}\s*\d{4})`
	PhonesWithExtsUSPattern = `(?i)(?:(?:\+?1\s*(?:[.-]\s*)?)?(?:\(\s*(?:[2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9])\s*\)|(?:[2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9]))\s*(?:[.-]\s*)?)?(?:[2-9]1[02-9]|[2-9][02-9]1|[2-9][02-9]{2})\s*(?:[.-]\s*)?(?:[0-9]{4})(?:\s*(?:#|x\.?|ext\.?|extension)\s*(?:\d+)?)`
	StreetAddressUSPattern  = `(?i)\d{1,4}\s+[a-z\s]+?\s+(?:street|st|avenue|ave|road|rd|highway|hwy|square|sq|trail|trl|drive|dr|court|ct|park|parkway|pkwy|circle|cir|boulevard|blvd)\b`
	ZipCodeUSPattern        = `\b\d{5}(?:[-\s]\d{4})?\b`
	PoBoxUSPattern          = `(?i)P\.? ?O\.? Box \d+`
	SSNUSPattern            = `(?:\d{3}-\d{2}-\d{4})`

	// International postal code patterns
	PostalCodeUKPattern     = `(?i)\b([A-Z]{1,2}\d[A-Z\d]?\s?\d[A-Z]{2})\b`
	PostalCodeFrancePattern = `\b(?:0[1-9]|[1-8]\d|9[0-8])\d{3}\b`
	PostalCodeSpainPattern  = `\b(?:0[1-9]|[1-4]\d|5[0-2])\d{3}\b`
	PostalCodeItalyPattern  = `\b(?:0[0-9]|[1-9]\d)\d{3}\b`

	// International street address patterns
	StreetAddressUKPattern     = `(?i)\b\d{1,4}[a-z]?\s+[a-z\-]+(?:\s+[a-z\-]+)*\s+(?:street|st|road|rd|lane|ln|avenue|ave|place|pl|square|sq|crescent|cres|close|cl|way|drive|dr|court|ct|terrace|ter|gardens|gdns|mews|hill|park|green|common|grove|rise|view|walk|bridge|manor|vale|row|circus|gate|heights|fields|meadow|cottage|house|villa|lodge|chambers|buildings|flats|towers|hall)\b`
	StreetAddressFrancePattern = `(?i)\b\d{1,4}\s+(?:rue|avenue|boulevard|place|impasse|allée|cours|quai|passage|square|villa|cité|résidence|hameau|chemin|route|voie|esplanade|promenade|parvis|mail|galerie|sentier|traverse|venelle)\s+(?:de\s+)?(?:la\s+|le\s+|les\s+|du\s+|des\s+)?[a-z\-'éèàçôöùûîôâêë]+(?:\s+[a-z\-'éèàçôöùûîôâêë]+){0,2}`
	StreetAddressSpainPattern  = `(?i)\b\d{1,4}\s+(?:calle|avenida|plaza|paseo|ronda|travesía|glorieta|carretera|camino|vía|callejón|callejuela|costanilla|corredera|rambla|alameda|boulevard|pasaje)\s+(?:de\s+)?(?:la\s+|el\s+|los\s+|las\s+|del\s+|de\s+los\s+|de\s+las\s+)?[a-z\-'ñáéíóúü]+(?:\s+[a-z\-'ñáéíóúü]+){0,2}`
	StreetAddressItalyPattern  = `(?i)\b\d{1,4}\s+(?:via|viale|piazza|corso|largo|strada|vicolo|piazzale|lungotevere|circonvallazione|passeggiata|salita|discesa|scalinata|rampa)\s+(?:del\s+|della\s+|dei\s+|delle\s+|di\s+)?[a-z\-'àèéìíîòóùú]+(?:\s+[a-z\-'àèéìíîòóùú]+){0,2}`

	// International/generic patterns
	EmailPattern          = `(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`
	IPv4Pattern           = `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	IPv6Pattern           = `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`
	IPPattern             = IPv4Pattern + `|` + IPv6Pattern
	CreditCardPattern     = `\b(?:(?:\d{4}[\s-]?){3}\d{4}|\d{15,16})\b`
	VISACreditCardPattern = `4\d{3}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}`
	MCCreditCardPattern   = `5[1-5]\d{2}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}`
	BtcAddressPattern     = `\b[13][a-km-zA-HJ-NP-Z1-9]{25,34}\b`
	IBANPattern           = `\b[A-Z]{2}\d{2}[A-Z0-9]{4,}\d{7,}[A-Z0-9]*\b`
)

// Compiled regex patterns using MustCompile
var (
	// US-specific compiled patterns
	PhoneUSRegex          = regexp.MustCompile(PhoneUSPattern)
	PhonesWithExtsUSRegex = regexp.MustCompile(PhonesWithExtsUSPattern)
	StreetAddressUSRegex  = regexp.MustCompile(StreetAddressUSPattern)
	ZipCodeUSRegex        = regexp.MustCompile(ZipCodeUSPattern)
	PoBoxUSRegex          = regexp.MustCompile(PoBoxUSPattern)
	SSNUSRegex            = regexp.MustCompile(SSNUSPattern)

	// International postal code compiled patterns
	PostalCodeUKRegex     = regexp.MustCompile(PostalCodeUKPattern)
	PostalCodeFranceRegex = regexp.MustCompile(PostalCodeFrancePattern)
	PostalCodeSpainRegex  = regexp.MustCompile(PostalCodeSpainPattern)
	PostalCodeItalyRegex  = regexp.MustCompile(PostalCodeItalyPattern)

	// International street address compiled patterns
	StreetAddressUKRegex     = regexp.MustCompile(StreetAddressUKPattern)
	StreetAddressFranceRegex = regexp.MustCompile(StreetAddressFrancePattern)
	StreetAddressSpainRegex  = regexp.MustCompile(StreetAddressSpainPattern)
	StreetAddressItalyRegex  = regexp.MustCompile(StreetAddressItalyPattern)

	// International/generic compiled patterns
	EmailRegex          = regexp.MustCompile(EmailPattern)
	IPv4Regex           = regexp.MustCompile(IPv4Pattern)
	IPv6Regex           = regexp.MustCompile(IPv6Pattern)
	IPRegex             = regexp.MustCompile(IPPattern)
	CreditCardRegex     = regexp.MustCompile(CreditCardPattern)
	VISACreditCardRegex = regexp.MustCompile(VISACreditCardPattern)
	MCCreditCardRegex   = regexp.MustCompile(MCCreditCardPattern)
	BtcAddressRegex     = regexp.MustCompile(BtcAddressPattern)
	IBANRegex           = regexp.MustCompile(IBANPattern)
)

func match(text string, regex *regexp.Regexp) []string {
	parsed := regex.FindAllString(text, -1)
	if parsed == nil {
		return []string{}
	}
	return parsed
}

// matchWithIndices returns matches along with their start and end positions
func matchWithIndices(text string, regex *regexp.Regexp) [][]int {
	return regex.FindAllStringIndex(text, -1)
}

// extractContext extracts the context around a match, prioritizing full sentences over word count
func extractContext(text string, start, end int) string {
	// First try to find a complete sentence
	sentence := extractSentence(text, start, end)
	if sentence != "" {
		return strings.TrimSpace(sentence)
	}

	// Fallback to 8 words before and after
	return extractWordContext(text, start, end)
}

// extractSentence tries to extract a complete sentence containing the match
func extractSentence(text string, start, end int) string {
	// Find sentence boundaries (., !, ?, or start/end of text)
	sentenceStart := start
	sentenceEnd := end

	// Look backwards for sentence start
	for i := start - 1; i >= 0; i-- {
		char := text[i]
		if char == '.' || char == '!' || char == '?' {
			sentenceStart = i + 1
			break
		}
		if i == 0 {
			sentenceStart = 0
		}
	}

	// Look forwards for sentence end
	for i := end; i < len(text); i++ {
		char := text[i]
		if char == '.' || char == '!' || char == '?' {
			sentenceEnd = i + 1
			break
		}
		if i == len(text)-1 {
			sentenceEnd = len(text)
		}
	}

	// Skip whitespace at the beginning
	for sentenceStart < len(text) && unicode.IsSpace(rune(text[sentenceStart])) {
		sentenceStart++
	}

	if sentenceStart < sentenceEnd && sentenceEnd <= len(text) {
		return text[sentenceStart:sentenceEnd]
	}

	return ""
}

// extractWordContext extracts 8 words before and after the match
func extractWordContext(text string, start, end int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	// Find the word index that contains our match
	wordStart := -1
	wordEnd := -1
	currentPos := 0

	for i, word := range words {
		wordStartPos := strings.Index(text[currentPos:], word) + currentPos
		wordEndPos := wordStartPos + len(word)

		if wordStartPos <= start && start < wordEndPos {
			wordStart = i
		}
		if wordStartPos < end && end <= wordEndPos {
			wordEnd = i
		}

		currentPos = wordEndPos

		if wordStart != -1 && wordEnd != -1 {
			break
		}
	}

	if wordStart == -1 || wordEnd == -1 {
		return ""
	}

	// Extract 8 words before and after
	contextStart := max(0, wordStart-8)
	contextEnd := min(len(words), wordEnd+8+1)

	return strings.Join(words[contextStart:contextEnd], " ")
}

// US-specific convenience functions
var PhonesUS = func(text string) []string { return match(text, PhoneUSRegex) }
var PhonesWithExtsUS = func(text string) []string { return match(text, PhonesWithExtsUSRegex) }
var StreetAddressesUS = func(text string) []string { return match(text, StreetAddressUSRegex) }
var ZipCodesUS = func(text string) []string { return match(text, ZipCodeUSRegex) }
var PoBoxesUS = func(text string) []string { return match(text, PoBoxUSRegex) }
var SSNsUS = func(text string) []string { return match(text, SSNUSRegex) }

// International postal code convenience functions
var PostalCodesUK = func(text string) []string { return match(text, PostalCodeUKRegex) }
var PostalCodesFrance = func(text string) []string { return match(text, PostalCodeFranceRegex) }
var PostalCodesSpain = func(text string) []string { return match(text, PostalCodeSpainRegex) }
var PostalCodesItaly = func(text string) []string { return match(text, PostalCodeItalyRegex) }

// International street address convenience functions
var StreetAddressesUK = func(text string) []string { return match(text, StreetAddressUKRegex) }
var StreetAddressesFrance = func(text string) []string { return match(text, StreetAddressFranceRegex) }
var StreetAddressesSpain = func(text string) []string { return match(text, StreetAddressSpainRegex) }
var StreetAddressesItaly = func(text string) []string { return match(text, StreetAddressItalyRegex) }

// International/generic convenience functions
var Emails = func(text string) []string { return match(text, EmailRegex) }
var IPv4s = func(text string) []string { return match(text, IPv4Regex) }
var IPv6s = func(text string) []string { return match(text, IPv6Regex) }
var IPs = func(text string) []string { return match(text, IPRegex) }
var CreditCards = func(text string) []string { return match(text, CreditCardRegex) }
var VISACreditCards = func(text string) []string { return match(text, VISACreditCardRegex) }
var MCCreditCards = func(text string) []string { return match(text, MCCreditCardRegex) }
var BtcAddresses = func(text string) []string { return match(text, BtcAddressRegex) }
var IBANs = func(text string) []string { return match(text, IBANRegex) }

// Structured extraction functions that return value objects

// extractWithContext is a generic function for extracting PII with context and counting
func extractWithContext[T any](text string, regex *regexp.Regexp, createItem func(value string, context string) T, updateItem func(item *T, context string)) []T {
	indices := matchWithIndices(text, regex)
	itemMap := make(map[string]*T)

	for _, idx := range indices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		context := extractContext(text, start, end)

		if item, exists := itemMap[value]; exists {
			updateItem(item, context)
		} else {
			newItem := createItem(value, context)
			itemMap[value] = &newItem
		}
	}

	items := make([]T, 0, len(itemMap))
	for _, item := range itemMap {
		items = append(items, *item)
	}
	return items
}

// ExtractPhonesUS extracts US phone numbers as PiiEntity objects with context
func ExtractPhonesUS(text string) []PiiEntity {
	phones := extractWithContext(text, PhoneUSRegex,
		func(value, context string) Phone {
			return Phone{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "US",
			}
		},
		func(phone *Phone, context string) {
			phone.BasePii.IncrementCount()
			phone.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, phone := range phones {
		entities = append(entities, PiiEntity{
			Type:  PiiTypePhone,
			Value: phone,
		})
	}
	return entities
}

// ExtractEmails extracts email addresses as PiiEntity objects with context
func ExtractEmails(text string) []PiiEntity {
	emails := extractWithContext(text, EmailRegex,
		func(value, context string) Email {
			return Email{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
			}
		},
		func(email *Email, context string) {
			email.BasePii.IncrementCount()
			email.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, email := range emails {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeEmail,
			Value: email,
		})
	}
	return entities
}

// ExtractSSNsUS extracts US SSNs as PiiEntity objects with context
func ExtractSSNsUS(text string) []PiiEntity {
	ssns := extractWithContext(text, SSNUSRegex,
		func(value, context string) SSN {
			return SSN{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "US",
			}
		},
		func(ssn *SSN, context string) {
			ssn.BasePii.IncrementCount()
			ssn.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, ssn := range ssns {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeSSN,
			Value: ssn,
		})
	}
	return entities
}

// ExtractZipCodesUS extracts US zip codes as PiiEntity objects with context
func ExtractZipCodesUS(text string) []PiiEntity {
	zipCodes := extractWithContext(text, ZipCodeUSRegex,
		func(value, context string) ZipCode {
			return ZipCode{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "US",
			}
		},
		func(zipCode *ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, zipCode := range zipCodes {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractStreetAddressesUS extracts US street addresses as PiiEntity objects with context
func ExtractStreetAddressesUS(text string) []PiiEntity {
	addresses := extractWithContext(text, StreetAddressUSRegex,
		func(value, context string) StreetAddress {
			return StreetAddress{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "US",
			}
		},
		func(address *StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, address := range addresses {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// ExtractCreditCards extracts credit cards as PiiEntity objects with context
func ExtractCreditCards(text string) []PiiEntity {
	cardMap := make(map[string]*CreditCard)

	// Check for VISA cards
	visaIndices := matchWithIndices(text, VISACreditCardRegex)
	for _, idx := range visaIndices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		context := extractContext(text, start, end)

		if card, exists := cardMap[value]; exists {
			card.BasePii.IncrementCount()
			card.BasePii.AddContext(context)
		} else {
			cardMap[value] = &CreditCard{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Type: "visa",
			}
		}
	}

	// Check for MasterCard
	mcIndices := matchWithIndices(text, MCCreditCardRegex)
	for _, idx := range mcIndices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		context := extractContext(text, start, end)

		if card, exists := cardMap[value]; exists {
			card.BasePii.IncrementCount()
			card.BasePii.AddContext(context)
		} else {
			cardMap[value] = &CreditCard{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Type: "mastercard",
			}
		}
	}

	// Check for generic cards (not already found as VISA/MC)
	genericIndices := matchWithIndices(text, CreditCardRegex)
	for _, idx := range genericIndices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		context := extractContext(text, start, end)

		// Skip if already found as VISA or MC
		if _, exists := cardMap[value]; !exists {
			cardMap[value] = &CreditCard{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Type: "generic",
			}
		}
	}

	var entities []PiiEntity
	for _, card := range cardMap {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeCreditCard,
			Value: *card,
		})
	}
	return entities
}

// ExtractIPAddresses extracts IP addresses as PiiEntity objects with context
func ExtractIPAddresses(text string) []PiiEntity {
	ipMap := make(map[string]*IPAddress)

	// Extract IPv4
	ipv4Indices := matchWithIndices(text, IPv4Regex)
	for _, idx := range ipv4Indices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		context := extractContext(text, start, end)

		if ip, exists := ipMap[value]; exists {
			ip.BasePii.IncrementCount()
			ip.BasePii.AddContext(context)
		} else {
			ipMap[value] = &IPAddress{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Version: "ipv4",
			}
		}
	}

	// Extract IPv6
	ipv6Indices := matchWithIndices(text, IPv6Regex)
	for _, idx := range ipv6Indices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		context := extractContext(text, start, end)

		if ip, exists := ipMap[value]; exists {
			ip.BasePii.IncrementCount()
			ip.BasePii.AddContext(context)
		} else {
			ipMap[value] = &IPAddress{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Version: "ipv6",
			}
		}
	}

	var entities []PiiEntity
	for _, ip := range ipMap {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeIPAddress,
			Value: *ip,
		})
	}
	return entities
}

// ExtractBtcAddresses extracts Bitcoin addresses as PiiEntity objects with context
func ExtractBtcAddresses(text string) []PiiEntity {
	btcAddresses := extractWithContext(text, BtcAddressRegex,
		func(value, context string) BtcAddress {
			return BtcAddress{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
			}
		},
		func(btc *BtcAddress, context string) {
			btc.BasePii.IncrementCount()
			btc.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, btcAddress := range btcAddresses {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeBtcAddress,
			Value: btcAddress,
		})
	}
	return entities
}

// ExtractPoBoxesUS extracts US P.O. Boxes as PiiEntity objects with context
func ExtractPoBoxesUS(text string) []PiiEntity {
	poBoxes := extractWithContext(text, PoBoxUSRegex,
		func(value, context string) PoBox {
			return PoBox{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "US",
			}
		},
		func(poBox *PoBox, context string) {
			poBox.BasePii.IncrementCount()
			poBox.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, poBox := range poBoxes {
		entities = append(entities, PiiEntity{
			Type:  PiiTypePoBox,
			Value: poBox,
		})
	}
	return entities
}

// ExtractIBANs extracts IBANs as PiiEntity objects with context
func ExtractIBANs(text string) []PiiEntity {
	ibans := extractWithContext(text, IBANRegex,
		func(value, context string) IBAN {
			country := ""
			if len(value) >= 2 {
				country = value[:2]
			}
			return IBAN{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: country,
			}
		},
		func(iban *IBAN, context string) {
			iban.BasePii.IncrementCount()
			iban.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, iban := range ibans {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeIBAN,
			Value: iban,
		})
	}
	return entities
}

// International PII extraction functions

// ExtractPostalCodesUK extracts UK postal codes as PiiEntity objects with context
func ExtractPostalCodesUK(text string) []PiiEntity {
	postalCodes := extractWithContext(text, PostalCodeUKRegex,
		func(value, context string) ZipCode {
			return ZipCode{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "UK",
			}
		},
		func(zipCode *ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractPostalCodesFrance extracts France postal codes as PiiEntity objects with context
func ExtractPostalCodesFrance(text string) []PiiEntity {
	postalCodes := extractWithContext(text, PostalCodeFranceRegex,
		func(value, context string) ZipCode {
			return ZipCode{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "France",
			}
		},
		func(zipCode *ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractPostalCodesSpain extracts Spain postal codes as PiiEntity objects with context
func ExtractPostalCodesSpain(text string) []PiiEntity {
	postalCodes := extractWithContext(text, PostalCodeSpainRegex,
		func(value, context string) ZipCode {
			return ZipCode{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Spain",
			}
		},
		func(zipCode *ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractPostalCodesItaly extracts Italy postal codes as PiiEntity objects with context
func ExtractPostalCodesItaly(text string) []PiiEntity {
	postalCodes := extractWithContext(text, PostalCodeItalyRegex,
		func(value, context string) ZipCode {
			return ZipCode{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Italy",
			}
		},
		func(zipCode *ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractStreetAddressesUK extracts UK street addresses as PiiEntity objects with context
func ExtractStreetAddressesUK(text string) []PiiEntity {
	addresses := extractWithContext(text, StreetAddressUKRegex,
		func(value, context string) StreetAddress {
			return StreetAddress{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "UK",
			}
		},
		func(address *StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, address := range addresses {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// ExtractStreetAddressesFrance extracts France street addresses as PiiEntity objects with context
func ExtractStreetAddressesFrance(text string) []PiiEntity {
	addresses := extractWithContext(text, StreetAddressFranceRegex,
		func(value, context string) StreetAddress {
			return StreetAddress{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "France",
			}
		},
		func(address *StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, address := range addresses {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// ExtractStreetAddressesSpain extracts Spain street addresses as PiiEntity objects with context
func ExtractStreetAddressesSpain(text string) []PiiEntity {
	addresses := extractWithContext(text, StreetAddressSpainRegex,
		func(value, context string) StreetAddress {
			return StreetAddress{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Spain",
			}
		},
		func(address *StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, address := range addresses {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// ExtractStreetAddressesItaly extracts Italy street addresses as PiiEntity objects with context
func ExtractStreetAddressesItaly(text string) []PiiEntity {
	addresses := extractWithContext(text, StreetAddressItalyRegex,
		func(value, context string) StreetAddress {
			return StreetAddress{
				BasePii: BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Italy",
			}
		},
		func(address *StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []PiiEntity
	for _, address := range addresses {
		entities = append(entities, PiiEntity{
			Type:  PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}
