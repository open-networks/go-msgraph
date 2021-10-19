package msgraph

import (
	"encoding/json"
	"fmt"
	"time"
)

// Device represents one device of ms graph
type Device struct {
	ID                            string    `json:"id,omitempty"`
	DeletedDateTime               time.Time `json:"deletedDateTime,omitempty"`
	AccountEnabled                bool      `json:"accountEnabled,omitempty"`
	ApproximateLastSignInDateTime time.Time `json:"approximateLastSignInDateTime,omitempty"`
	ComplianceExpirationDateTime  time.Time `json:"complianceExpirationDateTime,omitempty"`
	DeviceId                      string    `json:"deviceId,omitempty"`
	DisplayName                   string    `json:"displayName,omitempty"`
	IsCompliant                   bool      `json:"isCompliant,omitempty"`
	IsManaged                     bool      `json:"isManaged,omitempty"`
	Manufacturer                  string    `json:"manufacturer,omitempty"`
	MdmAppId                      string    `json:"mdmAppId,omitempty"`
	Model                         string    `json:"model,omitempty"`
	OnPremisesLastSyncDateTime    time.Time `json:"onPremisesLastSyncDateTime,omitempty"`
	OnPremisesSyncEnabled         bool      `json:"onPremisesSyncEnabled,omitempty"`
	OperatingSystem               string    `json:"operatingSystem,omitempty"`
	OperatingSystemVersion        string    `json:"operatingSystemVersion,omitempty"`
	PhysicalIds                   []string  `json:"physicalIds,omitempty"`
	ProfileType                   string    `json:"profileType,omitempty"`
	RegistrationDateTime          time.Time `json:"registrationDateTime,omitempty"`
	SystemLabels                  []string  `json:"systemLabels,omitempty"`
	TrustType                     string    `json:"trustType,omitempty"`
	MemberOf                      []Member  `json:"memberOf,omitempty"`

	graphClient *GraphClient // the graphClient that called the group
}

func (d Device) String() string {
	return fmt.Sprintf("Device(ID: \"%v\", DeletedDateTime: \"%v\" AccountEnabled: \"%v\", ApproximateLastSignInDateTime: \"%v\", ComplianceExpirationDateTime: \"%v\", "+
		"DeviceId: \"%v\", DisplayName: \"%v\", IsCompliant: \"%v\" IsManaged: \"%v\", Manufacturer: \"%v\", MdmAppId: \"%v\", Model: \"%v\", OnPremisesLastSyncDateTime: \"%v\", "+
		"OnPremisesSyncEnabled: \"%v\", OperatingSystem: \"%v\", OperatingSystemVersion: \"%v\", PhysicalIds: \"%v\", ProfileType: \"%v\", RegistrationDateTime: \"%v\", "+
		"SystemLabels: \"%v\", TrustType: \"%v\", DirectAPIConnection: %v)",
		d.ID, d.DeletedDateTime, d.AccountEnabled, d.ApproximateLastSignInDateTime, d.ComplianceExpirationDateTime, d.DeviceId, d.DisplayName, d.IsCompliant,
		d.IsManaged, d.Manufacturer, d.MdmAppId, d.Model, d.OnPremisesLastSyncDateTime, d.OnPremisesSyncEnabled, d.OperatingSystem,
		d.OperatingSystemVersion, d.PhysicalIds, d.ProfileType, d.RegistrationDateTime, d.SystemLabels, d.TrustType, d.graphClient != nil)
}

// setGraphClient sets the graphClient instance in this instance and all child-instances (if any)
func (d *Device) setGraphClient(gC *GraphClient) {
	d.graphClient = gC
}

