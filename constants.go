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

// LoginBaseURL represents the basic url used to acquire a token for the msgraph api
const LoginBaseURL string = "https://login.microsoftonline.com"

// BaseURL represents the URL used to perform all ms graph API-calls
const BaseURL string = "https://graph.microsoft.com"

// APIVersion represents the APIVersion of msgraph used by this implementation
const APIVersion string = "v1.0"

// MaxPageSize is the maximum Page size for an API-call. This will be rewritten to use paging some day. Currently limits environments to 999 entries (e.g. Users, CalendarEvents etc.)
const MaxPageSize int = 999

var (
	// ErrFindUser is returned on any func that tries to find a user with the given parameters that can not be found
	ErrFindUser = errors.New("unable to find user")
)
