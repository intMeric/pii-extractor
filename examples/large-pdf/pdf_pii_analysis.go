package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	piiextractor "github.com/intMeric/pii-extractor"
	"github.com/ledongthuc/pdf"
)

func main() {
	// Get PDF file path
	pdfPath := "goog-10-q-q1-2025.pdf"
	if len(os.Args) > 1 {
		pdfPath = os.Args[1]
	}

	// Check if file exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		log.Fatalf("❌ PDF file not found: %s", pdfPath)
	}

	fmt.Printf("🔍 Analyzing PDF: %s\n", filepath.Base(pdfPath))
	fmt.Println("=" + strings.Repeat("=", 60))

	// Extract text from PDF
	fmt.Println("📄 Extracting text from PDF...")
	startExtract := time.Now()

	text, err := extractTextFromPDF(pdfPath)
	if err != nil {
		log.Fatalf("❌ Error extracting text from PDF: %v", err)
	}

	extractTime := time.Since(startExtract)

	// Display text statistics
	wordCount := len(strings.Fields(text))
	lineCount := len(strings.Split(text, "\n"))

	fmt.Printf("📊 PDF Text Statistics:\n")
	fmt.Printf("   • Size: %d characters (~%.1f KB)\n", len(text), float64(len(text))/1024)
	fmt.Printf("   • Words: %d\n", wordCount)
	fmt.Printf("   • Lines: %d\n", lineCount)
	fmt.Printf("   • Extraction time: %v\n", extractTime)
	fmt.Println()

	// Create optimized RegexExtractor
	extractor := piiextractor.NewDefaultRegexExtractor()

	// Perform PII extraction with timing
	fmt.Println("🔍 Performing PII extraction with optimized RegexExtractor...")
	startPII := time.Now()

	result, err := extractor.Extract(text)
	if err != nil {
		log.Fatalf("❌ Error extracting PII: %v", err)
	}

	piiTime := time.Since(startPII)

	// Display results summary
	fmt.Printf("🎯 PII Extraction Results:\n")
	fmt.Printf("   • Total entities found: %d\n", result.Total)
	fmt.Printf("   • Extraction time: %v\n", piiTime)
	fmt.Printf("   • Processing rate: %.2f chars/ms\n", float64(len(text))/float64(piiTime.Milliseconds()))
	fmt.Printf("   • Entities per second: %.2f\n", float64(result.Total)/piiTime.Seconds())
	fmt.Println()

	// Display breakdown by type
	fmt.Printf("📋 PII Types Found:\n")
	if result.Total == 0 {
		fmt.Println("   • No PII entities detected")
	} else {
		for piiType, count := range result.Stats {
			fmt.Printf("   • %s: %d\n", piiType.String(), count)
		}
	}
	fmt.Println()

	// Show examples of each type found (limit to 3 per type)
	if result.Total > 0 {
		fmt.Println("💡 Sample PII Entities Found:")

		// Group by type for better display
		typeExamples := make(map[string][]string)
		typeContexts := make(map[string][]string)

		for _, entity := range result.Entities {
			typeName := entity.Type.String()
			if len(typeExamples[typeName]) < 3 { // Limit to 3 examples per type
				typeExamples[typeName] = append(typeExamples[typeName], entity.GetValue())

				// Get context (truncate if too long)
				context := ""
				if contexts := entity.GetContexts(); len(contexts) > 0 {
					context = contexts[0]
					if len(context) > 80 {
						context = context[:77] + "..."
					}
				}
				typeContexts[typeName] = append(typeContexts[typeName], context)
			}
		}

		for typeName, examples := range typeExamples {
			fmt.Printf("\n   🔸 %s:\n", strings.ToUpper(typeName))
			for i, example := range examples {
				fmt.Printf("     %d. %s\n", i+1, example)
				if context := typeContexts[typeName][i]; context != "" {
					fmt.Printf("        Context: %s\n", context)
				}
			}
		}
	}
}

func extractTextFromPDF(pdfPath string) (string, error) {
	// Open PDF file
	file, reader, err := pdf.Open(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer file.Close()

	var textBuilder strings.Builder
	totalPages := reader.NumPage()

	fmt.Printf("   • Processing %d pages...\n", totalPages)

	// Extract text from each page
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		// Extract text content
		pageText, err := page.GetPlainText(nil) // Pass nil for default font map
		if err != nil {
			// Continue with other pages if one fails
			fmt.Printf("   • Warning: Could not extract text from page %d: %v\n", pageNum, err)
			continue
		}

		textBuilder.WriteString(pageText)
		textBuilder.WriteString("\n") // Add page separator

		// Progress indicator for large PDFs
		if pageNum%10 == 0 || pageNum == totalPages {
			fmt.Printf("   • Processed %d/%d pages\n", pageNum, totalPages)
		}
	}

	text := textBuilder.String()
	if len(text) == 0 {
		return "", fmt.Errorf("no text content found in PDF")
	}

	return text, nil
}
