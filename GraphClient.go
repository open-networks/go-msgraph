// Package msgraph is a go lang implementation of the Microsoft Graph API
//
// See: https://developer.microsoft.com/en-us/graph/docs/concepts/overview
package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	odataSearchParamKey = "$search"
	odataFilterParamKey = "$filter"
	odataSelectParamKey = "$select"
)

// GraphClient represents a msgraph API connection instance.
//
// An instance can also be json-unmarshalled and will immediately be initialized, hence a Token will be
// grabbed. If grabbing a token fails the JSON-Unmarshal returns an error.
type GraphClient struct {
	apiCall sync.Mutex // lock it when performing an API-call to synchronize it

	TenantID      string // See https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal#get-tenant-id
	ApplicationID string // See https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal#get-application-id-and-authentication-key
	ClientSecret  string // See https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal#get-application-id-and-authentication-key

	token Token // the current token to be used

	// azureADAuthEndpoint is used for this instance of GraphClient. For available endpoints see https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	azureADAuthEndpoint string
	// serviceRootEndpoint is the basic API-url used for this instance of GraphClient, namely Microsoft Graph service root endpoints. For available endpoints see https://docs.microsoft.com/en-us/graph/deployments#microsoft-graph-and-graph-explorer-service-root-endpoints.
	serviceRootEndpoint string
}

func (g *GraphClient) String() string {
	var firstPart, lastPart string
	if len(g.ClientSecret) > 4 { // if ClientSecret is not initialized prevent a panic slice out of bounds
		firstPart = g.ClientSecret[0:3]
		lastPart = g.ClientSecret[len(g.ClientSecret)-3:]
	}
	return fmt.Sprintf("GraphClient(TenantID: %v, ApplicationID: %v, ClientSecret: %v...%v, Token validity: [%v - %v])",
		g.TenantID, g.ApplicationID, firstPart, lastPart, g.token.NotBefore, g.token.ExpiresOn)
}

// NewGraphClient creates a new GraphClient instance with the given parameters
// and grabs a token. Returns an error if the token cannot be initialized. The
// default ms graph API global endpoint is used.
//
// This method does not have to be used to create a new GraphClient. If not used, the default global ms Graph API endpoint is used.
func NewGraphClient(tenantID, applicationID, clientSecret string) (*GraphClient, error) {
	return NewGraphClientWithCustomEndpoint(tenantID, applicationID, clientSecret, AzureADAuthEndpointGlobal, ServiceRootEndpointGlobal)
}

// NewGraphClientCustomEndpoint creates a new GraphClient instance with the
// given parameters and tries to get a valid token. All available public endpoints
// for azureADAuthEndpoint and serviceRootEndpoint are available via msgraph.azureADAuthEndpoint*  and msgraph.ServiceRootEndpoint*
//
// For available endpoints from Microsoft, see documentation:
//   * Authentication Endpoints: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
//   * Service Root Endpoints: https://docs.microsoft.com/en-us/graph/deployments#microsoft-graph-and-graph-explorer-service-root-endpoints
//
// Returns an error if the token cannot be initialized. This func does not have
// to be used to create a new GraphClient.
func NewGraphClientWithCustomEndpoint(tenantID, applicationID, clientSecret string, azureADAuthEndpoint string, serviceRootEndpoint string) (*GraphClient, error) {
	g := GraphClient{
		TenantID:            tenantID,
		ApplicationID:       applicationID,
		ClientSecret:        clientSecret,
		azureADAuthEndpoint: azureADAuthEndpoint,
		serviceRootEndpoint: serviceRootEndpoint,
	}
	g.apiCall.Lock()         // lock because we will refresh the token
	defer g.apiCall.Unlock() // unlock after token refresh
	return &g, g.refreshToken()
}

// makeSureURLsAreSet ensures that the two fields g.azureADAuthEndpoint and g.serviceRootEndpoint
// of the graphClient are set and therefore not empty. If they are currently empty
// they will be set to the constants AzureADAuthEndpointGlobal and ServiceRootEndpointGlobal.
func (g *GraphClient) makeSureURLsAreSet() {
	if g.azureADAuthEndpoint == "" { // If AzureADAuthEndpoint is not set, use the global endpoint
		g.azureADAuthEndpoint = AzureADAuthEndpointGlobal
	}
	if g.serviceRootEndpoint == "" { // If ServiceRootEndpoint is not set, use the global endpoint
		g.serviceRootEndpoint = ServiceRootEndpointGlobal
	}
}

