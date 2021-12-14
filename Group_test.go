package msgraph

import (
	"fmt"
	"testing"
)

func GetTestGroup(t *testing.T) Group {
	t.Helper()
	groups, err := graphClient.ListGroups()
	if err != nil {
		t.Fatalf("Cannot GraphClient.ListGroups(): %v", err)
	}
	groupTest, err := groups.GetByDisplayName(msGraphExistingGroupDisplayName)
	if err != nil {
		t.Fatalf("Cannot groups.GetByDisplayName(%v): %v", msGraphExistingGroupDisplayName, err)
	}
	return groupTest
}

func TestGroup_ListMembers(t *testing.T) {
	groupTest := GetTestGroup(t)

	tests := []struct {
		name    string
		g       Group
		want    Users
		wantErr bool
	}{
		{
			name:    "GraphClient created Group",
			g:       groupTest,
			want:    Users{User{UserPrincipalName: msGraphExistingUserPrincipalInGroup}},
			wantErr: false,
		}, {
			name:    "Not GraphClient created Group",
			g:       Group{DisplayName: "Test"},
			want:    Users{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ListMembers()
			if (err != nil) != tt.wantErr {
				t.Errorf("Group.ListMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var found bool
			for _, searchObj := range tt.want {
				for _, checkObj := range got {
					found = found || searchObj.UserPrincipalName == checkObj.UserPrincipalName
				}
			}
			if !found && len(tt.want) > 0 {
				t.Errorf("GraphClient.ListGroups() = %v, searching for one of %v", got, tt.want)
			}
		})
	}
}

func TestGroup_String(t *testing.T) {
	testGroup := GetTestGroup(t)

	tests := []struct {
		name string
		g    Group
		want string
	}{
		{
			name: "Test All Groups",
			g:    testGroup,
			want: fmt.Sprintf("Group(ID: \"%v\", Description: \"%v\" DisplayName: \"%v\", CreatedDateTime: \"%v\", GroupTypes: \"%v\", Mail: \"%v\", MailEnabled: \"%v\", MailNickname: \"%v\", OnPremisesLastSyncDateTime: \"%v\", OnPremisesSecurityIdentifier: \"%v\", OnPremisesSyncEnabled: \"%v\", ProxyAddresses: \"%v\", SecurityEnabled \"%v\", Visibility: \"%v\", DirectAPIConnection: %v)",
				testGroup.ID, testGroup.Description, testGroup.DisplayName, testGroup.CreatedDateTime, testGroup.GroupTypes, testGroup.Mail, testGroup.MailEnabled, testGroup.MailNickname, testGroup.OnPremisesLastSyncDateTime, testGroup.OnPremisesSecurityIdentifier, testGroup.OnPremisesSyncEnabled, testGroup.ProxyAddresses, testGroup.SecurityEnabled, testGroup.Visibility, testGroup.graphClient != nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.String(); got != tt.want {
				t.Errorf("Group.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_ListTransitiveMembers(t *testing.T) {
	testGroup := GetTestGroup(t)

	tests := []struct {
		name    string
		g       Group
		wantErr bool
	}{
		{
			name:    "GraphClient created Group",
			g:       testGroup,
			wantErr: false,
		}, {
			name:    "Not GraphClient created Group",
			g:       Group{DisplayName: "Test not GraphClient sourced"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ListTransitiveMembers()
			if (err != nil) != tt.wantErr {
				t.Errorf("Group.ListTransitiveMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("Group.ListTransitiveMembers() = %v, len(%d), want at least one member of that group", got, len(got))
			}
		})
	}
}

func TestGroup_GetMemberGroupsAsStrings(t *testing.T) {
	testGroup := GetTestGroup(t)

	tests := []struct {
		name    string
		g       Group
		opts    []GetQueryOption
		wantErr bool
	}{
		{
			name:    "Test group func GetMembershipGroupsAsStrings",
			g:       testGroup,
			wantErr: false,
		}, {
			name:    "Test group func GetMembershipGroupsAsStrings - no securityGroupsEnabeledF",
			g:       testGroup,
			wantErr: false,
		},
		{
			name:    "Group not initialized by GraphClient",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.GetMemberGroupsAsStrings(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Group.GetMemberGroupsAsStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("Group.GetMemberGroupsAsStrings() = %v, len(%d), want at least one value", got, len(got))
			}
		})
	}
}
