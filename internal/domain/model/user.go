package model

type Username string

type UserID int

type Password string

type NumCoins int

type PasswordHash []byte

type User struct {
	ID UserID

	Username     Username
	PasswordHash PasswordHash

	Coins NumCoins
}

func NewUsername(s string) (Username, error) {
	// TODO: Add validation
	return Username(s), nil
}

func NewPassword(s string) (Password, error) {
	// TODO: Add validation
	return Password(s), nil
}

func NewNumCoins(n int) (NumCoins, error) {
	// TODO: Add validation
	return NumCoins(n), nil
}
