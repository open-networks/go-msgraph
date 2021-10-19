package msgraph

import (
	"fmt"
	"testing"
)

func GetTestDevice(t *testing.T) Device {
	t.Helper()
	devices, err := graphClient.ListDevices()
	if err != nil {
		t.Fatalf("Cannot GraphClient.ListDevice(): %v", err)
	}
	deviceTest, err := devices.GetByDisplayName(msGraphExistingDeviceDisplayName)
	if err != nil {
		t.Fatalf("Cannot devices.GetByDisplayName(%v): %v", msGraphExistingDeviceDisplayName, err)
	}
	return deviceTest
}

func TestDevice_String(t *testing.T) {
	testDevice := GetTestDevice(t)

	tests := []struct {
		name string
		g    Device
		want string
	}{
		{
			name: "Test All Devices",
			g:    testDevice,
			want: fmt.Sprintf("Device(ID: \"%v\", DeletedDateTime: \"%v\" AccountEnabled: \"%v\", ApproximateLastSignInDateTime: \"%v\", ComplianceExpirationDateTime: \"%v\", "+
				"DeviceId: \"%v\", DisplayName: \"%v\", IsCompliant: \"%v\" IsManaged: \"%v\", Manufacturer: \"%v\", MdmAppId: \"%v\", Model: \"%v\", OnPremisesLastSyncDateTime: \"%v\", "+
				"OnPremisesSyncEnabled: \"%v\", OperatingSystem: \"%v\", OperatingSystemVersion: \"%v\", PhysicalIds: \"%v\", ProfileType: \"%v\", RegistrationDateTime: \"%v\", "+
				"SystemLabels: \"%v\", TrustType: \"%v\", DirectAPIConnection: %v)",
				testDevice.ID, testDevice.DeletedDateTime, testDevice.AccountEnabled, testDevice.ApproximateLastSignInDateTime, testDevice.ComplianceExpirationDateTime, testDevice.DeviceId,
				testDevice.DisplayName, testDevice.IsCompliant, testDevice.IsManaged, testDevice.Manufacturer, testDevice.MdmAppId, testDevice.Model, testDevice.OnPremisesLastSyncDateTime,
				testDevice.OnPremisesSyncEnabled, testDevice.OperatingSystem, testDevice.OperatingSystemVersion, testDevice.PhysicalIds, testDevice.ProfileType, testDevice.RegistrationDateTime,
				testDevice.SystemLabels, testDevice.TrustType, testDevice.graphClient != nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.String(); got != tt.want {
				t.Errorf("Device.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
