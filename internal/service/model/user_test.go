package model_test

import (
	"strings"
	"testing"

	"github.com/cerfical/merchshop/internal/service/model"
	"github.com/stretchr/testify/assert"
)

func TestNewUsername(t *testing.T) {
	tests := []struct {
		Name  string
		Input string
		Want  model.Username
		Err   assert.ErrorAssertionFunc
	}{
		{"ok", "test_user", "test_user", assert.NoError},

		{"reject_empty", "", "", func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "not") &&
				assert.ErrorContains(t, err, "empty")
		}},

		{"reject_specials", "test/user", "", func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "only") &&
				assert.ErrorContains(t, err, "allowed")
		}},

		{"reject_spaces", "test user", "", func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "only") &&
				assert.ErrorContains(t, err, "allowed")
		}},

		{"reject_too_short", "t", "", func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "at least")
		}},

		{"reject_too_long", strings.Repeat("test_user", 100), "", func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "no more")
		}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := model.NewUsername(test.Input)
			test.Err(t, err)

			assert.Equal(t, test.Want, got)
		})
	}
}

func TestNewPassword(t *testing.T) {
	tests := []struct {
		Name  string
		Input string
		Want  model.Password
		Err   assert.ErrorAssertionFunc
	}{
		{"ok", "test_password", "test_password", assert.NoError},

		{"reject_empty", "", "", func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "not") &&
				assert.ErrorContains(t, err, "empty")
		}},

		{"reject_non_printables", "test\npassword", "", func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "only") &&
				assert.ErrorContains(t, err, "allowed")
		}},

		{"reject_too_short", "t", "", func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "at least")
		}},

		{"reject_too_long", strings.Repeat("test_password", 100), "", func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "no more")
		}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := model.NewPassword(test.Input)
			test.Err(t, err)

			assert.Equal(t, test.Want, got)
		})
	}
}

func TestNewNumCoins(t *testing.T) {
	tests := []struct {
		Name  string
		Input int
		Want  model.NumCoins
		Err   assert.ErrorAssertionFunc
	}{
		{"ok", 9, 9, assert.NoError},

		{"reject_zero", 0, 0, func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "not") &&
				assert.ErrorContains(t, err, "zero")
		}},

		{"reject_negative", -9, 0, func(t assert.TestingT, err error, _ ...any) bool {
			return assert.ErrorContains(t, err, "not") &&
				assert.ErrorContains(t, err, "negative")
		}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := model.NewNumCoins(test.Input)
			test.Err(t, err)

			assert.Equal(t, test.Want, got)
		})
	}
}
