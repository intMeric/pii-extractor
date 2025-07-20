# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Language Requirements

All README files, code comments, documentation, and commit messages must be written in English.

## Project Overview

This is a Go-based PII (Personally Identifiable Information) extractor tool. The project is in early development stages with minimal code structure currently in place.

## Development Commands

Since this is a Go project, use standard Go commands:

- **Build**: `go build`
- **Run**: `go run .`
- **Test**: `go test ./...`
- **Format**: `go fmt ./...`
- **Vet**: `go vet ./...`
- **Tidy dependencies**: `go mod tidy`

## Architecture

The project is currently minimal with only a go.mod file defining the module as `github.com/intMeric/pii-extractor` using Go 1.23.0. The architecture will evolve as the PII extraction functionality is implemented.