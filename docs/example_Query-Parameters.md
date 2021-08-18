# Query Parameters

Support for the following three query parameters has been added:

* `$select` - only return the specified fields of the object. This reduces the used bandwidth and therefore improves performance
* `$search` - search with `ConsistencyLevel` set to `eventual`
* `$filter` - filter results server-side and only return matching results

See [Query Parameters Documentation](https://docs.microsoft.com/en-us/graph/query-parameters) from Microsoft.

They can be passed as a parameter to all `Get` and `List` functions with the following helper functions:

* `msgraph.GetWithSelect("displayName")`
* `msgraph.ListWithSelect("displayName,createdDateTime")`
* ``msgraph.ListWithSearch(`"displayName:alice"`)``
* `msgraph.ListWithFilter("displayName eq 'bob')`

## Example

````go
// initialize GraphClient
graphClient, err := msgraph.NewGraphClient("<TenantID>", "<ApplicationID>", "<ClientSecret>")
if err != nil {
    fmt.Println("Credentials are probably wrong or system time is not synced: ", err)
}

// Get only displayName of the user
user, err := graphClient.GetUser("alice@contoso.com", msgraph.GetWithSelect("displayName"))
// Get a list of users using search, that containing "alice" in their displayName:
searchUser := "alice"
users, err := graphClient.ListUsers(msgraph.ListWithSearch(fmt.Sprintf(`"displayName:%s"`, searchUser)))
searchUser = "Alice Contoso"
// Get a list of users using filter, that contain "alice" in their displayName:
users, err := graphClient.ListUsers(msgraph.ListWithFilter(fmt.Sprintf("displayName eq '%s'", searchUser)))

// The two above examples can also be combined, for example search and select:
users, err := graphClient.ListUsers(
	msgraph.ListWithSelect("displayName"),
	msgraph.ListWithSearch(fmt.Sprintf(`"displayName:%s"`, searchUser)),
)

// Last, but not least, a context can also be added to the query:
users, err := graphClient.ListUsers(
	msgraph.ListWithSelect("displayName"),
	msgraph.ListWithSearch(fmt.Sprintf("displayName:%s", searchUser)),
	msgraph.ListWithContext(ctx.Background()),
)
````
