package msgraph

import (
	"math/rand"
	"testing"
)

func Test_supportedTimeZones_GetTimeZoneByAlias(t *testing.T) {
	testuser := GetTestUser(t)
	timezones, _ := testuser.getTimeZoneChoices(compileGetQueryOptions(nil))

	randomTimezone := timezones.Value[rand.Intn(len(timezones.Value))]
	_, err := timezones.GetTimeZoneByAlias(randomTimezone.Alias)
	if err != nil {
		t.Errorf("Cannot get timeZone with Alias %v, err: %v", randomTimezone.Alias, err)
	}
	_, err = timezones.GetTimeZoneByAlias("This is a non existing timezone")
	if err == nil {
		t.Errorf("Tried to get a non existing timezone, expected an error, but got nil")
	}
}

func Test_supportedTimeZones_GetTimeZoneByDisplayName(t *testing.T) {
	testuser := GetTestUser(t)
	timezones, _ := testuser.getTimeZoneChoices(compileGetQueryOptions(nil))

	randomTimezone := timezones.Value[rand.Intn(len(timezones.Value))]
	_, err := timezones.GetTimeZoneByDisplayName(randomTimezone.DisplayName)
	if err != nil {
		t.Errorf("Cannot get timeZone with DisplayName %v, err: %v", randomTimezone.DisplayName, err)
	}
	_, err = timezones.GetTimeZoneByDisplayName("This is a non existing timezone")
	if err == nil {
		t.Errorf("Tried to get a non existing timezone, expected an error, but got nil")
	}
}
