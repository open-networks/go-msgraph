package msgraph

import (
	"fmt"
	"testing"
	"time"
)

// GetTestUser returns a valid User instance for testing. Will issue a t.Fatalf if the
// user cannot be loaded
func GetTestUser(t *testing.T) User {
	t.Helper()
	userToTest, errUserToTest := graphClient.GetUser(msGraphExistingUserPrincipalInGroup)
	if errUserToTest != nil {
		t.Fatalf("Cannot find user %v for Testing", msGraphExistingUserPrincipalInGroup)
	}
	return userToTest
}

func TestUser_getTimeZoneChoices(t *testing.T) {
	userToTest := GetTestUser(t)
	if skipCalendarTests {
		t.Skip("Skipping due to missing 'MSGraphExistingCalendarsOfUser' value")
	}

	tests := []struct {
		name    string
		u       User
		want    []supportedTimeZones
		wantErr bool
	}{
		{
			name:    "Test all",
			u:       userToTest,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.u.getTimeZoneChoices(compileGetQueryOptions([]GetQueryOption{}))
			if (err != nil) != tt.wantErr {
				t.Errorf("User.getTimeZoneChoices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//t.Errorf("Printing for sure")
		})
	}
}

func TestUser_ListCalendars(t *testing.T) {
	if skipCalendarTests {
		t.Skip("Skipping due to missing 'MSGraphExistingCalendarsOfUser' value")
	}
	// testing for ErrNotGraphClientSourced
	notGraphClientSourcedUser := User{ID: "none"}
	_, err := notGraphClientSourcedUser.ListCalendars()
	if err != ErrNotGraphClientSourced {
		t.Errorf("Expected error \"ErrNotGraphClientSourced\", but got: %v", err)
	}
	// continue with normal tests
	userToTest := GetTestUser(t)

	var wantedCalendars []Calendar
	for _, calendarName := range msGraphExistingCalendarsOfUser {
		wantedCalendars = append(wantedCalendars, Calendar{Name: calendarName})
	}

	tests := []struct {
		name    string
		u       User
		want    Calendars
		wantErr bool
	}{
		{
			name:    "List All Calendars",
			u:       userToTest,
			want:    wantedCalendars,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.ListCalendars()
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.ListUserCalendars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		Outer:
			for _, searchCalendar := range tt.want {
				for _, toCompare := range got {
					if searchCalendar.Name == toCompare.Name {
						continue Outer
					}
				}
				t.Errorf("GraphClient.ListUserCalendars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_ListCalendarView(t *testing.T) {
	if skipCalendarTests {
		t.Skip("Skipping due to missing 'MSGraphExistingCalendarsOfUser' value")
	}
	userToTest := GetTestUser(t)

	type args struct {
		startdate time.Time
		enddate   time.Time
	}
	tests := []struct {
		name    string
		u       User
		args    args
		wantErr bool
	}{
		{
			name:    "Existing User",
			u:       userToTest,
			args:    args{startdate: time.Now(), enddate: time.Now().Add(7 * 24 * time.Hour)},
			wantErr: false,
		}, {
			name:    "User not initialized by GraphClient",
			u:       User{UserPrincipalName: "testuser"},
			args:    args{startdate: time.Now(), enddate: time.Now().Add(7 * 24 * time.Hour)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.ListCalendarView(tt.args.startdate, tt.args.enddate)
			//fmt.Println("Got User.ListCalendarview(): ", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.ListCalendarView() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(got) == 0 {
				t.Errorf("User.ListCalendarView() len is 0, but want > 0, Got: %v", got)
			}
		})
	}
}

func TestUser_String(t *testing.T) {
	u := GetTestUser(t)
	tt := struct {
		name string
		u    *User
		want string
	}{
		name: "Test user func String",
		u:    &u,
		want: fmt.Sprintf("User(ID: \"%v\", BusinessPhones: \"%v\", DisplayName: \"%v\", GivenName: \"%v\", JobTitle: \"%v\", "+
			"Mail: \"%v\", MobilePhone: \"%v\", PreferredLanguage: \"%v\", Surname: \"%v\", UserPrincipalName: \"%v\", "+
			"ActivePhone: \"%v\", DirectAPIConnection: %v)",
			u.ID, u.BusinessPhones, u.DisplayName, u.GivenName, u.JobTitle, u.Mail, u.MobilePhone, u.PreferredLanguage, u.Surname,
			u.UserPrincipalName, u.activePhone, u.graphClient != nil),
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := tt.u.String(); got != tt.want {
			t.Errorf("User.String() = %v, want %v", got, tt.want)
		}
	})
}

func TestUser_GetMemberGroupsAsStrings(t *testing.T) {
	u := GetTestUser(t)
	tests := []struct {
		name                   string
		u                      User
		securityGroupsEnabeled bool
		opt                    GetQueryOption
		wantErr                bool
	}{
		{
			name:                   "Test user func GetMembershipGroupsAsStrings",
			u:                      u,
			securityGroupsEnabeled: true,
			opt:                    GetWithContext(nil),
			wantErr:                false,
		}, {
			name:                   "Test user func GetMembershipGroupsAsStrings - no securityGroupsEnabeledF",
			u:                      u,
			securityGroupsEnabeled: false,
			opt:                    GetWithContext(nil),
			wantErr:                false,
		},
		{
			name:                   "User not initialized by GraphClient",
			securityGroupsEnabeled: true,
			opt:                    GetWithContext(nil),
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.GetMemberGroupsAsStrings(tt.securityGroupsEnabeled, tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.GetMemberGroupsAsStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("User.GetMemberGroupsAsStrings() = %v, len(%d), want at least one value", got, len(got))
			}
		})
	}
}

func TestUser_UpdateUser(t *testing.T) {
	// testing for ErrNotGraphClientSourced
	notGraphClientSourcedUser := User{ID: "none"}
	err := notGraphClientSourcedUser.UpdateUser(User{})
	if err != ErrNotGraphClientSourced {
		t.Errorf("Expected error \"ErrNotGraphClientSourced\", but got: %v", err)
	}

	// continue with normal tests
	testuser := createUnitTestUser(t)

	targetedCompanyName := "go-msgraph unit test suite UpdateUser" + randomString(25)
	testuser.UpdateUser(User{CompanyName: targetedCompanyName})
	getUser, err := graphClient.GetUser(testuser.ID, GetWithSelect("id,userPrincipalName,displayName,givenName,companyName"))
	if err != nil {
		testuser.DeleteUser()
		t.Errorf("Cannot perform User.UpdateUser, error: %v", err)
	}
	if getUser.CompanyName != targetedCompanyName {
		testuser.DeleteUser()
		t.Errorf("Performed User.UpdateUser, but CompanyName is still \"%v\" instead of wanted \"%v\"", getUser.CompanyName, targetedCompanyName)
	}
	err = testuser.DeleteUser()
	if err != nil {
		t.Errorf("Could not User.DeleteUser() (for %v) after User.UpdateUser tests: %v", testuser, err)
	}
}

func TestUser_GetShortName(t *testing.T) {
	// test a normal case
	testuser := User{UserPrincipalName: "dumpty@contoso.com"}
	if testuser.GetShortName() != "dumpty" {
		t.Errorf("user.GetShortName() should return \"dumpty\", but returns: %v", testuser.GetShortName())
	}
	// test a case that actually should never happen... but we all know murphy
	testuser = User{UserPrincipalName: "alice"}
	if testuser.GetShortName() != "alice" {
		t.Errorf("user.GetShortName() should return \"alice\", but returns: %v", testuser.GetShortName())
	}
}

func TestUser_GetFullName(t *testing.T) {
	testuser := User{GivenName: "Bob", Surname: "Rabbit"}
	wanted := fmt.Sprintf("%v %v", testuser.GivenName, testuser.Surname)
	if testuser.GetFullName() != wanted {
		t.Errorf("user.GetFullName() should return \"%v\", but returns: \"%v\"", wanted, testuser.GetFullName())
	}
}

func TestUser_PrettySimpleString(t *testing.T) {
	testuser := User{GivenName: "Bob", Surname: "Rabbit", Mail: "bob.rabbit@contoso.com", MobilePhone: "+1 23456789"}
	wanted := fmt.Sprintf("{ %v (%v) (%v) }", testuser.GetFullName(), testuser.Mail, testuser.GetActivePhone())
	if testuser.PrettySimpleString() != wanted {
		t.Errorf("user.GetFullName() should return \"%v\", but returns: \"%v\"", wanted, testuser.PrettySimpleString())
	}
}
