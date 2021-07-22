# Golang Microsoft Graph API implementation
[![godoc](https://godoc.org/github.com/open-networks/go-msgraph?status.svg)](https://godoc.org/github.com/open-networks/go-msgraph)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-networks/go-msgraph)](https://goreportcard.com/report/github.com/open-networks/go-msgraph)
[![codebeat badge](https://codebeat.co/badges/9d93c0c6-a981-42d3-97a7-bb48c296257f)](https://codebeat.co/projects/github-com-open-networks-go-msgraph-master)
[![codecov](https://codecov.io/gh/open-networks/go-msgraph/branch/master/graph/badge.svg)](https://codecov.io/gh/open-networks/go-msgraph)

go-msgraph is a go lang implementation of the Microsoft Graph API. See https://developer.microsoft.com/en-us/graph/docs/concepts/overview

## General
This implementation has been written to get various user, group and calendar details out of a Microsoft Azure Active Directory. Currently only READ-access is implemented, but you are welcome to add WRITE-support to it & backmerge it

## Features
working & tested:
- list users, groups, calendars, calendarevents
- automatically grab & refresh token for API-access
- json-load the GraphClient struct & initialize it
- set timezone for full-day CalendarEvent

in progress:
- implement paging to load huge data-sets, currently limitted to one page, 999 entries

planned:
- add further support for mail, personal contacts (outlook), devices and apps, files etc. See https://developer.microsoft.com/en-us/graph/docs/concepts/v1-overview

## Example
To get your credentials visit: https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal
````go
// initialize GraphClient manually
graphClient, err := msgraph.NewGraphClient("<TenantID>", "<ApplicationID>", "<ClientSecret>")
if err != nil {
    fmt.Println("Credentials are probably wrong or system time is not synced: ", err)
}

// initialize the GraphClient via JSON-Load. Please do proper error-handling (!)
// Specify JSON-Fields TenantID, ApplicationID and ClientSecret 
fileContents, err := ioutil.ReadFile("./msgraph-credentials.json")
var graphClient msgraph.GraphClient
err = json.Unmarshal(fileContents, &graphClient)

// List all users
users, err := graphClient.ListUsers()
// Gets all the detailled information about a user identified by it's ID or userPrincipalName
user, err := graphClient.GetUser("humpty@contoso.com") 
// List all groups
groups, err := graphClient.ListGroups()
// List all members of a group.
groupMembers, err := groups[0].ListMembers()
// Lists all Calendars of a user
calendars, err := user.ListCalendars()

// Let all full-day calendar events that are loaded from ms graph be set to timezone Europe/Vienna:
// Standard is time.Local
msgraph.FullDayEventTimeZone, _ = time.LoadLocation("Europe/Vienna")

// Lists all CalendarEvents of the given userPrincipalName/ID that starts/ends within the the next 7 days
startTime := time.Now()
endTime := time.Now().Add(time.Hour * 24 * 7)
events, err := graphClient.ListCalendarView("alice@contoso.com", startTime, endTime)
````

## Installation

### Using *go get*

    $ go get github.com/open-networks/go-msgraph

You can use `go get -u` to update the package.

## Documentation

For docs, see http://godoc.org/github.com/open-networks/go-msgraph or run:

    $ godoc github.com/open-networks/go-msgraph



