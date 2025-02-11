package api

import "time"

// AuthConfig configures API authentication.
type AuthConfig struct {
	// Token configures the API token generation process.
	Token struct {
		// Secret contains a secret to be used in signing of JWT tokens.
		Secret []byte

		// Lifetime optionally specifies the time period after which the token will become 'expired', i.e., invalid.
		Lifetime time.Duration
	}
}
