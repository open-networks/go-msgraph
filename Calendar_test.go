package msgraph

import (
	"fmt"
	"testing"
)

func TestCalendar_String(t *testing.T) {
	testCalendars := GetTestListCalendars(t)
	if skipCalendarTests {
		t.Skip("Skipping due to missing 'MSGraphExistingCalendarsOfUser' value")
	}

	for _, testCalendar := range testCalendars {
		tt := struct {
			name string
			c    *Calendar
			want string
		}{
			name: "Test first calendar",
			c:    &testCalendar,
			want: fmt.Sprintf("Calendar(ID: \"%v\", Name: \"%v\", canEdit: \"%v\", canShare: \"%v\", canViewPrivateItems: \"%v\", ChangeKey: \"%v\", "+
				"Owner: \"%v\")", testCalendar.ID, testCalendar.Name, testCalendar.CanEdit, testCalendar.CanShare, testCalendar.CanViewPrivateItems, testCalendar.ChangeKey, testCalendar.Owner),
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Calendar.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
