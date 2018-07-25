package msgraph

import (
	"encoding/json"
	"fmt"
	"time"
)

// ResponseStatus represents the response status for an Attendee to a CalendarEvent or just for a CalendarEvent
type ResponseStatus struct {
	Response string    // status of the response, may be organizer, accepted, declined etc.
	Time     time.Time // represents the time when the response was performed
}

func (s ResponseStatus) String() string {
	return fmt.Sprintf("Response: %s, Time: %s", s.Response, s.Time.Format(time.RFC3339Nano))
}

// Equal compares the ResponseStatus to the other Response status and returns true
// if the Response and time is equal
func (s ResponseStatus) Equal(other ResponseStatus) bool {
	return s.Response == other.Response && s.Time.Equal(other.Time)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (s *ResponseStatus) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Response  string `json:"response"`
		Timestamp string `json:"time"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	if tmp.Response == "" {
		return fmt.Errorf("Response-field is empty")
	}

	s.Response = tmp.Response
	s.Time, err = time.Parse(time.RFC3339Nano, tmp.Timestamp) // the timeZone is normally ALWAYS UTC, microsoft converts time date & time to that, but it does not matter here
	if err != nil {
		return fmt.Errorf("Can not parse timestamp with RFC3339Nano: %v", err)
	}

	return nil
}
