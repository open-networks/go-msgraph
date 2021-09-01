package msgraph

// MessageClassification is the importance of a user's message either inferred or overridden.
type MessageClassification string

var (
	MessageClassificationFocused MessageClassification = "focused"
	MessageClassificationOther   MessageClassification = "other"
)
