# Golang Microsoft Graph API implementation

[![Latest Release](https://img.shields.io/github/v/release/open-networks/go-msgraph)](https://github.com/open-networks/go-msgraph/releases)
[![Github Actions](https://github.com/open-networks/go-msgraph/actions/workflows/go.yml/badge.svg)](https://github.com/open-networks/go-msgraph/actions)
[![godoc](https://godoc.org/github.com/open-networks/go-msgraph?status.svg)](https://godoc.org/github.com/open-networks/go-msgraph)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-networks/go-msgraph)](https://goreportcard.com/report/github.com/open-networks/go-msgraph)
[![codebeat badge](https://codebeat.co/badges/9d93c0c6-a981-42d3-97a7-bb48c296257f)](https://codebeat.co/projects/github-com-open-networks-go-msgraph-master)
[![codecov](https://codecov.io/gh/open-networks/go-msgraph/branch/master/graph/badge.svg)](https://codecov.io/gh/open-networks/go-msgraph)
[![MIT License](https://img.shields.io/github/license/open-networks/go-msgraph)](LICENSE)

`go-msgraph` is an incomplete go lang implementation of the Microsoft Graph API. See [Overview of Microsoft Graph](https://docs.microsoft.com/en-us/graph/overview)

## General

This implementation has been written to get various user, group and calendar details out of a Microsoft Azure Active Directory and create / update the same.

## :warning: Deprecation warning

This code was created as part of a software project developed by Open Networks GmbH, Austria (the one in Europe :wink:) in the year 2017. The company has been bought by Bechtle AG, Germany in 2021 and the software project is planned to be replaced by end of 2022. Furthermore, the employee mainly working on this project also left the company by April 2022 and may only contribute in his leisure time.

Back in the days, there was no official support from Microsoft for Go lang. This support has been added in 2021, see [Issue #25](https://github.com/open-networks/go-msgraph/issues/25).

Therefore, I strongly advise you to use the new official implementation from Microsoft at [Golang MSGraph SDK by Microsoft](https://github.com/microsoftgraph/msgraph-sdk-go).

## Features

working & tested:

- list users, groups, calendars, calendarevents
- automatically grab & refresh token for API-access
- json-load the GraphClient struct & initialize it
- set timezone for full-day CalendarEvent
- use `$select`, `$search` and `$filter` when querying data
- `context`-aware API calls, can be cancelled.
- loading huge data sets with paging, thanks to PR #20 - [@Goorsky123](https://github.com/Goorsky123)

planned:

- String func that only prints non-empty values of an object, e.g. User
- add further support for mail, personal contacts (outlook), devices and apps, files etc. See [https://docs.microsoft.com/en-us/graph/overview](https://docs.microsoft.com/en-us/graph/overview)

## Example

To get your credentials to access the Microsoft Graph API visit: [Register an application with Azure AD and create a service principal](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#register-an-application-with-azure-ad-and-create-a-service-principal)

More examples can be found at the [docs](docs/). Here's a brief summary of some of the most common API-queries, ready to copy'n'paste:

````go
// initialize GraphClient manually
graphClient, err := msgraph.NewGraphClient("<TenantID>", "<ApplicationID>", "<ClientSecret>")
if err != nil {
    fmt.Println("Credentials are probably wrong or system time is not synced: ", err)
}

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

## Versioning & backwards compatibility

This project uses [Semantic versioning](https://semver.org/) with all tags prefixed with a `v`. Altough currently the case, I cannot promise to really keep everything backwards compatible for the 0.x version. If a 1.x version of this repository is ever released with enough API-calls implemented, I will keep this promise for sure. Any Breaking changes will be marked as such in the release notes of each release.

## Installation

I recommend to use [go modules](https://blog.golang.org/using-go-modules) and always use the latest tagged [release](https://github.com/open-networks/go-msgraph/releases). You may directly download the source code there, but the preffered way to install and update is with `go get`:

```shell
# Initially install
go get github.com/open-networks/go-msgraph
# Update
go get -u github.com/open-networks/go-msgraph
go mod tidy
```

## Documentation

There is some example code placed in the [docs/](docs/) folder. The code itself is pretty well documented with comments, hence see [http://godoc.org/github.com/open-networks/go-msgraph](http://godoc.org/github.com/open-networks/go-msgraph) or run:

```shell
godoc github.com/open-networks/go-msgraph
```

## License

[MIT](LICENSE)
