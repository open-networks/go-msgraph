# Contributing to Atom

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines for contributing to go-msgraph. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## How Can I Contribute?

You can either report bugs, create feature requests as an issue so that the community can work on it, or you are very welcome to implement the code yourself and submit it as a pull-request.

## Suggested Development Environment

The program is as of 08/2021 developed with [Visual Studio Code](https://code.visualstudio.com/) with several plugins and [docker](https://docs.docker.com/). With that setup it's easy to update to the latest docker version and not have any golang version installed locally.

Install the following to use this:

1. [Visual Studio Code](https://code.visualstudio.com/) with the following Plugins:
    * [Go](https://marketplace.visualstudio.com/items?itemName=golang.Go)
    * [Markdown All in One](https://marketplace.visualstudio.com/items?itemName=yzhang.markdown-all-in-one)
    * [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
2. [Docker Engine](https://docs.docker.com/install)

The repository contains a `.devcontainer/devcontainer.json` that automatically configures this docker container used for development.

## Code Testing

If you want to run `go test` locally, you *must* set the following environment variables:

* `MSGraphTenantID`: Microsoft Graph API TenantID
* `MSGraphApplicationID`: Microsoft Graph Application ID
* `MSGraphClientSecret`: Microsoft Graph ClientSecret
* `MSGraphExistingGroupDisplayName`: The name of an existing Group in Azure Active Directory
* `MSGraphExistingGroupDisplayNameNumRes`: The number of results when searching for the above group. E.g. the group "users" will return three results upon searching, if the following groups exist: "users", "users-sales", "technicians-users-only"
* `MSGraphExistingUserPrincipalInGroup`: a `userPrincipalName` that is member of the group above, format: `alice@contoso.com`

Furthermore, the following environment variables are *optional*:
* `MSGraphAzureADAuthEndpoint`: Defaults to `msgraph.AzureADAuthEndpoint`, hence https://login.microsoftonline.com. Set this environment variable to use e.g. a US endpoint
* `MSGraphServiceRootEndpoint`: Defaults to `msgraph.ServiceRootEndpoint`, hence https://graph.microsoft.com. Set this environment variable to use e.g. a US endpoint
* `MSGraphExistingCalendarsOfUser`: Existing calendars of the above users, separated by a comma (`,`). This is optional, tests are skipped if not present. In case you dont use Office 365 / Mailboxes.

This may be done locally, or in `.devcontainer/devcontainer.json` if you use my suggested environment, but be careful not to add the changes to the commit (!):

````json
{
// ...
	"settings": {
		"go.toolsManagement.checkForUpdates": "local",
		"go.useLanguageServer": true,
		"go.gopath": "/go",
		"go.goroot": "/usr/local/go",
		"go.testEnvVars": {
            "MSGraphTenantID": "67dce6ac-xxxx-xxxx-xxxx-0807c45243a7",
            "MSGraphApplicationID": "1b99ac3b-xxxx-xxxx-xxxx-6f7998277091",
            "MSGraphClientSecret": "PZ.Wzfbxxxxxxxxxxxx2oe++TOid/YVG",
            "MSGraphExistingGroupDisplayName": "technicians",
            "MSGraphExistingGroupDisplayNameNumRes" : "1",
            "MSGraphExistingUserPrincipalInGroup": "alice@contoso.com",
            "MSGraphAzureADAuthEndpoint": "https://login.microsoftonline.com",
            "MSGraphServiceRootEndpoint": "https://graph.microsoft.com",
            "MSGraphExistingCalendarsOfUser" : "Calendar,Birthdays",
		},
	},
// ...
}
````

If you use Visual Studio Code without `Remote Contaienrs`, you may place it in `.vscode/settings.json`:

````json
{
    "go.testEnvVars": {
        "MSGraphTenantID": "67dce6ac-xxxx-xxxx-xxxx-0807c45243a7",
        "MSGraphApplicationID": "1b99ac3b-xxxx-xxxx-xxxx-6f7998277091",
        "MSGraphClientSecret": "PZ.Wzfbxxxxxxxxxxxx2oe++TOid/YVG",
        "MSGraphExistingGroupDisplayName": "technicians",
        "MSGraphExistingGroupDisplayNameNumRes" : "1",
        "MSGraphExistingUserPrincipalInGroup": "alice@contoso.com",
        "MSGraphExistingCalendarsOfUser" : "Calendar,Birthdays",
    },
}
````

## Code Syleguide

* code is formatted with `go fmt`
