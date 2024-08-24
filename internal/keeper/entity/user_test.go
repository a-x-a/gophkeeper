package entity_test

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-x-a/gophkeeper/internal/keeper/entity"
	"github.com/a-x-a/gophkeeper/internal/util/gophtest"
)

func TestUserWithFromContext(t *testing.T) {
	expected := entity.User{
		ID:       uuid.NewV4(),
		Username: gophtest.Username,
	}

	ctx := expected.WithContext(context.Background())

	require.Equal(t, expected, *entity.UserFromContext(ctx))
}

func TestUserFromCleanContext(t *testing.T) {
	require.Nil(t, entity.UserFromContext(context.Background()))
}
