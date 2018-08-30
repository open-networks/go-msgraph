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

// ListMembers lists all members of the group, returns it as a Users-Instance. Only works
// if the Group has been loaded via a graphClient.
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/group_list_members
func (g Group) ListMembers() (Users, error) {
	if g.graphClient == nil {
		return nil, ErrNotGraphClientSourced
	}
	resource := fmt.Sprintf("/groups/%v/members", g.ID)

	var marsh struct {
		Users Users `json:"value"`
	}
	marsh.Users.setGraphClient(g.graphClient)
	return marsh.Users, g.graphClient.makeGETAPICall(resource, nil, &marsh)
}

func (g Group) String() string {
	return fmt.Sprintf("Group(ID: \"%v\", DisplayName: \"%v\", Mail: \"%v\", MailEnabled: \"%v\", MailNickname: \"%v\", DirectAPIConnection: %v)",
		g.ID, g.DisplayName, g.Mail, g.MailEnabled, g.MailNickname, g.graphClient != nil)
}
