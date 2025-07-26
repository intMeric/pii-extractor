# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Language Requirements

All README files, code comments, documentation, and commit messages must be written in English.

## Claude update

When architecture change (create/delete file etc...) update CLAUDE.md

## Project Overview

This is a Go-based PII (Personally Identifiable Information) extractor tool that provides comprehensive detection and extraction of sensitive data from text. The project includes both basic regex-based extraction and advanced LLM-powered validation capabilities.

### Key Features

- **Multi-country Support**: Extracts PII for US, UK, France, Spain, and Italy
- **Comprehensive PII Types**: Postal codes, street addresses, emails, phone numbers, SSNs, credit cards, IP addresses, IBANs, Bitcoin addresses
- **LLM Validation**: Optional validation using OpenAI, Anthropic, Gemini, Mistral, or Ollama models
- **Context-aware Extraction**: Provides surrounding context for each detected PII entity
- **Type-safe Value Objects**: Structured data with country metadata and occurrence counts

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

- **PiiExtractor Interface**: Main abstraction for PII extraction
- **RegexExtractor**: High-performance regex-based implementation
- **ValidatedExtractor**: LLM-enhanced validation wrapper
- **Value Objects**: Type-safe representations (Email, Phone, Address, etc.)
- **Validation System**: Configurable LLM validation with multiple providers

### File Structure

- `interface.go`: Core interfaces and type definitions
- `models.go`: PII value object implementations
- `regex.go`: Regex patterns and extraction functions
- `regex_extractor.go`: Main regex-based extractor
- `validated_extractor.go`: LLM validation wrapper
- `validation.go`: LLM validation implementation
- `*_test.go`: Comprehensive test suites

### International PII Support

The extractor supports country-specific formats:

- **UK**: Postal codes (SW1A 1AA), addresses with alphanumeric house numbers (221B Baker Street)
- **France**: Metropolitan and DOM-TOM postal codes (75001, 97110), street addresses (123 rue de la Paix)
- **Spain**: Mainland and island postal codes (28013, 35001), addresses (123 Calle Mayor)
- **Italy**: All valid postal codes (00186, 20100), addresses (123 Via del Corso)
- **US**: ZIP codes, SSNs, phone numbers, street addresses
