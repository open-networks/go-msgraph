package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// User represents a user from the ms graph API
type User struct {
	ID                string            `json:"id,omitempty"`
	BusinessPhones    []string          `json:"businessPhones,omitempty"`
	DisplayName       string            `json:"displayName,omitempty"`
	GivenName         string            `json:"givenName,omitempty"`
	JobTitle          string            `json:"jobTitle,omitempty"`
	Mail              string            `json:"mail,omitempty"`
	MobilePhone       string            `json:"mobilePhone,omitempty"`
	PreferredLanguage string            `json:"preferredLanguage,omitempty"`
	Surname           string            `json:"surname,omitempty"`
	UserPrincipalName string            `json:"userPrincipalName,omitempty"`
	AccountEnabled    bool              `json:"accountEnabled,omitempty"`
	AssignedLicenses  []AssignedLicense `json:"assignedLicenses,omitempty"`
	CompanyName       string            `json:"companyName,omitempty"`
	Department        string            `json:"department,omitempty"`
	MailNickname      string            `json:"mailNickname,omitempty"`
	PasswordProfile   PasswordProfile   `json:"passwordProfile,omitempty"`

	activePhone string       // private cache for the active phone number
	graphClient *GraphClient // the graphClient that called the user
}

type AssignedLicense struct {
	DisabledPlans []string `json:"disabledPlans,omitempty"`
	SkuID         string   `json:"skuId,omitempty"`
}

type PasswordProfile struct {
	ForceChangePasswordNextSignIn        bool   `json:"forceChangePasswordNextSignIn,omitempty"`
	ForceChangePasswordNextSignInWithMfa bool   `json:"forceChangePasswordNextSignInWithMfa,omitempty"`
	Password                             string `json:"password,omitempty"`
}

func (u User) String() string {
	return fmt.Sprintf("User(ID: \"%v\", BusinessPhones: \"%v\", DisplayName: \"%v\", GivenName: \"%v\", "+
		"JobTitle: \"%v\", Mail: \"%v\", MobilePhone: \"%v\", PreferredLanguage: \"%v\", Surname: \"%v\", "+
		"UserPrincipalName: \"%v\", ActivePhone: \"%v\", DirectAPIConnection: %v)",
		u.ID, u.BusinessPhones, u.DisplayName, u.GivenName, u.JobTitle, u.Mail, u.MobilePhone, u.PreferredLanguage, u.Surname,
		u.UserPrincipalName, u.activePhone, u.graphClient != nil)
}

// setGraphClient sets the graphClient instance in this instance and all child-instances (if any)
func (u *User) setGraphClient(gC *GraphClient) {
	u.graphClient = gC
}

// ListCalendars returns all calendars associated to that user.
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user_list_calendars
func (u User) ListCalendars(opts ...ListQueryOption) (Calendars, error) {
	if u.graphClient == nil {
		return Calendars{}, ErrNotGraphClientSourced
	}
	resource := fmt.Sprintf("/users/%v/calendars", u.ID)

	var marsh struct {
		Calendars Calendars `json:"value"`
	}
	err := u.graphClient.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	marsh.Calendars.setGraphClient(u.graphClient)
	return marsh.Calendars, err
}

// ListCalendarView returns the CalendarEvents of the given user within the specified
// start- and endDateTime. The calendar used is the default calendar of the user.
// Returns an error if the user it not GraphClient sourced or if there is any error
// during the API-call.
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user_list_calendarview
func (u User) ListCalendarView(startDateTime, endDateTime time.Time, opts ...ListQueryOption) (CalendarEvents, error) {
	if u.graphClient == nil {
		return CalendarEvents{}, ErrNotGraphClientSourced
	}

	if len(globalSupportedTimeZones.Value) == 0 {
		var err error
		// TODO: this is a dirty fix, because opts could contain other things than a context, e.g. select
		// parameters. This could produce unexpected outputs and therefore break the globalSupportedTimeZones variable.
		globalSupportedTimeZones, err = u.getTimeZoneChoices(compileListQueryOptions(opts))
		if err != nil {
			return CalendarEvents{}, err
		}
	}

	resource := fmt.Sprintf("/users/%v/calendar/calendarview", u.ID)

	// set GET-Params for start and end time
	var reqOpt = compileListQueryOptions(opts)
	reqOpt.queryValues.Add("startdatetime", startDateTime.Format("2006-01-02T00:00:00"))
	reqOpt.queryValues.Add("enddatetime", endDateTime.Format("2006-01-02T00:00:00"))

	var calendarEvents CalendarEvents
	return calendarEvents, u.graphClient.makeGETAPICall(resource, reqOpt, &calendarEvents)
}

