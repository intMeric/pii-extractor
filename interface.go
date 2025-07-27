package piiextractor

import (
	"github.com/intMeric/pii-extractor/extractors"
	hybridExtractor "github.com/intMeric/pii-extractor/extractors/hybrid"
	llmExtractor "github.com/intMeric/pii-extractor/extractors/llm"
	regexExtractor "github.com/intMeric/pii-extractor/extractors/regex"
	"github.com/intMeric/pii-extractor/pii"
)

// Re-export types from pii package for convenience
type PiiType = pii.PiiType
type PiiEntity = pii.PiiEntity
type PiiExtractionResult = pii.PiiExtractionResult
type ValidationStats = pii.ValidationStats
type ValidationResult = pii.ValidationResult

// Re-export PII value types
type Pii = pii.Pii
type BasePii = pii.BasePii
type Phone = pii.Phone
type Email = pii.Email
type SSN = pii.SSN
type ZipCode = pii.ZipCode
type StreetAddress = pii.StreetAddress
type PoBox = pii.PoBox
type CreditCard = pii.CreditCard
type IPAddress = pii.IPAddress
type BtcAddress = pii.BtcAddress
type IBAN = pii.IBAN

// Re-export constants
const (
	PiiTypePhone         = pii.PiiTypePhone
	PiiTypeEmail         = pii.PiiTypeEmail
	PiiTypeSSN           = pii.PiiTypeSSN
	PiiTypeZipCode       = pii.PiiTypeZipCode
	PiiTypePoBox         = pii.PiiTypePoBox
	PiiTypeStreetAddress = pii.PiiTypeStreetAddress
	PiiTypeCreditCard    = pii.PiiTypeCreditCard
	PiiTypeIPAddress     = pii.PiiTypeIPAddress
	PiiTypeBtcAddress    = pii.PiiTypeBtcAddress
	PiiTypeIBAN          = pii.PiiTypeIBAN
)

// Re-export extractors types for convenience
type ExtractionMethod = extractors.ExtractionMethod
type ExtractorConfig = extractors.ExtractorConfig
type PiiExtractor = extractors.PiiExtractor

// Re-export hybrid types for convenience
type ValidationConfig = hybridExtractor.ValidationConfig
type LLMProvider = hybridExtractor.LLMProvider
type ValidatedExtractor = hybridExtractor.ValidatedExtractor
type EnsembleExtractor = hybridExtractor.EnsembleExtractor

// Re-export extraction methods
const (
	MethodRegex  = extractors.MethodRegex
	MethodLLM    = extractors.MethodLLM
	MethodML     = extractors.MethodML
	MethodHybrid = extractors.MethodHybrid
)

// Re-export LLM providers
const (
	ProviderOpenAI    = hybridExtractor.ProviderOpenAI
	ProviderMistral   = hybridExtractor.ProviderMistral
	ProviderGemini    = hybridExtractor.ProviderGemini
	ProviderOllama    = hybridExtractor.ProviderOllama
	ProviderAnthropic = hybridExtractor.ProviderAnthropic
)

// Modern constructor functions

// NewRegexExtractor creates a new regex-based PII extractor
func NewRegexExtractor(config *ExtractorConfig) PiiExtractor {
	return regexExtractor.NewExtractor(config)
}

// NewDefaultRegexExtractor creates a regex extractor with default settings
func NewDefaultRegexExtractor() PiiExtractor {
	return regexExtractor.NewDefaultExtractor()
}

// NewLLMExtractor creates a new LLM-based PII extractor
func NewLLMExtractor(provider llmExtractor.Provider, model string, config *ExtractorConfig) (PiiExtractor, error) {
	return llmExtractor.NewExtractor(provider, model, config)
}

// NewEnsembleExtractor creates a new ensemble extractor that combines multiple extractors
func NewEnsembleExtractor(extractors ...PiiExtractor) *hybridExtractor.EnsembleExtractor {
	return hybridExtractor.NewEnsembleExtractor(extractors...)
}

// NewValidatedExtractor creates a new validated extractor that combines any base extractor with LLM validation
func NewValidatedExtractor(baseExtractor PiiExtractor, config *hybridExtractor.ValidationConfig) (*hybridExtractor.ValidatedExtractor, error) {
	return hybridExtractor.NewValidatedExtractor(baseExtractor, config)
}

// DefaultValidationConfig returns a default configuration for LLM validation
func DefaultValidationConfig() *ValidationConfig {
	return hybridExtractor.DefaultValidationConfig()
}

// Registry functions

// Register adds an extractor to the global registry
func Register(name string, extractor PiiExtractor) error {
	return extractors.Register(name, extractor)
}

// Get retrieves an extractor from the global registry
func Get(name string) (PiiExtractor, error) {
	return extractors.Get(name)
}

// GetByMethod returns all extractors that use the specified method
func GetByMethod(method ExtractionMethod) []PiiExtractor {
	return extractors.GetByMethod(method)
}

// List returns all registered extractor names
func List() []string {
	return extractors.List()
}

// Utility functions

// NewPiiExtractionResult creates a new extraction result
var NewPiiExtractionResult = pii.NewPiiExtractionResult

// PII constructors
var NewEmail = pii.NewEmail
var NewPhoneUS = pii.NewPhoneUS
var NewPhone = pii.NewPhone
var NewSSN = pii.NewSSN
var NewZipCode = pii.NewZipCode
var NewStreetAddress = pii.NewStreetAddress
var NewPoBox = pii.NewPoBox
var NewCreditCard = pii.NewCreditCard
var NewIPAddress = pii.NewIPAddress
var NewBtcAddress = pii.NewBtcAddress
var NewIBAN = pii.NewIBAN

// GetTypedValue performs a safe type assertion for PII values
func GetTypedValue[T Pii](entity PiiEntity) (T, bool) {
	return pii.GetTypedValue[T](entity)
}
