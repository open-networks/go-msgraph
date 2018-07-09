package msgraph

import "strings"

// Groups represents multiple Group-instances.
type Groups []Group

func (g Groups) String() string {
	var groups = make([]string, len(g))
	for i, calendar := range g {
		groups[i] = calendar.String()
	}
	return "Groups(" + strings.Join(groups, " | ") + ")"
}
