package msgraph

import (
	"strings"
)

// Calendars represents an array of Calendar instances combined with some helper-functions
//
// See: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/resources/calendar
type Calendars []Calendar

func (c Calendars) String() string {
	var calendars = make([]string, len(c))
	for i, calendar := range c {
		calendars[i] = calendar.String()
	}
	return "Calendars(" + strings.Join(calendars, " | ") + ")"
}

// setGraphClient sets the graphClient instance in this instance and all child-instances (if any)
func (c Calendars) setGraphClient(gC *GraphClient) Calendars {
	for i := range c {
		c[i].setGraphClient(gC)
	}
	return c
}

// GetByName returns the calendar obj of that array whose DisplayName matches
// the given name. Returns an ErrFindCalendar if no calendar exists that matches the given
// name.
func (c Calendars) GetByName(name string) (Calendar, error) {
	for _, calendar := range c {
		if calendar.Name == name {
			return calendar, nil
		}
	}
	return Calendar{}, ErrFindCalendar
}
