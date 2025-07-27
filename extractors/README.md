# PII Extractors

This package provides a flexible, extensible architecture for PII (Personally Identifiable Information) extraction using multiple methods.

## Architecture

The extractors package is organized around a common interface that allows different extraction methods to be used interchangeably:

```
extractors/
├── interface.go          # Common PiiExtractor interface
├── registry.go          # Extractor registration and discovery
├── regex/               # Regex-based extraction
│   ├── extractor.go     # RegexExtractor implementation
│   └── patterns/        # Regex patterns by country
│       ├── common.go    # International patterns
│       ├── us.go        # US-specific patterns
│       └── ...          # Other countries
├── llm/                 # LLM-based extraction
│   ├── extractor.go     # LLMExtractor implementation
│   ├── providers/       # LLM provider implementations
│   └── prompts/         # Extraction prompts
└── hybrid/              # Combination methods
    ├── ensemble.go      # EnsembleExtractor for combining methods
    └── validator.go     # Cross-validation utilities
```

## Usage

### Basic Usage

```go
import "github.com/intMeric/pii-extractor/extractors"

// Create a regex extractor
regexExtractor := regex.NewDefaultExtractor()

// Extract PII from text
result, err := regexExtractor.Extract("Contact me at john@example.com")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d entities\n", result.Total)
```

### Using the Registry

```go
import (
    "github.com/intMeric/pii-extractor/extractors"
    "github.com/intMeric/pii-extractor/extractors/regex"
)

// Register extractors
extractors.Register("regex", regex.NewDefaultExtractor())

// Get and use an extractor
extractor, err := extractors.Get("regex")
if err != nil {
    log.Fatal(err)
}

result, err := extractor.Extract(text)
```

### LLM-based Extraction (Future)

```go
import "github.com/intMeric/pii-extractor/extractors/llm"

// Create LLM extractor
config := &extractors.ExtractorConfig{
    Method: extractors.MethodLLM,
    Options: map[string]interface{}{
        "api_key": "your-openai-key",
        "temperature": 0.1,
    },
}

llmExtractor := llm.NewExtractor(llm.ProviderOpenAI, "gpt-4", config)
result, err := llmExtractor.Extract(text)
```

### Ensemble/Hybrid Extraction

```go
import "github.com/intMeric/pii-extractor/extractors/hybrid"

// Combine multiple extractors
regexExtractor := regex.NewDefaultExtractor()
llmExtractor := llm.NewExtractor(llm.ProviderOpenAI, "gpt-4", config)

ensemble := hybrid.NewEnsembleExtractor(regexExtractor, llmExtractor).
    WithStrategy(hybrid.StrategyUnion).
    WithValidation(hybrid.ValidationBasic)

result, err := ensemble.Extract(text)
```

## Extending with New Methods

To add a new extraction method:

1. Create a new package under `extractors/`
2. Implement the `PiiExtractor` interface
3. Register your extractor in the registry

Example:

```go
package mymethod

import "github.com/intMeric/pii-extractor/extractors"

type MyExtractor struct {
    name string
}

func (m *MyExtractor) Extract(text string) (*PiiExtractionResult, error) {
    // Your implementation here
}

func (m *MyExtractor) GetMethod() extractors.ExtractionMethod {
    return "mymethod"
}

// Implement other interface methods...
```

## Configuration

Each extractor can be configured using the `ExtractorConfig` struct:

```go
config := &extractors.ExtractorConfig{
    Method: extractors.MethodRegex,
    Countries: []string{"US", "FR", "UK"}, // Only extract for these countries
    Types: []PiiType{PiiTypeEmail, PiiTypePhone}, // Only extract these types
    Options: map[string]interface{}{
        "api_key": "...",
        "temperature": 0.1,
    },
}
```

## Future Enhancements

- Machine Learning-based extractors
- Cloud API integrations (AWS Comprehend, Google DLP, etc.)
- Real-time streaming extraction
- Custom validation rules
- Performance optimization with caching
- Async/batch processing support