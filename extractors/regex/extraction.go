package regex

import (
	"regexp"
	"github.com/intMeric/pii-extractor/pii"
	patterns "github.com/intMeric/pii-extractor/extractors/regex/patterns"
)

// extractWithContext is a generic function for extracting PII with context and counting
func extractWithContext[T any](text string, regexPattern *regexp.Regexp, createItem func(value string, context string) T, updateItem func(item *T, context string)) []T {
	indices := patterns.MatchWithIndices(text, regexPattern)
	if len(indices) == 0 {
		return []T{}
	}
	
	// Pre-size map based on expected unique matches (typically 70-80% of total matches are unique)
	expectedUnique := len(indices)*4/5 + 1
	itemMap := make(map[string]*T, expectedUnique)

	// Use context cache only if we have many matches (>= 10) to amortize the cost
	var contextCache *patterns.ContextCache
	if len(indices) >= 10 {
		contextCache = patterns.NewContextCache(text)
	}

	for _, idx := range indices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		
		var context string
		if contextCache != nil {
			context = contextCache.ExtractContext(start, end)
		} else {
			context = patterns.ExtractContext(text, start, end)
		}

		if item, exists := itemMap[value]; exists {
			updateItem(item, context)
		} else {
			newItem := createItem(value, context)
			itemMap[value] = &newItem
		}
	}

	// Pre-allocate result slice with exact size
	items := make([]T, 0, len(itemMap))
	for _, item := range itemMap {
		items = append(items, *item)
	}
	return items
}

// =============================================================================
// US-SPECIFIC EXTRACTION FUNCTIONS
// =============================================================================

