# Example GraphClient initialization

## Getting, listing, filtering users
 
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

````

## Create a user

````go
// example for Create-func:
user, err := graphClient.CreateUser(
    msgraph.User{
        AccountEnabled:    true,
        DisplayName:       "Rabbit",
        MailNickname:      "The rabbit",
        UserPrincipalName: "rabbit@contoso.com",
        PasswordProfile:   PasswordProfile{Password: "SecretCarrotBasedPassphrase"},
    },
    msgraph.CreateWithContext(ctx)
)
````

## Update a user

````go
// first, get the user:
user, err := graphClient.GetUser("rabbit@contoso.com")
// then create a user object, only set the fields you want to change.
err := user.UpdateUser(msgraph.User{DisplayName: "Rabbit 2.0"}, msgraph.UpdateWithContext(ctx))
// Hint 1: UpdateWithContext is optional
// Hint 2: you cannot disable a user that way, please user user.Disable
// Hint 3: after updating the account, you have to use GetUser("...") again.

// disable acccount
err := user.DisableAccount()
// enable user account again
err := user.UpdateUser(User{AccountEnabled: true})
// delete a user, use with caution!
err := user.DeleteUser()
````