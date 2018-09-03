package msgraph

import (
	"encoding/json"
	"fmt"
)

// Calendar represents a single calendar of a user
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/resources/calendar
type Calendar struct {
	ID                  string // The group's unique identifier. Read-only.
	Name                string // The calendar name.
	CanEdit             bool   // True if the user can write to the calendar, false otherwise. This property is true for the user who created the calendar. This property is also true for a user who has been shared a calendar and granted write access.
	CanShare            bool   // True if the user has the permission to share the calendar, false otherwise. Only the user who created the calendar can share it.
	CanViewPrivateItems bool   // True if the user can read calendar items that have been marked private, false otherwise.
	ChangeKey           string // Identifies the version of the calendar object. Every time the calendar is changed, changeKey changes as well. This allows Exchange to apply changes to the correct version of the object. Read-only.

	Owner EmailAddress // If set, this represents the user who created or added the calendar. For a calendar that the user created or added, the owner property is set to the user. For a calendar shared with the user, the owner property is set to the person who shared that calendar with the user.

	graphClient *GraphClient // the graphClient that created this instance
}

func (c Calendar) String() string {
	return fmt.Sprintf("Calendar(ID: \"%v\", Name: \"%v\", canEdit: \"%v\", canShare: \"%v\", canViewPrivateItems: \"%v\", ChangeKey: \"%v\", "+
		"Owner: \"%v\")", c.ID, c.Name, c.CanEdit, c.CanShare, c.CanViewPrivateItems, c.ChangeKey, c.Owner)
}

// setGraphClient sets the graphClient instance in this instance and all child-instances (if any)
func (c *Calendar) setGraphClient(graphClient *GraphClient) {
	c.graphClient = graphClient
	c.Owner.setGraphClient(graphClient)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (c *Calendar) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID                  string `json:"id"`                  // the calendars ID
		Name                string `json:"name"`                // the name of the calendar
		CanShare            bool   `json:"canShare"`            // true if the current account can shares this calendar
		CanViewPrivateItems bool   `json:"canViewPrivateItems"` // true if the current account can view private entries
		CanEdit             bool   `json:"canEdit"`             // true if the current account can edit the calendar
		ChangeKey           string `json:"changeKey"`
		Owner               EmailAddress
	}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	c.ID = tmp.ID
	c.Name = tmp.Name
	c.CanEdit = tmp.CanEdit
	c.CanShare = tmp.CanShare
	c.CanViewPrivateItems = tmp.CanViewPrivateItems
	c.ChangeKey = tmp.ChangeKey

	c.Owner = tmp.Owner

	return nil
}
