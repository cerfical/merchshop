package model

import "errors"

var (
	ErrAuthFail       = errors.New("auth failed")
	ErrUserNotExist   = errors.New("user doesn't exist")
)
