package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cerfical/merchshop/internal/service/model"
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

func (s *Storage) CreateUser(un model.Username, passwd model.PasswordHash, coins model.NumCoins) (*model.User, error) {
	row := s.db.QueryRow(`
			INSERT INTO users(username, password_hash)
			VALUES ($1, $2)
			ON CONFLICT (username) DO NOTHING
			RETURNING id`,
		un, passwd,
	)

	u := model.User{
		Username:     un,
		PasswordHash: passwd,
		Coins:        coins,
	}

	if err := row.Scan(&u.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserExist
		}
		return nil, err
	}

	return &u, nil
}

func (s *Storage) GetUserByUsername(un model.Username) (*model.User, error) {
	var u model.User
	row := s.db.QueryRow(`
			SELECT id, username, password_hash, coins
			FROM users
			WHERE username=$1`,
		un,
	)

	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Coins); errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrUserNotExist
	}

	return &u, nil
}

func (s *Storage) TransferCoins(from model.UserID, to model.UserID, amount model.NumCoins) (err error) {
	// TODO: For now, assume that provided user ids are always valid and users cannot be deleted

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if txErr := tx.Rollback(); err == nil && txErr != nil && !errors.Is(txErr, sql.ErrTxDone) {
			err = txErr
		}
	}()

	// Withdraw coins from one user
	res, err := tx.Exec(`
			UPDATE users
			SET coins = coins - $2
			WHERE id = $1
				AND coins >= $2`,
		from, amount,
	)

	// Check if any rows were updated
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrNotEnoughCoins
	}

	// Transfer the coins to another user
	_, err = tx.Exec(`
			UPDATE users
			SET coins = coins + $2
			WHERE id = $1`,
		to, amount,
	)
	if err != nil {
		return err
	}

	// Record the transfer transaction
	_, err = tx.Exec(`
			INSERT INTO coin_transactions(from_user_id, to_user_id, amount)
			VALUES ($1, $2, $3)`,
		from, to, amount,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) Close() error {
	return s.db.Close()
}
