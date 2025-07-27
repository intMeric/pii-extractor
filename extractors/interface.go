package extractors

import (
	"github.com/intMeric/pii-extractor/pii"
)

// ExtractionMethod represents the method used for PII extraction
type ExtractionMethod string

const (
	MethodRegex  ExtractionMethod = "regex"
	MethodLLM    ExtractionMethod = "llm"
	MethodML     ExtractionMethod = "ml"
	MethodHybrid ExtractionMethod = "hybrid"
)

// String returns the string representation of the extraction method
func (m ExtractionMethod) String() string {
	return string(m)
}

// PiiExtractor defines the interface that all PII extractors must implement
type PiiExtractor interface {
	// Extract performs PII extraction on the given text and returns all found entities
	Extract(text string) (*pii.PiiExtractionResult, error)
	
	// ExtractByType extracts only specific types of PII from the text
	ExtractByType(text string, piiType pii.PiiType) ([]pii.PiiEntity, error)
	
	// GetSupportedTypes returns the list of PII types this extractor can handle
	GetSupportedTypes() []pii.PiiType
	
	// GetMethod returns the extraction method used by this extractor
	GetMethod() ExtractionMethod
	
	// GetName returns a human-readable name for this extractor
	GetName() string
}

// ExtractorConfig represents configuration options for extractors
type ExtractorConfig struct {
	// Method specifies the extraction method to use
	Method ExtractionMethod `json:"method"`
	
	// Options contains method-specific configuration
	Options map[string]any `json:"options,omitempty"`
	
	// Countries specifies which countries to extract PII for (empty = all)
	Countries []string `json:"countries,omitempty"`
	
	// Types specifies which PII types to extract (empty = all)
	Types []pii.PiiType `json:"types,omitempty"`
}