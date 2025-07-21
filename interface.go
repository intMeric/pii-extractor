package piiextractor

// PiiExtractor defines the main interface for extracting PII from text
type PiiExtractor interface {
	Extract(text string) ([]PiiEntity, error)
}

// PiiEntity represents a single PII item found in text
type PiiEntity struct {
	Type  string `json:"type"`  // e.g., "phone", "email", "ssn", "credit_card", etc.
	Value Pii    `json:"value"` // The actual PII value object
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
