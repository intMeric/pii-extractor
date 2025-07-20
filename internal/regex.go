package internal

import "regexp"

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

// US-specific convenience functions
func PhonesUS(text string) []string {
	return match(text, PhoneUSRegex)
}

func PhonesWithExtsUS(text string) []string {
	return match(text, PhonesWithExtsUSRegex)
}

func StreetAddressesUS(text string) []string {
	return match(text, StreetAddressUSRegex)
}

func ZipCodesUS(text string) []string {
	return match(text, ZipCodeUSRegex)
}

func PoBoxesUS(text string) []string {
	return match(text, PoBoxUSRegex)
}

func SSNsUS(text string) []string {
	return match(text, SSNUSRegex)
}

// International/generic convenience functions
func Emails(text string) []string {
	return match(text, EmailRegex)
}

func IPv4s(text string) []string {
	return match(text, IPv4Regex)
}

func IPv6s(text string) []string {
	return match(text, IPv6Regex)
}

func IPs(text string) []string {
	return match(text, IPRegex)
}

func CreditCards(text string) []string {
	return match(text, CreditCardRegex)
}

func VISACreditCards(text string) []string {
	return match(text, VISACreditCardRegex)
}

func MCCreditCards(text string) []string {
	return match(text, MCCreditCardRegex)
}

func BtcAddresses(text string) []string {
	return match(text, BtcAddressRegex)
}

func IBANs(text string) []string {
	return match(text, IBANRegex)
}

// Structured extraction functions that return value objects

// ExtractPhonesUS extracts US phone numbers as Phone value objects
func ExtractPhonesUS(text string) []Phone {
	matches := match(text, PhoneUSRegex)
	phones := make([]Phone, len(matches))
	for i, match := range matches {
		phones[i] = NewPhoneUS(match)
	}
	return phones
}

// ExtractEmails extracts email addresses as Email value objects
func ExtractEmails(text string) []Email {
	matches := match(text, EmailRegex)
	emails := make([]Email, len(matches))
	for i, match := range matches {
		emails[i] = NewEmail(match)
	}
	return emails
}

// ExtractSSNsUS extracts US SSNs as SSN value objects
func ExtractSSNsUS(text string) []SSN {
	matches := match(text, SSNUSRegex)
	ssns := make([]SSN, len(matches))
	for i, match := range matches {
		ssns[i] = NewSSNUS(match)
	}
	return ssns
}

// ExtractZipCodesUS extracts US zip codes as ZipCode value objects
func ExtractZipCodesUS(text string) []ZipCode {
	matches := match(text, ZipCodeUSRegex)
	zipCodes := make([]ZipCode, len(matches))
	for i, match := range matches {
		zipCodes[i] = NewZipCodeUS(match)
	}
	return zipCodes
}

// ExtractStreetAddressesUS extracts US street addresses as StreetAddress value objects
func ExtractStreetAddressesUS(text string) []StreetAddress {
	matches := match(text, StreetAddressUSRegex)
	addresses := make([]StreetAddress, len(matches))
	for i, match := range matches {
		addresses[i] = NewStreetAddressUS(match)
	}
	return addresses
}

// ExtractCreditCards extracts credit cards as CreditCard value objects
func ExtractCreditCards(text string) []CreditCard {
	var cards []CreditCard
	
	// Check for VISA cards
	visaMatches := match(text, VISACreditCardRegex)
	for _, match := range visaMatches {
		cards = append(cards, NewCreditCard(match, "visa"))
	}
	
	// Check for MasterCard
	mcMatches := match(text, MCCreditCardRegex)
	for _, match := range mcMatches {
		cards = append(cards, NewCreditCard(match, "mastercard"))
	}
	
	// Check for generic cards (not VISA/MC)
	genericMatches := match(text, CreditCardRegex)
	for _, match := range genericMatches {
		// Skip if already found as VISA or MC
		isVisa := false
		isMC := false
		for _, visa := range visaMatches {
			if visa == match {
				isVisa = true
				break
			}
		}
		for _, mc := range mcMatches {
			if mc == match {
				isMC = true
				break
			}
		}
		if !isVisa && !isMC {
			cards = append(cards, NewCreditCard(match, "generic"))
		}
	}
	
	return cards
}

// ExtractIPAddresses extracts IP addresses as IPAddress value objects
func ExtractIPAddresses(text string) []IPAddress {
	var ips []IPAddress
	
	// Extract IPv4
	ipv4Matches := match(text, IPv4Regex)
	for _, match := range ipv4Matches {
		ips = append(ips, NewIPv4(match))
	}
	
	// Extract IPv6
	ipv6Matches := match(text, IPv6Regex)
	for _, match := range ipv6Matches {
		ips = append(ips, NewIPv6(match))
	}
	
	return ips
}

// ExtractBtcAddresses extracts Bitcoin addresses as BtcAddress value objects
func ExtractBtcAddresses(text string) []BtcAddress {
	matches := match(text, BtcAddressRegex)
	addresses := make([]BtcAddress, len(matches))
	for i, match := range matches {
		addresses[i] = NewBtcAddress(match)
	}
	return addresses
}

// ExtractIBANs extracts IBANs as IBAN value objects
func ExtractIBANs(text string) []IBAN {
	matches := match(text, IBANRegex)
	ibans := make([]IBAN, len(matches))
	for i, match := range matches {
		ibans[i] = NewIBAN(match)
	}
	return ibans
}