// ExtractPhonesUS extracts US phone numbers as PiiEntity objects with context
func ExtractPhonesUS(text string) []pii.PiiEntity {
	phones := extractWithContext(text, patterns.PhoneUSRegex,
		func(value, context string) pii.Phone {
			return pii.Phone{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "US",
			}
		},
		func(phone *pii.Phone, context string) {
			phone.BasePii.IncrementCount()
			phone.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, phone := range phones {
		// Filter out credit card false positives
		if !isCreditCardFalsePositive(phone.BasePii.Value) {
			entities = append(entities, pii.PiiEntity{
				Type:  pii.PiiTypePhone,
				Value: phone,
			})
		}
	}
	return entities
}

// isCreditCardFalsePositive checks if a phone number is actually part of a credit card or other PII
func isCreditCardFalsePositive(value string) bool {
	// Count digits and extract them without creating strings
	var digitChars [20]byte // Pre-allocated buffer for digits (max 20 digits)
	digitCount := 0
	
	for i := 0; i < len(value) && digitCount < 20; i++ {
		c := value[i]
		if c >= '0' && c <= '9' {
			digitChars[digitCount] = c
			digitCount++
		}
	}
	
	// If it's 14+ digits, it's likely a credit card or IBAN
	if digitCount >= 14 {
		return true
	}
	
	// Check if it looks like part of a credit card (starts with common prefixes)
	if digitCount >= 4 {
		// Check common credit card prefixes without string allocation
		prefix := digitChars[:4]
		if (prefix[0] == '4' && prefix[1] == '1' && prefix[2] == '1' && prefix[3] == '1') ||
		   (prefix[0] == '4' && prefix[1] == '0' && prefix[2] == '0' && prefix[3] == '0') ||
		   (prefix[0] == '5' && prefix[1] == '5' && prefix[2] == '5' && prefix[3] == '5') ||
		   (prefix[0] == '5' && prefix[1] == '1' && prefix[2] == '0' && prefix[3] == '5') ||
		   (prefix[0] == '1' && prefix[1] == '1' && prefix[2] == '1' && prefix[3] == '1') ||
		   (prefix[0] == '1' && prefix[1] == '2' && prefix[2] == '3' && prefix[3] == '4') {
			return true
		}
	}
	
	// Check if it's all the same digits (likely test data)
	if digitCount >= 4 {
		firstDigit := digitChars[0]
		allSame := true
		for i := 1; i < digitCount; i++ {
			if digitChars[i] != firstDigit {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}
	
	// Filter out sequences like 1111-1111 which are parts of credit cards
	if digitCount == 8 {
		// Check if all digits are the same
		if digitChars[0] == '1' {
			allOnes := true
			for i := 1; i < 8; i++ {
				if digitChars[i] != '1' {
					allOnes = false
					break
				}
			}
			if allOnes {
				return true
			}
		}
		
		// Check if first 4 digits match last 4 digits
		firstHalfSame := true
		for i := 0; i < 4; i++ {
			if digitChars[i] != digitChars[i+4] {
				firstHalfSame = false
				break
			}
		}
		if firstHalfSame {
			return true
		}
	}
	
	return false
}

// ExtractSSNsUS extracts US SSNs as PiiEntity objects with context
func ExtractSSNsUS(text string) []pii.PiiEntity {
	ssns := extractWithContext(text, patterns.SSNUSRegex,
		func(value, context string) pii.SSN {
			return pii.SSN{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "US",
			}
		},
		func(ssn *pii.SSN, context string) {
			ssn.BasePii.IncrementCount()
			ssn.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, ssn := range ssns {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeSSN,
			Value: ssn,
		})
	}
	return entities
}

// ExtractZipCodesUS extracts US zip codes as PiiEntity objects with context
func ExtractZipCodesUS(text string) []pii.PiiEntity {
	zipCodes := extractWithContext(text, patterns.ZipCodeUSRegex,
		func(value, context string) pii.ZipCode {
			return pii.ZipCode{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "US",
			}
		},
		func(zipCode *pii.ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, zipCode := range zipCodes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractStreetAddressesUS extracts US street addresses as PiiEntity objects with context
func ExtractStreetAddressesUS(text string) []pii.PiiEntity {
	addresses := extractWithContext(text, patterns.StreetAddressUSRegex,
		func(value, context string) pii.StreetAddress {
			return pii.StreetAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "US",
			}
		},
		func(address *pii.StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, address := range addresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// ExtractPoBoxesUS extracts US P.O. Boxes as PiiEntity objects with context
func ExtractPoBoxesUS(text string) []pii.PiiEntity {
	poBoxes := extractWithContext(text, patterns.PoBoxUSRegex,
		func(value, context string) pii.PoBox {
			return pii.PoBox{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "US",
			}
		},
		func(poBox *pii.PoBox, context string) {
			poBox.BasePii.IncrementCount()
			poBox.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, poBox := range poBoxes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypePoBox,
			Value: poBox,
		})
	}
	return entities
}

// =============================================================================
// INTERNATIONAL/GENERIC EXTRACTION FUNCTIONS
// =============================================================================

// ExtractEmails extracts email addresses as PiiEntity objects with context
func ExtractEmails(text string) []pii.PiiEntity {
	emails := extractWithContext(text, patterns.EmailRegex,
		func(value, context string) pii.Email {
			return pii.Email{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
			}
		},
		func(email *pii.Email, context string) {
			email.BasePii.IncrementCount()
			email.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, email := range emails {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeEmail,
			Value: email,
		})
	}
	return entities
}

// ExtractCreditCards extracts credit cards as PiiEntity objects with context
func ExtractCreditCards(text string) []pii.PiiEntity {
	// Estimate capacity based on typical credit card density in text
	estimatedCards := len(text)/2000 + 5 // ~1 card per 2000 chars
	cardMap := make(map[string]*pii.CreditCard, estimatedCards)

	// Check for VISA cards
	visaIndices := patterns.MatchWithIndices(text, patterns.VISACreditCardRegex)
	
	// Create context cache only if we have enough total matches to justify it
	var contextCache *patterns.ContextCache
	mcIndices := patterns.MatchWithIndices(text, patterns.MCCreditCardRegex)
	genericIndices := patterns.MatchWithIndices(text, patterns.CreditCardRegex)
	totalMatches := len(visaIndices) + len(mcIndices) + len(genericIndices)
	
	if totalMatches >= 5 {
		contextCache = patterns.NewContextCache(text)
	}

	for _, idx := range visaIndices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		var context string
		if contextCache != nil {
			context = contextCache.ExtractContext(start, end)
		} else {
			context = patterns.ExtractContext(text, start, end)
		}

		if card, exists := cardMap[value]; exists {
			card.BasePii.IncrementCount()
			card.BasePii.AddContext(context)
		} else {
			cardMap[value] = &pii.CreditCard{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Type: "visa",
			}
		}
	}

	// Check for MasterCard
	for _, idx := range mcIndices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		var context string
		if contextCache != nil {
			context = contextCache.ExtractContext(start, end)
		} else {
			context = patterns.ExtractContext(text, start, end)
		}

		if card, exists := cardMap[value]; exists {
			card.BasePii.IncrementCount()
			card.BasePii.AddContext(context)
		} else {
			cardMap[value] = &pii.CreditCard{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Type: "mastercard",
			}
		}
	}

	// Check for generic cards (not already found as VISA/MC)
	for _, idx := range genericIndices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		var context string
		if contextCache != nil {
			context = contextCache.ExtractContext(start, end)
		} else {
			context = patterns.ExtractContext(text, start, end)
		}

		// Skip if already found as VISA or MC
		if _, exists := cardMap[value]; !exists {
			cardMap[value] = &pii.CreditCard{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Type: "generic",
			}
		}
	}

	entities := make([]pii.PiiEntity, 0, len(cardMap))
	for _, card := range cardMap {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeCreditCard,
			Value: *card,
		})
	}
	return entities
}

// ExtractIPAddresses extracts IP addresses as PiiEntity objects with context
func ExtractIPAddresses(text string) []pii.PiiEntity {
	// Estimate capacity based on typical IP density in text
	estimatedIPs := len(text)/1500 + 3 // ~1 IP per 1500 chars
	ipMap := make(map[string]*pii.IPAddress, estimatedIPs)

	// Extract IPv4
	ipv4Indices := patterns.MatchWithIndices(text, patterns.IPv4Regex)
	ipv6Indices := patterns.MatchWithIndices(text, patterns.IPv6Regex)
	
	// Create context cache only if we have enough total matches to justify it
	var contextCache *patterns.ContextCache
	totalMatches := len(ipv4Indices) + len(ipv6Indices)
	
	if totalMatches >= 5 {
		contextCache = patterns.NewContextCache(text)
	}

	for _, idx := range ipv4Indices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		var context string
		if contextCache != nil {
			context = contextCache.ExtractContext(start, end)
		} else {
			context = patterns.ExtractContext(text, start, end)
		}

		if ip, exists := ipMap[value]; exists {
			ip.BasePii.IncrementCount()
			ip.BasePii.AddContext(context)
		} else {
			ipMap[value] = &pii.IPAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Version: "ipv4",
			}
		}
	}

	// Extract IPv6
	for _, idx := range ipv6Indices {
		start, end := idx[0], idx[1]
		value := text[start:end]
		var context string
		if contextCache != nil {
			context = contextCache.ExtractContext(start, end)
		} else {
			context = patterns.ExtractContext(text, start, end)
		}

		if ip, exists := ipMap[value]; exists {
			ip.BasePii.IncrementCount()
			ip.BasePii.AddContext(context)
		} else {
			ipMap[value] = &pii.IPAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Version: "ipv6",
			}
		}
	}

	entities := make([]pii.PiiEntity, 0, len(ipMap))
	for _, ip := range ipMap {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeIPAddress,
			Value: *ip,
		})
	}
	return entities
}

