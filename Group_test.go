package msgraph

import (
	"testing"
)

func TestGroup_ListMembers(t *testing.T) {
	TestEnvironmentVariablesPresent(t) // checks the env-variables and failsNow if any is missing

	groups, _ := graphClient.ListGroups()
	groupTest, err := groups.GetByDisplayName(msGraphExistingGroupDisplayName) // if that failed we dont have to check the err because the test will fail anyway
	if err != nil {
		t.Fatalf("Group \"%v\" not found in %v", msGraphExistingGroupDisplayName, groups)
	}

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
