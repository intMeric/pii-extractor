package pii

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

// ValidationResult contains the result of LLM validation
type ValidationResult struct {
	Valid      bool    `json:"valid"`
	Confidence float64 `json:"confidence"`
	Reasoning  string  `json:"reasoning,omitempty"`
	Provider   string  `json:"provider"`
	Model      string  `json:"model"`
}

// ValidationStats contains statistics about LLM validation results
type ValidationStats struct {
	TotalValidated    int     `json:"total_validated"`
	ValidCount        int     `json:"valid_count"`
	InvalidCount      int     `json:"invalid_count"`
	AverageConfidence float64 `json:"average_confidence"`
	Provider          string  `json:"provider,omitempty"`
	Model             string  `json:"model,omitempty"`
}

// Pii interface that all PII value objects must implement
type Pii interface {
	String() string
	GetValue() string
	GetContexts() []string
	GetCount() int
}

// BasePii provides common functionality for all PII types
type BasePii struct {
	Value    string   `json:"value"`
	Contexts []string `json:"contexts"`
	Count    int      `json:"count"`
}

// String returns the string representation of the PII
func (p BasePii) String() string {
	return p.Value
}

// GetValue returns the string value
func (p BasePii) GetValue() string {
	return p.Value
}

// GetContexts returns the contexts where this PII was found
func (p BasePii) GetContexts() []string {
	return p.Contexts
}

// GetCount returns how many times this PII value was found
func (p BasePii) GetCount() int {
	return p.Count
}

// IncrementCount increases the occurrence count
func (p *BasePii) IncrementCount() {
	p.Count++
}

// AddContext adds a new context if it doesn't already exist
func (p *BasePii) AddContext(context string) {
	for _, existingContext := range p.Contexts {
		if existingContext == context {
			return
		}
	}
	p.Contexts = append(p.Contexts, context)
}

// PII value objects

// Phone represents a phone number
type Phone struct {
	BasePii
	Country string `json:"country,omitempty"`
}

// Email represents an email address
type Email struct {
	BasePii
}

// SSN represents a Social Security Number
type SSN struct {
	BasePii
	Country string `json:"country,omitempty"`
}

// ZipCode represents a ZIP/postal code
type ZipCode struct {
	BasePii
	Country string `json:"country,omitempty"`
}

// StreetAddress represents a street address
type StreetAddress struct {
	BasePii
	Country string `json:"country,omitempty"`
}

// PoBox represents a P.O. Box
type PoBox struct {
	BasePii
	Country string `json:"country,omitempty"`
}

// CreditCard represents a credit card number
type CreditCard struct {
	BasePii
	Type string `json:"type,omitempty"` // visa, mastercard, etc.
}

// IPAddress represents an IP address
type IPAddress struct {
	BasePii
	Version string `json:"version,omitempty"` // ipv4, ipv6
}

// BtcAddress represents a Bitcoin address
type BtcAddress struct {
	BasePii
}

// IBAN represents an International Bank Account Number
type IBAN struct {
	BasePii
	Country string `json:"country,omitempty"`
}

// Constructor functions for PII types

// NewEmail creates a new Email PII value
func NewEmail(value string) Email {
	return Email{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
	}
}

// NewPhoneUS creates a new US Phone PII value
func NewPhoneUS(value string) Phone {
	return Phone{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
		Country: "US",
	}
}

// NewPhone creates a new Phone PII value with specified country
func NewPhone(value, country string) Phone {
	return Phone{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
		Country: country,
	}
}

// NewSSN creates a new SSN PII value
func NewSSN(value string) SSN {
	return SSN{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
		Country: "US",
	}
}

// NewZipCode creates a new ZipCode PII value
func NewZipCode(value, country string) ZipCode {
	return ZipCode{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
		Country: country,
	}
}

// NewStreetAddress creates a new StreetAddress PII value
func NewStreetAddress(value, country string) StreetAddress {
	return StreetAddress{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
		Country: country,
	}
}

// NewPoBox creates a new PoBox PII value
func NewPoBox(value, country string) PoBox {
	return PoBox{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
		Country: country,
	}
}

// NewCreditCard creates a new CreditCard PII value
func NewCreditCard(value, cardType string) CreditCard {
	return CreditCard{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
		Type: cardType,
	}
}

