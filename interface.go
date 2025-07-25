package piiextractor

// PiiType represents the type of PII entity
type PiiType int

const (
	PiiTypePhone PiiType = iota
	PiiTypeEmail
	PiiTypeSSN
	PiiTypeZipCode
	PiiTypePoBox
	PiiTypeStreetAddress
	PiiTypeCreditCard
	PiiTypeIPAddress
	PiiTypeBtcAddress
	PiiTypeIBAN
)

// String returns the string representation of the PII type
func (p PiiType) String() string {
	switch p {
	case PiiTypePhone:
		return "phone"
	case PiiTypeEmail:
		return "email"
	case PiiTypeSSN:
		return "ssn"
	case PiiTypeZipCode:
		return "zip_code"
	case PiiTypePoBox:
		return "po_box"
	case PiiTypeStreetAddress:
		return "street_address"
	case PiiTypeCreditCard:
		return "credit_card"
	case PiiTypeIPAddress:
		return "ip_address"
	case PiiTypeBtcAddress:
		return "btc_address"
	case PiiTypeIBAN:
		return "iban"
	default:
		return "unknown"
	}
}

// PiiExtractionResult represents the result of a PII extraction operation
type PiiExtractionResult struct {
	Entities []PiiEntity          `json:"entities"`
	Stats    map[PiiType]int      `json:"stats"`
	Total    int                  `json:"total"`
}

// PiiExtractor defines the main interface for extracting PII from text
type PiiExtractor interface {
	Extract(text string) (*PiiExtractionResult, error)
}

// PiiEntity represents a single PII item found in text
type PiiEntity struct {
	Type  PiiType `json:"type"`  // The type of PII (phone, email, ssn, etc.)
	Value Pii     `json:"value"` // The actual PII value object
}

// GetTypedValue performs a safe type assertion for the PII value
func GetTypedValue[T Pii](entity PiiEntity) (T, bool) {
	if value, ok := entity.Value.(T); ok {
		return value, true
	}
	var zero T
	return zero, false
}

// String returns a string representation of the PII entity
func (p PiiEntity) String() string {
	if p.Value != nil {
		return p.Value.String()
	}
	return ""
}

// GetValue returns the string value of the PII
func (p PiiEntity) GetValue() string {
	if p.Value != nil {
		return p.Value.GetValue()
	}
	return ""
}

// GetContexts returns the contexts where this PII was found
func (p PiiEntity) GetContexts() []string {
	if p.Value != nil {
		return p.Value.GetContexts()
	}
	return []string{}
}

// GetCount returns how many times this PII value was found
func (p PiiEntity) GetCount() int {
	if p.Value != nil {
		return p.Value.GetCount()
	}
	return 0
}

// Convenience methods for common type assertions

// AsPhone attempts to cast the value to a Phone
func (p PiiEntity) AsPhone() (Phone, bool) {
	return GetTypedValue[Phone](p)
}

// AsEmail attempts to cast the value to an Email
func (p PiiEntity) AsEmail() (Email, bool) {
	return GetTypedValue[Email](p)
}

// AsSSN attempts to cast the value to an SSN
func (p PiiEntity) AsSSN() (SSN, bool) {
	return GetTypedValue[SSN](p)
}

// AsCreditCard attempts to cast the value to a CreditCard
func (p PiiEntity) AsCreditCard() (CreditCard, bool) {
	return GetTypedValue[CreditCard](p)
}

// AsIPAddress attempts to cast the value to an IPAddress
func (p PiiEntity) AsIPAddress() (IPAddress, bool) {
	return GetTypedValue[IPAddress](p)
}

// Convenience type checking methods

// IsPhone returns true if the entity is a phone number
func (p PiiEntity) IsPhone() bool {
	return p.Type == PiiTypePhone
}

// IsEmail returns true if the entity is an email address
func (p PiiEntity) IsEmail() bool {
	return p.Type == PiiTypeEmail
}

// IsSSN returns true if the entity is a Social Security Number
func (p PiiEntity) IsSSN() bool {
	return p.Type == PiiTypeSSN
}

