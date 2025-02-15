package postgres

import (
	"fmt"
	"strings"
)

type Config struct {
	Host string
	Port string
	Name string

	User     string
	Password string

	Migrations struct {
		Dir string
	}
}

func (c *Config) ConnString() string {
	options := []struct {
		name, val string
	}{
		{"host", c.Host},
		{"port", c.Port},
		{"user", c.User},
		{"password", c.Password},
		{"database", c.Name},
		{"sslmode", "disable"},
	}

	var s []string
	for _, option := range options {
		if option.val == "" {
			continue
		}
		s = append(s, fmt.Sprintf("%s='%s'", option.name, option.val))
	}
	return strings.Join(s, " ")
}

func NewConfig() *Config {
	// Setup Postgres defaults
	return &Config{
		Host: "localhost",
		Port: "5432",
		Name: "postgres",
		User: "postgres",
	}
}
