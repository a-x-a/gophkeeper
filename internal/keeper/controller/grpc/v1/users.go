package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-x-a/gophkeeper/internal/keeper/entity"
	"github.com/a-x-a/gophkeeper/internal/keeper/usecase"
	"github.com/a-x-a/gophkeeper/pkg/goph"
)

// UsersServer provides implementation of the Users API.
type UsersServer struct {
	goph.UnimplementedUsersServiceServer

	usersUseCase usecase.Users
}

// NewUsersServer initializes and creates new UsersServer.
func NewUsersServer(users usecase.Users) *UsersServer {
	return &UsersServer{usersUseCase: users}
}

// Register creates new user.
func (s UsersServer) Register(
	ctx context.Context,
	req *goph.RegisterRequest,
) (*goph.RegisterResponse, error) {
	username := req.GetUsername()
	key := req.GetSecurityKey()

	if details, ok := validateCredentials(username, key); !ok {
		st := composeBadRequestError(details)

		return nil, st.Err()
	}

	accessToken, err := s.usersUseCase.Register(ctx, username, key)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return nil, status.Errorf(codes.AlreadyExists, entity.ErrUserExists.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &goph.RegisterResponse{AccessToken: accessToken.String()}, nil
}
