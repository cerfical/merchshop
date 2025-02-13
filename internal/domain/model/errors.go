package model

import "errors"

var (
	ErrAuthFail = errors.New("auth failed")
	ErrNotExist = errors.New("entity doesn't exist")
)
