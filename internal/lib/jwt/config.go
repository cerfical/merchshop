package jwt

import "time"

type TokenConfig struct {
	Secret   []byte
	Lifetime time.Duration
}
