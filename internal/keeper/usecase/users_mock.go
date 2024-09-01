package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/a-x-a/gophkeeper/internal/keeper/entity"
)

var _ Users = (*UsersUseCaseMock)(nil)

type UsersUseCaseMock struct {
	mock.Mock
}

func (m *UsersUseCaseMock) Register(
	ctx context.Context,
	username, securityKey string,
) (entity.AccessToken, error) {
	args := m.Called(ctx, username, securityKey)

	return args.Get(0).(entity.AccessToken), args.Error(1)
}
