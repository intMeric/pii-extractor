# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Language Requirements

All README files, code comments, documentation, and commit messages must be written in English.

## Claude update

When architecture change (create/delete file etc...) update CLAUDE.md

## Project Overview

This is a Go-based PII (Personally Identifiable Information) extractor library that provides high-accuracy detection and extraction of sensitive data from text across multiple countries. The project features intelligent deduplication, context extraction, and optional LLM validation.

### Key Features (v0.0.1)

- **Multi-country Support**: Extracts PII for US, UK, France, Spain, and Italy
- **Smart Deduplication**: Automatically merges duplicate entities and consolidates contexts
- **High Accuracy**: Improved regex patterns to minimize false positives
- **Context Extraction**: Captures surrounding sentences or 8 words before/after for context
- **Comprehensive PII Types**: Emails, phone numbers, SSNs, postal codes, street addresses, P.O. boxes, credit cards, IP addresses, IBANs, Bitcoin addresses
- **LLM Validation**: Optional validation using OpenAI, Anthropic, Gemini, Mistral, or Ollama models
- **Type-safe API**: Clean interface with re-exports and structured value objects

## Development Commands

Since this is a Go project, use standard Go commands:

- **Build**: `go build`
- **Run**: `go run .`
- **Test**: `go test ./...`
- **Format**: `go fmt ./...`
- **Vet**: `go vet ./...`
- **Tidy dependencies**: `go mod tidy`

## Architecture

The project follows a modular architecture with clear separation of concerns:

### Core Components

- **PiiExtractor Interface**: Main abstraction for PII extraction (`extractors/interface.go`)
- **RegexExtractor**: High-performance regex-based implementation with deduplication
- **ValidatedExtractor**: LLM-enhanced validation wrapper
- **LLMExtractor**: Pure LLM-based extraction
- **EnsembleExtractor**: Combines multiple extractors
- **Value Objects**: Type-safe representations with smart merging capabilities
- **Registry System**: Global extractor registry for reusable configurations

### File Structure (Updated v0.0.1)

```
pii-extractor/
├── interface.go                     # Main API with re-exports
├── pii/
│   └── types.go                    # PII value objects with deduplication logic
├── extractors/
│   ├── interface.go                # Core extractor interfaces
│   ├── registry.go                 # Extractor registry system
│   ├── regex/
│   │   ├── extractor.go           # Main regex-based extractor
│   │   ├── extraction.go          # Extraction logic with context handling
│   │   └── patterns/              # Country-specific regex patterns
│   │       ├── common.go          # Global patterns and context extraction
│   │       ├── us.go              # US-specific patterns (improved)
│   │       ├── uk.go              # UK postal codes and addresses
│   │       ├── fr.go              # France postal codes and addresses
│   │       ├── es.go              # Spain postal codes and addresses
│   │       └── it.go              # Italy postal codes and addresses
│   ├── llm/                       # LLM-based extraction
│   └── hybrid/                    # Validation and ensemble extractors
├── examples/
│   ├── basic/                     # Simple usage examples
│   └── regex-with-llm-cross-val/  # Advanced validation examples
└── README.md                      # Comprehensive documentation
```

### Key Improvements in v0.0.1

- **Smart Deduplication**: `deduplicateEntities()` function merges duplicate PII with context consolidation
- **Improved Regex**: Enhanced US phone pattern to reduce false positives from credit cards
- **Better Context**: Fixed word-based context extraction for accurate surrounding text
- **Country Unification**: When merging entities with different countries, sets to empty string
- **Type-safe Merging**: Context merging respects type-specific fields (country, card type, etc.)

### International PII Support

The extractor supports country-specific formats with high accuracy:

- **US**: Phone numbers (555) 123-4567, SSNs 123-45-6789, ZIP codes 10001, street addresses, P.O. boxes
- **UK**: Postal codes SW1A 1AA, street addresses 221B Baker Street
- **France**: Metropolitan/DOM-TOM postal codes 75001/97110, street addresses 123 rue de la Paix
- **Spain**: Mainland/island postal codes 28013/35001, street addresses 123 Calle Mayor
- **Italy**: All postal codes 00186/20100, street addresses 123 Via del Corso

## Version History

### v0.0.2 - Enhanced Accuracy Release

- Released with tag `v0.0.2`
- Advanced false positive filtering for phone numbers
- Enhanced detection to prevent credit card/IBAN segments being classified as phones
- Improved US phone pattern for better test coverage
- Fixed all unit tests to pass consistently

### v0.0.1 - Initial Release

- Released with tag `v0.0.1`
- Smart deduplication system implemented
- False positive reduction in phone number detection
- Multi-country support for 5 countries
- Context extraction improvements
- Comprehensive README and documentation

## Installation

Users can install the stable version:

```bash
go get github.com/intMeric/pii-extractor
```

## Testing Strategy

The library should be tested with:

- Multi-country text samples
- Edge cases for false positives
- Deduplication scenarios
- Context extraction accuracy
- LLM validation workflows

## Important Notes for Development

- Always test changes with `examples/basic/basic_usage.go` to ensure no regressions
- Phone regex pattern is critical - changes may reintroduce false positives
- Deduplication logic in `pii/types.go` handles complex merging scenarios
- Context extraction prioritizes full sentences over word counts
