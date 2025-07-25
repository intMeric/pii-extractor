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

	entities = append(entities, ExtractPhonesUS(text)...)
	entities = append(entities, ExtractEmails(text)...)
	entities = append(entities, ExtractSSNsUS(text)...)
	entities = append(entities, ExtractZipCodesUS(text)...)
	entities = append(entities, ExtractStreetAddressesUS(text)...)
	entities = append(entities, ExtractPoBoxesUS(text)...)
	entities = append(entities, ExtractCreditCards(text)...)
	entities = append(entities, ExtractIPAddresses(text)...)
	entities = append(entities, ExtractBtcAddresses(text)...)
	entities = append(entities, ExtractIBANs(text)...)

	return NewPiiExtractionResult(entities), nil
}