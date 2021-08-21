# Context awareness

Context awareness has been implemented and can be passed to all `Get` and `List` queries of the `GraphClient`. For a general documentation about `context` see [context package documentation](https://pkg.go.dev/context).

Two functions have been created to implement context awareness:

* `msgraph.GetWithContext(<context>)`
* `msgraph.ListWithContext(<context>)`
* `msgraph.CreateWithContext(<context>)`
* `msgraph.PatchWithContext(<context>)`

The result of those two functions must be passed as parameter to the respective `Get` or `List` function.

## Example

````go
// initialize GraphClient
graphClient, err := msgraph.NewGraphClient("<TenantID>", "<ApplicationID>", "<ClientSecret>")
if err != nil {
    fmt.Println("Credentials are probably wrong or system time is not synced: ", err)
}
// create new context:
ctx, cancel := context.WithCancel(context.Background())
// example for Get-func:
user, err := graphClient.GetUser("dumpty@contoso.com", msgraph.GetWithContext(ctx))
// example for List-func:
users, err := graphClient.ListUsers(msgraph.ListWithContext(ctx))
````

Note: the use of a context is optional. If no context is given, the context `context.Background()` will automatically be used for all API-calls.