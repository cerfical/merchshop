package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewStorage(config *Config) (*Storage, error) {
	connStr := makeConnString(config)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	// Check if the database connection is actually established
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db, config.Migrations}, nil
}

func makeConnString(config *Config) string {
	c := []struct {
		key, val string
	}{
		{"host", config.Host},
		{"port", config.Port},
		{"user", config.User},
		{"password", config.Password},
		{"database", config.Name},
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
	return connStr
}

type Storage struct {
	db         *sql.DB
	migrations string
}

func (s *Storage) UpMigrations() error {
	return s.migrate(func(m *migrate.Migrate) error {
		return m.Up()
	})
}

func (s *Storage) DownMigrations() error {
	return s.migrate(func(m *migrate.Migrate) error {
		return m.Down()
	})
}

func (s *Storage) migrate(f func(*migrate.Migrate) error) error {
	// Create a database migration driver from the existing DB instance
	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Load migrations from the local filesystem
	migrationsPath := fmt.Sprintf("file://%v", s.migrations)
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return err
	}

	// Apply the migrations if there are any changes
	if err := f(m); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (s *Storage) PutUser(u *model.User) (*model.User, error) {
	// TODO: Unnecessary update?
	row := s.db.QueryRow(`
			INSERT INTO users(username, password_hash)
			VALUES($1, $2)
			ON CONFLICT (username) DO UPDATE
				SET username=EXCLUDED.username
			RETURNING *`,
		u.Username,
		u.PasswordHash,
	)

	var uu model.User
	if err := row.Scan(&uu.ID, &uu.Username, &uu.PasswordHash); err != nil {
		return nil, err
	}
	return &uu, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
