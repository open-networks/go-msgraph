package msgraph

import (
	"testing"
	"time"
)

func TestAttendees_String(t *testing.T) {
	tests := []struct {
		name string
		a    Attendees
		want string
	}{
		{
			name: "test eq",
			a:    Attendees{testAttendee1, testAttendee2},
			want: "Attendees(Name: Testname, Type: attendee, E-mail: testname@contoso.com, ResponseStatus: Response: accepted, Time: " +
				testAttendee1.ResponseStatus.Time.Format(time.RFC3339Nano) +
				" | Name: Testuser, Type: attendee, E-mail: testuser@contoso.com, ResponseStatus: Response: declined, Time: " +
				testAttendee2.ResponseStatus.Time.Format(time.RFC3339Nano) + ")",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("Attendees.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestAttendees_Equal(t *testing.T) {
	type args struct {
		other Attendees
	}
	tests := []struct {
		name string
		a    Attendees
		args args
		want bool
	}{
		{
			name: "Equal, different order",
			a:    Attendees{testAttendee1, testAttendee2},
			args: args{other: Attendees{testAttendee2, testAttendee1}},
			want: true,
		}, {
			name: "non-equal, missing one",
			a:    Attendees{testAttendee1, testAttendee2},
			args: args{other: Attendees{testAttendee1}},
			want: false,
		}, {
			name: "non-equal, too many",
			a:    Attendees{testAttendee1},
			args: args{other: Attendees{testAttendee2, testAttendee1}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(tt.args.other); got != tt.want {
				t.Errorf("Attendees.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
