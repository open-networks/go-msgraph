package msgraph

import "fmt"

// Group represents one group
type Group struct {
	ID           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Mail         string `json:"mail"`
	MailEnabled  bool   `json:"mailEnabled"`
	MailNickname string `json:"mailNickname"`
}

func (g Group) String() string {
	return fmt.Sprintf("Group(ID: \"%v\", DisplayName: \"%v\", Mail: \"%v\", MailEnabled: \"%v\", MailNickname: \"%v\")",
		g.ID, g.DisplayName, g.Mail, g.MailEnabled, g.MailNickname)
}
