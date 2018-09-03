package msgraph

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// User represents a user from the ms graph API
type User struct {
	ID                string   `json:"id"`
	BusinessPhones    []string `json:"businessPhones"`
	DisplayName       string   `json:"displayName"`
	GivenName         string   `json:"givenName"`
	Mail              string   `json:"mail"`
	MobilePhone       string   `json:"mobilePhone"`
	PreferredLanguage string   `json:"preferredLanguage"`
	Surname           string   `json:"surname"`
	UserPrincipalName string   `json:"userPrincipalName"`

	activePhone string       // private cache for the active phone number
	graphClient *GraphClient // the graphClient that called the user
}

func (u *User) String() string {
	return fmt.Sprintf("User(ID: \"%v\", BusinessPhones: \"%v\", DisplayName: \"%v\", GivenName: \"%v\", "+
		"Mail: \"%v\", MobilePhone: \"%v\", PreferredLanguage: \"%v\", Surname: \"%v\", UserPrincipalName: \"%v\", "+
		"ActivePhone: \"%v\", DirectAPIConnection: %v)",
		u.ID, u.BusinessPhones, u.DisplayName, u.GivenName, u.Mail, u.MobilePhone, u.PreferredLanguage, u.Surname,
		u.UserPrincipalName, u.activePhone, u.graphClient != nil)
}

// setGraphClient sets the graphClient instance in this instance and all child-instances (if any)
func (u *User) setGraphClient(gC *GraphClient) {
	u.graphClient = gC
}

// ListCalendars returns all calendars associated to that user.
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user_list_calendars
func (u User) ListCalendars() (Calendars, error) {
	if u.graphClient == nil {
		return Calendars{}, ErrNotGraphClientSourced
	}
	resource := fmt.Sprintf("/users/%v/calendars", u.ID)

	var marsh struct {
		Calendars Calendars `json:"value"`
	}
	err := u.graphClient.makeGETAPICall(resource, nil, &marsh)
	marsh.Calendars.setGraphClient(u.graphClient)
	return marsh.Calendars, err
}

// ListCalendarView returns the CalendarEvents of the given user within the specified
// start- and endDateTime. The calendar used is the default calendar of the user.
// Returns an error if the user it not GraphClient sourced or if there is any error
// during the API-call.
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user_list_calendarview
func (u User) ListCalendarView(startDateTime, endDateTime time.Time) (CalendarEvents, error) {
	if u.graphClient == nil {
		return CalendarEvents{}, ErrNotGraphClientSourced
	}

	if len(globalSupportedTimeZones.Value) == 0 {
		var err error
		globalSupportedTimeZones, err = u.getTimeZoneChoices()
		if err != nil {
			return CalendarEvents{}, err
		}
	}

	resource := fmt.Sprintf("/users/%v/calendar/calendarview", u.ID)

	// set GET-Params for start and end time
	getParams := url.Values{}
	getParams.Add("startdatetime", startDateTime.Format("2006-01-02T00:00:00"))
	getParams.Add("enddatetime", endDateTime.Format("2006-01-02T00:00:00"))

	var calendarEvents CalendarEvents
	return calendarEvents, u.graphClient.makeGETAPICall(resource, getParams, &calendarEvents)
}

// getTimeZoneChoices grabs all supported time zones from microsoft for this user.
// This should actually be the same for every user. Only used internally by this
// msgraph package.
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/outlookuser_supportedtimezones
func (u User) getTimeZoneChoices() (supportedTimeZones, error) {
	var ret supportedTimeZones
	err := u.graphClient.makeGETAPICall(fmt.Sprintf("/users/%s/outlook/supportedTimeZones", u.ID), nil, &ret)
	return ret, err
}

// GetActivePhone returns the space-trimmed active phone-number of the user. The active
// phone number is either the MobilePhone number or the first business-Phone number
func (u *User) GetActivePhone() string {
	if u.activePhone != "" { // use cached value if any
		return u.activePhone
	}
	// no cached active phone number, evaluate & cache it now:
	u.activePhone = strings.Replace(u.MobilePhone, " ", "", -1)
	if u.activePhone == "" && len(u.BusinessPhones) > 0 {
		u.activePhone = strings.Replace(u.BusinessPhones[0], " ", "", -1)
	}
	return u.activePhone
}

// GetShortName returns the first part of UserPrincipalName before the @. If there
// is no @, then just the UserPrincipalName will be returned
func (u User) GetShortName() string {
	supn := strings.Split(u.UserPrincipalName, "@")
	if len(supn) != 2 {
		return u.UserPrincipalName
	}
	return strings.ToUpper(supn[0])
}

// GetFullName returns the full name in that format: <firstname> <lastname>
func (u User) GetFullName() string {
	return fmt.Sprintf("%v %v", u.GivenName, u.Surname)
}

// PrettySimpleString returns the User-instance simply formatted for logging purposes: {FullName (email) (activePhone)}
func (u User) PrettySimpleString() string {
	return fmt.Sprintf("{ %v (%v) (%v) }", u.GetFullName(), u.Mail, u.GetActivePhone())
}

// Equal returns wether the user equals the other User by comparing every property
// of the user including the ID
func (u User) Equal(other User) bool {
	var equalBool = true
	for i := 0; i < len(u.BusinessPhones) && i < len(other.BusinessPhones); i++ {
		equalBool = equalBool && u.BusinessPhones[i] == other.BusinessPhones[i]
	}
	equalBool = equalBool && len(u.BusinessPhones) == len(other.BusinessPhones)
	return equalBool && u.ID == other.ID && u.DisplayName == other.DisplayName && u.GivenName == other.GivenName &&
		u.Mail == other.Mail && u.MobilePhone == other.MobilePhone && u.PreferredLanguage == other.PreferredLanguage &&
		u.Surname == other.Surname && u.UserPrincipalName == other.UserPrincipalName
}
