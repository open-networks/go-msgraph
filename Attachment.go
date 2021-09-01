package msgraph

import "time"

// Attachment represents content related to a event, message, or group post.
//
// See https://docs.microsoft.com/en-us/graph/api/resources/Attachment
type Attachment struct {
	ContentType          string    `json:"contentType,omitempty"`
	ID                   string    `json:"id,omitempty"`
	IsInline             bool      `json:"isInline,omitempty"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime,omitempty"`
	Name                 string    `json:"name,omitempty"`
	Size                 int32     `json:"size,omitempty"`
}
