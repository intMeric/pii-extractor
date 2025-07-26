package main

import (
	"encoding/json"
	"fmt"
	"log"

	piiextractor "github.com/intMeric/pii-extractor"
)

func main() {
	// Create a new RegexExtractor instance
	extractor := piiextractor.NewRegexExtractor()

	// Sample text containing various PII types
	text := `
	Hello, my name is John Doe. You can reach me at john.doe@example.com 
	or call me at (555) 123-4567. My home address is 123 Main Street, 
	New York, NY 10001. 
	
	For business purposes, my SSN is 123-45-6789 and you can send mail to 
	P.O. Box 456. My credit card number is 4111-1111-1111-1111.
	
	Server details:
	- IP Address: 192.168.1.100
	- Bitcoin wallet: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
	- Bank account (IBAN): GB82WEST12345698765432
	`

	// Extract PII from the text
	result, err := extractor.Extract(text)
	if err != nil {
		log.Fatalf("Error extracting PII: %v", err)
	}

	// Display summary
	fmt.Printf("Found %d PII entities:\n", result.Total)
	fmt.Printf("Types found: %v\n\n", result.Stats)

	for i, entity := range result.Entities {
		fmt.Printf("--- Entity %d ---\n", i+1)
		fmt.Printf("Type: %s\n", entity.Type)
		fmt.Printf("Value: %s\n", entity.GetValue())
		fmt.Printf("Count: %d\n", entity.GetCount())

		contexts := entity.GetContexts()
		if len(contexts) > 0 {
			fmt.Printf("Context: %s\n", contexts[0])
		}

		// Demonstrate type-specific casting
		switch entity.Type {
		case piiextractor.PiiTypeEmail:
			if email, ok := entity.AsEmail(); ok {
				fmt.Printf("Email domain: %s\n", getEmailDomain(email.GetValue()))
			}
		case piiextractor.PiiTypePhone:
			if phone, ok := entity.AsPhone(); ok {
				fmt.Printf("Phone country: %s\n", phone.Country)
			}
		case piiextractor.PiiTypeCreditCard:
			if cc, ok := entity.AsCreditCard(); ok {
				fmt.Printf("Card type: %s\n", cc.Type)
			}
		case piiextractor.PiiTypeIPAddress:
			if ip, ok := entity.AsIPAddress(); ok {
				fmt.Printf("IP version: %s\n", ip.Version)
			}
		}
		fmt.Println()
	}

	// Demonstrate new result methods
	fmt.Println("--- PiiExtractionResult Methods ---")
	fmt.Printf("Emails found: %d\n", len(result.GetEmails()))
	fmt.Printf("Phones found: %d\n", len(result.GetPhones()))
	fmt.Printf("Has credit cards: %t\n", result.HasType(piiextractor.PiiTypeCreditCard))
	fmt.Printf("Is empty: %t\n\n", result.IsEmpty())

	// Demonstrate JSON serialization
	fmt.Println("--- JSON Output ---")
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

// Helper function to extract domain from email
func getEmailDomain(email string) string {
	for i := len(email) - 1; i >= 0; i-- {
		if email[i] == '@' {
			return email[i+1:]
		}
	}
	return ""
}
