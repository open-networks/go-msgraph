package msgraph

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// get graph client config from environment
var (
	// Microsoft Graph tenant ID
	msGraphTenantID = os.Getenv("MSGraphTenantID")
	// Microsoft Graph Application ID
	msGraphApplicationID = os.Getenv("MSGraphApplicationID")
	// Microsoft Graph Client Secret
	msGraphClientSecret = os.Getenv("MSGraphClientSecret")
	// a valid groupdisplayname from msgraph, e.g. technicians@contoso.com
	msGraphExistingGroupDisplayName = os.Getenv("MSGraphExistingGroupDisplayName")
	// a valid userprincipalname in the above group, e.g. felix@contoso.com
	msGraphExistingUserPrincipalInGroup = os.Getenv("MSGraphExistingUserPrincipalInGroup")
	// valid calendar names that belong to the above user, sepearated by a colon (","). e.g.: "Kalender,Feiertage in Österreich,Geburtstage"
	msGraphExistingCalendarsOfUser = strings.Split(os.Getenv("MSGraphExistingCalendarsOfUser"), ",")
)

// the graphclient used to perform all tests
var graphClient, _ = NewGraphClient(msGraphTenantID, msGraphApplicationID, msGraphClientSecret)

func TestEnvironmentVariablesPresent(t *testing.T) {
	if msGraphTenantID == "" {
		t.Fatal("Environment Variable for Tenant ID named <MSGraphTenantID> is mising!")
	}
	if msGraphApplicationID == "" {
		t.Fatal("Environment Variable for Application ID named <MSGraphApplicationID> is mising!")
	}
	if msGraphClientSecret == "" {
		t.Fatal("Environment Variable for Client Secret named <MSGraphClientSecret> is mising!")
	}
	if msGraphExistingGroupDisplayName == "" {
		t.Fatal("Environment Variable for an existing group's displayName named <MSGraphExistingGroupDisplayName> is mising!")
	}
	if msGraphExistingUserPrincipalInGroup == "" {
		t.Fatal("Environment Variable for an existing user in the group named <MSGraphExistingUserPrincipalInGroup> is missing!")
	}
	if msGraphExistingCalendarsOfUser[0] == "" {
		t.Fatal("Environment Variable for existing calendars of the given user named <MSGraphExistingCalendarsOfUser> is missing!")
	}
}

