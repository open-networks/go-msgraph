package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

// get graph client config from environment
var (
	// Optional: Azure AD Authentication Endpoint, defaults to msgraph.AzureADAuthEndpointGlobal: https://login.microsoftonline.com
	msGraphAzureADAuthEndpoint string
	// Optional: Azure AD Authentication Endpoint, defaults to msgraph.ServiceRootEndpointGlobal: https://graph.microsoft.com
	msGraphServiceRootEndpoint string
	// Microsoft Graph tenant ID
	msGraphTenantID string
	// Microsoft Graph Application ID
	msGraphApplicationID string
	// Microsoft Graph Client Secret
	msGraphClientSecret string
	// a valid groupdisplayname from msgraph, e.g. technicians@contoso.com
	msGraphExistingGroupDisplayName string
	// a valid userprincipalname in the above group, e.g. felix@contoso.com
	msGraphExistingUserPrincipalInGroup string
	// valid calendar names that belong to the above user, seperated by a colon (","). e.g.: "Kalender,Feiertage in Österreich,Geburtstage"
	msGraphExistingCalendarsOfUser []string
	// the number of expected results when searching for the msGraphExistingGroupDisplayName with $search or $filter
	msGraphExistingGroupDisplayNameNumRes uint64
	// a domain-name for unit tests to create a user or other objects, e.g. contoso.com - omit the @
	msGraphDomainNameForCreateTests string
	// the graphclient used to perform all tests
	graphClient *GraphClient
	// marker if the calendar tests should be skipped - set if msGraphExistingCalendarsOfUser is empty
	skipCalendarTests bool
)

func getEnvOrPanic(key string) string {
	var val = os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Expected %s to be set, but is empty", key))
	}
	return val
}

func TestMain(m *testing.M) {
	msGraphTenantID = getEnvOrPanic("MSGraphTenantID")
	msGraphApplicationID = getEnvOrPanic("MSGraphApplicationID")
	msGraphClientSecret = getEnvOrPanic("MSGraphClientSecret")
	msGraphExistingGroupDisplayName = getEnvOrPanic("MSGraphExistingGroupDisplayName")
	msGraphExistingUserPrincipalInGroup = getEnvOrPanic("MSGraphExistingUserPrincipalInGroup")
	msGraphDomainNameForCreateTests = getEnvOrPanic("MSGraphDomainNameForCreateTests")

	if msGraphAzureADAuthEndpoint = os.Getenv("MSGraphAzureADAuthEndpoint"); msGraphAzureADAuthEndpoint == "" {
		msGraphAzureADAuthEndpoint = AzureADAuthEndpointGlobal
	}
	if msGraphServiceRootEndpoint = os.Getenv("MSGraphServiceRootEndpoint"); msGraphServiceRootEndpoint == "" {
		msGraphServiceRootEndpoint = ServiceRootEndpointGlobal
	}
	if msGraphExistingCalendarsOfUser = strings.Split(os.Getenv("MSGraphExistingCalendarsOfUser"), ","); msGraphExistingCalendarsOfUser[0] == "" {
		fmt.Println("Skipping calendar tests due to missing 'MSGraphExistingCalendarsOfUser' value")
		skipCalendarTests = true
	}

	var err error
	msGraphExistingGroupDisplayNameNumRes, err = strconv.ParseUint(os.Getenv("MSGraphExistingGroupDisplayNameNumRes"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Environment variable \"MSGraphExistingGroupDisplayNameNumRes\" seems to be invalid, cannot be parsed to unsigned integer: %v", err))
	}

	graphClient, err = NewGraphClientWithCustomEndpoint(msGraphTenantID, msGraphApplicationID, msGraphClientSecret, msGraphAzureADAuthEndpoint, msGraphServiceRootEndpoint)
	if err != nil {
		panic(fmt.Sprintf("Cannot initialize a new GraphClient, error: %v", err))
	}

	rand.Seed(time.Now().UnixNano())

	os.Exit(m.Run())
}

