package msgraph

// SingleValueLegacyExtendedProperty is an extended that contains a single value.
//
// See https://docs.microsoft.com/en-us/graph/api/resources/SingleValueLegacyExtendedProperty
type SingleValueLegacyExtendedProperty struct {
	ID    string `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
}
