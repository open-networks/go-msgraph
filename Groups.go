package msgraph

import (
	"strings"
)

// Groups represents multiple Group-instances.
type Groups []Group

func (g Groups) String() string {
	var groups = make([]string, len(g))
	for i, calendar := range g {
		groups[i] = calendar.String()
	}
	return "Groups(" + strings.Join(groups, " | ") + ")"
}

// setGraphClient sets the GraphClient within that particular instance. Hence it's directly created by GraphClient
func (g Groups) setGraphClient(gC *GraphClient) Groups {
	for i := range g {
		g[i].setGraphClient(gC)
	}
	return g
}

// GetByDisplayName returns the Group obj of that array whose DisplayName matches
// the given name. Returns an ErrFindGroup if no group exists that matches the given
// DisplayName.
func (g Groups) GetByDisplayName(displayName string) (Group, error) {
	for _, group := range g {
		if group.DisplayName == displayName {
			return group, nil
		}
	}
	return Group{}, ErrFindGroup
}
