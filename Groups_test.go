package msgraph

import "testing"

func TestGroups_String(t *testing.T) {
	testGroup := GetTestGroup(t)
	groups := Groups{testGroup}
	wanted := "Groups(" + testGroup.String() + ")"
	if wanted != groups.String() {
		t.Errorf("Groups.String() result: %v, wanted: %v", testGroup, wanted)
	}
}
