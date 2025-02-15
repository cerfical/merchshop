package postgres_test

import (
	"testing"

	"github.com/cerfical/merchshop/internal/lib/postgres"
	"github.com/stretchr/testify/assert"
)

func TestConfig_ConnString(t *testing.T) {
	tests := []struct {
		Name   string
		Config *postgres.Config
		Want   string
	}{
		{"empty_config", &postgres.Config{}, "sslmode='disable'"},
		{"default_config", postgres.NewConfig(), "host='localhost' port='5432' user='postgres' database='postgres' sslmode='disable'"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Want, test.Config.ConnString())
		})
	}
}