func randomString(n int) string {
	var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func createUnitTestUser(t *testing.T) User {
	t.Helper()
	rndstring := randomString(32)
	user, err := graphClient.CreateUser(User{
		AccountEnabled:    true,
		DisplayName:       "go-msgraph unit-test generated user " + time.Now().Format("2006-01-02") + " - random " + rndstring,
		MailNickname:      "go-msgraph.unit-test.generated." + rndstring,
		UserPrincipalName: "go-msgraph.unit-test.generated." + rndstring + "@" + msGraphDomainNameForCreateTests,
		PasswordProfile:   PasswordProfile{Password: randomString(32)},
	})
	if err != nil {
		t.Errorf("Cannot create a new User for unit tests: %v", err)
	}
	return user
}

func TestNewGraphClient(t *testing.T) {
	if msGraphAzureADAuthEndpoint != AzureADAuthEndpointGlobal || msGraphServiceRootEndpoint != ServiceRootEndpointGlobal {
		t.Skip("Skipping TestNewGraphClient because the endpoint is not the default - global - endpoint")
	}
	type args struct {
		tenantID      string
		applicationID string
		clientSecret  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "GraphClient from Environment-variables",
			args:    args{tenantID: msGraphTenantID, applicationID: msGraphApplicationID, clientSecret: msGraphClientSecret},
			wantErr: false,
		}, {
			name:    "GraphClient fail - wrong tenant ID",
			args:    args{tenantID: "wrong tenant id", applicationID: msGraphApplicationID, clientSecret: msGraphClientSecret},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong application ID",
			args:    args{tenantID: msGraphTenantID, applicationID: "wrong application id", clientSecret: msGraphClientSecret},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong client secret",
			args:    args{tenantID: msGraphTenantID, applicationID: msGraphApplicationID, clientSecret: "wrong client secret"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGraphClient(tt.args.tenantID, tt.args.applicationID, tt.args.clientSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGraphClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewGraphClientWithCustomEndpoint(t *testing.T) {
	type args struct {
		tenantID            string
		applicationID       string
		clientSecret        string
		azureADAuthEndpoint string
		serviceRootEndpoint string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "GraphClient from Environment-variables",
			args:    args{tenantID: msGraphTenantID, applicationID: msGraphApplicationID, clientSecret: msGraphClientSecret, azureADAuthEndpoint: msGraphAzureADAuthEndpoint, serviceRootEndpoint: msGraphServiceRootEndpoint},
			wantErr: false,
		}, {
			name:    "GraphClient fail - wrong tenant ID",
			args:    args{tenantID: "wrong tenant id", applicationID: msGraphApplicationID, clientSecret: msGraphClientSecret, azureADAuthEndpoint: msGraphAzureADAuthEndpoint, serviceRootEndpoint: msGraphServiceRootEndpoint},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong application ID",
			args:    args{tenantID: msGraphTenantID, applicationID: "wrong application id", clientSecret: msGraphClientSecret, azureADAuthEndpoint: msGraphAzureADAuthEndpoint, serviceRootEndpoint: msGraphServiceRootEndpoint},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong client secret",
			args:    args{tenantID: msGraphTenantID, applicationID: msGraphApplicationID, clientSecret: "wrong client secret", azureADAuthEndpoint: msGraphAzureADAuthEndpoint, serviceRootEndpoint: msGraphServiceRootEndpoint},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong Azure AD Authentication Endpoint",
			args:    args{tenantID: msGraphTenantID, applicationID: msGraphApplicationID, clientSecret: msGraphClientSecret, azureADAuthEndpoint: "completely invalid URL", serviceRootEndpoint: msGraphServiceRootEndpoint},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong Azure AD Authentication Endpoint",
			args:    args{tenantID: msGraphTenantID, applicationID: msGraphApplicationID, clientSecret: msGraphClientSecret, azureADAuthEndpoint: "https://iguess-this-does-not.exist.in.this.dksfowe3834myksweroqiwer.world", serviceRootEndpoint: msGraphServiceRootEndpoint},
			wantErr: true,
		}, {
			name:    "GraphClient fail - wrong Service Root Endpoint",
			args:    args{tenantID: msGraphTenantID, applicationID: msGraphApplicationID, clientSecret: msGraphClientSecret, azureADAuthEndpoint: msGraphAzureADAuthEndpoint, serviceRootEndpoint: "invalid URL"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGraphClientWithCustomEndpoint(tt.args.tenantID, tt.args.applicationID, tt.args.clientSecret, tt.args.azureADAuthEndpoint, tt.args.serviceRootEndpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGraphClientWithCustomEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGraphClient_ListUsers(t *testing.T) {
	tests := []struct {
		name    string
		g       *GraphClient
		want    User
		wantErr bool
	}{
		{
			name:    fmt.Sprintf("List all Users, check for user %v", msGraphExistingUserPrincipalInGroup),
			g:       graphClient,
			want:    User{UserPrincipalName: msGraphExistingUserPrincipalInGroup},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ListUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("GraphClient.ListUsers() len = 0, want more than 0: %v", got)
			}
			isGraphClientInitializd := true
			found := false
			for _, user := range got {
				isGraphClientInitializd = isGraphClientInitializd && user.graphClient != nil
				found = found || user.UserPrincipalName == tt.want.UserPrincipalName
			}
			if !found {
				t.Errorf("GraphClient.ListUsers() user %v not found, users: %v", tt.want.UserPrincipalName, got)
			}
			if !isGraphClientInitializd {
				t.Errorf("GraphClient.ListUsers() graphClient is nil, but was initialized from GraphClient")
			}
		})
	}
}

func TestGraphClient_ListGroups(t *testing.T) {
	tests := []struct {
		name    string
		g       *GraphClient
		opts    []ListQueryOption
		want    Group
		wantErr bool
	}{
		{
			name:    fmt.Sprintf("Test if Group %v is present", msGraphExistingGroupDisplayName),
			g:       graphClient,
			want:    Group{DisplayName: msGraphExistingGroupDisplayName},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ListGroups(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.ListGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			found := false
			isGraphClientInitialized := true
			for _, checkObj := range got {
				found = found || tt.want.DisplayName == checkObj.DisplayName
				isGraphClientInitialized = isGraphClientInitialized && checkObj.graphClient != nil
			}
			if !found {
				t.Errorf("GraphClient.ListGroups() = %v, searching for one of %v", got, tt.want)
			}
			if !isGraphClientInitialized {
				t.Errorf("GraphClient.ListGroups() graphClient is nil, but was initialized from GraphClient")
			}
		})
	}
}

func TestGraphClient_ListGroupsWithSelect(t *testing.T) {
	tests := []struct {
		name              string
		g                 *GraphClient
		opts              []ListQueryOption
		want              Group
		wantErr           bool
		wantZeroFields    []string
		wantNonZeroFields []string
	}{
		{
			name:    fmt.Sprintf("Test if Group %v is present and contains only specified fields", msGraphExistingGroupDisplayName),
			g:       graphClient,
			want:    Group{DisplayName: msGraphExistingGroupDisplayName},
			wantErr: false,
			opts: []ListQueryOption{
				ListWithSelect("displayName,createdDateTime"),
			},
			wantZeroFields:    []string{"ID"},
			wantNonZeroFields: []string{"DisplayName", "CreatedDateTime"},
		},
		{
			name:    fmt.Sprintf("Test if Group %v is present and contains only specified fields with context", msGraphExistingGroupDisplayName),
			g:       graphClient,
			want:    Group{DisplayName: msGraphExistingGroupDisplayName},
			wantErr: false,
			opts: []ListQueryOption{
				ListWithSelect("displayName,createdDateTime"),
				ListWithContext(context.Background()),
			},
			wantZeroFields:    []string{"ID"},
			wantNonZeroFields: []string{"DisplayName", "CreatedDateTime"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ListGroups(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.ListGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			found := false
			isGraphClientInitialized := true

			for _, checkObj := range got {
				found = found || tt.want.DisplayName == checkObj.DisplayName

				assertZeroFields(t, checkObj, tt.wantZeroFields, tt.wantNonZeroFields)
				isGraphClientInitialized = isGraphClientInitialized && checkObj.graphClient != nil
			}
			if !found {
				t.Errorf("GraphClient.ListGroups() = %v, searching for one of %v", got, tt.want)
			}

			if !isGraphClientInitialized {
				t.Errorf("GraphClient.ListGroups() graphClient is nil, but was initialized from GraphClient")
			}
		})
	}
}

func TestGraphClient_ListGroupsWithSearchAndFilter(t *testing.T) {
	tests := []struct {
		name    string
		g       *GraphClient
		opts    []ListQueryOption
		want    Group
		wantErr bool
	}{
		{
			name:    fmt.Sprintf("Test if Group %v is present when using searchQuery", msGraphExistingGroupDisplayName),
			g:       graphClient,
			want:    Group{DisplayName: msGraphExistingGroupDisplayName},
			wantErr: false,
			opts: []ListQueryOption{
				ListWithSearch(fmt.Sprintf(`"displayName:%s"`, msGraphExistingGroupDisplayName)),
				ListWithFilter(fmt.Sprintf("displayName eq '%s'", msGraphExistingGroupDisplayName)),
			},
		},
		{
			name:    fmt.Sprintf("Test if Group %v is present when using searchQuery with context", msGraphExistingGroupDisplayName),
			g:       graphClient,
			want:    Group{DisplayName: msGraphExistingGroupDisplayName},
			wantErr: false,
			opts: []ListQueryOption{
				ListWithSearch(fmt.Sprintf(`"displayName:%s"`, msGraphExistingGroupDisplayName)),
				ListWithFilter(fmt.Sprintf("displayName eq '%s'", msGraphExistingGroupDisplayName)),
				ListWithContext(context.Background()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ListGroups(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.ListGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			found := false
			isGraphClientInitialized := true

			if len(got) != int(msGraphExistingGroupDisplayNameNumRes) {
				t.Errorf("GraphClient.ListGroups(): Did not find expected number of results. Wanted: %d, got: %d", msGraphExistingGroupDisplayNameNumRes, len(got))
			}

			for _, checkObj := range got {
				found = found || tt.want.DisplayName == checkObj.DisplayName
				isGraphClientInitialized = isGraphClientInitialized && checkObj.graphClient != nil
			}
			if !found {
				t.Errorf("GraphClient.ListGroups() = %v, searching for one of %v", got, tt.want)
			}

			if !isGraphClientInitialized {
				t.Errorf("GraphClient.ListGroups() graphClient is nil, but was initialized from GraphClient")
			}
		})
	}
}

func TestGraphClient_ListGroupsWithSelectAndFilter(t *testing.T) {
	tests := []struct {
		name              string
		g                 *GraphClient
		opts              []ListQueryOption
		want              Group
		wantErr           bool
		wantZeroFields    []string
		wantNonZeroFields []string
	}{
		{
			name:    fmt.Sprintf("Test if Group %v is present when using searchQuery", msGraphExistingGroupDisplayName),
			g:       graphClient,
			want:    Group{DisplayName: msGraphExistingGroupDisplayName},
			wantErr: false,
			opts: []ListQueryOption{
				ListWithSelect("displayName,createdDateTime"),
				ListWithFilter(fmt.Sprintf("displayName eq '%s'", msGraphExistingGroupDisplayName)),
			},
			wantZeroFields:    []string{"ID"},
			wantNonZeroFields: []string{"DisplayName", "CreatedDateTime"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.ListGroups(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.ListGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			found := false
			isGraphClientInitialized := true

			if len(got) != int(msGraphExistingGroupDisplayNameNumRes) {
				t.Errorf("GraphClient.ListGroups(): Did not find expected number of results. Wanted: %d, got: %d", msGraphExistingGroupDisplayNameNumRes, len(got))
			}

			for _, checkObj := range got {
				found = found || tt.want.DisplayName == checkObj.DisplayName

				assertZeroFields(t, checkObj, tt.wantZeroFields, tt.wantNonZeroFields)
				isGraphClientInitialized = isGraphClientInitialized && checkObj.graphClient != nil
			}
			if !found {
				t.Errorf("GraphClient.ListGroups() = %v, searching for one of %v", got, tt.want)
			}

			if !isGraphClientInitialized {
				t.Errorf("GraphClient.ListGroups() graphClient is nil, but was initialized from GraphClient")
			}
		})
	}
}

func assertZeroFields(tb testing.TB, v interface{}, zeroFieldNames []string, nonZeroFieldNames []string) {
	tb.Helper()

	var (
		jsonBytes  []byte
		err        error
		mappedData = make(map[string]interface{})
	)

	if jsonBytes, err = json.Marshal(v); err != nil {
		tb.Fatalf("json.Marshal() error = %v", err)
	}

	if err = json.Unmarshal(jsonBytes, &mappedData); err != nil {
		tb.Fatalf("json.Unmarshal() error = %v", err)
	}

	for _, fieldName := range nonZeroFieldNames {
		if isZeroValue(v, fieldName, mappedData) {
			tb.Fatalf("Expected field %s to have non zero value", fieldName)
		}
	}

	for _, fieldName := range zeroFieldNames {
		if !isZeroValue(v, fieldName, mappedData) {
			tb.Fatalf("Expected field %s to have zero value but got %v", fieldName, mappedData[fieldName])
		}
	}
}

func isZeroValue(v interface{}, prop string, m map[string]interface{}) bool {
	// get value of 'v' if it's a reference
	underlying := reflect.Indirect(reflect.ValueOf(v))
	// if v is nil pointer return zero straight away
	if underlying.IsZero() {
		return true
	}

	// check if property has a IsZero() bool func e.g. for time.Time
	if zeroable, hasIsZero := underlying.FieldByName(prop).Interface().(interface{ IsZero() bool }); hasIsZero {
		return zeroable.IsZero()
	}

	return reflect.ValueOf(m[prop]).IsZero()
}

func TestGraphClient_GetUser(t *testing.T) {
	type args struct {
		identifier string
	}
	tests := []struct {
		name    string
		g       *GraphClient
		args    args
		want    User
		wantErr bool
	}{
		{
			name:    fmt.Sprintf("Test if user %v is present", msGraphExistingUserPrincipalInGroup),
			g:       graphClient,
			args:    args{identifier: msGraphExistingUserPrincipalInGroup},
			want:    User{UserPrincipalName: msGraphExistingUserPrincipalInGroup},
			wantErr: false,
		}, {
			name:    "Test if non-existing user produces err",
			g:       graphClient,
			args:    args{identifier: "ThisUserwillNotExistForSure@contoso.com"},
			want:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.GetUser(tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.UserPrincipalName != tt.want.UserPrincipalName {
				t.Errorf("GraphClient.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphClient_GetGroup(t *testing.T) {
	tests := []struct {
		name    string
		g       *GraphClient
		opts    []GetQueryOption
		want    Group
		wantErr bool
	}{
		{
			name:    fmt.Sprintf("Test if Group %v is present and GetGroup-able", msGraphExistingGroupDisplayName),
			g:       graphClient,
			want:    Group{DisplayName: msGraphExistingGroupDisplayName},
			wantErr: false,
		},
		{
			name:    fmt.Sprintf("Test if Group %v is present and GetGroup-able with context", msGraphExistingGroupDisplayName),
			g:       graphClient,
			opts:    []GetQueryOption{GetWithContext(context.Background())},
			want:    Group{DisplayName: msGraphExistingGroupDisplayName},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allGroups, err := tt.g.ListGroups()
			if err != nil { // check if groups can be listed
				t.Fatalf("GraphClient.ListGroups(): cannot list groups: %v", err)
			}
			targetGroup, err := allGroups.GetByDisplayName(tt.want.DisplayName)
			if err != nil { // check if the group to be tested is in the list
				t.Fatalf("Groups.GetByDisplayName(): cannot find group %v in %v, err: %v", tt.want.DisplayName, allGroups, err)
			}
			got, err := tt.g.GetGroup(targetGroup.ID) // actually execute the test we want to test
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.GetGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(got.DisplayName == tt.want.DisplayName) {
				t.Errorf("GraphClient.GetGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphClient_CreateAndDeleteUser(t *testing.T) {
	var rndstring = randomString(32)
	tests := []struct {
		name    string
		g       *GraphClient
		want    User
		wantErr bool
	}{
		{
			name: "Create new User",
			g:    graphClient,
			want: User{
				AccountEnabled:    true,
				DisplayName:       "go-msgraph unit-test generated user " + time.Now().Format("2006-01-02") + " - random " + rndstring,
				MailNickname:      "go-msgraph.unit-test.generated." + rndstring,
				UserPrincipalName: "go-msgraph.unit-test.generated." + rndstring + "@" + msGraphDomainNameForCreateTests,
				PasswordProfile:   PasswordProfile{Password: randomString(32)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test CreateUser
			got, err := tt.g.CreateUser(tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GraphClient.CreateUser() result: %v\n", got)
			// test DisableAccount
			err = got.DisableAccount()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.DisableAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
			// get user again to compare AccountEnabled field
			got, err = tt.g.GetUser(got.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got.AccountEnabled == true {
				t.Errorf("User.DisableAccount() did not work, AccountEnabled is still true")
			}
			// Delete user again
			err = got.DeleteUser()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetQueryOptions_Context(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		ctxSetup     func(tb testing.TB) context.Context
		wantDeadline bool
	}{
		{
			name:         "do not set a context explicitly expect background context",
			wantDeadline: false,
		},
		{
			name: "set background context",
			ctxSetup: func(testing.TB) context.Context {
				return context.Background()
			},
			wantDeadline: false,
		},
		{
			name: "set context with timeout",
			ctxSetup: func(tb testing.TB) context.Context {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				tb.Cleanup(cancel)
				return ctx
			},
			wantDeadline: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var opts []GetQueryOption
			if tt.ctxSetup != nil {
				ctx := tt.ctxSetup(t)
				opts = append(opts, GetWithContext(ctx))
			}

			var compiledOpts = compileGetQueryOptions(opts)
			var effectiveCtx = compiledOpts.Context()
			if _, ok := effectiveCtx.Deadline(); ok != tt.wantDeadline {
				t.Errorf("wantDeadline = %t but got %t", tt.wantDeadline, ok)
			}
		})
	}
}

func TestListQueryOptions_Context(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		ctxSetup     func(tb testing.TB) context.Context
		wantDeadline bool
	}{
		{
			name:         "do not set a context explicitly expect background context",
			wantDeadline: false,
		},
		{
			name: "set background context",
			ctxSetup: func(testing.TB) context.Context {
				return context.Background()
			},
			wantDeadline: false,
		},
		{
			name: "set context with timeout",
			ctxSetup: func(tb testing.TB) context.Context {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				tb.Cleanup(cancel)
				return ctx
			},
			wantDeadline: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var opts []ListQueryOption
			if tt.ctxSetup != nil {
				ctx := tt.ctxSetup(t)
				opts = append(opts, ListWithContext(ctx))
			}

			var compiledOpts = compileListQueryOptions(opts)
			var effectiveCtx = compiledOpts.Context()
			if _, ok := effectiveCtx.Deadline(); ok != tt.wantDeadline {
				t.Errorf("wantDeadline = %t but got %t", tt.wantDeadline, ok)
			}
		})
	}
}

func TestGetQueryOptions_Values(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		opts       []GetQueryOption
		wantValues string
	}{
		{
			name: "add $select",
			opts: []GetQueryOption{GetWithSelect("displayName")},
			wantValues: url.Values{
				"$select": []string{"displayName"},
			}.Encode(),
		},
		{
			name: "Select multiple values",
			opts: []GetQueryOption{GetWithSelect("displayName,createdDateTime")},
			wantValues: url.Values{
				"$select": []string{"displayName,createdDateTime"},
			}.Encode(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var compiledOpts = compileGetQueryOptions(tt.opts)
			if encodedValues := compiledOpts.Values().Encode(); tt.wantValues != encodedValues {
				unescapedWant, _ := url.PathUnescape(tt.wantValues)
				unescapedGot, _ := url.PathUnescape(encodedValues)

				t.Errorf("Expected values %s but got %s", unescapedWant, unescapedGot)
			}
		})
	}
}

func TestListQueryOptions_ValuesAndHeaders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		opts        []ListQueryOption
		wantValues  string
		wantHeaders map[string]string
	}{
		{
			name: "add $select",
			opts: []ListQueryOption{ListWithSelect("displayName")},
			wantValues: url.Values{
				"$select": []string{"displayName"},
			}.Encode(),
		},
		{
			name: "Select multiple values",
			opts: []ListQueryOption{ListWithSelect("displayName,createdDateTime")},
			wantValues: url.Values{
				"$select": []string{"displayName,createdDateTime"},
			}.Encode(),
		},
		{
			name: "Add $filter",
			opts: []ListQueryOption{ListWithFilter("displayName eq SomeGroupName")},
			wantValues: url.Values{
				"$filter": []string{"displayName eq SomeGroupName"},
			}.Encode(),
		},
		{
			name: "Add $search",
			opts: []ListQueryOption{ListWithSearch("displayName:hello")},
			wantValues: url.Values{
				"$search": []string{"displayName:hello"},
			}.Encode(),
			wantHeaders: map[string]string{
				"ConsistencyLevel": "eventual",
			},
		},
		{
			name: "Add $search and $filter",
			opts: []ListQueryOption{
				ListWithSearch("displayName:hello"),
				ListWithFilter("displayName eq 'hello world'"),
			},
			wantValues: url.Values{
				"$search": []string{"displayName:hello"},
				"$filter": []string{"displayName eq 'hello world'"},
			}.Encode(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var compiledOpts = compileListQueryOptions(tt.opts)
			if encodedValues := compiledOpts.Values().Encode(); tt.wantValues != encodedValues {
				unescapedWant, _ := url.PathUnescape(tt.wantValues)
				unescapedGot, _ := url.PathUnescape(encodedValues)
				t.Errorf("Expected values %s but got %s", unescapedWant, unescapedGot)
				return
			}
			if tt.wantHeaders != nil {
				var actualHeaders = compiledOpts.Headers()
				for key, wantValue := range tt.wantHeaders {
					if got := actualHeaders.Get(key); got != wantValue {
						t.Errorf("Expected %s for header %s but got %s", wantValue, key, got)
					}
				}
			}
		})
	}
}

func TestGraphClient_UnmarshalJSON(t *testing.T) {

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "All correct",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\", \"AzureADAuthEndpoint\": \"%v\",\"ServiceRootEndpoint\": \"%v\"}", msGraphTenantID, msGraphApplicationID, msGraphClientSecret, msGraphAzureADAuthEndpoint, msGraphServiceRootEndpoint))},
			wantErr: false,
		}, {
			name:    "JSON-syntax error",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\", \"AzureADAuthEndpoint\": \"%v\",\"ServiceRootEndpoint\": \"%v\"", msGraphTenantID, msGraphApplicationID, msGraphClientSecret, msGraphAzureADAuthEndpoint, msGraphServiceRootEndpoint))},
			wantErr: true,
		}, {
			name:    "TenantID incorrect",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\", \"AzureADAuthEndpoint\": \"%v\",\"ServiceRootEndpoint\": \"%v\"}", "wrongtenant", msGraphApplicationID, msGraphClientSecret, msGraphAzureADAuthEndpoint, msGraphServiceRootEndpoint))},
			wantErr: true,
		}, {
			name:    "TenantID empty",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\", \"AzureADAuthEndpoint\": \"%v\",\"ServiceRootEndpoint\": \"%v\"}", "", msGraphApplicationID, msGraphClientSecret, msGraphAzureADAuthEndpoint, msGraphServiceRootEndpoint))},
			wantErr: true,
		}, {
			name:    "ApplicationID incorrect",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\", \"AzureADAuthEndpoint\": \"%v\",\"ServiceRootEndpoint\": \"%v\"}", msGraphTenantID, "wrongapplication", msGraphClientSecret, msGraphAzureADAuthEndpoint, msGraphServiceRootEndpoint))},
			wantErr: true,
		}, {
			name:    "ApplicationID empty",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\", \"AzureADAuthEndpoint\": \"%v\",\"ServiceRootEndpoint\": \"%v\"}", msGraphTenantID, "", msGraphClientSecret, msGraphAzureADAuthEndpoint, msGraphServiceRootEndpoint))},
			wantErr: true,
		}, {
			name:    "ClientSecret incorrect",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\", \"AzureADAuthEndpoint\": \"%v\",\"ServiceRootEndpoint\": \"%v\"}", msGraphTenantID, msGraphApplicationID, "wrongclientsecret", msGraphAzureADAuthEndpoint, msGraphServiceRootEndpoint))},
			wantErr: true,
		}, {
			name:    "ClientSecret empty",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\", \"AzureADAuthEndpoint\": \"%v\",\"ServiceRootEndpoint\": \"%v\"}", msGraphTenantID, msGraphApplicationID, "", msGraphAzureADAuthEndpoint, msGraphServiceRootEndpoint))},
			wantErr: true,
		},
	}

	// only test empty endpoints if the endpoint is the default - global - endpoint. Otherwise the default values do not apply.
	if msGraphAzureADAuthEndpoint == AzureADAuthEndpointGlobal && msGraphServiceRootEndpoint == ServiceRootEndpointGlobal {
		tests = append(tests, struct {
			name    string
			args    args
			wantErr bool
		}{
			name:    "Empty AzureADAuthEndpoint",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\",\"ServiceRootEndpoint\": \"%v\"}", msGraphTenantID, msGraphApplicationID, msGraphClientSecret, msGraphServiceRootEndpoint))},
			wantErr: false,
		})
		tests = append(tests, struct {
			name    string
			args    args
			wantErr bool
		}{
			name:    "Empty ServiceRootEndpoint",
			args:    args{data: []byte(fmt.Sprintf("{\"TenantID\": \"%v\", \"ApplicationID\": \"%v\",\"ClientSecret\": \"%v\",\"AzureADAuthEndpoint\": \"%v\"}", msGraphTenantID, msGraphApplicationID, msGraphClientSecret, msGraphAzureADAuthEndpoint))},
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unmarshalTest GraphClient
			if err := unmarshalTest.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("GraphClient.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGraphClient_String(t *testing.T) {
	if fmt.Sprintf("GraphClient(TenantID: %v, ApplicationID: %v, ClientSecret: %v...%v, Token validity: [%v - %v])",
		graphClient.TenantID, graphClient.ApplicationID, graphClient.ClientSecret[0:3], graphClient.ClientSecret[len(graphClient.ClientSecret)-3:], graphClient.token.NotBefore, graphClient.token.ExpiresOn) != graphClient.String() {
		t.Errorf("GraphClient.String(): String function failed")
	}
}
