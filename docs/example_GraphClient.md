# Example GraphClient initialization

To get your credentials to access the Microsoft Graph API visit: [Register an application with Azure AD and create a service principal](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#register-an-application-with-azure-ad-and-create-a-service-principal)

## Manual initialization of a GraphClient

````go
// initialize GraphClient manually with Global Authentication and Service endpoint
graphClient, err := msgraph.NewGraphClient("<TenantID>", "<ApplicationID>", "<ClientSecret>")
if err != nil {
    fmt.Println("Credentials are probably wrong or system time is not synced: ", err)
}
// initialize GraphClient manually with US Government L4 endpoint.
graphClient, err := msgraph.NewGraphClient("<TenantID>", "<ApplicationID>", "<ClientSecret>", msgraph.AzureADauthEndpointUSGov, msgraph.ServiceRootEndpointUSGovL4)
if err != nil {
    fmt.Println("Credentials are probably wrong or system time is not synced: ", err)
}
````

All of the available endpoints are created as `const` variables:

* `msgraph.AzureADauthEndpoint<Global,USGov,China,Germany>`
* `msgraph.ServiceRootEndpoint<Global,USGovL4,USGovL5,China,Germany>`

The Microsoft documentation for all available service endpoints can be found here:

* Azure AD authentication endpoints: https://docs.microsoft.com/en-us/azure/active-directory/develop/authentication-national-cloud#azure-ad-authentication-endpoints
* Serivce Root Endpoints: https://docs.microsoft.com/en-us/graph/deployments#microsoft-graph-and-graph-explorer-service-root-endpoints.


## JSON initialize the Graphclient

The GraphClient can be initilized directly via a JSON-file, also nested in other objects. The GraphClient will immediately initialize upon `json.Unmarshal`, and therefore check if the credentials are valid and a valid token can be aquired. If this fails, the `json.Unmarshal` will return an error.

example contents of the json-file `./msgraph-credentials.json`:
````json
{
  "TenantID": "67dce6ac-xxxx-xxxx-xxxx-0807c45243a7",
  "ApplicationID": "1b99ac3b-xxxx-xxxx-xxxx-6f7998277091",
  "ClientSecret": "PZ.Wzfbxxxxxxxxxxxx2oe++TOid/YVG",
  "AzureADAuthEndpoint": "https://login.microsoftonline.com", // this is optional
  "ServiceRootEndpoint" : "https://graph.microsoft.com" // this is optional
}
````

*Hint*: `AzureADAuthEndpoint` and `ServiceRootEndpoint` are optional and default to the two `Global` endpoints: `msgraph.AzureADAuthEndpointGlobal` and `msgraph.ServiceRootEndpointGlobal`

Example to initialize the `GraphClient` with the json file:

````go
// initialize the GraphClient via JSON-Load. Please do proper error-handling (!)
// Specify JSON-Fields TenantID, ApplicationID and ClientSecret
fileContents, err := ioutil.ReadFile("./msgraph-credentials.json")
if (err != nil) {
    fmt.println("Cannot read ./msgraph-credentials.json: ", err)
}
var graphClient msgraph.GraphClient
err = json.Unmarshal(fileContents, &graphClient)
if (err != nil) {
    fmt.println("Could not initialize GraphClient from JSON: ", err)
}
````

## Other options

I could think about an initialization directly with a `yaml` file, or via enviroment variables. If you need this in your code, please feel free to implement it and open a pull-request.