// ListMembers - Get a list of the domain's direct members.
func (d Device) ListMembers(opts ...ListQueryOption) (Devices, error) {
	if d.graphClient == nil {
		return nil, ErrNotGraphClientSourced
	}
	resource := fmt.Sprintf("/devices/%v/members", d.ID)

	var marsh struct {
		Devices Devices `json:"value"`
	}
	marsh.Devices.setGraphClient(d.graphClient)
	return marsh.Devices, d.graphClient.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (d *Device) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID                            string   `json:"id,omitempty"`
		DeletedDateTime               string   `json:"deletedDateTime,omitempty"`
		AccountEnabled                bool     `json:"accountEnabled,omitempty"`
		ApproximateLastSignInDateTime string   `json:"approximateLastSignInDateTime,omitempty"`
		ComplianceExpirationDateTime  string   `json:"complianceExpirationDateTime,omitempty"`
		DeviceId                      string   `json:"deviceId,omitempty"`
		DisplayName                   string   `json:"displayName,omitempty"`
		IsCompliant                   bool     `json:"isCompliant,omitempty"`
		IsManaged                     bool     `json:"isManaged,omitempty"`
		Manufacturer                  string   `json:"manufacturer,omitempty"`
		MdmAppId                      string   `json:"mdmAppId,omitempty"`
		Model                         string   `json:"model,omitempty"`
		OnPremisesLastSyncDateTime    string   `json:"onPremisesLastSyncDateTime,omitempty"`
		OnPremisesSyncEnabled         bool     `json:"onPremisesSyncEnabled,omitempty"`
		OperatingSystem               string   `json:"operatingSystem,omitempty"`
		OperatingSystemVersion        string   `json:"operatingSystemVersion,omitempty"`
		PhysicalIds                   []string `json:"physicalIds,omitempty"`
		ProfileType                   string   `json:"profileType,omitempty"`
		RegistrationDateTime          string   `json:"registrationDateTime,omitempty"`
		SystemLabels                  []string `json:"systemLabels,omitempty"`
		TrustType                     string   `json:"trustType,omitempty"`
		MemberOf                      []Member `json:"memberOf,omitempty"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	d.ID = tmp.ID
	d.DeletedDateTime, err = time.Parse(time.RFC3339, tmp.DeletedDateTime)
	if err != nil && tmp.DeletedDateTime != "" {
		return fmt.Errorf("cannot parse DeletedDateTime %v with RFC3339: %v", tmp.DeletedDateTime, err)
	}
	d.AccountEnabled = tmp.AccountEnabled
	d.ApproximateLastSignInDateTime, err = time.Parse(time.RFC3339, tmp.ApproximateLastSignInDateTime)
	if err != nil && tmp.ApproximateLastSignInDateTime != "" {
		return fmt.Errorf("cannot parse ApproximateLastSignInDateTime %v with RFC3339: %v", tmp.ApproximateLastSignInDateTime, err)
	}
	d.ComplianceExpirationDateTime, err = time.Parse(time.RFC3339, tmp.ComplianceExpirationDateTime)
	if err != nil && tmp.ComplianceExpirationDateTime != "" {
		return fmt.Errorf("cannot parse ApproximateLastSignInDateTime %v with RFC3339: %v", tmp.ComplianceExpirationDateTime, err)
	}
	d.DeviceId = tmp.DeviceId
	d.OnPremisesSyncEnabled = tmp.OnPremisesSyncEnabled
	d.DisplayName = tmp.DisplayName
	d.IsCompliant = tmp.IsCompliant
	d.IsManaged = tmp.IsManaged
	d.Manufacturer = tmp.Manufacturer
	d.MdmAppId = tmp.MdmAppId
	d.Model = tmp.Model
	d.OnPremisesLastSyncDateTime, err = time.Parse(time.RFC3339, tmp.OnPremisesLastSyncDateTime)
	if err != nil && tmp.OnPremisesLastSyncDateTime != "" {
		return fmt.Errorf("cannot parse OnPremisesLastSyncDateTime %v with RFC3339: %v", tmp.OnPremisesLastSyncDateTime, err)
	}
	d.OnPremisesSyncEnabled = tmp.OnPremisesSyncEnabled
	d.OperatingSystem = tmp.OperatingSystem
	d.OperatingSystemVersion = tmp.OperatingSystemVersion
	d.PhysicalIds = tmp.PhysicalIds
	d.ProfileType = tmp.ProfileType
	d.RegistrationDateTime, err = time.Parse(time.RFC3339, tmp.RegistrationDateTime)
	if err != nil && tmp.RegistrationDateTime != "" {
		return fmt.Errorf("cannot parse OnPremisesLastSyncDateTime %v with RFC3339: %v", tmp.RegistrationDateTime, err)
	}
	d.SystemLabels = tmp.SystemLabels
	d.TrustType = tmp.TrustType
	d.MemberOf = tmp.MemberOf
	return nil
}
