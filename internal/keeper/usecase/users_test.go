package usecase_test

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-x-a/gophkeeper/internal/keeper/entity"
	"github.com/a-x-a/gophkeeper/internal/keeper/repo"
	"github.com/a-x-a/gophkeeper/internal/keeper/usecase"
	"github.com/a-x-a/gophkeeper/internal/util/gophtest"
)

func doRegisterUser(t *testing.T, repoErr error) (entity.AccessToken, error) {
	t.Helper()

	m := &repo.UsersRepoMock{}
	m.On(
		"Register",
		mock.Anything,
		gophtest.Username,
		gophtest.SecurityKey,
	).
		Return(uuid.NewV4(), repoErr)

	sat := usecase.NewUsersUseCase(gophtest.Secret, m)
	token, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	m.AssertExpectations(t)

	return token, err
}

func TestRegisterUser(t *testing.T) {
	token, err := doRegisterUser(t, nil)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestRegisterUserFailsIfUserExists(t *testing.T) {
	_, err := doRegisterUser(t, entity.ErrUserExists)

	require.Error(t, err)
}
