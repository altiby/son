package domain

import "errors"

const (
	adminRoleID = "admin"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID         string
	Role       string
	FirstName  string
	SecondName string
	Birthdate  string
	Biography  string
	City       string
}

func (u User) IsAdmin() bool {
	return u.Role == adminRoleID
}
