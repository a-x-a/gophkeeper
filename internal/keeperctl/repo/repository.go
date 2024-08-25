package repo

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/a-x-a/gophkeeper/internal/keeperctl/infra/grpcconn"
	"github.com/a-x-a/gophkeeper/pkg/goph"
)

type Auth interface {
	Login(ctx context.Context, username, securityKey string) (string, error)
}

type Secrets interface {
	Push(
		ctx context.Context,
		token, name string,
		kind goph.DataKind,
		description, payload []byte,
	) (uuid.UUID, error)

	List(ctx context.Context, token string) ([]*goph.Secret, error)
	Get(ctx context.Context, token string, id uuid.UUID) (*goph.Secret, []byte, error)

	Update(
		ctx context.Context,
		token string,
		id uuid.UUID,
		name string,
		description []byte,
		noDescription bool,
		data []byte,
	) error

	Delete(ctx context.Context, token string, id uuid.UUID) error
}

type Users interface {
	Register(ctx context.Context, username, securityKey string) (string, error)
}

// Repositories is a collection of data repositories.
type Repositories struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of data repositories.
func New(conn *grpcconn.Connection) *Repositories {
	c := conn.Instance()

	return &Repositories{
		Auth:    NewAuthRepo(goph.NewAuthServiceClient(c)),
		Secrets: NewSecretsRepo(goph.NewSecretsServiceClient(c)),
		Users:   NewUsersRepo(goph.NewUsersServiceClient(c)),
	}
}
