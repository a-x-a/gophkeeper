package repo

import (
	"context"
	"fmt"

	"github.com/a-x-a/gophkeeper/internal/keeperctl/entity"
	"github.com/a-x-a/gophkeeper/pkg/goph"
)

var _ Users = (*UsersRepo)(nil)

// UsersRepo is facade to operations regarding Keeper.
type UsersRepo struct {
	client goph.UsersServiceClient
}

// NewUsersRepo creates and initializes UsersRepo object.
func NewUsersRepo(client goph.UsersServiceClient) *UsersRepo {
	return &UsersRepo{client}
}

// Register creates a new user.
func (r *UsersRepo) Register(
	ctx context.Context,
	username, securityKey string,
) (string, error) {
	req := &goph.RegisterRequest{
		Username:    username,
		SecurityKey: securityKey,
	}

	resp, err := r.client.Register(ctx, req)
	if err != nil {
		return "", fmt.Errorf("UsersRepo - Register - r.client.Register: %w", entity.NewRequestError(err))
	}

	return resp.GetAccessToken(), nil
}
