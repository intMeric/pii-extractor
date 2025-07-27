package regex

import (
	"slices"
	
	"github.com/intMeric/pii-extractor/extractors"
	"github.com/intMeric/pii-extractor/pii"
)

// RegexExtractor implements PII extraction using regular expressions
type RegexExtractor struct {
	name      string
	countries []string
	types     []pii.PiiType
}

// NewExtractor creates a new regex-based PII extractor
func NewExtractor(config *extractors.ExtractorConfig) *RegexExtractor {
	extractor := &RegexExtractor{
		name: "regex-extractor",
	}

	if config != nil {
		if config.Countries != nil {
			extractor.countries = config.Countries
		}
		if config.Types != nil {
			extractor.types = config.Types
		}
	}

	return extractor
}

// NewDefaultExtractor creates a regex extractor with default settings
func NewDefaultExtractor() *RegexExtractor {
	return NewExtractor(nil)
}

// Extract performs PII extraction on the given text
func (r *RegexExtractor) Extract(text string) (*pii.PiiExtractionResult, error) {
	var allEntities []pii.PiiEntity

	// If specific types are configured, extract only those
	if len(r.types) > 0 {
		for _, piiType := range r.types {
			entities, err := r.ExtractByType(text, piiType)
			if err != nil {
				return nil, err
			}
			allEntities = append(allEntities, entities...)
		}
	} else {
		// Extract all types
		allEntities = append(allEntities, ExtractEmails(text)...)
		allEntities = append(allEntities, ExtractCreditCards(text)...)
		allEntities = append(allEntities, ExtractIPAddresses(text)...)
		allEntities = append(allEntities, ExtractBtcAddresses(text)...)
		allEntities = append(allEntities, ExtractIBANs(text)...)

		// US-specific extractions
		if r.shouldExtractForCountry("US") {
			allEntities = append(allEntities, ExtractPhonesUS(text)...)
			allEntities = append(allEntities, ExtractSSNsUS(text)...)
			allEntities = append(allEntities, ExtractZipCodesUS(text)...)
			allEntities = append(allEntities, ExtractStreetAddressesUS(text)...)
			allEntities = append(allEntities, ExtractPoBoxesUS(text)...)
		}

		// UK-specific extractions
		if r.shouldExtractForCountry("UK") {
			allEntities = append(allEntities, ExtractPostalCodesUK(text)...)
			allEntities = append(allEntities, ExtractStreetAddressesUK(text)...)
		}

		// France-specific extractions
		if r.shouldExtractForCountry("France") {
			allEntities = append(allEntities, ExtractPostalCodesFrance(text)...)
			allEntities = append(allEntities, ExtractStreetAddressesFrance(text)...)
		}

		// Spain-specific extractions
		if r.shouldExtractForCountry("Spain") {
			allEntities = append(allEntities, ExtractPostalCodesSpain(text)...)
			allEntities = append(allEntities, ExtractStreetAddressesSpain(text)...)
		}

		// Italy-specific extractions
		if r.shouldExtractForCountry("Italy") {
			allEntities = append(allEntities, ExtractPostalCodesItaly(text)...)
			allEntities = append(allEntities, ExtractStreetAddressesItaly(text)...)
		}
	}

	return pii.NewPiiExtractionResult(allEntities), nil
}

// ExtractByType extracts only specific types of PII from the text
func (r *RegexExtractor) ExtractByType(text string, piiType pii.PiiType) ([]pii.PiiEntity, error) {
	switch piiType {
	case pii.PiiTypeEmail:
		return ExtractEmails(text), nil
	case pii.PiiTypeCreditCard:
		return ExtractCreditCards(text), nil
	case pii.PiiTypeIPAddress:
		return ExtractIPAddresses(text), nil
	case pii.PiiTypeBtcAddress:
		return ExtractBtcAddresses(text), nil
	case pii.PiiTypeIBAN:
		return ExtractIBANs(text), nil
	case pii.PiiTypePhone:
		if r.shouldExtractForCountry("US") {
			return ExtractPhonesUS(text), nil
		}
	case pii.PiiTypeSSN:
		if r.shouldExtractForCountry("US") {
			return ExtractSSNsUS(text), nil
		}
	case pii.PiiTypeZipCode:
		var entities []pii.PiiEntity
		if r.shouldExtractForCountry("US") {
			entities = append(entities, ExtractZipCodesUS(text)...)
		}
		if r.shouldExtractForCountry("UK") {
			entities = append(entities, ExtractPostalCodesUK(text)...)
		}
		if r.shouldExtractForCountry("France") {
			entities = append(entities, ExtractPostalCodesFrance(text)...)
		}
		if r.shouldExtractForCountry("Spain") {
			entities = append(entities, ExtractPostalCodesSpain(text)...)
		}
		if r.shouldExtractForCountry("Italy") {
			entities = append(entities, ExtractPostalCodesItaly(text)...)
		}
		return entities, nil
	case pii.PiiTypeStreetAddress:
		var entities []pii.PiiEntity
		if r.shouldExtractForCountry("US") {
			entities = append(entities, ExtractStreetAddressesUS(text)...)
		}
		if r.shouldExtractForCountry("UK") {
			entities = append(entities, ExtractStreetAddressesUK(text)...)
		}
		if r.shouldExtractForCountry("France") {
			entities = append(entities, ExtractStreetAddressesFrance(text)...)
		}
		if r.shouldExtractForCountry("Spain") {
			entities = append(entities, ExtractStreetAddressesSpain(text)...)
		}
		if r.shouldExtractForCountry("Italy") {
			entities = append(entities, ExtractStreetAddressesItaly(text)...)
		}
		return entities, nil
	case pii.PiiTypePoBox:
		if r.shouldExtractForCountry("US") {
			return ExtractPoBoxesUS(text), nil
		}
	}

	return []pii.PiiEntity{}, nil
}

// shouldExtractForCountry checks if extraction should be performed for a specific country
func (r *RegexExtractor) shouldExtractForCountry(country string) bool {
	// If no countries specified, extract for all
	if len(r.countries) == 0 {
		return true
	}

	// Check if the country is in the allowed list
	return slices.Contains(r.countries, country)
}

// GetSupportedTypes returns the list of PII types this extractor can handle
func (r *RegexExtractor) GetSupportedTypes() []pii.PiiType {
	return []pii.PiiType{
		pii.PiiTypePhone,
		pii.PiiTypeEmail,
		pii.PiiTypeSSN,
		pii.PiiTypeZipCode,
		pii.PiiTypePoBox,
		pii.PiiTypeStreetAddress,
		pii.PiiTypeCreditCard,
		pii.PiiTypeIPAddress,
		pii.PiiTypeBtcAddress,
		pii.PiiTypeIBAN,
	}
}

// GetMethod returns the extraction method used by this extractor
func (r *RegexExtractor) GetMethod() extractors.ExtractionMethod {
	return extractors.MethodRegex
}

// GetName returns a human-readable name for this extractor
func (r *RegexExtractor) GetName() string {
	return r.name
}

// GetCountries returns the list of countries this extractor is configured for
func (r *RegexExtractor) GetCountries() []string {
	return r.countries
}

// GetTypes returns the list of PII types this extractor is configured for
func (r *RegexExtractor) GetTypes() []pii.PiiType {
	return r.types
}
