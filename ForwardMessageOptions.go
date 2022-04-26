package msgraph

// Structure to reply to messages.
//
// See https://docs.microsoft.com/en-us/graph/api/message-reply
type ForwardMessageOptions struct {
	Comment      string      `json:"comment,omitempty"`      // A comment to add.
	ToRecipients []Recipient `json:"toRecipients,omitempty"` // Properties to be updated when repling.
}
