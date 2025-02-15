package postgres_test

import (
	"database/sql"
	"testing"

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

	t.Require().NoError(t.storage.MigrateUp())
}

func (t *StorageIntegrationTest) TearDownSuite() {
	t.Require().NoError(t.storage.MigrateDown())
	t.Require().NoError(t.db.Close())
	t.Require().NoError(t.storage.Close())
}

func (t *StorageIntegrationTest) TestCreateUser() {
	u, err := t.storage.CreateUser("test_new_user", model.PasswordHash{1, 2, 3}, 9)
	t.Require().NoError(err)

	var uu model.User
	err = t.db.QueryRow(`
		SELECT id, username, password_hash, coins
		FROM users
		WHERE username = 'test_new_user'`,
	).Scan(&uu.ID, &uu.Username, &uu.PasswordHash, &uu.Coins)

	t.Require().NoError(err)
	t.Require().Equal(u, &uu)
}

func (t *StorageIntegrationTest) TestGetUser() {
	var id model.UserID
	err := t.db.QueryRow(`
		INSERT INTO users (username, password_hash, coins)
		VALUES ('test_existing_user', '\001\002\003', 9)
		RETURNING id`,
	).Scan(&id)
	t.Require().NoError(err)

	u, err := t.storage.GetUser("test_existing_user")
	t.Require().NoError(err)
	t.Require().Equal(u, &model.User{
		ID:           id,
		Username:     "test_existing_user",
		PasswordHash: model.PasswordHash{1, 2, 3},
		Coins:        9,
	})
}

func (t *StorageIntegrationTest) TestPurchaseMerch() {
	// Create a buyer user
	var buyer model.UserID
	err := t.db.QueryRow(`
		INSERT INTO users (username, password_hash, coins)
		VALUES ('test_buyer', '', 9)
		RETURNING id`,
	).Scan(&buyer)
	t.Require().NoError(err)

	err = t.storage.PurchaseMerch(buyer, &model.MerchItem{Kind: "shirt", Price: 8})
	t.Require().NoError(err)

	// Check that buyer's coin balance was decreased by the merch price
	var coins model.NumCoins
	err = t.db.QueryRow(`
		SELECT coins FROM users
		WHERE username = 'test_buyer'`,
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
		JOIN users ON user_id = id
		WHERE username = 'test_buyer'`,
	).Scan(&merch, &quantity)

	t.Require().NoError(err)
	t.Require().Equal(merch, "shirt")
	t.Require().Equal(quantity, 1)
}

func (t *StorageIntegrationTest) TestTransferCoins() {
	// Create users to participate in coin transfer
	_, err := t.db.Exec(`
		INSERT INTO users (username, password_hash, coins)
		VALUES
			('test_sender', '', 9),
			('test_recipient', '', 10)`,
	)
	t.Require().NoError(err)

	var sender, recipient model.UserID
	err = t.db.QueryRow("SELECT id FROM users WHERE username = 'test_sender'").Scan(&sender)
	t.Require().NoError(err)
	err = t.db.QueryRow("SELECT id FROM users WHERE username = 'test_recipient'").Scan(&recipient)
	t.Require().NoError(err)

	err = t.storage.TransferCoins(sender, recipient, 8)
	t.Require().NoError(err)

	// Check that the correct amount of coins was withdrawn from the sender
	var coins int
	err = t.db.QueryRow(`
		SELECT coins
		FROM users
		WHERE username = 'test_sender'`,
	).Scan(&coins)

	t.Require().NoError(err)
	t.Require().Equal(1, coins)

	// And the correct amount was deposited to the recipient
	err = t.db.QueryRow(`
		SELECT coins
		FROM users
		WHERE username = 'test_recipient'`,
	).Scan(&coins)

	t.Require().NoError(err)
	t.Require().Equal(18, coins)
}
