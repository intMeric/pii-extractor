# PII Extractor

A comprehensive Go-based library for extracting and identifying Personally Identifiable Information (PII) from text data. This tool uses regex patterns to detect various types of PII and optionally validates findings using Large Language Models (LLMs) for improved accuracy.

## Features

- **Multi-type PII Detection**: Supports detection of 10+ PII types including:
  - Phone numbers (US format)
  - Email addresses
  - Social Security Numbers (US format)
  - ZIP codes (US format)
  - Street addresses (US format)
  - P.O. Box addresses
  - Credit card numbers
  - IP addresses (IPv4/IPv6)
  - Bitcoin addresses
  - IBAN numbers

- **LLM Validation**: Optional validation using multiple LLM providers:
  - OpenAI (GPT models)
  - Anthropic (Claude models)
  - Google Gemini
  - Mistral AI
  - Ollama (local models)

- **Structured Results**: Returns detailed extraction results with:
  - Entity counts and statistics
  - Context information for each finding
  - Validation confidence scores
  - Type-safe value objects

- **Type Safety**: Full Go type safety with convenience methods and type assertions

## Installation

```bash
go get github.com/intMeric/pii-extractor
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    pii "github.com/intMeric/pii-extractor"
)

func main() {
    // Create a regex-based extractor
    extractor := pii.NewRegexExtractor()
    
    // Extract PII from text
    text := "Contact John at john@example.com or call (555) 123-4567"
    result, err := extractor.Extract(text)
    if err != nil {
        log.Fatal(err)
    }
    
    // Print results
    fmt.Printf("Found %d PII entities:\n", result.Total)
    for _, entity := range result.Entities {
        fmt.Printf("- %s: %s\n", entity.Type.String(), entity.GetValue())
    }
}
```

### With LLM Validation

```go
package main

import (
    "fmt"
    "log"
    
    pii "github.com/intMeric/pii-extractor"
)

func main() {
    // Configure validation
    config := pii.DefaultValidationConfig()
    config.Enabled = true
    config.Provider = pii.ProviderOpenAI
    config.APIKey = "your-api-key"
    
    // Create a validated extractor with configuration
    baseExtractor := pii.NewRegexExtractor()
    extractor, err := pii.NewValidatedExtractor(baseExtractor, config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Extract and validate PII using the configured validation
    text := "My email is john@example.com and my phone is (555) 123-4567"
    result, err := extractor.ExtractWithValidation(text)
    if err != nil {
        log.Fatal(err)
    }
    
    // Print validation results
    for _, entity := range result.Entities {
        if entity.IsValidated() {
            fmt.Printf("%s: %s (Valid: %t, Confidence: %.2f)\n", 
                entity.Type.String(), 
                entity.GetValue(), 
                entity.IsValid(), 
                entity.GetValidationConfidence())
        }
    }
}
```

### Advanced Usage

To change validation settings, create a new extractor with different configuration:

```go
// For different validation settings, create a new extractor
newConfig := &pii.ValidationConfig{
    Enabled: true,
    Provider: pii.ProviderMistral,
    APIKey: "different-api-key",
    MinConfidence: 0.9,
}

newExtractor, err := pii.NewValidatedExtractor(baseExtractor, newConfig)
result, err := newExtractor.ExtractWithValidation(text)
```

## API Reference

### Core Interfaces

- `PiiExtractor`: Basic extraction interface
- `ValidatedPiiExtractor`: Extended interface with LLM validation
- `PiiExtractionResult`: Structured results with statistics and utilities

### Supported PII Types

- `PiiTypePhone`: Phone numbers
- `PiiTypeEmail`: Email addresses  
- `PiiTypeSSN`: Social Security Numbers
- `PiiTypeZipCode`: ZIP/postal codes
- `PiiTypeStreetAddress`: Street addresses
- `PiiTypePoBox`: P.O. Box addresses
- `PiiTypeCreditCard`: Credit card numbers
- `PiiTypeIPAddress`: IP addresses
- `PiiTypeBtcAddress`: Bitcoin addresses
- `PiiTypeIBAN`: International Bank Account Numbers

### Value Objects

Each PII type has a corresponding value object with:
- `GetValue()`: Returns the raw string value
- `GetContexts()`: Returns surrounding text contexts
- `GetCount()`: Returns occurrence count
- Type-specific fields (e.g., country, card type)

## Examples

The `examples/` directory contains:
- `basic/`: Simple usage examples
- `regex-with-llm-cross-val/`: Advanced validation examples

## Development

### Commands

- **Build**: `go build`
- **Test**: `go test ./...`
- **Format**: `go fmt ./...`
- **Vet**: `go vet ./...`
- **Tidy**: `go mod tidy`

### Dependencies

- [gollm](https://github.com/teilomillet/gollm): LLM integration library
- Go 1.23.0 or later

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]