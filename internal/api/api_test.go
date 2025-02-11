package api_test

import (
	"net/http"
	"testing"

	"github.com/cerfical/merchshop/internal/api"
	"github.com/cerfical/merchshop/internal/mocks"
	"github.com/cerfical/merchshop/internal/model"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

var user = model.User{
	ID:       7,
	Name:     "testuser",
	Password: "123",
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITest))
}

type APITest struct {
	suite.Suite

	expect *httpexpect.Expect
	users  *mocks.UserStore
}

func (t *APITest) SetupSubTest() {
	t.users = mocks.NewUserStore(t.T())

	api := api.New(&api.AuthConfig{}, t.users, nil)
	t.expect = httpexpect.WithConfig(httpexpect.Config{
		TestName: t.T().Name(),
		BaseURL:  "/api/auth",
		Reporter: httpexpect.NewAssertReporter(t.T()),
		Client: &http.Client{
			Transport: httpexpect.NewBinder(api),
		},
	})
}

func (t *APITest) TestAuth() {
	tests := []struct {
		name    string
		passwd  string
		status  int
		matcher func(*httpexpect.Object)
	}{
		{"valid_password", "123", http.StatusOK, func(r *httpexpect.Object) {
			r.Value("token").
				String()
		}},

		{"invalid_password", "124", http.StatusUnauthorized, func(r *httpexpect.Object) {
			r.Value("errors").
				String()
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			t.users.EXPECT().
				GetUser(&model.UserCreds{
					Name:     user.Name,
					Password: test.passwd,
				}).
				Return(&user, nil)

			e := t.expect.POST("").
				WithJSON(&AuthRequest{
					Username: user.Name,
					Password: test.passwd,
				}).
				Expect()

			e.Status(test.status)
			test.matcher(e.JSON().
				Object())
		})
	}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
