package msgraph

// EmailAddress represents an emailAddress instance as microsoft.graph.EmailAddress. This is used at
// various positions, for example in CalendarEvents for attenees, owners, organizers or in Calendar
// for the owner.
//
// Short: The name and email address of a contact or message recipient.
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/resources/emailaddress
type EmailAddress struct {
	Address string `json:"address"` // The email address of the person or entity.
	Name    string `json:"name"`    // The display name of the person or entity.

	graphClient *GraphClient // the initiator of this EMailAddress Instance
}

func (e EmailAddress) String() string {
	return e.Name + "<" + e.Address + ">"
}

// setGraphClient sets the graphClient instance in this instance and all child-instances (if any)
func (e *EmailAddress) setGraphClient(graphClient *GraphClient) {
	e.graphClient = graphClient
}

// GetUser tries to get the real User-Instance directly from msgraph
// identified by the e-mail address of the user. This should normally
// be the userPrincipalName anyways. Returns an error if any from GraphClient.
func (e EmailAddress) GetUser() (User, error) {
	return e.graphClient.GetUser(e.Address)
}