// refreshToken refreshes the current Token. Grabs a new one and saves it within the GraphClient instance
func (g *GraphClient) refreshToken() error {
	g.makeSureURLsAreSet()
	if g.TenantID == "" {
		return fmt.Errorf("tenant ID is empty")
	}
	resource := fmt.Sprintf("/%v/oauth2/token", g.TenantID)
	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("client_id", g.ApplicationID)
	data.Add("client_secret", g.ClientSecret)
	data.Add("resource", g.serviceRootEndpoint)

	u, err := url.ParseRequestURI(g.azureADAuthEndpoint)
	if err != nil {
		return fmt.Errorf("unable to parse URI: %v", err)
	}

	u.Path = resource
	req, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(data.Encode()))

	if err != nil {
		return fmt.Errorf("HTTP Request Error: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	var newToken Token
	err = g.performRequest(req, &newToken) // perform the prepared request
	if err != nil {
		return fmt.Errorf("error on getting msgraph Token: %v", err)
	}
	g.token = newToken
	return err
}

// makeGETAPICall performs an API-Call to the msgraph API.
func (g *GraphClient) makeGETAPICall(apiCall string, reqParams getRequestParams, v interface{}) error {
	return g.makeAPICall(apiCall, http.MethodGet, reqParams, nil, v)
}

// makeGETAPICall performs an API-Call to the msgraph API.
func (g *GraphClient) makePOSTAPICall(apiCall string, reqParams getRequestParams, body io.Reader, v interface{}) error {
	return g.makeAPICall(apiCall, http.MethodPost, reqParams, body, v)
}

// makePATCHAPICall performs an API-Call to the msgraph API.
func (g *GraphClient) makePATCHAPICall(apiCall string, reqParams getRequestParams, body io.Reader, v interface{}) error {
	return g.makeAPICall(apiCall, http.MethodPatch, reqParams, body, v)
}

// makeDELETEAPICall performs an API-Call to the msgraph API.
func (g *GraphClient) makeDELETEAPICall(apiCall string, reqParams getRequestParams, v interface{}) error {
	return g.makeAPICall(apiCall, http.MethodDelete, reqParams, nil, v)
}

// makeAPICall performs an API-Call to the msgraph API. This func uses sync.Mutex to synchronize all API-calls.
//
// Parameter httpMethod may be http.MethodGet, http.MethodPost or http.MethodPatch
//
// Parameter body may be nil to not provide any content - e.g. when using a http GET request.
func (g *GraphClient) makeAPICall(apiCall string, httpMethod string, reqParams getRequestParams, body io.Reader, v interface{}) error {
	g.makeSureURLsAreSet()
	g.apiCall.Lock()
	defer g.apiCall.Unlock() // unlock when the func returns
	// Check token
	if g.token.WantsToBeRefreshed() { // Token not valid anymore?
		err := g.refreshToken()
		if err != nil {
			return err
		}
	}

	reqURL, err := url.ParseRequestURI(g.serviceRootEndpoint)
	if err != nil {
		return fmt.Errorf("unable to parse URI %v: %v", g.serviceRootEndpoint, err)
	}

	// Add Version to API-Call, the leading slash is always added by the calling func
	reqURL.Path = "/" + APIVersion + apiCall

	req, err := http.NewRequestWithContext(reqParams.Context(), httpMethod, reqURL.String(), body)
	if err != nil {
		return fmt.Errorf("HTTP request error: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", g.token.GetAccessToken())

	for key, vals := range reqParams.Headers() {
		for idx := range vals {
			req.Header.Add(key, vals[idx])
		}
	}

	var getParams = reqParams.Values()

	if httpMethod == http.MethodGet {
		// TODO: Improve performance with using $skip & paging instead of retrieving all results with $top
		// TODO: MaxPageSize is currently 999, if there are any time more than 999 entries this will make the program unpredictable... hence start to use paging (!)
		getParams.Add("$top", strconv.Itoa(MaxPageSize))
	}
	req.URL.RawQuery = getParams.Encode() // set query parameters

	return g.performRequest(req, v)
}

// performRequest performs a pre-prepared http.Request and does the proper error-handling for it.
// does a json.Unmarshal into the v interface{} and returns the error of it if everything went well so far.
func (g *GraphClient) performRequest(req *http.Request, v interface{}) error {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP response error: %v of http.Request: %v", err, req.URL)
	}
	defer resp.Body.Close() // close body when func returns

	body, err := ioutil.ReadAll(resp.Body) // read body first to append it to the error (if any)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// Hint: this will mostly be the case if the tenant ID cannot be found, the Application ID cannot be found or the clientSecret is incorrect.
		// The cause will be described in the body, hence we have to return the body too for proper error-analysis
		return fmt.Errorf("StatusCode is not OK: %v. Body: %v ", resp.StatusCode, string(body))
	}

	//fmt.Println("Body: ", string(body))

	if err != nil {
		return fmt.Errorf("HTTP response read error: %v of http.Request: %v", err, req.URL)
	}

	// no content returned when http PATCH or DELETE is used, e.g. User.DeleteUser()
	if req.Method == http.MethodDelete || req.Method == http.MethodPatch {
		return nil
	}
	return json.Unmarshal(body, &v) // return the error of the json unmarshal
}

