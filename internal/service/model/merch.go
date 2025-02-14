package model

type MerchKind string

type MerchItem struct {
	Kind  MerchKind
	Price NumCoins
}

func NewMerchItem(s string) (*MerchItem, error) {
	if price, ok := merchKinds[s]; ok {
		return &MerchItem{
			Kind:  MerchKind(s),
			Price: price,
		}, nil
	}

	return nil, ErrMerchNotExist
}

var merchKinds = map[string]NumCoins{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}
