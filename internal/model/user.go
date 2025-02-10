package model

// UserCreds specifies user credentials required to authenticate the user.
type UserCreds struct {
	Name     string
	Password string
}

// User identifies a single user.
type User struct {
	ID       int
	Name     string
	Password string
}

// UserStore provide storage to store user credentials.
type UserStore interface {
	// GetUser retrieves a user from the store or, if it doesn't exist, creates one with the provided credentials.
	GetUser(u *UserCreds) (*User, error)
}
