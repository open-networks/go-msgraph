package msgraph

import (
	"time"
)

// Message represents a single email message.
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/resources/Message
type Message struct {
	BCCRecipients                 []Recipient                         `json:"bccRecipients,omitempty"`
	Body                          *ItemBody                           `json:"body,omitempty"`
	BodyPreview                   string                              `json:"bodyPreview,omitempty"`
	Categories                    []string                            `json:"categories,omitempty"`
	CCRecipients                  []Recipient                         `json:"ccRecipients,omitempty"`
	ChangeKey                     string                              `json:"changeKey,omitempty"`
	ConversationID                string                              `json:"conversationId,omitempty"`
	ConversationIndex             string                              `json:"conversationIndex,omitempty"`
	CreatedDateTime               time.Time                           `json:"createdDateTime,omitempty"`
	Flag                          *FollowupFlag                       `json:"flag,omitempty"`
	From                          *Recipient                          `json:"from,omitempty"`
	HasAttachments                *bool                               `json:"hasAttachments,omitempty"`
	ID                            string                              `json:"id,omitempty"`
	Importance                    MessageImportance                   `json:"importance,omitempty"`
	InferenceClassification       MessageClassification               `json:"inferenceClassification,omitempty"`
	InternetMessageHeaders        []InternetMessageHeader             `json:"internetMessageHeaders,omitempty"`
	InternetMessageID             string                              `json:"internetMessageId,omitempty"`
	IsDeliveryReceiptRequested    *bool                               `json:"isDeliveryReceiptRequested,omitempty"`
	IsDraft                       *bool                               `json:"isDraft,omitempty"`
	IsRead                        *bool                               `json:"isRead,omitempty"`
	IsReadReceiptRequested        *bool                               `json:"isReadReceiptRequested,omitempty"`
	LastModifiedDateTime          time.Time                           `json:"lastModifiedDateTime,omitempty"`
	ParentFolderId                string                              `json:"parentFolderId,omitempty"`
	ReceivedDateTime              time.Time                           `json:"receivedDateTime,omitempty"`
	ReplyTo                       []Recipient                         `json:"replyTo,omitempty"`
	Sender                        *Recipient                          `json:"sender,omitempty"`
	SentDateTime                  time.Time                           `json:"sentDateTime,omitempty"`
	Subject                       string                              `json:"subject,omitempty"`
	ToRecipients                  []Recipient                         `json:"toRecipients,omitempty"`
	UniqueBody                    *ItemBody                           `json:"uniqueBody,omitempty"`
	WebLink                       string                              `json:"webLink,omitempty"`
	Attachments                   []Attachment                        `json:"attachments,omitempty"`
	Extensions                    []Extension                         `json:"extensions,omitempty"`
	MultiValueExtendedProperties  []MultiValueLegacyExtendedProperty  `json:"multiValueExtendedProperties,omitempty"`
	SingleValueExtendedProperties []SingleValueLegacyExtendedProperty `json:"singleValueExtendedProperties,omitempty"`
}
