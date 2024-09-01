package usecase

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/a-x-a/gophkeeper/internal/keeper/config"
	"github.com/a-x-a/gophkeeper/internal/keeper/entity"
	"github.com/a-x-a/gophkeeper/internal/keeper/repo"
	"github.com/a-x-a/gophkeeper/pkg/goph"
)

type Auth interface {
	Login(ctx context.Context, username, securityKey string) (entity.AccessToken, error)
}

type Secrets interface {
	Create(
		ctx context.Context,
		owner uuid.UUID,
		name string,
		kind goph.DataKind,
		metadata, data []byte,
	) (uuid.UUID, error)

	List(ctx context.Context, owner uuid.UUID) ([]entity.Secret, error)
	Get(ctx context.Context, owner, id uuid.UUID) (*entity.Secret, error)

	Update(
		ctx context.Context,
		owner, id uuid.UUID,
		changed []string,
		name string,
		metadata []byte,
		data []byte,
	) error

	Delete(ctx context.Context, owner, id uuid.UUID) error
}

type Users interface {
	Register(ctx context.Context, username, securityKey string) (entity.AccessToken, error)
}

// UseCases is a collection of business logic use cases.
type UseCases struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of business logic use cases.
func New(cfg *config.Config, repos *repo.Repositories) *UseCases {
	return &UseCases{
		Auth:    NewAuthUseCase(cfg.Secret, repos.Users),
		Secrets: NewSecretsUseCase(repos.Secrets),
		Users:   NewUsersUseCase(cfg.Secret, repos.Users),
	}
}
