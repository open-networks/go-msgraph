package msgraph

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// CalendarEvents represents multiple events of a Calendar. The amount of entries is determined by the timespan that is used to load the Calendar
type CalendarEvents []CalendarEvent

func (c CalendarEvents) String() string {
	var events = make([]string, len(c))
	for i, calendarEvent := range c {
		events[i] = calendarEvent.String()
	}
	return fmt.Sprintf("CalendarEvents(%v)", strings.Join(events, ", "))
}

// PrettySimpleString returns all Calendar Events in a readable format, mostly used for logging purposes
func (c CalendarEvents) PrettySimpleString() string {
	var events = make([]string, len(c))
	for i, calendarEvent := range c {
		events[i] = fmt.Sprintf("{ %v (%v) [%v - %v] }", calendarEvent.Subject, calendarEvent.GetFirstAttendee().Name, calendarEvent.StartTime, calendarEvent.EndTime)
	}
	return fmt.Sprintf("CalendarEvents(%v)", strings.Join(events, ", "))
}

// SortByStartDateTime sorts the array in this CalendarEvents instance
func (c CalendarEvents) SortByStartDateTime() {
	sort.Slice(c, func(i, j int) bool { return c[i].StartTime.Before(c[j].StartTime) })
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library. The only
// purpose of this overwrite is to immediately sort the []CalendarEvent by StartDateTime
func (c *CalendarEvents) UnmarshalJSON(data []byte) error {
	tmp := struct {
		CalendarEvents []CalendarEvent `json:"value"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return fmt.Errorf("Can not UnmarshalJSON: %v | Data: %v", err, string(data))
	}

	*c = tmp.CalendarEvents // re-assign the

	//c.CalendarEvents = tmp.CalendarEvents
	c.SortByStartDateTime()
	return nil
}
