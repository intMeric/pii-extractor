package internal

// Phone represents a phone number value object
type Phone struct {
	Value     string
	Country   string
	Contexts  []string
	Count     int
}

// Email represents an email address value object
type Email struct {
	Value    string
	Contexts []string
	Count    int
}

// SSN represents a Social Security Number value object
type SSN struct {
	Value    string
	Country  string
	Contexts []string
	Count    int
}

// ZipCode represents a postal code value object
type ZipCode struct {
	Value    string
	Country  string
	Contexts []string
	Count    int
}

// PoBox represents a Post Office Box value object
type PoBox struct {
	Value    string
	Country  string
	Contexts []string
	Count    int
}

// StreetAddress represents a street address value object
type StreetAddress struct {
	Value    string
	Country  string
	Contexts []string
	Count    int
}

// CreditCard represents a credit card number value object
type CreditCard struct {
	Value    string
	Type     string // "visa", "mastercard", "generic"
	Contexts []string
	Count    int
}

// IPAddress represents an IP address value object
type IPAddress struct {
	Value    string
	Version  string // "ipv4", "ipv6"
	Contexts []string
	Count    int
}

// BtcAddress represents a Bitcoin address value object
type BtcAddress struct {
	Value    string
	Contexts []string
	Count    int
}

// IBAN represents an International Bank Account Number value object
type IBAN struct {
	Value    string
	Country  string
	Contexts []string
	Count    int
}

// NewPhoneUS creates a new US phone number value object
func NewPhoneUS(value string) Phone {
	return Phone{Value: value, Country: "US", Contexts: []string{}, Count: 1}
}

// NewPhone creates a new phone number value object
func NewPhone(value, country string) Phone {
	return Phone{Value: value, Country: country, Contexts: []string{}, Count: 1}
}

// NewEmail creates a new email value object
func NewEmail(value string) Email {
	return Email{Value: value, Contexts: []string{}, Count: 1}
}

// NewSSNUS creates a new US SSN value object
func NewSSNUS(value string) SSN {
	return SSN{Value: value, Country: "US", Contexts: []string{}, Count: 1}
}

// NewZipCodeUS creates a new US zip code value object
func NewZipCodeUS(value string) ZipCode {
	return ZipCode{Value: value, Country: "US", Contexts: []string{}, Count: 1}
}

// NewPoBoxUS creates a new US P.O. Box value object
func NewPoBoxUS(value string) PoBox {
	return PoBox{Value: value, Country: "US", Contexts: []string{}, Count: 1}
}

// NewStreetAddressUS creates a new US street address value object
func NewStreetAddressUS(value string) StreetAddress {
	return StreetAddress{Value: value, Country: "US", Contexts: []string{}, Count: 1}
}

// NewCreditCard creates a new credit card value object
func NewCreditCard(value, cardType string) CreditCard {
	return CreditCard{Value: value, Type: cardType, Contexts: []string{}, Count: 1}
}

// NewIPv4 creates a new IPv4 address value object
func NewIPv4(value string) IPAddress {
	return IPAddress{Value: value, Version: "ipv4", Contexts: []string{}, Count: 1}
}

// NewIPv6 creates a new IPv6 address value object
func NewIPv6(value string) IPAddress {
	return IPAddress{Value: value, Version: "ipv6", Contexts: []string{}, Count: 1}
}

// NewBtcAddress creates a new Bitcoin address value object
func NewBtcAddress(value string) BtcAddress {
	return BtcAddress{Value: value, Contexts: []string{}, Count: 1}
}

// NewIBAN creates a new IBAN value object
func NewIBAN(value string) IBAN {
	country := ""
	if len(value) >= 2 {
		country = value[:2]
	}
	return IBAN{Value: value, Country: country, Contexts: []string{}, Count: 1}
}

// String methods for display
func (p Phone) String() string {
	return p.Value
}

func (e Email) String() string {
	return e.Value
}

func (s SSN) String() string {
	return s.Value
}

func (z ZipCode) String() string {
	return z.Value
}

func (p PoBox) String() string {
	return p.Value
}

func (s StreetAddress) String() string {
	return s.Value
}

func (c CreditCard) String() string {
	return c.Value
}

func (i IPAddress) String() string {
	return i.Value
}

func (b BtcAddress) String() string {
	return b.Value
}

func (i IBAN) String() string {
	return i.Value
}