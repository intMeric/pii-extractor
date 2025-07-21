package piiextractor

// Pii is a common interface for all PII types
type Pii interface {
	GetValue() string
	GetContexts() []string
	GetCount() int
	String() string
}

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

// Pii interface implementations

// Phone methods
func (p Phone) GetValue() string {
	return p.Value
}

func (p Phone) GetContexts() []string {
	return p.Contexts
}

func (p Phone) GetCount() int {
	return p.Count
}

func (p Phone) String() string {
	return p.Value
}

// Email methods
func (e Email) GetValue() string {
	return e.Value
}

func (e Email) GetContexts() []string {
	return e.Contexts
}

func (e Email) GetCount() int {
	return e.Count
}

func (e Email) String() string {
	return e.Value
}

// SSN methods
func (s SSN) GetValue() string {
	return s.Value
}

func (s SSN) GetContexts() []string {
	return s.Contexts
}

func (s SSN) GetCount() int {
	return s.Count
}

func (s SSN) String() string {
	return s.Value
}

// ZipCode methods
func (z ZipCode) GetValue() string {
	return z.Value
}

func (z ZipCode) GetContexts() []string {
	return z.Contexts
}

func (z ZipCode) GetCount() int {
	return z.Count
}

func (z ZipCode) String() string {
	return z.Value
}

// PoBox methods
func (p PoBox) GetValue() string {
	return p.Value
}

func (p PoBox) GetContexts() []string {
	return p.Contexts
}

func (p PoBox) GetCount() int {
	return p.Count
}

func (p PoBox) String() string {
	return p.Value
}

// StreetAddress methods
func (s StreetAddress) GetValue() string {
	return s.Value
}

func (s StreetAddress) GetContexts() []string {
	return s.Contexts
}

func (s StreetAddress) GetCount() int {
	return s.Count
}

func (s StreetAddress) String() string {
	return s.Value
}

// CreditCard methods
func (c CreditCard) GetValue() string {
	return c.Value
}

func (c CreditCard) GetContexts() []string {
	return c.Contexts
}

func (c CreditCard) GetCount() int {
	return c.Count
}

func (c CreditCard) String() string {
	return c.Value
}

// IPAddress methods
func (i IPAddress) GetValue() string {
	return i.Value
}

func (i IPAddress) GetContexts() []string {
	return i.Contexts
}

func (i IPAddress) GetCount() int {
	return i.Count
}

func (i IPAddress) String() string {
	return i.Value
}

// BtcAddress methods
func (b BtcAddress) GetValue() string {
	return b.Value
}

func (b BtcAddress) GetContexts() []string {
	return b.Contexts
}

func (b BtcAddress) GetCount() int {
	return b.Count
}

func (b BtcAddress) String() string {
	return b.Value
}

// IBAN methods
func (i IBAN) GetValue() string {
	return i.Value
}

func (i IBAN) GetContexts() []string {
	return i.Contexts
}

func (i IBAN) GetCount() int {
	return i.Count
}

func (i IBAN) String() string {
	return i.Value
}