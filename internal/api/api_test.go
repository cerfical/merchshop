package api_test

import (
	"net/http"
	"testing"

	"github.com/cerfical/merchshop/internal/api"
	"github.com/cerfical/merchshop/internal/mocks"
	"github.com/cerfical/merchshop/internal/service/auth"
	"github.com/cerfical/merchshop/internal/service/model"
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
					AuthUser(
						model.Username("testuser"),
						model.Password("123"),
					).
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
					AuthUser(
						model.Username("testuser"),
						model.Password("123"),
					).
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

func (t *APITest) TestInfo() {
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
					AuthToken(auth.Token("123321")).
					Return("testuser", nil)
				t.coinService.EXPECT().
					GetCoinBalance(model.Username("testuser")).
					Return(9, nil)
			},

			Builder: func(r *httpexpect.Request) {
				r.WithHeader("Authorization", "Bearer 123321")
				r.WithMatcher(func(r *httpexpect.Response) {
					r.JSON().Object().
						Value("coins").IsEqual(9)
				})
			},
			Status: http.StatusOK,
		},

		{
			Name: "auth_fail",
			Setup: func() {
				t.authService.EXPECT().
					AuthToken(auth.Token("124421")).
					Return("", model.ErrAuthFail)
			},

			Builder: func(r *httpexpect.Request) {
				r.WithHeader("Authorization", "Bearer 124421")
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

			e := t.expect.Builder(test.Builder).GET("/info").
				Expect()
			e.Status(test.Status)
		})
	}
}
