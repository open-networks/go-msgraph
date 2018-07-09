package msgraph

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// CalendarEvents represents multiple events of a Calendar. The amount of entries is determined by the timespan that is used to load the Calendar
type CalendarEvents struct {
	CalendarEvents []CalendarEvent `json:"value"`
}

func (c *CalendarEvents) String() string {
	return fmt.Sprintf("CalendarEvents(%v)", c.CalendarEvents)
}

// PrettySimpleString returns all Calendar Events in a readable format, mostly used for logging purposes
func (c *CalendarEvents) PrettySimpleString() string {
	var events = make([]string, len(c.CalendarEvents))
	for i, calendarEvent := range c.CalendarEvents {
		events[i] = fmt.Sprintf("{ %v (%v) [%v - %v] }", calendarEvent.Subject, calendarEvent.GetFirstAttendee().Name, calendarEvent.StartTime, calendarEvent.EndTime)
		//start.Format("02.01.2006 15:04:05"), end.Format("02.01.2006 15:04:05")))
	}
	return strings.Join(events, ", ")
}

// SortByStartDateTime sorts the array in this CalendarEvents instance
func (c *CalendarEvents) SortByStartDateTime() {
	sort.Slice(c.CalendarEvents, func(i, j int) bool { return c.CalendarEvents[i].StartTime.Before(c.CalendarEvents[j].StartTime) })
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

	c.CalendarEvents = tmp.CalendarEvents
	c.SortByStartDateTime()
	return nil
}
