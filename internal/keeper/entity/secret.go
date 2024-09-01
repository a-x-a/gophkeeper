package entity

import (
	"errors"

	uuid "github.com/satori/go.uuid"

	"github.com/a-x-a/gophkeeper/pkg/goph"
)

var (
	ErrSecretNotFound     = errors.New("secret not found")
	ErrSecretExists       = errors.New("secret already exists")
	ErrSecretNameConflict = errors.New("secret with such name already exists")
)

// Secret represents full secret info stored in the service.
type Secret struct {
	ID       uuid.UUID `db:"secret_id"`
	Name     string
	Kind     goph.DataKind
	Metadata []byte
	Data     []byte
}
