package repo

import (
	"context"
	"fmt"

	"github.com/a-x-a/gophkeeper/internal/keeperctl/entity"
	"github.com/a-x-a/gophkeeper/pkg/goph"
)

var _ Auth = (*AuthRepo)(nil)

// AuthRepo is facade to operations regarding authentication in Keeper.
type AuthRepo struct {
	client goph.AuthServiceClient
}

// NewAuthRepo creates and initializes AuthRepo object.
func NewAuthRepo(client goph.AuthServiceClient) *AuthRepo {
	return &AuthRepo{client}
}

// Login authenticates user in the Keeper service.
func (r *AuthRepo) Login(
	ctx context.Context,
	username, securityKey string,
) (string, error) {
	req := &goph.LoginRequest{
		Username:    username,
		SecurityKey: securityKey,
	}

	resp, err := r.client.Login(ctx, req)
	if err != nil {
		return "", fmt.Errorf("AuthRepo - Login - r.client.Login: %w", entity.NewRequestError(err))
	}

	return resp.GetAccessToken(), nil
}
