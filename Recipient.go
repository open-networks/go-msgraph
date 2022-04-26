package msgraph

// Represents a user that has sent or been sent an event, message, or group post.
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/resources/Recipient
type Recipient struct {
	EmailAddress `json:"emailAddress,omitempty"`
}
