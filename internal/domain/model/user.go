package model

func NewUserCreds(username, passwd string) (UserCreds, error) {
	un, err := NewUsername(username)
	if err != nil {
		return UserCreds{}, err
	}

	pd, err := NewPassword(passwd)
	if err != nil {
		return UserCreds{}, err
	}

	return UserCreds{
		Username: un,
		Password: pd,
	}, nil
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

type User struct {
	ID UserID

	Username     Username
	PasswordHash Hash

	Coins NumCoins
}

type UserCreds struct {
	Username Username
	Password Password
}

type Username string

type UserID int

type Password string

type NumCoins int

type Hash []byte
