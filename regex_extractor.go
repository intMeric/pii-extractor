package piiextractor

// RegexExtractor is a simple implementation of PiiExtractor using regex patterns
type RegexExtractor struct{}

// NewRegexExtractor creates a new regex-based PII extractor
func NewRegexExtractor() *RegexExtractor {
	return &RegexExtractor{}
}

// Extract implements the PiiExtractor interface
func (r *RegexExtractor) Extract(text string) (*PiiExtractionResult, error) {
	var entities []PiiEntity

	// US-specific extractions
	entities = append(entities, ExtractPhonesUS(text)...)
	entities = append(entities, ExtractSSNsUS(text)...)
	entities = append(entities, ExtractZipCodesUS(text)...)
	entities = append(entities, ExtractStreetAddressesUS(text)...)
	entities = append(entities, ExtractPoBoxesUS(text)...)

	// International postal codes
	entities = append(entities, ExtractPostalCodesUK(text)...)
	entities = append(entities, ExtractPostalCodesFrance(text)...)
	entities = append(entities, ExtractPostalCodesSpain(text)...)
	entities = append(entities, ExtractPostalCodesItaly(text)...)

	// International street addresses
	entities = append(entities, ExtractStreetAddressesUK(text)...)
	entities = append(entities, ExtractStreetAddressesFrance(text)...)
	entities = append(entities, ExtractStreetAddressesSpain(text)...)
	entities = append(entities, ExtractStreetAddressesItaly(text)...)

	// Generic/international extractions
	entities = append(entities, ExtractEmails(text)...)
	entities = append(entities, ExtractCreditCards(text)...)
	entities = append(entities, ExtractIPAddresses(text)...)
	entities = append(entities, ExtractBtcAddresses(text)...)
	entities = append(entities, ExtractIBANs(text)...)

	return NewPiiExtractionResult(entities), nil
}
