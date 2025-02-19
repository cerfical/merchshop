package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/cerfical/merchshop/internal/config"
	"github.com/cerfical/merchshop/internal/lib/postgres"
	"github.com/cerfical/merchshop/internal/service/model"
	"github.com/stretchr/testify/suite"
)

func TestIntegrationStorage(t *testing.T) {
	suite.Run(t, new(StorageIntegrationTest))
}

type StorageIntegrationTest struct {
	suite.Suite

	storage *postgres.Storage
	db      *sql.DB
}

func (t *StorageIntegrationTest) SetupSuite() {
	config := config.MustLoad(nil)

	var err error
	t.storage, err = postgres.NewStorage(&config.DB)
	t.Require().NoError(err)

	t.db, err = sql.Open("pgx", config.DB.ConnString())
	t.Require().NoError(err)
}

func (t *StorageIntegrationTest) TearDownSuite() {
	t.Require().NoError(t.db.Close())
	t.Require().NoError(t.storage.Close())
}

func (t *StorageIntegrationTest) TestCreateUser() {
	newUser := newUser("new_user")

	u, err := t.storage.CreateUser(context.Background(), newUser, model.PasswordHash{1, 2, 3}, 9)
	t.Require().NoError(err)

	var uu model.User
	err = t.db.QueryRow(`
			SELECT id, username, password_hash, coins
			FROM users
			WHERE username = $1`,
		newUser,
	).Scan(&uu.ID, &uu.Username, &uu.PasswordHash, &uu.Coins)

	t.Require().NoError(err)
	t.Require().Equal(u, &uu)
}

func (t *StorageIntegrationTest) TestGetUser() {
	user := newUser("user")

	var userID model.UserID
	err := t.db.QueryRow(`
			INSERT INTO users (username, password_hash, coins)
			VALUES ($1, '\001\002\003', 9)
			RETURNING id`,
		user,
	).Scan(&userID)
	t.Require().NoError(err)

	u, err := t.storage.GetUser(context.Background(), user)
	t.Require().NoError(err)
	t.Require().Equal(u, &model.User{
		ID:           userID,
		Username:     user,
		PasswordHash: model.PasswordHash{1, 2, 3},
		Coins:        9,
	})
}

func (t *StorageIntegrationTest) TestPurchaseMerch() {
	buyer := newUser("buyer")

	// Create a buyer user
	var buyerID model.UserID
	err := t.db.QueryRow(`
			INSERT INTO users (username, password_hash, coins)
			VALUES ($1, '', 9)
			RETURNING id`,
		buyer,
	).Scan(&buyerID)
	t.Require().NoError(err)

	err = t.storage.PurchaseMerch(context.Background(), buyerID, &model.MerchItem{Kind: "shirt", Price: 8})
	t.Require().NoError(err)

	// Check that buyer's coin balance was decreased by the merch price
	var coins model.NumCoins
	err = t.db.QueryRow(`
			SELECT coins FROM users
			WHERE id = $1`,
		buyerID,
	).Scan(&coins)

	t.Require().NoError(err)
	t.Require().Equal(model.NumCoins(1), coins)

	// And buyer's inventory was replenished by the recently bought merch item
	var (
		merch    string
		quantity int
	)
	err = t.db.QueryRow(`
			SELECT merch, quantity
			FROM user_inventories
			WHERE user_id = $1`,
		buyerID,
	).Scan(&merch, &quantity)

	t.Require().NoError(err)
	t.Require().Equal(merch, "shirt")
	t.Require().Equal(quantity, 1)
}

func (t *StorageIntegrationTest) TestTransferCoins() {
	// Create users to perform coin transfer
	sender := newUser("sender")
	recipient := newUser("recipient")

	_, err := t.db.Exec(`
			INSERT INTO users (username, password_hash, coins)
			VALUES ($1, '', 9), ($2, '', 10)`,
		sender, recipient,
	)
	t.Require().NoError(err)

	var senderID, recipientID model.UserID
	err = t.db.QueryRow("SELECT id FROM users WHERE username = $1", sender).
		Scan(&senderID)
	t.Require().NoError(err)

	err = t.db.QueryRow("SELECT id FROM users WHERE username = $1", recipient).
		Scan(&recipientID)
	t.Require().NoError(err)

	err = t.storage.TransferCoins(context.Background(), senderID, recipientID, 8)
	t.Require().NoError(err)

	// Check that the correct amount of coins was withdrawn from the sender
	var coins int
	err = t.db.QueryRow(`
			SELECT coins
			FROM users
			WHERE id = $1`,
		senderID,
	).Scan(&coins)

	t.Require().NoError(err)
	t.Require().Equal(1, coins)

	// And the correct amount was deposited to the recipient
	err = t.db.QueryRow(`
			SELECT coins
			FROM users
			WHERE id = $1`,
		recipientID,
	).Scan(&coins)

	t.Require().NoError(err)
	t.Require().Equal(18, coins)
}

func newUser(s string) model.Username {
	return model.Username(fmt.Sprintf("test_storage_%s_%d", s, time.Now().UnixNano()))
}
