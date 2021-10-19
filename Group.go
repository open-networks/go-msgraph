package msgraph

import (
	"encoding/json"
	"fmt"
	"time"
)

// Group represents one group of ms graph
//
// See: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/group_get
type Group struct {
	ID                           string    `json:"id,omitempty"`
	Description                  string    `json:"description,omitempty"`
	DisplayName                  string    `json:"displayName,omitempty"`
	CreatedDateTime              time.Time `json:"createdDateTime,omitempty"`
	GroupTypes                   []string  `json:"groupTypes,omitempty"`
	Mail                         string    `json:"mail,omitempty"`
	MailEnabled                  bool      `json:"mailEnabled,omitempty"`
	MailNickname                 string    `json:"mailNickname,omitempty"`
	OnPremisesLastSyncDateTime   time.Time `json:"onPremisesLastSyncDateTime,omitempty"` // defaults to 0001-01-01 00:00:00 +0000 UTC if there's none
	OnPremisesSecurityIdentifier string    `json:"onPremisesSecurityIdentifier,omitempty"`
	OnPremisesSyncEnabled        bool      `json:"onPremisesSyncEnabled,omitempty"`
	ProxyAddresses               []string  `json:"proxyAddresses,omitempty"`
	SecurityEnabled              bool      `json:"securityEnabled,omitempty"`
	Visibility                   string    `json:"visibility,omitempty"`
	MemberOf                     []Member  `json:"memberOf,omitempty"`

	graphClient *GraphClient // the graphClient that called the group
}

func (g Group) String() string {
	return fmt.Sprintf("Group(ID: \"%v\", Description: \"%v\" DisplayName: \"%v\", CreatedDateTime: \"%v\", GroupTypes: \"%v\", Mail: \"%v\", MailEnabled: \"%v\", MailNickname: \"%v\", OnPremisesLastSyncDateTime: \"%v\", OnPremisesSecurityIdentifier: \"%v\", OnPremisesSyncEnabled: \"%v\", ProxyAddresses: \"%v\", SecurityEnabled \"%v\", Visibility: \"%v\", DirectAPIConnection: %v)",
		g.ID, g.Description, g.DisplayName, g.CreatedDateTime, g.GroupTypes, g.Mail, g.MailEnabled, g.MailNickname, g.OnPremisesLastSyncDateTime, g.OnPremisesSecurityIdentifier, g.OnPremisesSyncEnabled, g.ProxyAddresses, g.SecurityEnabled, g.Visibility, g.graphClient != nil)
}

// setGraphClient sets the graphClient instance in this instance and all child-instances (if any)
func (g *Group) setGraphClient(gC *GraphClient) {
	g.graphClient = gC
}

// ListMembers - Get a list of the group's direct members. A group can have users,
// contacts, and other groups as members. This operation is not transitive. This
// method will currently ONLY return User-instances of members
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/group_list_members
func (g Group) ListMembers(opts ...ListQueryOption) (Users, error) {
	if g.graphClient == nil {
		return nil, ErrNotGraphClientSourced
	}
	resource := fmt.Sprintf("/groups/%v/members", g.ID)

	var marsh struct {
		Users Users `json:"value"`
	}
	marsh.Users.setGraphClient(g.graphClient)
	return marsh.Users, g.graphClient.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (g *Group) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID                           string   `json:"id"`
		Description                  string   `json:"description"`
		DisplayName                  string   `json:"displayName"`
		CreatedDateTime              string   `json:"createdDateTime"`
		GroupTypes                   []string `json:"groupTypes"`
		Mail                         string   `json:"mail"`
		MailEnabled                  bool     `json:"mailEnabled"`
		MailNickname                 string   `json:"mailNickname"`
		OnPremisesLastSyncDateTime   string   `json:"onPremisesLastSyncDateTime"`
		OnPremisesSecurityIdentifier string   `json:"onPremisesSecurityIdentifier"`
		OnPremisesSyncEnabled        bool     `json:"onPremisesSyncEnabled"`
		ProxyAddresses               []string `json:"proxyAddresses"`
		SecurityEnabled              bool     `json:"securityEnabled"`
		Visibility                   string   `json:"visibility"`
		MemberOf                     []Member `json:"memberOf"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	g.ID = tmp.ID
	g.Description = tmp.Description
	g.DisplayName = tmp.DisplayName
	g.CreatedDateTime, err = time.Parse(time.RFC3339, tmp.CreatedDateTime)
	if err != nil && tmp.CreatedDateTime != "" {
		return fmt.Errorf("cannot parse CreatedDateTime %v with RFC3339: %v", tmp.CreatedDateTime, err)
	}
	g.GroupTypes = tmp.GroupTypes
	g.Mail = tmp.Mail
	g.MailEnabled = tmp.MailEnabled
	g.MailNickname = tmp.MailNickname
	g.OnPremisesLastSyncDateTime, err = time.Parse(time.RFC3339, tmp.OnPremisesLastSyncDateTime)
	if err != nil && tmp.OnPremisesLastSyncDateTime != "" {
		return fmt.Errorf("cannot parse OnPremisesLastSyncDateTime %v with RFC3339: %v", tmp.OnPremisesLastSyncDateTime, err)
	}
	g.OnPremisesSecurityIdentifier = tmp.OnPremisesSecurityIdentifier
	g.OnPremisesSyncEnabled = tmp.OnPremisesSyncEnabled
	g.ProxyAddresses = tmp.ProxyAddresses
	g.SecurityEnabled = tmp.SecurityEnabled
	g.Visibility = tmp.Visibility
	g.MemberOf = tmp.MemberOf

	return nil
}
