package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cerfical/merchshop/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(cfg *Config) (*Storage, error) {
	connStr, err := makeConnString(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL
		)
	`); err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

func makeConnString(cfg *Config) (string, error) {
	c := []struct {
		key, val string
	}{
		{"host", cfg.Host},
		{"port", cfg.Port},
		{"user", cfg.User},
		{"password", cfg.Password},
		{"database", cfg.Name},
		{"sslmode", "disable"},
	}

	var options []string
	for _, cc := range c {
		if cc.val == "" {
			continue
		}
		options = append(options, fmt.Sprintf("%v='%v'", cc.key, cc.val))
	}

	connStr := strings.Join(options, " ")
	return connStr, nil
}

type Storage struct {
	db *sql.DB
}

func (s *Storage) GetUser(uc *model.UserCreds) (*model.User, error) {
	// TODO: Unnecessary update?
	row := s.db.QueryRow(`
		INSERT INTO users(name, password_hash) VALUES($1, $2)
		ON CONFLICT (name) DO UPDATE SET
			name=EXCLUDED.name
		RETURNING *`,
		uc.Name,
		uc.PasswordHash,
	)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.PasswordHash); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