// ExtractBtcAddresses extracts Bitcoin addresses as PiiEntity objects with context
func ExtractBtcAddresses(text string) []pii.PiiEntity {
	btcAddresses := extractWithContext(text, patterns.BtcAddressRegex,
		func(value, context string) pii.BtcAddress {
			return pii.BtcAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
			}
		},
		func(btc *pii.BtcAddress, context string) {
			btc.BasePii.IncrementCount()
			btc.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, btcAddress := range btcAddresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeBtcAddress,
			Value: btcAddress,
		})
	}
	return entities
}

// ExtractIBANs extracts IBANs as PiiEntity objects with context
func ExtractIBANs(text string) []pii.PiiEntity {
	ibans := extractWithContext(text, patterns.IBANRegex,
		func(value, context string) pii.IBAN {
			country := ""
			if len(value) >= 2 {
				country = value[:2]
			}
			return pii.IBAN{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: country,
			}
		},
		func(iban *pii.IBAN, context string) {
			iban.BasePii.IncrementCount()
			iban.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, iban := range ibans {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeIBAN,
			Value: iban,
		})
	}
	return entities
}

// =============================================================================
// INTERNATIONAL POSTAL CODES & ADDRESSES
// =============================================================================

// --- UK PII ---

// ExtractPostalCodesUK extracts UK postal codes as PiiEntity objects with context
func ExtractPostalCodesUK(text string) []pii.PiiEntity {
	postalCodes := extractWithContext(text, patterns.PostalCodeUKRegex,
		func(value, context string) pii.ZipCode {
			return pii.ZipCode{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "UK",
			}
		},
		func(zipCode *pii.ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractStreetAddressesUK extracts UK street addresses as PiiEntity objects with context
func ExtractStreetAddressesUK(text string) []pii.PiiEntity {
	addresses := extractWithContext(text, patterns.StreetAddressUKRegex,
		func(value, context string) pii.StreetAddress {
			return pii.StreetAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "UK",
			}
		},
		func(address *pii.StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, address := range addresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// --- France PII ---

// ExtractPostalCodesFrance extracts France postal codes as PiiEntity objects with context
func ExtractPostalCodesFrance(text string) []pii.PiiEntity {
	postalCodes := extractWithContext(text, patterns.PostalCodeFranceRegex,
		func(value, context string) pii.ZipCode {
			return pii.ZipCode{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "France",
			}
		},
		func(zipCode *pii.ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractStreetAddressesFrance extracts France street addresses as PiiEntity objects with context
func ExtractStreetAddressesFrance(text string) []pii.PiiEntity {
	addresses := extractWithContext(text, patterns.StreetAddressFranceRegex,
		func(value, context string) pii.StreetAddress {
			return pii.StreetAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "France",
			}
		},
		func(address *pii.StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, address := range addresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// --- Spain PII ---

// ExtractPostalCodesSpain extracts Spain postal codes as PiiEntity objects with context
func ExtractPostalCodesSpain(text string) []pii.PiiEntity {
	postalCodes := extractWithContext(text, patterns.PostalCodeSpainRegex,
		func(value, context string) pii.ZipCode {
			return pii.ZipCode{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Spain",
			}
		},
		func(zipCode *pii.ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractStreetAddressesSpain extracts Spain street addresses as PiiEntity objects with context
func ExtractStreetAddressesSpain(text string) []pii.PiiEntity {
	addresses := extractWithContext(text, patterns.StreetAddressSpainRegex,
		func(value, context string) pii.StreetAddress {
			return pii.StreetAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Spain",
			}
		},
		func(address *pii.StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, address := range addresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// --- Italy PII ---

// ExtractPostalCodesItaly extracts Italy postal codes as PiiEntity objects with context
func ExtractPostalCodesItaly(text string) []pii.PiiEntity {
	postalCodes := extractWithContext(text, patterns.PostalCodeItalyRegex,
		func(value, context string) pii.ZipCode {
			return pii.ZipCode{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Italy",
			}
		},
		func(zipCode *pii.ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractStreetAddressesItaly extracts Italy street addresses as PiiEntity objects with context
func ExtractStreetAddressesItaly(text string) []pii.PiiEntity {
	addresses := extractWithContext(text, patterns.StreetAddressItalyRegex,
		func(value, context string) pii.StreetAddress {
			return pii.StreetAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Italy",
			}
		},
		func(address *pii.StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, address := range addresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// =============================================================================
// NEW COUNTRIES EXTRACTION FUNCTIONS
// =============================================================================

// --- Germany PII ---

// ExtractPostalCodesGermany extracts Germany postal codes as PiiEntity objects with context
func ExtractPostalCodesGermany(text string) []pii.PiiEntity {
	postalCodes := extractWithContext(text, patterns.PostalCodeGermanyRegex,
		func(value, context string) pii.ZipCode {
			return pii.ZipCode{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Germany",
			}
		},
		func(zipCode *pii.ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractPhonesGermany extracts Germany phone numbers as PiiEntity objects with context
func ExtractPhonesGermany(text string) []pii.PiiEntity {
	phones := extractWithContext(text, patterns.PhoneGermanyRegex,
		func(value, context string) pii.Phone {
			return pii.Phone{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Germany",
			}
		},
		func(phone *pii.Phone, context string) {
			phone.BasePii.IncrementCount()
			phone.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, phone := range phones {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypePhone,
			Value: phone,
		})
	}
	return entities
}

// ExtractStreetAddressesGermany extracts Germany street addresses as PiiEntity objects with context
func ExtractStreetAddressesGermany(text string) []pii.PiiEntity {
	addresses := extractWithContext(text, patterns.StreetAddressGermanyRegex,
		func(value, context string) pii.StreetAddress {
			return pii.StreetAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Germany",
			}
		},
		func(address *pii.StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, address := range addresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// --- China PII ---

// ExtractPostalCodesChina extracts China postal codes as PiiEntity objects with context
func ExtractPostalCodesChina(text string) []pii.PiiEntity {
	postalCodes := extractWithContext(text, patterns.PostalCodeChinaRegex,
		func(value, context string) pii.ZipCode {
			return pii.ZipCode{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "China",
			}
		},
		func(zipCode *pii.ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractPhonesChina extracts China phone numbers as PiiEntity objects with context
func ExtractPhonesChina(text string) []pii.PiiEntity {
	phones := extractWithContext(text, patterns.PhoneChinaRegex,
		func(value, context string) pii.Phone {
			return pii.Phone{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "China",
			}
		},
		func(phone *pii.Phone, context string) {
			phone.BasePii.IncrementCount()
			phone.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, phone := range phones {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypePhone,
			Value: phone,
		})
	}
	return entities
}

// ExtractStreetAddressesChina extracts China street addresses as PiiEntity objects with context
func ExtractStreetAddressesChina(text string) []pii.PiiEntity {
	addresses := extractWithContext(text, patterns.StreetAddressChinaRegex,
		func(value, context string) pii.StreetAddress {
			return pii.StreetAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "China",
			}
		},
		func(address *pii.StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, address := range addresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// --- India PII ---

// ExtractPostalCodesIndia extracts India postal codes as PiiEntity objects with context
func ExtractPostalCodesIndia(text string) []pii.PiiEntity {
	postalCodes := extractWithContext(text, patterns.PostalCodeIndiaRegex,
		func(value, context string) pii.ZipCode {
			return pii.ZipCode{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "India",
			}
		},
		func(zipCode *pii.ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractPhonesIndia extracts India phone numbers as PiiEntity objects with context
func ExtractPhonesIndia(text string) []pii.PiiEntity {
	phones := extractWithContext(text, patterns.PhoneIndiaRegex,
		func(value, context string) pii.Phone {
			return pii.Phone{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "India",
			}
		},
		func(phone *pii.Phone, context string) {
			phone.BasePii.IncrementCount()
			phone.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, phone := range phones {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypePhone,
			Value: phone,
		})
	}
	return entities
}

// ExtractStreetAddressesIndia extracts India street addresses as PiiEntity objects with context
func ExtractStreetAddressesIndia(text string) []pii.PiiEntity {
	addresses := extractWithContext(text, patterns.StreetAddressIndiaRegex,
		func(value, context string) pii.StreetAddress {
			return pii.StreetAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "India",
			}
		},
		func(address *pii.StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, address := range addresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// --- Arabic Countries PII ---

// ExtractPostalCodesArabic extracts Arabic countries postal codes as PiiEntity objects with context
func ExtractPostalCodesArabic(text string) []pii.PiiEntity {
	postalCodes := extractWithContext(text, patterns.PostalCodeArabicRegex,
		func(value, context string) pii.ZipCode {
			return pii.ZipCode{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Arabic",
			}
		},
		func(zipCode *pii.ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractPhonesArabic extracts Arabic countries phone numbers as PiiEntity objects with context
func ExtractPhonesArabic(text string) []pii.PiiEntity {
	phones := extractWithContext(text, patterns.PhoneArabicRegex,
		func(value, context string) pii.Phone {
			return pii.Phone{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Arabic",
			}
		},
		func(phone *pii.Phone, context string) {
			phone.BasePii.IncrementCount()
			phone.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, phone := range phones {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypePhone,
			Value: phone,
		})
	}
	return entities
}

// ExtractStreetAddressesArabic extracts Arabic countries street addresses as PiiEntity objects with context
func ExtractStreetAddressesArabic(text string) []pii.PiiEntity {
	addresses := extractWithContext(text, patterns.StreetAddressArabicRegex,
		func(value, context string) pii.StreetAddress {
			return pii.StreetAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Arabic",
			}
		},
		func(address *pii.StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, address := range addresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}

// --- Russia PII ---

// ExtractPostalCodesRussia extracts Russia postal codes as PiiEntity objects with context
func ExtractPostalCodesRussia(text string) []pii.PiiEntity {
	postalCodes := extractWithContext(text, patterns.PostalCodeRussiaRegex,
		func(value, context string) pii.ZipCode {
			return pii.ZipCode{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Russia",
			}
		},
		func(zipCode *pii.ZipCode, context string) {
			zipCode.BasePii.IncrementCount()
			zipCode.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, zipCode := range postalCodes {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeZipCode,
			Value: zipCode,
		})
	}
	return entities
}

// ExtractPhonesRussia extracts Russia phone numbers as PiiEntity objects with context
func ExtractPhonesRussia(text string) []pii.PiiEntity {
	phones := extractWithContext(text, patterns.PhoneRussiaRegex,
		func(value, context string) pii.Phone {
			return pii.Phone{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Russia",
			}
		},
		func(phone *pii.Phone, context string) {
			phone.BasePii.IncrementCount()
			phone.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, phone := range phones {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypePhone,
			Value: phone,
		})
	}
	return entities
}

// ExtractStreetAddressesRussia extracts Russia street addresses as PiiEntity objects with context
func ExtractStreetAddressesRussia(text string) []pii.PiiEntity {
	addresses := extractWithContext(text, patterns.StreetAddressRussiaRegex,
		func(value, context string) pii.StreetAddress {
			return pii.StreetAddress{
				BasePii: pii.BasePii{
					Value:    value,
					Contexts: []string{context},
					Count:    1,
				},
				Country: "Russia",
			}
		},
		func(address *pii.StreetAddress, context string) {
			address.BasePii.IncrementCount()
			address.BasePii.AddContext(context)
		})

	var entities []pii.PiiEntity
	for _, address := range addresses {
		entities = append(entities, pii.PiiEntity{
			Type:  pii.PiiTypeStreetAddress,
			Value: address,
		})
	}
	return entities
}