// NewIPAddress creates a new IPAddress PII value
func NewIPAddress(value, version string) IPAddress {
	return IPAddress{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
		Version: version,
	}
}

// NewBtcAddress creates a new BtcAddress PII value
func NewBtcAddress(value string) BtcAddress {
	return BtcAddress{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
	}
}

// NewIBAN creates a new IBAN PII value
func NewIBAN(value, country string) IBAN {
	return IBAN{
		BasePii: BasePii{
			Value:    value,
			Contexts: []string{},
			Count:    1,
		},
		Country: country,
	}
}

// PiiEntity represents a single PII item found in text
type PiiEntity struct {
	Type       PiiType           `json:"type"`                 // The type of PII (phone, email, ssn, etc.)
	Value      Pii               `json:"value"`                // The actual PII value object
	Validation *ValidationResult `json:"validation,omitempty"` // Optional LLM validation result
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

// AsZipCode attempts to cast the value to a ZipCode
func (p PiiEntity) AsZipCode() (ZipCode, bool) {
	return GetTypedValue[ZipCode](p)
}

// AsStreetAddress attempts to cast the value to a StreetAddress
func (p PiiEntity) AsStreetAddress() (StreetAddress, bool) {
	return GetTypedValue[StreetAddress](p)
}

// AsPoBox attempts to cast the value to a PoBox
func (p PiiEntity) AsPoBox() (PoBox, bool) {
	return GetTypedValue[PoBox](p)
}

// AsBtcAddress attempts to cast the value to a BtcAddress
func (p PiiEntity) AsBtcAddress() (BtcAddress, bool) {
	return GetTypedValue[BtcAddress](p)
}

// AsIBAN attempts to cast the value to an IBAN
func (p PiiEntity) AsIBAN() (IBAN, bool) {
	return GetTypedValue[IBAN](p)
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

// IsValidated returns true if the entity has been validated by an LLM
func (p PiiEntity) IsValidated() bool {
	return p.Validation != nil
}

// IsValid returns true if the entity is validated and marked as valid
func (p PiiEntity) IsValid() bool {
	return p.Validation != nil && p.Validation.Valid
}

// GetValidationConfidence returns the validation confidence score (0.0 if not validated)
func (p PiiEntity) GetValidationConfidence() float64 {
	if p.Validation != nil {
		return p.Validation.Confidence
	}
	return 0.0
}

// PiiExtractionResult represents the result of a PII extraction operation
type PiiExtractionResult struct {
	Entities        []PiiEntity      `json:"entities"`
	Stats           map[PiiType]int  `json:"stats"`
	Total           int              `json:"total"`
	ValidationStats *ValidationStats `json:"validation_stats,omitempty"` // Optional validation statistics
}

// NewPiiExtractionResult creates a new PiiExtractionResult from entities with deduplication
func NewPiiExtractionResult(entities []PiiEntity) *PiiExtractionResult {
	// Deduplicate entities
	dedupedEntities := deduplicateEntities(entities)
	
	stats := make(map[PiiType]int)
	for _, entity := range dedupedEntities {
		stats[entity.Type]++
	}

	return &PiiExtractionResult{
		Entities: dedupedEntities,
		Stats:    stats,
		Total:    len(dedupedEntities),
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

// IsEmpty returns true if no PII entities were found
func (r *PiiExtractionResult) IsEmpty() bool {
	return r.Total == 0
}

// HasType returns true if the result contains entities of the specified type
func (r *PiiExtractionResult) HasType(piiType PiiType) bool {
	return r.Stats[piiType] > 0
}

// GetValidatedEntities returns only entities that have been validated by LLM
func (r *PiiExtractionResult) GetValidatedEntities() []PiiEntity {
	var result []PiiEntity
	for _, entity := range r.Entities {
		if entity.IsValidated() {
			result = append(result, entity)
		}
	}
	return result
}

// GetValidEntities returns only entities that are validated and marked as valid
func (r *PiiExtractionResult) GetValidEntities() []PiiEntity {
	var result []PiiEntity
	for _, entity := range r.Entities {
		if entity.IsValid() {
			result = append(result, entity)
		}
	}
	return result
}

// GetInvalidEntities returns only entities that are validated but marked as invalid
func (r *PiiExtractionResult) GetInvalidEntities() []PiiEntity {
	var result []PiiEntity
	for _, entity := range r.Entities {
		if entity.IsValidated() && !entity.IsValid() {
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

// International extraction convenience methods

// GetZipCodesByCountry returns all ZIP/postal code entities for a specific country
func (r *PiiExtractionResult) GetZipCodesByCountry(country string) []PiiEntity {
	var result []PiiEntity
	for _, entity := range r.GetZipCodes() {
		if zipCode, ok := entity.AsZipCode(); ok && zipCode.Country == country {
			result = append(result, entity)
		}
	}
	return result
}

// GetStreetAddressesByCountry returns all street address entities for a specific country
func (r *PiiExtractionResult) GetStreetAddressesByCountry(country string) []PiiEntity {
	var result []PiiEntity
	for _, entity := range r.GetStreetAddresses() {
		if address, ok := entity.AsStreetAddress(); ok && address.Country == country {
			result = append(result, entity)
		}
	}
	return result
}

// GetPhonesByCountry returns all phone entities for a specific country
func (r *PiiExtractionResult) GetPhonesByCountry(country string) []PiiEntity {
	var result []PiiEntity
	for _, entity := range r.GetPhones() {
		if phone, ok := entity.AsPhone(); ok && phone.Country == country {
			result = append(result, entity)
		}
	}
	return result
}

// Convenience methods for specific countries

// GetUKEntities returns all UK-specific PII entities (postal codes and addresses)
func (r *PiiExtractionResult) GetUKEntities() []PiiEntity {
	var result []PiiEntity
	result = append(result, r.GetZipCodesByCountry("UK")...)
	result = append(result, r.GetStreetAddressesByCountry("UK")...)
	return result
}

// GetFranceEntities returns all France-specific PII entities (postal codes and addresses)
func (r *PiiExtractionResult) GetFranceEntities() []PiiEntity {
	var result []PiiEntity
	result = append(result, r.GetZipCodesByCountry("France")...)
	result = append(result, r.GetStreetAddressesByCountry("France")...)
	return result
}

// GetSpainEntities returns all Spain-specific PII entities (postal codes and addresses)
func (r *PiiExtractionResult) GetSpainEntities() []PiiEntity {
	var result []PiiEntity
	result = append(result, r.GetZipCodesByCountry("Spain")...)
	result = append(result, r.GetStreetAddressesByCountry("Spain")...)
	return result
}

// GetItalyEntities returns all Italy-specific PII entities (postal codes and addresses)
func (r *PiiExtractionResult) GetItalyEntities() []PiiEntity {
	var result []PiiEntity
	result = append(result, r.GetZipCodesByCountry("Italy")...)
	result = append(result, r.GetStreetAddressesByCountry("Italy")...)
	return result
}

// GetUSEntities returns all US-specific PII entities (phones, SSNs, ZIP codes, addresses, P.O. boxes)
func (r *PiiExtractionResult) GetUSEntities() []PiiEntity {
	var result []PiiEntity
	result = append(result, r.GetPhonesByCountry("US")...)
	result = append(result, r.GetZipCodesByCountry("US")...)
	result = append(result, r.GetStreetAddressesByCountry("US")...)
	result = append(result, r.GetSSNs()...)    // SSNs are US-specific
	result = append(result, r.GetPoBoxes()...) // P.O. boxes are currently US-specific
	return result
}

// deduplicateEntities removes duplicate entities and merges their contexts
func deduplicateEntities(entities []PiiEntity) []PiiEntity {
	entityMap := make(map[string]*PiiEntity)
	
	for _, entity := range entities {
		key := generateEntityKey(entity)
		
		if existing, exists := entityMap[key]; exists {
			// Merge contexts and update count
			mergeEntityContexts(existing, &entity)
		} else {
			// Create a copy to avoid modifying the original
			entityCopy := entity
			entityMap[key] = &entityCopy
		}
	}
	
	// Convert map back to slice
	result := make([]PiiEntity, 0, len(entityMap))
	for _, entity := range entityMap {
		result = append(result, *entity)
	}
	
	return result
}

// generateEntityKey creates a unique key for an entity based on type and value
func generateEntityKey(entity PiiEntity) string {
	return entity.Type.String() + ":" + entity.GetValue()
}

// mergeEntityContexts merges contexts from source entity into target entity
func mergeEntityContexts(target, source *PiiEntity) {
	if target.Value == nil || source.Value == nil {
		return
	}
	
	// Get the underlying value objects
	targetValue := target.Value
	sourceValue := source.Value
	
	// Merge contexts
	sourceContexts := sourceValue.GetContexts()
	
	switch tv := targetValue.(type) {
	case Phone:
		if sv, ok := sourceValue.(Phone); ok {
			// Unify countries if different - set to empty string if they differ
			if tv.Country != sv.Country && tv.Country != "" && sv.Country != "" {
				tv.Country = ""
			}
			// Add new contexts
			for _, context := range sourceContexts {
				tv.BasePii.AddContext(context)
			}
			tv.BasePii.Count += sv.BasePii.Count
			target.Value = tv
		}
	case Email:
		if sv, ok := sourceValue.(Email); ok {
			for _, context := range sourceContexts {
				tv.BasePii.AddContext(context)
			}
			tv.BasePii.Count += sv.BasePii.Count
			target.Value = tv
		}
	case SSN:
		if sv, ok := sourceValue.(SSN); ok {
			if tv.Country != sv.Country && tv.Country != "" && sv.Country != "" {
				tv.Country = ""
			}
			for _, context := range sourceContexts {
				tv.BasePii.AddContext(context)
			}
			tv.BasePii.Count += sv.BasePii.Count
			target.Value = tv
		}
	case ZipCode:
		if sv, ok := sourceValue.(ZipCode); ok {
			if tv.Country != sv.Country && tv.Country != "" && sv.Country != "" {
				tv.Country = ""
			}
			for _, context := range sourceContexts {
				tv.BasePii.AddContext(context)
			}
			tv.BasePii.Count += sv.BasePii.Count
			target.Value = tv
		}
	case StreetAddress:
		if sv, ok := sourceValue.(StreetAddress); ok {
			if tv.Country != sv.Country && tv.Country != "" && sv.Country != "" {
				tv.Country = ""
			}
			for _, context := range sourceContexts {
				tv.BasePii.AddContext(context)
			}
			tv.BasePii.Count += sv.BasePii.Count
			target.Value = tv
		}
	case PoBox:
		if sv, ok := sourceValue.(PoBox); ok {
			if tv.Country != sv.Country && tv.Country != "" && sv.Country != "" {
				tv.Country = ""
			}
			for _, context := range sourceContexts {
				tv.BasePii.AddContext(context)
			}
			tv.BasePii.Count += sv.BasePii.Count
			target.Value = tv
		}
	case CreditCard:
		if sv, ok := sourceValue.(CreditCard); ok {
			if tv.Type != sv.Type && tv.Type != "" && sv.Type != "" {
				tv.Type = "generic"
			}
			for _, context := range sourceContexts {
				tv.BasePii.AddContext(context)
			}
			tv.BasePii.Count += sv.BasePii.Count
			target.Value = tv
		}
	case IPAddress:
		if sv, ok := sourceValue.(IPAddress); ok {
			if tv.Version != sv.Version && tv.Version != "" && sv.Version != "" {
				tv.Version = ""
			}
			for _, context := range sourceContexts {
				tv.BasePii.AddContext(context)
			}
			tv.BasePii.Count += sv.BasePii.Count
			target.Value = tv
		}
	case BtcAddress:
		if sv, ok := sourceValue.(BtcAddress); ok {
			for _, context := range sourceContexts {
				tv.BasePii.AddContext(context)
			}
			tv.BasePii.Count += sv.BasePii.Count
			target.Value = tv
		}
	case IBAN:
		if sv, ok := sourceValue.(IBAN); ok {
			if tv.Country != sv.Country && tv.Country != "" && sv.Country != "" {
				tv.Country = ""
			}
			for _, context := range sourceContexts {
				tv.BasePii.AddContext(context)
			}
			tv.BasePii.Count += sv.BasePii.Count
			target.Value = tv
		}
	}
}