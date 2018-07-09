package msgraph

import (
	"strings"
)

// Calendars represents an array of Calendar instances
type Calendars struct {
	Calendars []Calendar `json:"value"` // A list of Calendars
}

func (c Calendars) String() string {
	var calendars = make([]string, len(c.Calendars))
	for i, calendar := range c.Calendars {
		calendars[i] = calendar.String()
	}
	return "Calendars(" + strings.Join(calendars, " | ") + ")"
}
