package msgraph

// FollowupFlag are flags on a message that display how it should be followed up on.
//
// See https://docs.microsoft.com/en-us/graph/api/resources/FollowupFlag
type FollowupFlag struct {
	CompletedDateTime DateTimeTimeZone `json:"completedDateTime,omitempty"`
	DueDateTime       DateTimeTimeZone `json:"dueDateTime,omitempty"`
	FlagStatus        string           `json:"flagStatus,omitempty"`
	StartDateTime     DateTimeTimeZone `json:"startDateTime,omitempty"`
}
