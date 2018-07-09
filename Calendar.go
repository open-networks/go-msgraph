package msgraph

import (
	"encoding/json"
	"fmt"
)

// Calendar represents a single calendar of a user
type Calendar struct {
	ID                  string // the calendars ID
	Name                string // the name of the calendar
	CanShare            bool   // true if the current account can shares this calendar
	CanViewPrivateItems bool   // true if the current account can view private entries
	CanEdit             bool   // true if the current account can edit the calendar
	OwnerName           string // onwer name determined by e-mail addr field, not reliable to identify the user
	OwnerEMailAddr      string // email addr of ofnwer, reliable to identify the user
}

func (c *Calendar) String() string {
	return fmt.Sprintf("Calendar(ID: \"%v\", Name: \"%v\", CanShare: \"%v\", CanViewPrivateItems: \"%v\", "+
		"canEdit: \"%v\", OwnerName: \"%v\", OnwerEMailAddr: \"%v\")", c.ID, c.Name, c.CanShare, c.CanViewPrivateItems, c.CanEdit, c.OwnerName, c.OwnerEMailAddr)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (c *Calendar) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID                  string `json:"id"`                  // the calendars ID
		Name                string `json:"name"`                // the name of the calendar
		CanShare            bool   `json:"canShare"`            // true if the current account can shares this calendar
		CanViewPrivateItems bool   `json:"canViewPrivateItems"` // true if the current account can view private entries
		CanEdit             bool   `json:"canEdit"`             // true if the current account can edit the calendar
		Owner               struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"owner"`
	}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	c.ID = tmp.ID
	c.Name = tmp.Name
	c.CanShare = tmp.CanShare
	c.CanViewPrivateItems = tmp.CanViewPrivateItems
	c.CanEdit = tmp.CanEdit
	c.OwnerName = tmp.Owner.Name
	c.OwnerEMailAddr = tmp.Owner.Address

	return nil

}
