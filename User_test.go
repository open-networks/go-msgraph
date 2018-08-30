package msgraph

import (
	"testing"
	"time"
)

// GetTestUser returns a valid User instance for testing. Will issue a t.Fatalf if the
// user can not be loaded
func GetTestUser(t *testing.T) User {
	TestEnvironmentVariablesPresent(t)
	userToTest, errUserToTest := graphClient.GetUser(msGraphExistingUserPrincipalInGroup)
	if errUserToTest != nil {
		t.Fatalf("Can not find user %v for Testing", msGraphExistingUserPrincipalInGroup)
	}
	return userToTest
}

func TestUser_ListCalendarView(t *testing.T) {
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

func TestUser_ListCalendars(t *testing.T) {
	userToTest := GetTestUser(t)

	tests := []struct {
		name    string
		u       User
		want    Calendars
		wantErr bool
	}{
		{
			name:    "List All Calendars",
			u:       userToTest,
			want:    Calendars{Calendar{Name: "Kalender"}, Calendar{Name: "Feiertage in Österreich"}, Calendar{Name: "Geburtstage"}},
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
