package msgraph

import (
	"reflect"
	"testing"
)

var (
	testUser1 = User{ID: "567898765678", UserPrincipalName: "test1@contoso.com", Mail: "test1@contoso.com", MobilePhone: "+1 12346789"}
	testUser2 = User{ID: "123456", UserPrincipalName: "test2@contoso.com", Mail: "test2@contoso.com", BusinessPhones: []string{"+9 9876543"}}
)

func TestUsers_GetUserByShortName(t *testing.T) {
	type args struct {
		shortName string
	}
	tests := []struct {
		name    string
		u       Users
		args    args
		want    User
		wantErr bool
	}{
		{
			name:    "Find user",
			u:       Users{testUser1, testUser2},
			args:    args{shortName: testUser1.GetShortName()},
			want:    testUser1,
			wantErr: false,
		}, {
			name:    "Missing user",
			u:       Users{testUser1, testUser2},
			args:    args{shortName: "forsurenotexistingusername@contoso.com"},
			want:    User{},
			wantErr: true,
		}, {
			name:    "faulty userprincipal",
			u:       Users{testUser1, testUser2},
			args:    args{shortName: "forsurenotexistingusername"},
			want:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.GetUserByShortName(tt.args.shortName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.GetUserByShortName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.GetUserByShortName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_GetUserByActivePhone(t *testing.T) {
	type args struct {
		activePhone string
	}
	tests := []struct {
		name    string
		u       Users
		args    args
		want    User
		wantErr bool
	}{
		{
			name:    "Find user with mobile phone",
			u:       Users{testUser1, testUser2},
			args:    args{activePhone: testUser1.GetActivePhone()},
			want:    testUser1,
			wantErr: false,
		}, {
			name:    "Find user with business phone",
			u:       Users{testUser1, testUser2},
			args:    args{activePhone: testUser2.GetActivePhone()},
			want:    testUser2,
			wantErr: false,
		}, {
			name:    "Missing user",
			u:       Users{testUser1, testUser2},
			args:    args{activePhone: "123"},
			want:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.GetUserByActivePhone(tt.args.activePhone)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.GetUserByActivePhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.GetUserByActivePhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_GetUserByMail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		u       Users
		args    args
		want    User
		wantErr bool
	}{
		{
			name:    "Find user",
			u:       Users{testUser1, testUser2},
			args:    args{email: testUser1.Mail},
			want:    testUser1,
			wantErr: false,
		}, {
			name:    "Missing user",
			u:       Users{testUser1, testUser2},
			args:    args{email: "forsurenotexistingusername@contoso.com"},
			want:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.GetUserByMail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.GetUserByMail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.GetUserByMail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_Equal(t *testing.T) {
	type args struct {
		other Users
	}
	tests := []struct {
		name string
		u    Users
		args args
		want bool
	}{
		{
			name: "Equal Users",
			u:    Users{testUser1, testUser2},
			args: args{other: Users{testUser2, testUser1}},
			want: true,
		}, {
			name: "non-equal Users",
			u:    Users{testUser1, testUser2},
			args: args{other: Users{testUser2}},
			want: false,
		}, {
			name: "too many users",
			u:    Users{testUser1},
			args: args{other: Users{testUser1, testUser2}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Equal(tt.args.other); got != tt.want {
				t.Errorf("Users.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
