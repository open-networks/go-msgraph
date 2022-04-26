package msgraph

// InternetMessageHeader is a key-value pair that represents an Internet message header that provides details of the network path taken from the message to the recipient. See RFC5322 for the definition.
//
// See https://docs.microsoft.com/en-us/graph/api/resources/InternetMessageHeader
type InternetMessageHeader struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
