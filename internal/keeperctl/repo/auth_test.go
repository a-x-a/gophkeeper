package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-x-a/gophkeeper/internal/keeperctl/repo"
	"github.com/a-x-a/gophkeeper/internal/util/gophtest"
	"github.com/a-x-a/gophkeeper/pkg/goph"
)

func newLoginRequest() *goph.LoginRequest {
	return &goph.LoginRequest{
		Username:    gophtest.Username,
		SecurityKey: gophtest.SecurityKey,
	}
}

func TestLogin(t *testing.T) {
	resp := &goph.LoginResponse{
		AccessToken: gophtest.AccessToken,
	}

	m := &goph.AuthClientMock{}
	m.On(
		"Login",
		mock.Anything,
		newLoginRequest(),
		mock.Anything,
	).
		Return(resp, nil)

	sat := repo.NewAuthRepo(m)
	token, err := sat.Login(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, token)
	m.AssertExpectations(t)
}

func TestLoginOnClientFailure(t *testing.T) {
	m := &goph.AuthClientMock{}
	m.On(
		"Login",
		mock.Anything,
		newLoginRequest(),
		mock.Anything,
	).
		Return(nil, gophtest.ErrUnexpected)

	sat := repo.NewAuthRepo(m)
	_, err := sat.Login(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.Error(t, err)
	m.AssertExpectations(t)
}
