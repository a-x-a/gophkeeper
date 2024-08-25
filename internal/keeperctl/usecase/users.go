package usecase

import (
	"context"
	"fmt"

	"github.com/a-x-a/gophkeeper/internal/keeperctl/entity"
	"github.com/a-x-a/gophkeeper/internal/keeperctl/repo"
)

var _ Users = (*UsersUseCase)(nil)

// UsersUseCase contains business logic related to users management.
type UsersUseCase struct {
	usersRepo repo.Users
}

// NewUsersUseCase create and initializes new UsersUseCase object.
func NewUsersUseCase(users repo.Users) *UsersUseCase {
	return &UsersUseCase{users}
}

// Register creates a new user.
func (uc *UsersUseCase) Register(
	ctx context.Context,
	username string,
	key entity.Key,
) (string, error) {
	securityKey := key.Hash()

	accessToken, err := uc.usersRepo.Register(ctx, username, securityKey)
	if err != nil {
		return "", fmt.Errorf("UsersUseCase - Register - uc.usersRepo.Register: %w", err)
	}

	return accessToken, nil
}
