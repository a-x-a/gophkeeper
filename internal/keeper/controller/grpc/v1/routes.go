// Package v1 implements v1 version of the gRPC API.
package v1

import (
	"google.golang.org/grpc"

	"github.com/a-x-a/gophkeeper/internal/keeper/usecase"
	"github.com/a-x-a/gophkeeper/pkg/goph"
)

// DefaultMaxMessageSize suggests limit for maximum length of gRPC message.
const DefaultMaxMessageSize = DefaultDataLimit + DefaultMetadataLimit + 2*DefaultMaxSecretNameLength

// RegisterRoutes injects new routes into the provided gRPC server.
func RegisterRoutes(server *grpc.Server, useCases *usecase.UseCases) {
	auth := NewAuthServer(useCases.Auth)
	goph.RegisterAuthServer(server, auth)

	secrets := NewSecretsServer(useCases.Secrets)
	goph.RegisterSecretsServer(server, secrets)

	users := NewUsersServer(useCases.Users)
	goph.RegisterUsersServer(server, users)
}
