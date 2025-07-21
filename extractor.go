package piiextractor

// RegexExtractor is a simple implementation of PiiExtractor using regex patterns
type RegexExtractor struct{}

// NewRegexExtractor creates a new regex-based PII extractor
func NewRegexExtractor() *RegexExtractor {
	return &RegexExtractor{}
}

// addEntities is a helper function to convert PII values to entities and append them
func addEntities[T Pii](entities *[]PiiEntity, piiType string, values []T) {
	for _, value := range values {
		*entities = append(*entities, PiiEntity{
			Type:  piiType,
			Value: value,
		})
	}
}

// Extract implements the PiiExtractor interface
func (r *RegexExtractor) Extract(text string) ([]PiiEntity, error) {
	var entities []PiiEntity

	addEntities(&entities, "phone", ExtractPhonesUS(text))
	addEntities(&entities, "email", ExtractEmails(text))
	addEntities(&entities, "ssn", ExtractSSNsUS(text))
	addEntities(&entities, "zip_code", ExtractZipCodesUS(text))
	addEntities(&entities, "street_address", ExtractStreetAddressesUS(text))
	addEntities(&entities, "po_box", ExtractPoBoxesUS(text))
	addEntities(&entities, "credit_card", ExtractCreditCards(text))
	addEntities(&entities, "ip_address", ExtractIPAddresses(text))
	addEntities(&entities, "btc_address", ExtractBtcAddresses(text))
	addEntities(&entities, "iban", ExtractIBANs(text))

	return entities, nil
}