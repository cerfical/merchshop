package model

const (
	ErrAuthFail     = Error("auth failed")
	ErrUserNotExist = Error("user doesn't exist")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
