package msgraph

import (
	"fmt"
)

// Domain represents one domain of ms graph
type Domain struct {
	ID                               string      `json:"id,omitempty"`
	AuthenticationType               string      `json:"authenticationType,omitempty"`
	IsAdminManaged                   bool        `json:"isAdminManaged,omitempty"`
	IsDefault                        bool        `json:"isDefault,omitempty"`
	IsInitial                        bool        `json:"isInitial,omitempty"`
	IsRoot                           bool        `json:"isRoot,omitempty"`
	IsVerified                       bool        `json:"isVerified,omitempty"`
	SupportedServices                []string    `json:"supportedServices,omitempty"`
	PasswordValidityPeriodInDays     *int        `json:"passwordValidityPeriodInDays,omitempty"`
	PasswordNotificationWindowInDays *int        `json:"passwordNotificationWindowInDays,omitempty"`
	State                            interface{} `json:"state,omitempty"`

	graphClient *GraphClient // the graphClient that called the group
}

func (d Domain) String() string {
	return fmt.Sprintf("Domain(ID: \"%v\", AuthenticationType: \"%v\" IsAdminManaged: \"%v\", IsDefault: \"%v\", IsInitial: \"%v\", IsRoot: \"%v\", "+
		"IsVerified: \"%v\", SupportedServices: \"%v\", PasswordValidityPeriodInDays: \"%v\", PasswordNotificationWindowInDays: \"%v\", DirectAPIConnection: %v)",
		d.ID, d.AuthenticationType, d.IsAdminManaged, d.IsDefault, d.IsInitial, d.IsRoot, d.IsVerified, d.SupportedServices, d.PasswordValidityPeriodInDays,
		d.PasswordNotificationWindowInDays, d.graphClient != nil)
}

// setGraphClient sets the graphClient instance in this instance and all child-instances (if any)
func (d *Domain) setGraphClient(gC *GraphClient) {
	d.graphClient = gC
}

// ListMembers - Get a list of the domain's direct members.
func (d Domain) ListMembers(opts ...ListQueryOption) (Users, error) {
	if d.graphClient == nil {
		return nil, ErrNotGraphClientSourced
	}
	resource := fmt.Sprintf("/domains/%v/members", d.ID)

	var marsh struct {
		Users Users `json:"value"`
	}
	marsh.Users.setGraphClient(d.graphClient)
	return marsh.Users, d.graphClient.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
}
