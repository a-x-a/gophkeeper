package v1

import (
	"context"
	"errors"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-x-a/gophkeeper/internal/keeper/entity"
	"github.com/a-x-a/gophkeeper/internal/keeper/usecase"
	"github.com/a-x-a/gophkeeper/pkg/goph"
)

// SecretsServer provides implementation of the Secrets API.
type SecretsServer struct {
	goph.UnimplementedSecretsServiceServer

	secretsUseCase usecase.Secrets
}

// NewSecretsServer initializes and creates new SecretsServer.
func NewSecretsServer(secrets usecase.Secrets) *SecretsServer {
	return &SecretsServer{secretsUseCase: secrets}
}

// Create creates new secret for a user.
func (s SecretsServer) Create(
	ctx context.Context,
	req *goph.CreateRequest,
) (*goph.CreateResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	if details, ok := validateCreateSecretReq(req); !ok {
		st := composeBadRequestError(details)

		return nil, st.Err()
	}

	id, err := s.secretsUseCase.Create(
		ctx,
		owner.ID,
		req.GetName(),
		req.GetKind(),
		req.GetMetadata(),
		req.GetData(),
	)
	if err != nil {
		if errors.Is(err, entity.ErrSecretExists) {
			return nil, status.Errorf(codes.AlreadyExists, entity.ErrSecretExists.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &goph.CreateResponse{Id: id.String()}, nil
}

// List retrieves list of the secrets stored a user.
func (s SecretsServer) List(
	ctx context.Context,
	_ *goph.ListRequest,
) (*goph.ListResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	data, err := s.secretsUseCase.List(ctx, owner.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rv := make([]*goph.Secret, 0, len(data))
	for _, val := range data {
		rv = append(rv, &goph.Secret{
			Id:       val.ID.String(),
			Name:     val.Name,
			Kind:     val.Kind,
			Metadata: val.Metadata,
		})
	}

	return &goph.ListResponse{Secrets: rv}, nil
}

// Get returns particular secret with data.
func (s SecretsServer) Get(
	ctx context.Context,
	req *goph.GetRequest,
) (*goph.GetResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	id, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	secret, err := s.secretsUseCase.Get(ctx, owner.ID, id)
	if err != nil {
		if errors.Is(err, entity.ErrSecretNotFound) {
			return nil, status.Errorf(codes.NotFound, entity.ErrSecretNotFound.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &goph.GetResponse{
		Secret: &goph.Secret{
			Id:       secret.ID.String(),
			Name:     secret.Name,
			Kind:     secret.Kind,
			Metadata: secret.Metadata,
		},
		Data: secret.Data,
	}, nil
}

// Update updates particular secret stored by a user.
func (s SecretsServer) Update(
	ctx context.Context,
	req *goph.UpdateRequest,
) (*goph.UpdateResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	id, details := validateUpdateSecretReq(req)
	if details != nil {
		st := composeBadRequestError(details)

		return nil, st.Err()
	}

	mask := req.GetUpdateMask()
	// Remove redundand paths.
	mask.Normalize()

	if err := s.secretsUseCase.Update(
		ctx,
		owner.ID,
		id,
		mask.GetPaths(),
		req.GetName(),
		req.GetMetadata(),
		req.GetData(),
	); err != nil {
		if errors.Is(err, entity.ErrSecretNotFound) {
			return nil, status.Errorf(codes.NotFound, entity.ErrSecretNotFound.Error())
		}

		if errors.Is(err, entity.ErrSecretNameConflict) {
			return nil, status.Errorf(codes.AlreadyExists, entity.ErrSecretNameConflict.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &goph.UpdateResponse{}, nil
}

// Delete removes particular secret stored by a user.
func (s SecretsServer) Delete(
	ctx context.Context,
	req *goph.DeleteRequest,
) (*goph.DeleteResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	id, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := s.secretsUseCase.Delete(ctx, owner.ID, id); err != nil {
		if errors.Is(err, entity.ErrSecretNotFound) {
			return nil, status.Errorf(codes.NotFound, entity.ErrSecretNotFound.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &goph.DeleteResponse{}, nil
}
