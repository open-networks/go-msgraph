package msgraph

import (
	"testing"
	"time"
)

var (
	// First Attendee used for the Unit tests
	testAttendee1 = Attendee{Name: "Testname", Email: "testname@contoso.com", Type: "attendee", ResponseStatus: ResponseStatus{Response: "accepted", Time: time.Now()}}
	// Second Attendee used for the Unit tests
	testAttendee2 = Attendee{Name: "Testuser", Email: "testuser@contoso.com", Type: "attendee", ResponseStatus: ResponseStatus{Response: "declined", Time: time.Now()}}
)

func TestAttendee_Equal(t *testing.T) {
	type args struct {
		other Attendee
	}
	tests := []struct {
		name string
		a    Attendee
		args args
		want bool
	}{
		{
			name: "All equal",
			a:    testAttendee1,
			args: args{other: testAttendee1},
			want: true,
		}, {
			name: "non-equal",
			a:    testAttendee1,
			args: args{other: testAttendee2},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(tt.args.other); got != tt.want {
				t.Errorf("Attendee.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttendee_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		a       *Attendee
		args    args
		wantErr bool
	}{
		{
			name: "All good",
			a:    &Attendee{},
			args: args{data: []byte(
				`{	"Type" : "organizer", 
					"Status": {"response" : "accepted", "time" : "` + time.Now().Format(time.RFC3339Nano) + `"},
					"emailAddress" : {
						"name" : "TestUserName",
						"address": "TestUserName@contoso.com"
					}}`)},
			wantErr: false,
		}, {
			name: "wrong json",
			a:    &Attendee{},
			args: args{data: []byte(
				`{	"Type" : "organizer", 
					"Status": {"response" : "accepted", "time" : "` + time.Now().Format(time.RFC3339Nano) + `"},
					"emailAddress" : {
						"name" : "TestUserName",
						"address": "TestUserName@contoso.com",
					}}`)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Attendee.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAttendee_String(t *testing.T) {
	tests := []struct {
		name string
		a    Attendee
		want string
	}{
		{
			name: "Test String-func",
			a:    testAttendee1,
			want: "Name: Testname, Type: attendee, E-mail: testname@contoso.com, ResponseStatus: Response: accepted, Time: " + testAttendee1.ResponseStatus.Time.Format(time.RFC3339Nano),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("Attendee.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
