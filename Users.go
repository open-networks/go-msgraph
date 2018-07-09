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
	var users []string
	for _, user := range u {
		users = append(users, user.String())
	}
	return fmt.Sprintf("Users(%v)", strings.Join(users, ", "))
}

// PrettySimpleString returns the whole []Users pretty simply formatted for logging purposes
func (u Users) PrettySimpleString() string {
	var users []string
	for _, user := range u {
		users = append(users, user.String())
	}
	return strings.Join(users, ", ")
}
