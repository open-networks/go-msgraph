package msgraph

// Member represents one member of ms graph
type Member struct {
	Type        string `json:"@odata.type,omitempty"`
	Id          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}
