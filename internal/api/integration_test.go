package api_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cerfical/merchshop/internal/api"
	"github.com/cerfical/merchshop/internal/config"
	"github.com/cerfical/merchshop/internal/lib/bcrypt"
	"github.com/cerfical/merchshop/internal/lib/jwt"
	"github.com/cerfical/merchshop/internal/lib/postgres"
	"github.com/cerfical/merchshop/internal/service/auth"
	"github.com/cerfical/merchshop/internal/service/coins"
	"github.com/cerfical/merchshop/internal/service/model"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APIIntegrationTest))
}

type APIIntegrationTest struct {
	suite.Suite

	storage *postgres.Storage
	expect  *httpexpect.Expect

	user, user1 string
	authToken   string
}

func (t *APIIntegrationTest) SetupSuite() {
	config := config.MustLoad(nil)

	var err error
	t.storage, err = postgres.NewStorage(&config.DB)
	t.Require().NoError(err)

	tokenAuth := jwt.NewTokenAuth(&config.API.Auth.Token)
	coins := coins.NewCoinService(t.storage)
	auth := auth.NewAuthService(tokenAuth, t.storage, bcrypt.NewHasher())

	t.expect = httpexpect.WithConfig(httpexpect.Config{
		TestName: t.T().Name(),
		BaseURL:  "/api/",
		Reporter: httpexpect.NewAssertReporter(t.T()),
		Client: &http.Client{
			Transport: httpexpect.NewBinder(api.NewHandler(auth, coins, nil)),
		},
	})

	// Create test users
	t.user = newUser("user")
	token, err := auth.AuthUser(model.Username(t.user), "12345678910")
	t.Require().NoError(err)

	t.user1 = newUser("user1")
	_, err = auth.AuthUser(model.Username(t.user1), "12345678910")
	t.Require().NoError(err)

	t.authToken = string(token)
}

func (t *APIIntegrationTest) TearDownSuite() {
	t.Require().NoError(t.storage.Close())
}

func (t *APIIntegrationTest) TestAuth() {
	tests := []struct {
		Name     string
		Username string
		Password string
		Status   int
	}{
		{
			Name:     "new_user",
			Username: newUser("new_user"),
			Password: "12345678910",
			Status:   http.StatusOK,
		},

		{
			Name:     "existing_user",
			Username: t.user,
			Password: "12345678910",
			Status:   http.StatusOK,
		},

		{
			Name:     "bad_password",
			Username: t.user,
			Password: "bad_12345678910",
			Status:   http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			e := t.expect.POST("/auth").
				WithJSON(map[string]string{
					"username": test.Username,
					"password": test.Password,
				}).
				Expect()
			e.Status(test.Status)
		})
	}
}

func (t *APIIntegrationTest) TestInfo() {
	tests := []struct {
		Name      string
		AuthToken string
		Status    int
	}{
		{
			Name:      "ok",
			AuthToken: t.authToken,
			Status:    http.StatusOK,
		},

		{
			Name:      "invalid_auth_token",
			AuthToken: "bad_token",
			Status:    http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			e := t.expect.GET("/info").
				WithHeader("Authorization", fmt.Sprintf("Bearer %s", test.AuthToken)).
				Expect()
			e.Status(test.Status)
		})
	}
}

func (t *APIIntegrationTest) TestSendCoin() {
	tests := []struct {
		Name      string
		ToUser    string
		AuthToken string
		Amount    int
		Status    int
	}{
		{
			Name:      "ok",
			ToUser:    t.user1,
			AuthToken: t.authToken,
			Amount:    9,
			Status:    http.StatusOK,
		},

		{
			Name:      "invalid_auth_token",
			ToUser:    t.user1,
			AuthToken: "bad_token",
			Amount:    9,
			Status:    http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			e := t.expect.POST("/sendCoin").
				WithHeader("Authorization", fmt.Sprintf("Bearer %s", test.AuthToken)).
				WithJSON(map[string]any{
					"toUser": test.ToUser,
					"amount": test.Amount,
				}).
				Expect()
			e.Status(test.Status)
		})
	}
}

func (t *APIIntegrationTest) TestBuyItem() {
	tests := []struct {
		Name      string
		AuthToken string
		MerchItem string
		Status    int
	}{
		{
			Name:      "ok",
			AuthToken: t.authToken,
			MerchItem: "t-shirt",
			Status:    http.StatusOK,
		},

		{
			Name:      "invalid_merch_item",
			AuthToken: t.authToken,
			MerchItem: "bad_merch",
			Status:    http.StatusBadRequest,
		},

		{
			Name:      "invalid_auth_token",
			AuthToken: "bad_token",
			MerchItem: "t-shirt",
			Status:    http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			e := t.expect.GET("/buy/{item}").WithPath("item", test.MerchItem).
				WithHeader("Authorization", fmt.Sprintf("Bearer %s", test.AuthToken)).
				Expect()

			e.Status(test.Status)
		})
	}
}

func newUser(s string) string {
	return fmt.Sprintf("test_api_%s_%d", s, time.Now().UnixNano())
}