func TestNewGraphClient(t *testing.T) {
	type args struct {
		tenantID      string
		applicationID string
		clientSecret  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "GraphClient from Environment-variables",
			args:    args{tenantID: msGraphTenantID, applicationID: msGraphApplicationID, clientSecret: msGraphClientSecret},
			wantErr: false,
		}, {
			name:    "GraphClient fail - wrong tenant ID",
			args:    args{tenantID: "wrong tenant id", applicationID: msGraphApplicationID, clientSecret: msGraphClientSecret},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong application ID",
			args:    args{tenantID: msGraphTenantID, applicationID: "wrong applicatio id", clientSecret: msGraphClientSecret},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong client secret",
			args:    args{tenantID: msGraphTenantID, applicationID: msGraphApplicationID, clientSecret: "wrong client secret"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGraphClient(tt.args.tenantID, tt.args.applicationID, tt.args.clientSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGraphClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGraphClient_ListUsers(t *testing.T) {
	tests := []struct {
		name    string
		g       *GraphClient
		want    User
		wantErr bool
	}{
		{
			name:    fmt.Sprintf("List all Users, check for user %v", msGraphExistingUserPrincipalInGroup),
			g:       graphClient,
			want:    User{UserPrincipalName: msGraphExistingUserPrincipalInGroup},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ListUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("GraphClient.ListUsers() len = 0, want more than 0: %v", got)
			}
			isGraphClientInitializd := true
			found := false
			for _, user := range got {
				isGraphClientInitializd = isGraphClientInitializd && user.graphClient != nil
				found = found || user.UserPrincipalName == tt.want.UserPrincipalName
			}
			if !found {
				t.Errorf("GraphClient.ListUsers() user %v not found, users: %v", tt.want.UserPrincipalName, got)
			}
			if !isGraphClientInitializd {
				t.Errorf("GraphClient.ListUsers() graphClient is nil, but was initialized from GraphClient")
			}
		})
	}
}

func TestGraphClient_ListGroups(t *testing.T) {
	tests := []struct {
		name    string
		g       *GraphClient
		want    Group
		wantErr bool
	}{
		{
			name:    fmt.Sprintf("Test if Group %v is present", msGraphExistingGroupDisplayName),
			g:       graphClient,
			want:    Group{DisplayName: msGraphExistingGroupDisplayName},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ListGroups()
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.ListGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			found := false
			isGraphClientInitializd := true
			for _, checkObj := range got {
				found = found || tt.want.DisplayName == checkObj.DisplayName
				isGraphClientInitializd = isGraphClientInitializd && checkObj.graphClient != nil
			}
			if !found {
				t.Errorf("GraphClient.ListGroups() = %v, searching for one of %v", got, tt.want)
			}
			if !isGraphClientInitializd {
				t.Errorf("GraphClient.ListGroups() graphClient is nil, but was initialized from GraphClient")
			}
		})
	}
}

func TestGraphClient_GetUser(t *testing.T) {
	type args struct {
		identifier string
	}
	tests := []struct {
		name    string
		g       *GraphClient
		args    args
		want    User
		wantErr bool
	}{
		{
			name:    fmt.Sprintf("Test if user %v is present", msGraphExistingUserPrincipalInGroup),
			g:       graphClient,
			args:    args{identifier: msGraphExistingUserPrincipalInGroup},
			want:    User{UserPrincipalName: msGraphExistingUserPrincipalInGroup},
			wantErr: false,
		}, {
			name:    "Test if non-existing user produces err",
			g:       graphClient,
			args:    args{identifier: "ThisUserwillNotExistForSure@contoso.com"},
			want:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.GetUser(tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.UserPrincipalName != tt.want.UserPrincipalName {
				t.Errorf("GraphClient.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphClient_GetGroup(t *testing.T) {
	tests := []struct {
		name    string
		g       *GraphClient
		want    Group
		wantErr bool
	}{
		{
			name:    fmt.Sprintf("Test if Group %v is presnt and GetGroup-able", msGraphExistingGroupDisplayName),
			g:       graphClient,
			want:    Group{DisplayName: msGraphExistingGroupDisplayName},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allGroups, err := tt.g.ListGroups()
			if err != nil { // check if groups can be listed
				t.Fatalf("GraphClient.ListGroups(): cannot list groups: %v", err)
			}
			targetGroup, err := allGroups.GetByDisplayName(tt.want.DisplayName)
			if err != nil { // check if the group to be tested is in the list
				t.Fatalf("Groups.GetByDisplayName(): cannot find group %v in %v, err: %v", tt.want.DisplayName, allGroups, err)
			}
			got, err := tt.g.GetGroup(targetGroup.ID) // actually execute the test we want to test
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.GetGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(got.DisplayName == tt.want.DisplayName) {
				t.Errorf("GraphClient.GetGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphClient_UnmarshalJSON(t *testing.T) {
	TestEnvironmentVariablesPresent(t) // check prerequisites

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "All correct",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\"}", msGraphTenantID, msGraphApplicationID, msGraphClientSecret))},
			wantErr: false,
		}, {
			name:    "JSON-syntax error",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\"", msGraphTenantID, msGraphApplicationID, msGraphClientSecret))},
			wantErr: true,
		}, {
			name:    "TenantID incorrect",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\"}", "wrongtenant", msGraphApplicationID, msGraphClientSecret))},
			wantErr: true,
		}, {
			name:    "TenantID empty",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\"}", "", msGraphApplicationID, msGraphClientSecret))},
			wantErr: true,
		}, {
			name:    "ApplicationID incorrect",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\"}", msGraphTenantID, "wrongapplication", msGraphClientSecret))},
			wantErr: true,
		}, {
			name:    "ApplicationID empty",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\"}", msGraphTenantID, "", msGraphClientSecret))},
			wantErr: true,
		}, {
			name:    "ClientSecret incorrect",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\"}", msGraphTenantID, msGraphApplicationID, "wrongclientsecret"))},
			wantErr: true,
		}, {
			name:    "ClientSecret empty",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\"}", msGraphTenantID, msGraphApplicationID, ""))},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unmarshalTest GraphClient
			if err := unmarshalTest.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
