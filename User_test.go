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
		want: fmt.Sprintf("User(ID: \"%v\", BusinessPhones: \"%v\", DisplayName: \"%v\", GivenName: \"%v\", "+
			"Mail: \"%v\", MobilePhone: \"%v\", PreferredLanguage: \"%v\", Surname: \"%v\", UserPrincipalName: \"%v\", "+
			"ActivePhone: \"%v\", DirectAPIConnection: %v)",
			u.ID, u.BusinessPhones, u.DisplayName, u.GivenName, u.Mail, u.MobilePhone, u.PreferredLanguage, u.Surname,
			u.UserPrincipalName, u.activePhone, u.graphClient != nil),
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := tt.u.String(); got != tt.want {
			t.Errorf("User.String() = %v, want %v", got, tt.want)
		}
	})
}

func TestUser_UpdateUser(t *testing.T) {
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
