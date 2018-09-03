package msgraph

import (
	"fmt"
	"strings"
)

// Users represents multiple Users, used in JSON unmarshal
type Users []User

// GetUserByShortName returns the first User object that has the given shortName.
// Will return an error ErrFindUser if the user can not be found
func (u Users) GetUserByShortName(shortName string) (User, error) {
	for _, user := range u {
		if user.GetShortName() == shortName {
			return user, nil
		}
	}
	return User{}, ErrFindUser
}

// GetUserByActivePhone returns the User-instance whose activeNumber equals the given phone number.
// Will return an error ErrFindUser if the user can not be found
func (u Users) GetUserByActivePhone(activePhone string) (User, error) {
	for _, user := range u {
		if user.GetActivePhone() == activePhone {
			return user, nil
		}
	}
	return User{}, ErrFindUser
}

// GetUserByMail returns the User-instance that e-mail address matches the given e-mail addr.
// Will return an error ErrFindUser if the user can not be found.
func (u Users) GetUserByMail(email string) (User, error) {
	for _, user := range u {
		if user.Mail == email {
			return user, nil
		}
	}
	return User{}, ErrFindUser
}

func (u Users) String() string {
	var strs = make([]string, len(u))
	for i, user := range u {
		strs[i] = user.String()
	}
	return fmt.Sprintf("Users(%v)", strings.Join(strs, ", "))
}

// PrettySimpleString returns the whole []Users pretty simply formatted for logging purposes
func (u Users) PrettySimpleString() string {
	var strs = make([]string, len(u))
	for i, user := range u {
		strs[i] = user.PrettySimpleString()
	}
	return fmt.Sprintf("Users(%v)", strings.Join(strs, ", "))
}

// setGraphClient sets the GraphClient within that particular instance. Hence it's directly created by GraphClient
func (u Users) setGraphClient(gC *GraphClient) Users {
	for i := range u {
		u[i].setGraphClient(gC)
	}
	return u
}

// Equal compares the Users to the other Users and returns true
// if the two given Users are equal. Otherwise returns false
func (u Users) Equal(other Users) bool {
Outer:
	for _, user := range u {
		for _, toCompare := range other {
			if user.Equal(toCompare) {
				continue Outer
			}
		}
		return false
	}
	return len(u) == len(other) // if we reach this, all users have been found, now return if len of the users are equal
}
