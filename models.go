package piiextractor

// Pii is a common interface for all PII types
type Pii interface {
	GetValue() string
	GetContexts() []string
	GetCount() int
	String() string
}

// BasePii contains the common fields for all PII types
type BasePii struct {
	Value    string
	Contexts []string
	Count    int
}

// GetValue returns the string value of the PII
func (b BasePii) GetValue() string {
	return b.Value
}

// GetContexts returns the contexts where this PII was found
func (b BasePii) GetContexts() []string {
	return b.Contexts
}

// GetCount returns how many times this PII value was found
func (b BasePii) GetCount() int {
	return b.Count
}

// String returns a string representation of the PII
func (b BasePii) String() string {
	return b.Value
}

// AddContext adds a context to the PII
func (b *BasePii) AddContext(context string) {
	b.Contexts = append(b.Contexts, context)
}

// IncrementCount increments the count of the PII
func (b *BasePii) IncrementCount() {
	b.Count++
}

// Phone represents a phone number value object
type Phone struct {
	BasePii
	Country string
}

// Email represents an email address value object
type Email struct {
	BasePii
}

// SSN represents a Social Security Number value object
type SSN struct {
	BasePii
	Country string
}

// ZipCode represents a postal code value object
type ZipCode struct {
	BasePii
	Country string
}

// PoBox represents a Post Office Box value object
type PoBox struct {
	BasePii
	Country string
}

// StreetAddress represents a street address value object
type StreetAddress struct {
	BasePii
	Country string
}

// CreditCard represents a credit card number value object
type CreditCard struct {
	BasePii
	Type string // "visa", "mastercard", "generic"
}

// IPAddress represents an IP address value object
type IPAddress struct {
	BasePii
	Version string // "ipv4", "ipv6"
}

// BtcAddress represents a Bitcoin address value object
type BtcAddress struct {
	BasePii
}

// IBAN represents an International Bank Account Number value object
type IBAN struct {
	BasePii
	Country string
}

// NewPhoneUS creates a new US phone number value object
func NewPhoneUS(value string) Phone {
	return Phone{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
		Country: "US",
	}
}

// NewPhone creates a new phone number value object
func NewPhone(value, country string) Phone {
	return Phone{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
		Country: country,
	}
}

// NewEmail creates a new email value object
func NewEmail(value string) Email {
	return Email{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
	}
}

// NewSSNUS creates a new US SSN value object
func NewSSNUS(value string) SSN {
	return SSN{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
		Country: "US",
	}
}

// NewZipCodeUS creates a new US zip code value object
func NewZipCodeUS(value string) ZipCode {
	return ZipCode{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
		Country: "US",
	}
}

// NewPoBoxUS creates a new US P.O. Box value object
func NewPoBoxUS(value string) PoBox {
	return PoBox{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
		Country: "US",
	}
}

// NewStreetAddressUS creates a new US street address value object
func NewStreetAddressUS(value string) StreetAddress {
	return StreetAddress{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
		Country: "US",
	}
}

// NewCreditCard creates a new credit card value object
func NewCreditCard(value, cardType string) CreditCard {
	return CreditCard{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
		Type:    cardType,
	}
}

// NewIPv4 creates a new IPv4 address value object
func NewIPv4(value string) IPAddress {
	return IPAddress{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
		Version: "ipv4",
	}
}

// NewIPv6 creates a new IPv6 address value object
func NewIPv6(value string) IPAddress {
	return IPAddress{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
		Version: "ipv6",
	}
}

// NewBtcAddress creates a new Bitcoin address value object
func NewBtcAddress(value string) BtcAddress {
	return BtcAddress{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
	}
}

// NewIBAN creates a new IBAN value object
func NewIBAN(value string) IBAN {
	country := ""
	if len(value) >= 2 {
		country = value[:2]
	}
	return IBAN{
		BasePii: BasePii{Value: value, Contexts: []string{}, Count: 1},
		Country: country,
	}
}
