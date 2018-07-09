package msgraph

import (
	"fmt"
	"strings"
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
	activePhone       string   // private cache for the active phone number
}

func (u *User) String() string {
	return fmt.Sprintf("User(ID: \"%v\", BusinessPhones: \"%v\", DisplayName: \"%v\", GivenName: \"%v\", "+
		"Mail: \"%v\", MobilePhone: \"%v\", PreferredLanguage: \"%v\", Surname: \"%v\", UserPrincipalName: \"%v\", "+
		"ActivePhone: \"%v\")",
		u.ID, u.BusinessPhones, u.DisplayName, u.GivenName, u.Mail, u.MobilePhone, u.PreferredLanguage, u.Surname,
		u.UserPrincipalName, u.activePhone)
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
func (u *User) GetShortName() string {
	supn := strings.Split(u.UserPrincipalName, "@")
	if len(supn) != 2 {
		return u.UserPrincipalName
	}
	return strings.ToUpper(supn[0])
}

// GetFullName returns the full name in that format: <firstname> <lastname>
func (u *User) GetFullName() string {
	return fmt.Sprintf("%v %v", u.GivenName, u.Surname)
}

// PrettySimpleString returns the User-instance simply formatted for logging purposes: {FullName (email) (activePhone)}
func (u *User) PrettySimpleString() string {
	return fmt.Sprintf("{ %v (%v) (%v) }", u.GetFullName(), u.Mail, u.GetActivePhone())
}
