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

// AuthServer provides implementation of the Auth API.
type AuthServer struct {
	goph.UnimplementedAuthServer
	authUseCase usecase.Auth
}

// NewAuthServer initializes and creates new AuthServer.
func NewAuthServer(auth usecase.Auth) *AuthServer {
	return &AuthServer{authUseCase: auth}
}

// Login authenticates a user in the service.
func (s AuthServer) Login(
	ctx context.Context,
	req *goph.LoginRequest,
) (*goph.LoginResponse, error) {
	username := req.GetUsername()
	key := req.GetSecurityKey()

	if details, ok := validateCredentials(username, key); !ok {
		st := composeBadRequestError(details)

		return nil, st.Err()
	}

	accessToken, err := s.authUseCase.Login(ctx, username, key)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &goph.LoginResponse{AccessToken: accessToken.String()}, nil
}
