package postgres

import (
	"database/sql"
	"errors"
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
			password TEXT NOT NULL
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

func (s *Storage) GetUser(uc *model.UserCreds) (user *model.User, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) && err == nil {
			err = rbErr
		}
	}()

	// TODO: Hash the stored passwords
	passwdHash := uc.Password
	row := tx.QueryRow("SELECT * FROM users WHERE name=$1", uc.Name)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Password); err != nil {
		// If the error is not caused by the absence of the user record
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		// Create a new user
		row := tx.QueryRow(`
			INSERT INTO users(name, password) VALUES($1, $2)
			RETURNING *`,
			uc.Name,
			passwdHash,
		)

		if err := row.Scan(&u.ID, &u.Name, &u.Password); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
