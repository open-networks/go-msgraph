package msgraph

// MultiValueLegacyExtendedProperty is an extended value that contains a slice of values.
type MultiValueLegacyExtendedProperty struct {
	ID    string   `json:"id,omitempty"`
	Value []string `json:"value,omitempty"`
}
