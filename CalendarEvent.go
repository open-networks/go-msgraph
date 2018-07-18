package msgraph

import (
	"encoding/json"
	"fmt"
	"time"
)

// CalendarEvent represents a single event within a calendar
type CalendarEvent struct {
	ID                    string
	CreatedDateTime       time.Time      // Creation time of the CalendarEvent, has the correct timezone set from OriginalStartTimeZone (json)
	LastModifiedDateTime  time.Time      // Last modified time of the CalendarEvent, has the correct timezone set from OriginalEndTimeZone (json)
	OriginalStartTimeZone *time.Location // The original start-timezone, is already integrated in the calendartimes. Caution: is UTC on full day events
	OriginalEndTimeZone   *time.Location // The original end-timezone, is already integrated in the calendartimes. Caution: is UTC on full day events
	ICalUID               string
	Subject               string
	Importance            string
	Sensitivity           string
	IsAllDay              bool   // true = full day event, otherwise false
	IsCancelled           bool   // calendar event has been cancelled but is still in the calendar
	IsOrganizer           bool   // true if the calendar owner is the organizer
	SeriesMasterID        string // the ID of the master-entry of this series-event if any
	ShowAs                string
	Type                  string
	ResponseStatus        ResponseStatus // how the calendar-owner responded to the event (normally "organizer" because support-calendar is the host)
	StartTime             time.Time      // starttime of the Event, correct timezone is set
	EndTime               time.Time      // endtime of the event, correct timezone is set

	Attendees      []Attendee // represents all attendees to this CalendarEvent
	OrganizerName  string     // the name of the organizer from the e-mail, not reliable to identify anyone
	OrganizerEMail string     // the e-mail address of the organizer, use this to identify the user
}

// GetEventTimes returns the begin-time and the end-time of the Calendar Event
func (c CalendarEvent) GetEventTimes() (time.Time, time.Time) {
	return c.StartTime, c.EndTime
}

// GetFirstAttendee returns the first Attendee that is not the organizer of the event from the Attendees array.
// If none is found then an Attendee with the Name of "None" will be returned.
func (c CalendarEvent) GetFirstAttendee() Attendee {
	for _, attendee := range c.Attendees {
		if attendee.Email != c.OrganizerEMail {
			return attendee
		}
	}

	return Attendee{Name: "None"}
}

