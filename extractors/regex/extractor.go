package regex

import (
	"runtime"
	"slices"
	"sync"
	
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
	// Pre-allocate slice with estimated capacity based on text length
	// Rough estimation: 1 PII entity per 200 characters
	estimatedCapacity := len(text)/200 + 10
	if estimatedCapacity > 1000 {
		estimatedCapacity = 1000 // Cap at reasonable maximum
	}
	allEntities := make([]pii.PiiEntity, 0, estimatedCapacity)

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
		// Collect all extraction operations and batch them
		var extractorFuncs []func(string) []pii.PiiEntity
		
		// Generic/International extractors
		extractorFuncs = append(extractorFuncs,
			ExtractEmails,
			ExtractCreditCards,
			ExtractIPAddresses,
			ExtractBtcAddresses,
			ExtractIBANs,
		)

		// Country-specific extractors
		if r.shouldExtractForCountry("US") {
			extractorFuncs = append(extractorFuncs,
				ExtractPhonesUS,
				ExtractSSNsUS,
				ExtractZipCodesUS,
				ExtractStreetAddressesUS,
				ExtractPoBoxesUS,
			)
		}

		if r.shouldExtractForCountry("UK") {
			extractorFuncs = append(extractorFuncs,
				ExtractPostalCodesUK,
				ExtractStreetAddressesUK,
			)
		}

		if r.shouldExtractForCountry("France") {
			extractorFuncs = append(extractorFuncs,
				ExtractPostalCodesFrance,
				ExtractStreetAddressesFrance,
			)
		}

		if r.shouldExtractForCountry("Spain") {
			extractorFuncs = append(extractorFuncs,
				ExtractPostalCodesSpain,
				ExtractStreetAddressesSpain,
			)
		}

		if r.shouldExtractForCountry("Italy") {
			extractorFuncs = append(extractorFuncs,
				ExtractPostalCodesItaly,
				ExtractStreetAddressesItaly,
			)
		}

		if r.shouldExtractForCountry("Germany") {
			extractorFuncs = append(extractorFuncs,
				ExtractPostalCodesGermany,
				ExtractPhonesGermany,
				ExtractStreetAddressesGermany,
			)
		}

		if r.shouldExtractForCountry("China") {
			extractorFuncs = append(extractorFuncs,
				ExtractPostalCodesChina,
				ExtractPhonesChina,
				ExtractStreetAddressesChina,
			)
		}

		if r.shouldExtractForCountry("India") {
			extractorFuncs = append(extractorFuncs,
				ExtractPostalCodesIndia,
				ExtractPhonesIndia,
				ExtractStreetAddressesIndia,
			)
		}

		if r.shouldExtractForCountry("Arabic") {
			extractorFuncs = append(extractorFuncs,
				ExtractPostalCodesArabic,
				ExtractPhonesArabic,
				ExtractStreetAddressesArabic,
			)
		}

		if r.shouldExtractForCountry("Russia") {
			extractorFuncs = append(extractorFuncs,
				ExtractPostalCodesRussia,
				ExtractPhonesRussia,
				ExtractStreetAddressesRussia,
			)
		}

		// Use parallel execution for large text or many extractors
		if len(text) > 10000 && len(extractorFuncs) > 8 {
			allEntities = r.executeExtractorsParallel(text, extractorFuncs, allEntities)
		} else {
			// Sequential execution for smaller workloads
			for _, extractorFunc := range extractorFuncs {
				entities := extractorFunc(text)
				if len(entities) > 0 {
					allEntities = append(allEntities, entities...)
				}
			}
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
		entities := make([]pii.PiiEntity, 0, 20) // Pre-allocate for typical phone count
		if r.shouldExtractForCountry("US") {
			entities = append(entities, ExtractPhonesUS(text)...)
		}
		if r.shouldExtractForCountry("Germany") {
			entities = append(entities, ExtractPhonesGermany(text)...)
		}
		if r.shouldExtractForCountry("China") {
			entities = append(entities, ExtractPhonesChina(text)...)
		}
		if r.shouldExtractForCountry("India") {
			entities = append(entities, ExtractPhonesIndia(text)...)
		}
		if r.shouldExtractForCountry("Arabic") {
			entities = append(entities, ExtractPhonesArabic(text)...)
		}
		if r.shouldExtractForCountry("Russia") {
			entities = append(entities, ExtractPhonesRussia(text)...)
		}
		return entities, nil
	case pii.PiiTypeSSN:
		if r.shouldExtractForCountry("US") {
			return ExtractSSNsUS(text), nil
		}
	case pii.PiiTypeZipCode:
		entities := make([]pii.PiiEntity, 0, 30) // Pre-allocate for typical postal code count
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
		if r.shouldExtractForCountry("Germany") {
			entities = append(entities, ExtractPostalCodesGermany(text)...)
		}
		if r.shouldExtractForCountry("China") {
			entities = append(entities, ExtractPostalCodesChina(text)...)
		}
		if r.shouldExtractForCountry("India") {
			entities = append(entities, ExtractPostalCodesIndia(text)...)
		}
		if r.shouldExtractForCountry("Arabic") {
			entities = append(entities, ExtractPostalCodesArabic(text)...)
		}
		if r.shouldExtractForCountry("Russia") {
			entities = append(entities, ExtractPostalCodesRussia(text)...)
		}
		return entities, nil
	case pii.PiiTypeStreetAddress:
		entities := make([]pii.PiiEntity, 0, 25) // Pre-allocate for typical address count
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
		if r.shouldExtractForCountry("Germany") {
			entities = append(entities, ExtractStreetAddressesGermany(text)...)
		}
		if r.shouldExtractForCountry("China") {
			entities = append(entities, ExtractStreetAddressesChina(text)...)
		}
		if r.shouldExtractForCountry("India") {
			entities = append(entities, ExtractStreetAddressesIndia(text)...)
		}
		if r.shouldExtractForCountry("Arabic") {
			entities = append(entities, ExtractStreetAddressesArabic(text)...)
		}
		if r.shouldExtractForCountry("Russia") {
			entities = append(entities, ExtractStreetAddressesRussia(text)...)
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

// executeExtractorsParallel runs extraction functions in parallel using worker pool
func (r *RegexExtractor) executeExtractorsParallel(text string, extractorFuncs []func(string) []pii.PiiEntity, initialEntities []pii.PiiEntity) []pii.PiiEntity {
	numWorkers := runtime.NumCPU()
	if numWorkers > len(extractorFuncs) {
		numWorkers = len(extractorFuncs)
	}
	
	// Create channels for work distribution
	jobs := make(chan func(string) []pii.PiiEntity, len(extractorFuncs))
	results := make(chan []pii.PiiEntity, len(extractorFuncs))
	
	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for extractorFunc := range jobs {
				entities := extractorFunc(text)
				results <- entities
			}
		}()
	}
	
	// Send jobs to workers
	go func() {
		for _, extractorFunc := range extractorFuncs {
			jobs <- extractorFunc
		}
		close(jobs)
	}()
	
	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Collect all results
	allEntities := initialEntities
	for entities := range results {
		if len(entities) > 0 {
			allEntities = append(allEntities, entities...)
		}
	}
	
	return allEntities
}
