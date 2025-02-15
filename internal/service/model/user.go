package model

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
