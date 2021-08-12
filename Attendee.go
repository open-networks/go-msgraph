package msgraph

import (
	"encoding/json"
	"fmt"
)

// Attendee struct represents an attendee for a CalendarEvent
type Attendee struct {
	Type           string         // the type of the invitation, e.g. required, optional etc.
	Name           string         // the name of the person, comes from the E-Mail Address - hence not a reliable name to search for
	Email          string         // the e-mail address of the person - use this to identify the user
	ResponseStatus ResponseStatus // the ResponseStatus for that particular Attendee for the CalendarEvent
}

func (a Attendee) String() string {
	return fmt.Sprintf("Name: %s, Type: %s, E-mail: %s, ResponseStatus: %v", a.Name, a.Type, a.Email, a.ResponseStatus)
}

// Equal compares the Attendee to the other Attendee and returns true
// if the two given Attendees are equal. Otherwise returns false
func (a Attendee) Equal(other Attendee) bool {
	return a.Type == other.Type && a.Name == other.Name && a.Email == other.Email && a.ResponseStatus.Equal(other.ResponseStatus)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (a *Attendee) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type         string         `json:"type"`
		Status       ResponseStatus `json:"status"`
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return fmt.Errorf("Attendee: %v", err.Error())
	}

	a.Type = tmp.Type
	a.Name = tmp.EmailAddress.Name
	a.Email = tmp.EmailAddress.Address
	a.ResponseStatus = tmp.Status

	return nil
}
