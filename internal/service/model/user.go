package model

import (
	"fmt"
	"regexp"
)

type Username string

type UserID int

type Password string

type NumCoins int

type NumItems int

type PasswordHash []byte

type Deposit struct {
	Amount NumCoins
	From   Username
}

type Withdrawal struct {
	Amount NumCoins
	To     Username
}

type InventoryItem struct {
	Merch    MerchKind
	Quantity NumItems
}

type User struct {
	ID UserID

	Username     Username
	PasswordHash PasswordHash

	Inventory   []InventoryItem
	Withdrawals []Withdrawal
	Deposits    []Deposit

	Coins NumCoins
}

// TODO: Hardcoded validation parameters?

var userRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func NewUsername(s string) (Username, error) {
	if err := validateString(s); err != nil {
		return "", err
	}

	if !userRegex.MatchString(s) {
		return "", Error("only letters, digits and underscores are allowed")
	}
	return Username(s), nil
}

var passwdRegex = regexp.MustCompile(`^[[:print:]]+$`)

func NewPassword(s string) (Password, error) {
	if err := validateString(s); err != nil {
		return "", err
	}

	if !passwdRegex.MatchString(s) {
		return "", Error("only printable characters are allowed")
	}
	return Password(s), nil
}

func validateString(s string) error {
	const minLen = 8
	const maxLen = 128

	switch {
	case s == "":
		return Error("must not be empty")
	case len(s) < minLen:
		return Error(fmt.Sprintf("must be at least %d characters long", minLen))
	case len(s) > maxLen:
		return Error(fmt.Sprintf("must be no more than %d characters long", maxLen))
	default:
		return nil
	}
}

func NewNumCoins(n int) (NumCoins, error) {
	switch {
	case n == 0:
		return 0, Error("must not be zero")
	case n < 0:
		return 0, Error("must not be negative")
	default:
		return NumCoins(n), nil
	}
}
