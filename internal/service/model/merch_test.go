package model_test

import (
	"testing"

	"github.com/cerfical/merchshop/internal/service/model"
	"github.com/stretchr/testify/assert"
)

func TestNewMerchItem(t *testing.T) {
	tests := []struct {
		Name  string
		Input string
		Want  model.MerchKind
		Err   assert.ErrorAssertionFunc
	}{
		{"ok", "t-shirt", "t-shirt", assert.NoError},
		{"invalid_merch_item", "T-shirt", "", func(t assert.TestingT, err error, args ...any) bool {
			return assert.ErrorIs(t, err, model.ErrMerchNotExist, args...)
		}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := model.NewMerchItem(test.Input)
			test.Err(t, err)

			if got != nil {
				assert.Equal(t, test.Want, got.Kind)
			}
		})
	}
}
