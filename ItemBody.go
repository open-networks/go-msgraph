package msgraph

// ItemBody represents the body of an item, either a message, event, or group post.
//
// See https://docs.microsoft.com/en-us/graph/api/resources/ItemBody
type ItemBody struct {
	Content     string `json:"content,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}
