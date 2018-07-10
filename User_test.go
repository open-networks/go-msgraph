package msgraph

import (
	"testing"
	"time"
)

func TestUser_ListCalendarView(t *testing.T) {
	TestEnvironmentVariablesPresent(t) // checks the env-variables and failsNow if any is missing

	userToTest, errUserToTest := graphClient.GetUser(msGraphExistingUserPrincipalInGroup)
	if errUserToTest != nil {
		t.Fatalf("Can not find user %v for Testing", msGraphExistingUserPrincipalInGroup)
	}

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
