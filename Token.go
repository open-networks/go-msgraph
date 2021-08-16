package msgraph

import (
	"encoding/json"
	"fmt"
	"time"
)

// Token struct holds the Microsoft Graph API authentication token used by GraphClient to authenticate API-requests to the ms graph API
type Token struct {
	TokenType   string    // should always be "Bearer" for msgraph API-calls
	NotBefore   time.Time // time when the access token starts to be valid
	ExpiresOn   time.Time // time when the access token expires
	Resource    string    // will most likely be https://graph.microsoft.*, hence the Service Root Endpoint
	AccessToken string    // the access-token itself
}

func (t Token) String() string {
	return fmt.Sprintf("Token {TokenType: \"%v\", NotBefore: \"%v\", ExpiresOn: \"%v\", "+
		"Resource: \"%v\", AccessToken: \"%v\"}",
		t.TokenType, t.ExpiresOn, t.NotBefore, t.Resource, t.AccessToken)
}

// GetAccessToken teturns the API access token in Bearer format representation ready to send to the API interface.
func (t Token) GetAccessToken() string {
	return fmt.Sprintf("%v %v", t.TokenType, t.AccessToken)
}

// IsValid returns true if the token is already valid and is still valid. Otherwise false.
//
// Hint: this is a wrapper for >>token.IsAlreadyValid() && token.IsStillValid()<<
func (t Token) IsValid() bool {
	return t.IsAlreadyValid() && t.IsStillValid()
}

// IsAlreadyValid returns true if the token is already valid, hence the
// NotBefore is before the current time. Otherwise false.
//
// Hint: The current time is determined by time.Now()
func (t Token) IsAlreadyValid() bool {
	return time.Now().After(t.NotBefore)
}

// IsStillValid returns true if the token is still valid, hence the current time is before ExpiresOn.
// Does NOT check it the token is yet valid or in the future.
//
// Hint: The current time is determined by time.Now()
func (t Token) IsStillValid() bool {
	return time.Now().Before(t.ExpiresOn)
}

// HasExpired returns true if the token has already expired.
//
// Hint: this is a wrapper for >>!token.IsStillValid()<<
func (t Token) HasExpired() bool {
	return !t.IsStillValid()
}

// WantsToBeRefreshed returns true if the token is already invalid or close to
// expire (10 second before ExpiresOn), otherwise false. time.Now() is used to
// determine the current time.
func (t Token) WantsToBeRefreshed() bool {
	return !t.IsValid() || time.Now().After(t.ExpiresOn.Add(-10*time.Second))
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library.
//
// Hint: the UnmarshalJSON also checks immediately if the token is valid, hence
// the current time.Now() is after NotBefore and before ExpiresOn
func (t *Token) UnmarshalJSON(data []byte) error {
	tmp := struct {
		TokenType   string `json:"token_type"`        // should normally be "Bearer"
		ExpiresOn   int64  `json:"expires_on,string"` // = UNIX timestamp, parse to int64 immediately
		NotBefore   int64  `json:"not_before,string"` // = UNIX timestamp, parse to int64 immediately
		Resource    string `json:"resource"`          // will typically be https://graph.microsoft.com or wherever it came from
		AccessToken string `json:"access_token"`      // the actual access token - veeery long string
		//ExpiresIn   string `json:"expires_in"` // not used
	}{}

	// unmarshal to tmp-struct, return if error
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("err on json.Unmarshal: %v | Data: %v", err, string(data))
	}

	t.TokenType = tmp.TokenType
	t.ExpiresOn = time.Unix(tmp.ExpiresOn, 0)
	t.NotBefore = time.Unix(tmp.NotBefore, 0)
	t.Resource = tmp.Resource
	t.AccessToken = tmp.AccessToken

	if t.HasExpired() {
		return fmt.Errorf("Access-Token ExpiresOn %v is before current system-time %v", t.ExpiresOn, time.Now())
	}
	if !t.IsAlreadyValid() {
		return fmt.Errorf("Access-Token NotBefore %v is after current system-time %v", t.NotBefore, time.Now())
	}

	return nil
}
