package msgraph

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// get graph client config from environment
var (
	msGraphTenantID                     = os.Getenv("MSGraphTenantID")
	msGraphApplicationID                = os.Getenv("MSGraphApplicationID")
	msGraphClientSecret                 = os.Getenv("MSGraphClientSecret")
	msGraphExistingGroupDisplayName     = os.Getenv("MSGraphExistingGroupDisplayName")
	msGraphExistingUserPrincipalInGroup = os.Getenv("MSGraphExistingUserPrincipalInGroup")
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

func TestGraphClient_ListGroups(t *testing.T) {
	tests := []struct {
		name    string
		g       *GraphClient
		want    Groups
		wantErr bool
	}{
		{
			name: "Test if one of the groups provided is present",
			g:    graphClient,
			want: Groups{
				Group{DisplayName: "Employees"},
				Group{DisplayName: "Finance"},
				Group{DisplayName: "Sales"},
				Group{DisplayName: "Admins"},
				Group{DisplayName: "Administrators"},
			},
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
			var found bool
			for _, searchObj := range tt.want {
				for _, checkObj := range got {
					found = found || searchObj.DisplayName == checkObj.DisplayName
				}
			}
			if !found {
				t.Errorf("GraphClient.ListGroups() = %v, searching for one of %v", got, tt.want)
			}
		})
	}
}

func TestGraphClient_UnmarshalJSON(t *testing.T) {
	TestEnvironmentVariablesPresent(t) // check prerequisites

	jsonString := fmt.Sprintf("{"+
		"\"TenantID\": \"%v\","+
		"\"ApplicationID\": \"%v\","+
		"\"ClientSecret\": \"%v\""+
		"}",
		msGraphTenantID, msGraphApplicationID, msGraphClientSecret,
	)
	var marsh GraphClient
	err := json.Unmarshal([]byte(jsonString), &marsh)
	if err != nil {
		t.Errorf("Error on JSON-Unmarshal: %v", err)
	}
}
