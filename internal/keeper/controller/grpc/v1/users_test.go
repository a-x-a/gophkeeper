package v1_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	v1 "github.com/a-x-a/gophkeeper/internal/keeper/controller/grpc/v1"
	"github.com/a-x-a/gophkeeper/internal/keeper/entity"
	"github.com/a-x-a/gophkeeper/internal/keeper/usecase"
	"github.com/a-x-a/gophkeeper/internal/util/gophtest"
	"github.com/a-x-a/gophkeeper/pkg/goph"
)

func TestRegisterUser(t *testing.T) {
	tt := []struct {
		name     string
		userName string
	}{
		{
			name:     "Register user",
			userName: gophtest.Username,
		},
		{
			name:     "Register user with long name",
			userName: strings.Repeat("#", v1.DefaultMaxUsernameLength),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newUseCasesMock()
			m.Users.(*usecase.UsersUseCaseMock).On(
				"Register",
				mock.Anything,
				tc.userName,
				gophtest.SecurityKey,
			).
				Return(entity.AccessToken(gophtest.AccessToken), nil)

			conn := createTestServer(t, m)

			req := &goph.RegisterRequest{
				Username:    tc.userName,
				SecurityKey: gophtest.SecurityKey,
			}

			client := goph.NewUsersServiceClient(conn)
			resp, err := client.Register(context.Background(), req)

			require.NoError(t, err)
			require.Equal(t, gophtest.AccessToken, resp.GetAccessToken())
			m.Users.(*usecase.UsersUseCaseMock).AssertExpectations(t)
		})
	}
}

func TestRegisterUserWithBadRequest(t *testing.T) {
	tt := []struct {
		name     string
		username string
		key      string
	}{
		{
			name:     "Register user fails if username is empty",
			username: "",
			key:      gophtest.SecurityKey,
		},
		{
			name:     "Register user fails if security key is empty",
			username: gophtest.Username,
			key:      "",
		},
		{
			name:     "Register user fails if username is too long",
			username: strings.Repeat("#", v1.DefaultMaxUsernameLength+1),
			key:      gophtest.SecurityKey,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conn := createTestServer(t, newUseCasesMock())

			req := &goph.RegisterRequest{
				Username:    tc.username,
				SecurityKey: tc.key,
			}

			client := goph.NewUsersServiceClient(conn)
			_, err := client.Register(context.Background(), req)

			requireEqualCode(t, codes.InvalidArgument, err)
		})
	}
}

func TestRegisterUserOnUseCaseFailure(t *testing.T) {
	tt := []struct {
		name       string
		useCaseErr error
		expected   codes.Code
	}{
		{
			name:       "Register user fails if user already exists",
			useCaseErr: entity.ErrUserExists,
			expected:   codes.AlreadyExists,
		},
		{
			name:       "Register user fails if something bad happened",
			useCaseErr: gophtest.ErrUnexpected,
			expected:   codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newUseCasesMock()
			m.Users.(*usecase.UsersUseCaseMock).On(
				"Register",
				mock.Anything,
				gophtest.Username,
				gophtest.SecurityKey,
			).
				Return(entity.AccessToken(""), tc.useCaseErr)

			conn := createTestServer(t, m)

			req := &goph.RegisterRequest{
				Username:    gophtest.Username,
				SecurityKey: gophtest.SecurityKey,
			}

			client := goph.NewUsersServiceClient(conn)
			_, err := client.Register(context.Background(), req)

			requireEqualCode(t, tc.expected, err)
			m.Users.(*usecase.UsersUseCaseMock).AssertExpectations(t)
		})
	}
}
