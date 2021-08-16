package msgraph

import (
	"errors"
	"time"
)

// FullDayEventTimeZone is used by CalendarEvent.UnmarshalJSON to set the timezone for full day events.
//
// That method json-unmarshal automatically sets the Begin/End Date to 00:00 with the correnct days then.
// This has to be done because Microsoft always sets the timezone to UTC for full day events. To work
// with that within your program is probably a bad idea, hence configure this as you need or
// probably even back to time.UTC
var FullDayEventTimeZone = time.Local

const (

	// Azure AD authentication endpoint "Global". Used to aquire a token for the ms graph API connection.
	//
	// Microsoft Documentation: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	AzureADAuthEndpointGlobal string = "https://login.microsoftonline.com"

	// Azure AD authentication endpoint "Germany". Used to aquire a token for the ms graph API connection.
	//
	// Microsoft Documentation: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	AzureADAuthEndpointGermany string = "https://login.microsoftonline.de"

	// Azure AD authentication endpoint "US Government". Used to aquire a token for the ms graph API connection.
	//
	// Microsoft Documentation: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	AzureADAuthEndpointUSGov string = "https://login.microsoftonline.us"

	// Azure AD authentication endpoint "China by 21 Vianet". Used to aquire a token for the ms graph API connection.
	//
	// Microsoft Documentation: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	AzureADAuthEndpointChina string = "https://login.partner.microsoftonline.cn"

	// ServiceRootEndpointGlobal represents the default Service Root Endpoint used to perform all ms graph
	// API-calls, hence the Service Root Endpoint.
	//
	// See https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
	ServiceRootEndpointGlobal string = "https://graph.microsoft.com"

	// Service Root Endpoint "US Government L4".
	//
	// See https://docs.microsoft.com/en-us/graph/deployments#microsoft-graph-and-graph-explorer-service-root-endpoints
	ServiceRootEndpointUSGovL4 string = "https://graph.microsoft.us"

	// Service Root Endpoint "US Government L5 (DOD)".
	//
	// See https://docs.microsoft.com/en-us/graph/deployments#microsoft-graph-and-graph-explorer-service-root-endpoints
	ServiceRootEndpointUSGovL5 string = "https://dod-graph.microsoft.us"

	// Service Root Endpoint "Germany".
	//
	// See https://docs.microsoft.com/en-us/graph/deployments#microsoft-graph-and-graph-explorer-service-root-endpoints
	ServiceRootEndpointGermany string = "https://graph.microsoft.de"

	// Service Root Endpoint "China operated by 21Vianet".
	//
	// See https://docs.microsoft.com/en-us/graph/deployments#microsoft-graph-and-graph-explorer-service-root-endpoints
	ServiceRootEndpointChina string = "https://microsoftgraph.chinacloudapi.cn"
)

// APIVersion represents the APIVersion of msgraph used by this implementation
const APIVersion string = "v1.0"

// MaxPageSize is the maximum Page size for an API-call. This will be rewritten to use paging some day. Currently limits environments to 999 entries (e.g. Users, CalendarEvents etc.)
const MaxPageSize int = 999

var (
	// ErrFindUser is returned on any func that tries to find a user with the given parameters that cannot be found
	ErrFindUser = errors.New("unable to find user")
	// ErrFindGroup is returned on any func that tries to find a group with the given parameters that cannot be found
	ErrFindGroup = errors.New("unable to find group")
	// ErrFindCalendar is returned on any func that tries to find a calendar with the given parameters that cannot be found
	ErrFindCalendar = errors.New("unable to find calendar")
	// ErrNotGraphClientSourced is returned if e.g. a ListMembers() is called but the Group has not been created by a graphClient query
	ErrNotGraphClientSourced = errors.New("instance is not created from a GraphClient API-Call, cannot directly get further information")
)