func (c *CalendarEvent) String() string {
	return fmt.Sprintf("CalendarEvent(ID: \"%v\", CreatedDateTime: \"%v\", LastModifiedDateTime: \"%v\", "+
		"ICalUId: \"%v\", Subject: \"%v\", "+
		"Importance: \"%v\", Sensitivity: \"%v\", IsAllDay: \"%v\", IsCancelled: \"%v\", "+
		"IsOrganizer: \"%v\", SeriesMasterId: \"%v\", ShowAs: \"%v\", Type: \"%v\", ResponseStatus: \"%v\", "+
		"Attendees: \"%v\", Organizer: \"%v\", Start: \"%v\", End: \"%v\")", c.ID, c.CreatedDateTime, c.LastModifiedDateTime,
		c.ICalUID, c.Subject, c.Importance,
		c.Sensitivity, c.IsAllDay, c.IsCancelled, c.IsOrganizer, c.SeriesMasterID, c.ShowAs,
		c.Type, c.ResponseStatus, c.Attendees, c.OrganizerName+" "+c.OrganizerEMail, c.StartTime, c.EndTime)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (c *CalendarEvent) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID                    string            `json:"id"`
		CreatedDateTime       string            `json:"createdDateTime"`
		LastModifiedDateTime  string            `json:"lastModifiedDateTime"`
		OriginalStartTimeZone string            `json:"originalStartTimeZone"`
		OriginalEndTimeZone   string            `json:"originalEndTimeZone"`
		ICalUID               string            `json:"iCalUId"`
		Subject               string            `json:"subject"`
		Importance            string            `json:"importance"`
		Sensitivity           string            `json:"sensitivity"`
		IsAllDay              bool              `json:"isAllDay"`
		IsCancelled           bool              `json:"isCancelled"`
		IsOrganizer           bool              `json:"isOrganizer"`
		SeriesMasterID        string            `json:"seriesMasterId"`
		ShowAs                string            `json:"showAs"`
		Type                  string            `json:"type"`
		ResponseStatus        ResponseStatus    `json:"responseStatus"`
		Start                 map[string]string `json:"start"`
		End                   map[string]string `json:"end"`
		Attendees             []Attendee        `json:"attendees"`
		Organizer             struct {
			EmailAddress struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			} `json:"emailAddress"`
		} `json:"organizer"`
	}{}

	var err error
	// unmarshal to tmp-struct, return if error
	if err = json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("Error on json.Unmarshal: %v | Data: %v", err, string(data))
	}

	c.ID = tmp.ID
	c.CreatedDateTime, err = time.Parse(time.RFC3339Nano, tmp.CreatedDateTime)
	if err != nil {
		return fmt.Errorf("Can not time.Parse with RFC3339Nano createdDateTime %v: %v", tmp.CreatedDateTime, err)
	}
	c.LastModifiedDateTime, err = time.Parse(time.RFC3339Nano, tmp.LastModifiedDateTime)
	if err != nil {
		return fmt.Errorf("Can not time.Parse with RFC3339Nano lastModifiedDateTime %v: %v", tmp.LastModifiedDateTime, err)
	}
	c.OriginalStartTimeZone, err = time.LoadLocation(mapTimeZoneStrings(tmp.OriginalStartTimeZone))
	if err != nil {
		return fmt.Errorf("Can not time.LoadLocation originalStartTimeZone %v: %v", tmp.OriginalStartTimeZone, err)
	}
	c.OriginalEndTimeZone, err = time.LoadLocation(mapTimeZoneStrings(tmp.OriginalEndTimeZone))
	if err != nil {
		return fmt.Errorf("Can not time.LoadLocation originalEndTimeZone %v: %v", tmp.OriginalEndTimeZone, err)
	}
	c.ICalUID = tmp.ICalUID
	c.Subject = tmp.Subject
	c.Importance = tmp.Importance
	c.Sensitivity = tmp.Sensitivity
	c.IsAllDay = tmp.IsAllDay
	c.IsCancelled = tmp.IsCancelled
	c.IsOrganizer = tmp.IsOrganizer
	c.SeriesMasterID = tmp.SeriesMasterID
	c.ShowAs = tmp.ShowAs
	c.Type = tmp.Type
	c.ResponseStatus = tmp.ResponseStatus
	c.Attendees = tmp.Attendees
	c.OrganizerName = tmp.Organizer.EmailAddress.Name
	c.OrganizerEMail = tmp.Organizer.EmailAddress.Address

	// Parse event start & endtime with timezone
	c.StartTime, err = parseTimeAndLocation(tmp.Start["dateTime"], tmp.Start["timeZone"]) // the timeZone is normally ALWAYS UTC, microsoft converts time date & time to that
	if err != nil {
		return fmt.Errorf("Can not parse start-dateTime %v AND timeZone %v: %v", tmp.Start["dateTime"], tmp.Start["timeZone"], err)
	}
	c.EndTime, err = parseTimeAndLocation(tmp.End["dateTime"], tmp.End["timeZone"]) // the timeZone is normally ALWAYS UTC, microsoft converts time date & time to that
	if err != nil {
		return fmt.Errorf("Can not parse end-dateTime %v AND timeZone %v: %v", tmp.End["dateTime"], tmp.End["timeZone"], err)
	}

	// Hint: OriginalStartTimeZone & end are UTC (set by microsoft) if it is a full-day event, this will be handled in the next section
	c.StartTime = c.StartTime.In(c.OriginalStartTimeZone) // move the StartTime to the orignal start-timezone
	c.EndTime = c.EndTime.In(c.OriginalEndTimeZone)       // move the EndTime to the orignal end-timezone

	// Now check if it's a full-day event, if yes, the event is UTC anyway. We need it to be accurate for the program to work
	// hence we set it to time.Local. It can later be manipulated by the program to a different timezone but the times also have
	// to be recalculated. E.g. we set it to UTC+2 hence it will start at 02:00 and end at 02:00, not 00:00 -> manually set to 00:00
	if c.IsAllDay && FullDayEventTimeZone != time.UTC {
		// set to local location
		c.StartTime = c.StartTime.In(FullDayEventTimeZone)
		c.EndTime = c.EndTime.In(FullDayEventTimeZone)
		// get offset in seconds
		_, startOffSet := c.StartTime.Zone()
		_, endOffSet := c.EndTime.Zone()
		// decrease time to 00:00 again
		c.StartTime = c.StartTime.Add(-1 * time.Second * time.Duration(startOffSet))
		c.EndTime = c.EndTime.Add(-1 * time.Second * time.Duration(endOffSet))
	}

	return nil
}

// parseTimeAndLocation is just a helper method to shorten the code in the Unmarshal json
func parseTimeAndLocation(timeToParse, locationToParse string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02T15:04:05.999999999", timeToParse)
	if err != nil {
		return time.Time{}, err
	}
	parsedTimeZone, err := time.LoadLocation(locationToParse)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime.In(parsedTimeZone), nil
}

// mapTimeZoneStrings maps various Timezones used by Microsoft to go-understandable timezones or returns the source-zone if no mapping is found
func mapTimeZoneStrings(timeZone string) string {
	switch timeZone {
	case "W. Europe Standard Time":
		return "Europe/Vienna"
	case "Central Europe Standard Time":
		return "CET"
	case "tzone://Microsoft/Utc":
		return "UTC"
	case "India Standard Time":
		return "IST"
	}
	// TODO: Make this list bigger / improve this list with more synonyms or make a more dynamic workaround
	// Hint: probably this could list all timezones and even enable to automatically determine the correct timezone: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/outlookuser_supportedtimezones
	return timeZone
}