// IsCreditCard returns true if the entity is a credit card number
func (p PiiEntity) IsCreditCard() bool {
	return p.Type == PiiTypeCreditCard
}

// IsIPAddress returns true if the entity is an IP address
func (p PiiEntity) IsIPAddress() bool {
	return p.Type == PiiTypeIPAddress
}

// IsZipCode returns true if the entity is a ZIP code
func (p PiiEntity) IsZipCode() bool {
	return p.Type == PiiTypeZipCode
}

// IsPoBox returns true if the entity is a P.O. Box
func (p PiiEntity) IsPoBox() bool {
	return p.Type == PiiTypePoBox
}

// IsStreetAddress returns true if the entity is a street address
func (p PiiEntity) IsStreetAddress() bool {
	return p.Type == PiiTypeStreetAddress
}

// IsBtcAddress returns true if the entity is a Bitcoin address
func (p PiiEntity) IsBtcAddress() bool {
	return p.Type == PiiTypeBtcAddress
}

// IsIBAN returns true if the entity is an IBAN
func (p PiiEntity) IsIBAN() bool {
	return p.Type == PiiTypeIBAN
}

// PiiExtractionResult utility methods
// NewPiiExtractionResult creates a new PiiExtractionResult from entities
func NewPiiExtractionResult(entities []PiiEntity) *PiiExtractionResult {
	stats := make(map[PiiType]int)
	for _, entity := range entities {
		stats[entity.Type]++
	}
	
	return &PiiExtractionResult{
		Entities: entities,
		Stats:    stats,
		Total:    len(entities),
	}
}

// GetEntitiesByType returns all entities of a specific type
func (r *PiiExtractionResult) GetEntitiesByType(piiType PiiType) []PiiEntity {
	var result []PiiEntity
	for _, entity := range r.Entities {
		if entity.Type == piiType {
			result = append(result, entity)
		}
	}
	return result
}

// GetPhones returns all phone entities
func (r *PiiExtractionResult) GetPhones() []PiiEntity {
	return r.GetEntitiesByType(PiiTypePhone)
}

// GetEmails returns all email entities
func (r *PiiExtractionResult) GetEmails() []PiiEntity {
	return r.GetEntitiesByType(PiiTypeEmail)
}

// GetSSNs returns all SSN entities
func (r *PiiExtractionResult) GetSSNs() []PiiEntity {
	return r.GetEntitiesByType(PiiTypeSSN)
}

// GetCreditCards returns all credit card entities
func (r *PiiExtractionResult) GetCreditCards() []PiiEntity {
	return r.GetEntitiesByType(PiiTypeCreditCard)
}

// GetIPAddresses returns all IP address entities
func (r *PiiExtractionResult) GetIPAddresses() []PiiEntity {
	return r.GetEntitiesByType(PiiTypeIPAddress)
}

// GetZipCodes returns all ZIP code entities
func (r *PiiExtractionResult) GetZipCodes() []PiiEntity {
	return r.GetEntitiesByType(PiiTypeZipCode)
}

// GetPoBoxes returns all P.O. Box entities
func (r *PiiExtractionResult) GetPoBoxes() []PiiEntity {
	return r.GetEntitiesByType(PiiTypePoBox)
}

// GetStreetAddresses returns all street address entities
func (r *PiiExtractionResult) GetStreetAddresses() []PiiEntity {
	return r.GetEntitiesByType(PiiTypeStreetAddress)
}

// GetBtcAddresses returns all Bitcoin address entities
func (r *PiiExtractionResult) GetBtcAddresses() []PiiEntity {
	return r.GetEntitiesByType(PiiTypeBtcAddress)
}

// GetIBANs returns all IBAN entities
func (r *PiiExtractionResult) GetIBANs() []PiiEntity {
	return r.GetEntitiesByType(PiiTypeIBAN)
}

// HasType returns true if the result contains entities of the specified type
func (r *PiiExtractionResult) HasType(piiType PiiType) bool {
	return r.Stats[piiType] > 0
}

// IsEmpty returns true if no PII entities were found
func (r *PiiExtractionResult) IsEmpty() bool {
	return r.Total == 0
}
