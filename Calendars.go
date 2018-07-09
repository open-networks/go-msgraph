package msgraph

import (
	"strings"
)

// Calendars represents an array of Calendar instances
type Calendars []Calendar

func (c Calendars) String() string {
	var calendars = make([]string, len(c))
	for i, calendar := range c {
		calendars[i] = calendar.String()
	}
	return "Calendars(" + strings.Join(calendars, " | ") + ")"
}
