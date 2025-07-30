# Large PDF Analysis Example

This example demonstrates the performance of the optimized RegexExtractor on large PDF documents.

## Features

- **PDF Text Extraction**: Extracts text content from multi-page PDF files
- **Performance Monitoring**: Measures extraction time and processing rates
- **Comprehensive Analysis**: Shows detailed PII statistics and examples
- **Parallel Processing Detection**: Indicates when optimizations are active

## Usage

```bash
# Install dependencies
go mod tidy

# Run analysis on the included PDF
go run pdf_analysis.go

# Or analyze a different PDF file
go run pdf_analysis.go path/to/your/document.pdf
```

- **Memory Optimization**: Pre-allocated data structures reduce GC pressure
- **Batch Processing**: Optimized entity collection and processing

