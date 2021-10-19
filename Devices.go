package msgraph

import (
	"strings"
)

// Devices represents multiple Device-instances and provides funcs to work with them.
type Devices []Device

func (d Devices) String() string {
	var devices = make([]string, len(d))
	for i, calendar := range d {
		devices[i] = calendar.String()
	}
	return "Devices(" + strings.Join(devices, " | ") + ")"
}

// setGraphClient sets the GraphClient within that particular instance. Hence it's directly created by GraphClient
func (d Devices) setGraphClient(gC *GraphClient) Devices {
	for i := range d {
		d[i].setGraphClient(gC)
	}
	return d
}

// GetByDisplayName returns the Device obj of that array whose DisplayName matches
// the given name. Returns an ErrFindDevice if no device exists that matches the given
// DisplayName.
func (d Devices) GetByDisplayName(displayName string) (Device, error) {
	for _, device := range d {
		if device.DisplayName == displayName {
			return device, nil
		}
	}
	return Device{}, ErrFindDevice
}
