package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cerfical/merchshop/internal/service/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewStorage(config *Config) (*Storage, error) {
	db, err := sql.Open("pgx", config.ConnString())
	if err != nil {
		return nil, err
	}

	// Check if the database connection is actually established
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Storage{db}, nil
}

type Storage struct {
	db *sql.DB
}

func (s *Storage) CreateUser(username model.Username, passwd model.PasswordHash, coins model.NumCoins) (*model.User, error) {
	row := s.db.QueryRow(`
			INSERT INTO users(username, password_hash, coins)
			VALUES ($1, $2, $3)
			ON CONFLICT (username) DO NOTHING
			RETURNING id`,
		username, passwd, coins,
	)

	u := model.User{
		Username:     username,
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

func (s *Storage) GetUser(username model.Username) (*model.User, error) {
	// TODO: Implement conditional fetching of fields
	// TODO: Is the transaction really necessary?

	tx, err := s.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	u, err := getUser(tx, username)
	if err != nil {
		return nil, fmt.Errorf("get user data: %w", err)
	}

	u.Inventory, err = getUserInventory(tx, u)
	if err != nil {
		return nil, fmt.Errorf("ger user inventory: %w", err)
	}

	u.Withdrawals, u.Deposits, err = getUserTransactions(tx, u)
	if err != nil {
		return nil, fmt.Errorf("ger user transactions: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}

func getUser(tx *sql.Tx, username model.Username) (*model.User, error) {
	u := model.User{
		Username: username,

		// Ensure the extracted data is not nil
		Deposits:    []model.Deposit{},
		Withdrawals: []model.Withdrawal{},
		Inventory:   []model.InventoryItem{},
	}

	row := tx.QueryRow(`
			SELECT id, password_hash, coins
			FROM users
			WHERE username=$1`,
		username,
	)

	if err := row.Scan(&u.ID, &u.PasswordHash, &u.Coins); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotExist
		}
		return nil, err
	}

	return &u, nil
}

func getUserInventory(tx *sql.Tx, u *model.User) (i []model.InventoryItem, _ error) {
	rows, err := tx.Query(`
			SELECT merch, quantity
			FROM user_inventories
			WHERE user_id = $1`,
		u.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.InventoryItem
		if err := rows.Scan(&item.Merch, &item.Quantity); err != nil {
			return nil, err
		}
		i = append(i, item)
	}

	return i, rows.Err()
}

func getUserTransactions(tx *sql.Tx, u *model.User) (w []model.Withdrawal, d []model.Deposit, _ error) {
	rows, err := tx.Query(`
			SELECT to_users.username, from_users.username, amount
			FROM coin_transactions
			JOIN users AS to_users ON to_users.id = to_user_id
			JOIN users AS from_users ON from_users.id = from_user_id
			WHERE to_user_id = $1
				OR from_user_id = $1`,
		u.ID,
	)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			from, to model.Username
			amount   model.NumCoins
		)
		if err := rows.Scan(&to, &from, &amount); err != nil {
			return nil, nil, err
		}

		if from == u.Username {
			w = append(w, model.Withdrawal{
				To:     to,
				Amount: amount,
			})
		} else {
			d = append(d, model.Deposit{
				From:   from,
				Amount: amount,
			})
		}
	}

	return w, d, rows.Err()
}

func (s *Storage) TransferCoins(from model.UserID, to model.UserID, amount model.NumCoins) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Withdraw coins from one user
	res, err := tx.Exec(`
			UPDATE users
			SET coins = coins - $2
			WHERE id = $1
				AND coins >= $2`,
		from, amount,
	)
	if err != nil {
		return fmt.Errorf("withdraw coins: %w", err)
	}

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
		return fmt.Errorf("deposit coins: %w", err)
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

func (s *Storage) PurchaseMerch(buyer model.UserID, m *model.MerchItem) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Withdraw the amount of coins needed to purchase the item
	res, err := tx.Exec(`
			UPDATE users
			SET coins = coins - $2
			WHERE id = $1
				AND coins >= $2`,
		buyer, m.Price,
	)
	if err != nil {
		return fmt.Errorf("purchase item: %w", err)
	}

	// Check if the purchase was successful, i.e. the user has the required amount of coins
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrNotEnoughCoins
	}

	// Create a record for the item in the user's inventory
	_, err = tx.Exec(`
			INSERT INTO user_inventories(user_id, merch, quantity)
			VALUES($1, $2, 1)
			ON CONFLICT (user_id, merch) DO UPDATE
			SET quantity = user_inventories.quantity + 1`,
		buyer, m.Kind,
	)
	if err != nil {
		return fmt.Errorf("update user inventory: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) Close() error {
	return s.db.Close()
}
