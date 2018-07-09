package msgraph

import (
	"fmt"
	"strings"
)

// Users represents multiple Users, used in JSON unmarshal
type Users struct {
	Users []User `json:"value"`
}

// GetUserByShortName returns the first User object that has the given shortName.
// Will return an error ErrFindUser if the user can not be found
func (u *Users) GetUserByShortName(shortName string) (*User, error) {
	for _, user := range u.Users {
		if user.GetShortName() == shortName {
			return &user, nil
		}
	}

	return nil, ErrFindUser
}

// GetUserByActivePhone returns the User-instance whose activeNumber equals the given phone number.
// Will return an error ErrFindUser if the user can not be found
func (u *Users) GetUserByActivePhone(activePhone string) (*User, error) {
	for _, user := range u.Users {
		if user.GetActivePhone() == activePhone {
			return &user, nil
		}
	}

	return nil, ErrFindUser
}

// GetUserByMail returns the User-instance that e-mail address matches the given e-mail addr.
// Will return an error ErrFindUser if the user can not be found.
func (u *Users) GetUserByMail(email string) (*User, error) {
	for _, user := range u.Users {
		if user.Mail == email {
			return &user, nil
		}
	}

	return nil, ErrFindUser
}

func (u *Users) String() string {
	return fmt.Sprintf("Users(%v)", u.Users)
}

// PrettySimpleString returns the whole []Users pretty simply formatted for logging purposes
func (u *Users) PrettySimpleString() string {
	var users []string
	for _, user := range u.Users {
		users = append(users, user.PrettySimpleString())
	}
	return strings.Join(users, ", ")
}
