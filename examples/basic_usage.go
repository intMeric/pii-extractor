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
	entities, err := extractor.Extract(text)
	if err != nil {
		log.Fatalf("Error extracting PII: %v", err)
	}

	// Display results
	fmt.Printf("Found %d PII entities:\n\n", len(entities))

	for i, entity := range entities {
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
		case "email":
			if email, ok := entity.AsEmail(); ok {
				fmt.Printf("Email domain: %s\n", getEmailDomain(email.GetValue()))
			}
		case "phone":
			if phone, ok := entity.AsPhone(); ok {
				fmt.Printf("Phone country: %s\n", phone.Country)
			}
		case "credit_card":
			if cc, ok := entity.AsCreditCard(); ok {
				fmt.Printf("Card type: %s\n", cc.Type)
			}
		case "ip_address":
			if ip, ok := entity.AsIPAddress(); ok {
				fmt.Printf("IP version: %s\n", ip.Version)
			}
		}
		fmt.Println()
	}

	// Demonstrate JSON serialization
	fmt.Println("--- JSON Output ---")
	jsonData, err := json.MarshalIndent(entities, "", "  ")
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