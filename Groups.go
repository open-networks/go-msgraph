package msgraph

import "strings"

// Groups represents multiple Group-instances.
type Groups struct {
	Groups []Group `json:"value"`
}

func (g Groups) String() string {
	var groups = make([]string, len(g.Groups))
	for i, calendar := range g.Groups {
		groups[i] = calendar.String()
	}
	return "Groups(" + strings.Join(groups, " | ") + ")"
}
