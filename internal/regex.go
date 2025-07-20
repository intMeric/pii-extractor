package internal

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

// ExtractPhonesUS extracts US phone numbers as Phone value objects with context
func ExtractPhonesUS(text string) []Phone {
	return extractWithContext(text, PhoneUSRegex,
		func(value, context string) Phone {
			return Phone{
				Value:    value,
				Country:  "US",
				Contexts: []string{context},
				Count:    1,
			}
		},
		func(phone *Phone, context string) {
			phone.Count++
			phone.Contexts = append(phone.Contexts, context)
		})
}

// ExtractEmails extracts email addresses as Email value objects with context
func ExtractEmails(text string) []Email {
	return extractWithContext(text, EmailRegex,
		func(value, context string) Email {
			return Email{
				Value:    value,
				Contexts: []string{context},
				Count:    1,
			}
		},
		func(email *Email, context string) {
			email.Count++
			email.Contexts = append(email.Contexts, context)
		})
}

// ExtractSSNsUS extracts US SSNs as SSN value objects with context
func ExtractSSNsUS(text string) []SSN {
	return extractWithContext(text, SSNUSRegex,
		func(value, context string) SSN {
			return SSN{
				Value:    value,
				Country:  "US",
				Contexts: []string{context},
				Count:    1,
			}
		},
		func(ssn *SSN, context string) {
			ssn.Count++
			ssn.Contexts = append(ssn.Contexts, context)
		})
}

// ExtractZipCodesUS extracts US zip codes as ZipCode value objects with context
func ExtractZipCodesUS(text string) []ZipCode {
	return extractWithContext(text, ZipCodeUSRegex,
		func(value, context string) ZipCode {
			return ZipCode{
				Value:    value,
				Country:  "US",
				Contexts: []string{context},
				Count:    1,
			}
		},
		func(zipCode *ZipCode, context string) {
			zipCode.Count++
			zipCode.Contexts = append(zipCode.Contexts, context)
		})
}

// ExtractStreetAddressesUS extracts US street addresses as StreetAddress value objects with context
func ExtractStreetAddressesUS(text string) []StreetAddress {
	return extractWithContext(text, StreetAddressUSRegex,
		func(value, context string) StreetAddress {
			return StreetAddress{
				Value:    value,
				Country:  "US",
				Contexts: []string{context},
				Count:    1,
			}
		},
		func(address *StreetAddress, context string) {
			address.Count++
			address.Contexts = append(address.Contexts, context)
		})
}

// ExtractCreditCards extracts credit cards as CreditCard value objects with context
func ExtractCreditCards(text string) []CreditCard {
	cardMap := make(map[string]*CreditCard)

	// Check for VISA cards
	visaIndices := matchWithIndices(text, VISACreditCardRegex)
	for _, idx := range visaIndices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		context := extractContext(text, start, end)

		if card, exists := cardMap[value]; exists {
			card.Count++
			card.Contexts = append(card.Contexts, context)
		} else {
			cardMap[value] = &CreditCard{
				Value:    value,
				Type:     "visa",
				Contexts: []string{context},
				Count:    1,
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
			card.Count++
			card.Contexts = append(card.Contexts, context)
		} else {
			cardMap[value] = &CreditCard{
				Value:    value,
				Type:     "mastercard",
				Contexts: []string{context},
				Count:    1,
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
				Value:    value,
				Type:     "generic",
				Contexts: []string{context},
				Count:    1,
			}
		}
	}

	cards := make([]CreditCard, 0, len(cardMap))
	for _, card := range cardMap {
		cards = append(cards, *card)
	}
	return cards
}

// ExtractIPAddresses extracts IP addresses as IPAddress value objects with context
func ExtractIPAddresses(text string) []IPAddress {
	ipMap := make(map[string]*IPAddress)

	// Extract IPv4
	ipv4Indices := matchWithIndices(text, IPv4Regex)
	for _, idx := range ipv4Indices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		context := extractContext(text, start, end)

		if ip, exists := ipMap[value]; exists {
			ip.Count++
			ip.Contexts = append(ip.Contexts, context)
		} else {
			ipMap[value] = &IPAddress{
				Value:    value,
				Version:  "ipv4",
				Contexts: []string{context},
				Count:    1,
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
			ip.Count++
			ip.Contexts = append(ip.Contexts, context)
		} else {
			ipMap[value] = &IPAddress{
				Value:    value,
				Version:  "ipv6",
				Contexts: []string{context},
				Count:    1,
			}
		}
	}

	ips := make([]IPAddress, 0, len(ipMap))
	for _, ip := range ipMap {
		ips = append(ips, *ip)
	}
	return ips
}

// ExtractBtcAddresses extracts Bitcoin addresses as BtcAddress value objects with context
func ExtractBtcAddresses(text string) []BtcAddress {
	return extractWithContext(text, BtcAddressRegex,
		func(value, context string) BtcAddress {
			return BtcAddress{
				Value:    value,
				Contexts: []string{context},
				Count:    1,
			}
		},
		func(btc *BtcAddress, context string) {
			btc.Count++
			btc.Contexts = append(btc.Contexts, context)
		})
}

// ExtractIBANs extracts IBANs as IBAN value objects with context
func ExtractIBANs(text string) []IBAN {
	return extractWithContext(text, IBANRegex,
		func(value, context string) IBAN {
			country := ""
			if len(value) >= 2 {
				country = value[:2]
			}
			return IBAN{
				Value:    value,
				Country:  country,
				Contexts: []string{context},
				Count:    1,
			}
		},
		func(iban *IBAN, context string) {
			iban.Count++
			iban.Contexts = append(iban.Contexts, context)
		})
}
