package msgraph

import "time"

// DateTimeTimeZone represents the date, time, and timezone for a given time.
//
// See https://docs.microsoft.com/en-us/graph/api/resources/DateTimeTimeZone?view=graph-rest-1.0
type DateTimeTimeZone struct {
	DateTime time.Time      `json:"dateTime,omitempty"`
	TimeZone *time.Location `json:"timeZone,omitempty"`
}
