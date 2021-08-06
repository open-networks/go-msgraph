package msgraph

import (
	"strings"
	"testing"
)

// GetTestListCalendars tries to get a valid User from the User-test cases
// and then uses ListCalendars() for this instance.
func GetTestListCalendars(t *testing.T) Calendars {
	t.Helper()
	testUser := GetTestUser(t)
	if skipCalendarTests {
		t.Skip("Skipping due to missing 'MSGraphExistingCalendarsOfUser' value")
	}
	testCalendars, err := testUser.ListCalendars()
	if err != nil {
		t.Fatalf("Cannot User.ListCalendars() for user %v: %v", testUser, err)
	}
	if len(testCalendars) == 0 {
		t.Fatalf("Cannot User.ListCalendars() for user %v, 0 calendars returned", testUser)
	}
	return testCalendars
}

func TestCalendars_GetByName(t *testing.T) {
	testCalendars := GetTestListCalendars(t)
	if skipCalendarTests {
		t.Skip("Skipping due to missing 'MSGraphExistingCalendarsOfUser' value")
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		c       Calendars
		args    args
		wantErr bool
	}{
		{
			name:    "Find valid Calendar",
			c:       testCalendars,
			args:    args{name: msGraphExistingCalendarsOfUser[0]},
			wantErr: false,
		}, {
			name:    "non-existing calendar",
			c:       testCalendars,
			args:    args{name: "DSfk9sdf89io23rasdfasfasdfasdf"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.c.GetByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calendars.GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCalendars_String(t *testing.T) {
	if skipCalendarTests {
		t.Skip("Skipping due to missing 'MSGraphExistingCalendarsOfUser' value")
	}
	testCalendars := GetTestListCalendars(t)

	// write custom string func
	var calendars = make([]string, len(testCalendars))
	for i, calendar := range testCalendars {
		calendars[i] = calendar.String()
	}
	want := "Calendars(" + strings.Join(calendars, " | ") + ")"

	t.Run("Test Calendars of TestUser", func(t *testing.T) {
		if got := testCalendars.String(); got != want {
			t.Errorf("Calendars.String() = %v, want %v", got, want)
		}
	})
}
