package msgraph

import (
	"strings"
)

// Attendees struct represents multiple Attendees for a CalendarEvent
type Attendees []Attendee

func (a Attendees) String() string {
	var attendees = make([]string, len(a))
	for i, attendee := range a {
		attendees[i] = attendee.String()
	}
	return "Attendees(" + strings.Join(attendees, " | ") + ")"
}

// Equal compares the Attendee to the other Attendee and returns true
// if the two given Attendees are equal. Otherwise returns false
func (a Attendees) Equal(other Attendees) bool {
Outer:
	for _, attendee := range a {
		for _, toCompare := range other {
			if attendee.Equal(toCompare) {
				continue Outer
			}
		}
		return false
	}
	return len(a) == len(other) // if we reach this, all attendees have been found, now return if the attendees are equally much
}