// ListUsers returns a list of all users
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user_list
func (g *GraphClient) ListUsers(opts ...ListQueryOption) (Users, error) {
	resource := "/users"
	var marsh struct {
		Users Users `json:"value"`
	}
	err := g.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	marsh.Users.setGraphClient(g)
	return marsh.Users, err
}

// ListGroups returns a list of all groups
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/group_list
func (g *GraphClient) ListGroups(opts ...ListQueryOption) (Groups, error) {
	resource := "/groups"

	var reqParams = compileListQueryOptions(opts)

	var marsh struct {
		Groups Groups `json:"value"`
	}
	err := g.makeGETAPICall(resource, reqParams, &marsh)
	marsh.Groups.setGraphClient(g)
	return marsh.Groups, err
}

// GetUser returns the user object associated to the given user identified by either
// the given ID or userPrincipalName
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user_get
func (g *GraphClient) GetUser(identifier string, opts ...GetQueryOption) (User, error) {
	resource := fmt.Sprintf("/users/%v", identifier)
	user := User{graphClient: g}
	err := g.makeGETAPICall(resource, compileGetQueryOptions(opts), &user)
	return user, err
}

// GetGroup returns the group object identified by the given groupID.
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/group_get
func (g *GraphClient) GetGroup(groupID string, opts ...GetQueryOption) (Group, error) {
	resource := fmt.Sprintf("/groups/%v", groupID)
	group := Group{graphClient: g}
	err := g.makeGETAPICall(resource, compileGetQueryOptions(opts), &group)
	return group, err
}

// CreateUser creates a new user given a user object and returns and updated object
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user-post-users
func (g *GraphClient) CreateUser(userInput User, opts ...CreateQueryOption) (User, error) {
	user := User{graphClient: g}
	bodyBytes, err := json.Marshal(userInput)
	if err != nil {
		return user, err
	}

	reader := bytes.NewReader(bodyBytes)
	err = g.makePOSTAPICall("/users", compileCreateQueryOptions(opts), reader, &user)

	return user, err
}

// SendMail sends a message.
//
// See https://docs.microsoft.com/en-us/graph/api/user-sendmail
func (g *GraphClient) SendMail(userID string, message Message, saveToSentItems bool) error {
	if g == nil {
		return ErrNotGraphClientSourced
	}

	bodyBytes, err := json.Marshal(struct {
		Message         Message `json:"message,omitempty"`
		SaveToSentItems bool    `json:"saveToSentItems,omitempty"`
	}{
		message,
		saveToSentItems,
	})
	if err != nil {
		return fmt.Errorf("could not marshal message body: %w", err)
	}

	reader := bytes.NewReader(bodyBytes)

	resource := fmt.Sprintf("/users/%v/sendMail", userID)

	return g.makePOSTAPICall(resource, compileCreateQueryOptions(nil), reader, nil)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library.
// This method additionally to loading the TenantID, ApplicationID and ClientSecret
// immediately gets a Token from msgraph (hence initialize this GraphAPI instance)
// and returns an error if any of the data provided is incorrect or the token cannot be acquired
func (g *GraphClient) UnmarshalJSON(data []byte) error {
	tmp := struct {
		TenantID            string
		ApplicationID       string
		ClientSecret        string
		AzureADAuthEndpoint string
		ServiceRootEndpoint string
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	g.TenantID = tmp.TenantID
	if g.TenantID == "" {
		return fmt.Errorf("TenantID is empty")
	}
	g.ApplicationID = tmp.ApplicationID
	if g.ApplicationID == "" {
		return fmt.Errorf("ApplicationID is empty")
	}
	g.ClientSecret = tmp.ClientSecret
	if g.ClientSecret == "" {
		return fmt.Errorf("ClientSecret is empty")
	}
	g.azureADAuthEndpoint = tmp.AzureADAuthEndpoint
	g.serviceRootEndpoint = tmp.ServiceRootEndpoint
	g.makeSureURLsAreSet()

	// get a token and return the error (if any)
	err = g.refreshToken()
	if err != nil {
		return fmt.Errorf("can't get Token: %v", err)
	}
	return nil
}