// getTimeZoneChoices grabs all supported time zones from microsoft for this user.
// This should actually be the same for every user. Only used internally by this
// msgraph package.
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/outlookuser_supportedtimezones
func (u User) getTimeZoneChoices(opts getRequestParams) (supportedTimeZones, error) {
	var ret supportedTimeZones
	err := u.graphClient.makeGETAPICall(fmt.Sprintf("/users/%s/outlook/supportedTimeZones", u.ID), opts, &ret)
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
	return supn[0]
}

// GetFullName returns the full name in that format: <firstname> <lastname>
func (u User) GetFullName() string {
	return fmt.Sprintf("%v %v", u.GivenName, u.Surname)
}

// GetMemberGroupsAsStrings returns a list of all group IDs the user is a member of.
// You can specify the securityGroupsEnabeled parameter to only return security group IDs.
//
// opts ...GetQueryOption - only msgraph.GetWithContext is supported.
//
// Reference: https://docs.microsoft.com/en-us/graph/api/directoryobject-getmembergroups?view=graph-rest-1.0&tabs=http
func (u User) GetMemberGroupsAsStrings(securityGroupsEnabeled bool, opts ...GetQueryOption) ([]string, error) {
	if u.graphClient == nil {
		return nil, ErrNotGraphClientSourced
	}
	return u.graphClient.getMemberGroups(u.ID, securityGroupsEnabeled, opts...)
}

// PrettySimpleString returns the User-instance simply formatted for logging purposes: {FullName (email) (activePhone)}
func (u User) PrettySimpleString() string {
	return fmt.Sprintf("{ %v (%v) (%v) }", u.GetFullName(), u.Mail, u.GetActivePhone())
}

// UpdateUser patches this user object. Note, only set the fields that should be changed.
//
// IMPORTANT: the user cannot be disabled (field AccountEnabled) this way, because the
// default value of a boolean is false - and hence will not be posted via json - omitempty
// is used. user func user.DisableAccount() instead.
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user-update
func (u User) UpdateUser(userInput User, opts ...UpdateQueryOption) error {
	if u.graphClient == nil {
		return ErrNotGraphClientSourced
	}
	resource := fmt.Sprintf("/users/%v", u.ID)

	bodyBytes, err := json.Marshal(userInput)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(bodyBytes)
	// Hint: API-call body does not return any data / no json object.
	err = u.graphClient.makePATCHAPICall(resource, compileUpdateQueryOptions(opts), reader, nil)
	return err
}

// DisableAccount disables the User-Account, hence sets the AccountEnabled-field to false.
// This function must be used instead of user.UpdateUser, because the AccountEnabled-field
// with json "omitempty" will never be sent when false. Without omitempty, the user account would
// always accidentially disabled upon an update of e.g. only "DisplayName"
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user-update
func (u User) DisableAccount(opts ...UpdateQueryOption) error {
	if u.graphClient == nil {
		return ErrNotGraphClientSourced
	}
	resource := fmt.Sprintf("/users/%v", u.ID)

	bodyBytes, err := json.Marshal(struct {
		AccountEnabled bool `json:"accountEnabled"`
	}{AccountEnabled: false})
	if err != nil {
		return err
	}

	reader := bytes.NewReader(bodyBytes)
	// Hint: API-call body does not return any data / no json object.
	err = u.graphClient.makePATCHAPICall(resource, compileUpdateQueryOptions(opts), reader, nil)
	return err
}

// DeleteUser deletes this user instance at the Microsoft Azure AD. Use with caution.
//
// Reference: https://docs.microsoft.com/en-us/graph/api/user-delete
func (u User) DeleteUser(opts ...DeleteQueryOption) error {
	if u.graphClient == nil {
		return ErrNotGraphClientSourced
	}
	resource := fmt.Sprintf("/users/%v", u.ID)

	// TODO: check return body, maybe there is some potential success or error message hidden in it?
	err := u.graphClient.makeDELETEAPICall(resource, compileDeleteQueryOptions(opts), nil)
	return err
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
		u.JobTitle == other.JobTitle && u.Mail == other.Mail && u.MobilePhone == other.MobilePhone &&
		u.PreferredLanguage == other.PreferredLanguage && u.Surname == other.Surname && u.UserPrincipalName == other.UserPrincipalName
}
