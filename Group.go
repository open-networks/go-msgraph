package msgraph

import "fmt"

// Group represents one group
type Group struct {
	ID           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Mail         string `json:"mail"`
	MailEnabled  bool   `json:"mailEnabled"`
	MailNickname string `json:"mailNickname"`

	graphClient *GraphClient // the graphClient that called the group
}

// ListMembers lists all members of that group
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/group_list_members
func (g Group) ListMembers() (Users, error) {
	if g.graphClient == nil {
		return nil, ErrNotGraphClientSourced
	}
	return g.graphClient.ListMembersOfGroup(g.ID)
}

func (g Group) String() string {
	return fmt.Sprintf("Group(ID: \"%v\", DisplayName: \"%v\", Mail: \"%v\", MailEnabled: \"%v\", MailNickname: \"%v\", DirectAPIConnection: %v)",
		g.ID, g.DisplayName, g.Mail, g.MailEnabled, g.MailNickname, g.graphClient != nil)
}
