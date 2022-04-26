package msgraph

// Structure to reply to messages.
//
// See https://docs.microsoft.com/en-us/graph/api/message-reply
type MessageReply struct {
	Comment string  `json:"comment,omitempty"` // A comment to add.
	Message Message `json:"message,omitempty"` // Properties to be updated when repling.
}
