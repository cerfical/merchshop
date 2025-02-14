package model

const (
	ErrAuthFail     = Error("auth failed")
	ErrUserNotExist = Error("user doesn't exist")
	ErrUserExist    = Error("user already exists")

	ErrNotEnoughCoins    = Error("insufficient coins")
	ErrSenderNotExist    = Error("sender doesn't exist")
	ErrRecipientNotExist = Error("recipient doesn't exist")
	ErrSenderIsRecipient = Error("same sender and recipient")

	ErrMerchNotExist = Error("merch doesn't exist")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
