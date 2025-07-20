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
