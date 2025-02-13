package api_test

import (
	"net/http"
	"testing"

	"github.com/cerfical/merchshop/internal/api"
	"github.com/cerfical/merchshop/internal/domain/model"
	"github.com/cerfical/merchshop/internal/mocks"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITest))
}

type APITest struct {
	suite.Suite

	expect      *httpexpect.Expect
	authService *mocks.AuthService
	coinService *mocks.CoinService
}

func (t *APITest) SetupSubTest() {
	t.authService = mocks.NewAuthService(t.T())
	t.coinService = mocks.NewCoinService(t.T())

	t.expect = httpexpect.WithConfig(httpexpect.Config{
		TestName: t.T().Name(),
		BaseURL:  "/api/",
		Reporter: httpexpect.NewAssertReporter(t.T()),
		Client: &http.Client{
			Transport: httpexpect.NewBinder(api.NewHandler(t.authService, t.coinService, nil)),
		},
	})
}

func (t *APITest) TestAuth() {
	tests := []struct {
		Name  string
		Setup func()

		Builder func(*httpexpect.Request)
		Status  int
	}{
		{
			Name: "ok",
			Setup: func() {
				t.authService.EXPECT().
					AuthUser(model.UserCreds{
						Username: "testuser",
						Password: "123",
					}).
					Return("123321", nil)
			},

			Builder: func(r *httpexpect.Request) {
				r.WithJSON(map[string]string{
					"username": "testuser",
					"password": "123",
				})

				r.WithMatcher(func(r *httpexpect.Response) {
					r.JSON().Object().
						Value("token").IsEqual("123321")
				})
			},
			Status: http.StatusOK,
		},

		{
			Name: "user_auth_fail",
			Setup: func() {
				t.authService.EXPECT().
					AuthUser(model.UserCreds{
						Username: "testuser",
						Password: "123",
					}).
					Return("", model.ErrAuthFail)
			},

			Builder: func(r *httpexpect.Request) {
				r.WithJSON(map[string]string{
					"username": "testuser",
					"password": "123",
				})

				r.WithMatcher(func(r *httpexpect.Response) {
					r.JSON().Object().
						Value("errors")
				})
			},
			Status: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func() {
			test.Setup()

			e := t.expect.Builder(test.Builder).POST("/auth").
				Expect()
			e.Status(test.Status)
		})
	}
}
