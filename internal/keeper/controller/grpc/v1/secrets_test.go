package v1_test

import (
	"context"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	v1 "github.com/a-x-a/gophkeeper/internal/keeper/controller/grpc/v1"
	"github.com/a-x-a/gophkeeper/internal/keeper/entity"
	"github.com/a-x-a/gophkeeper/internal/keeper/usecase"
	"github.com/a-x-a/gophkeeper/internal/util/gophtest"
	"github.com/a-x-a/gophkeeper/pkg/goph"
)

func doListSecrets(
	t *testing.T,
	mockRV []entity.Secret,
	mockErr error,
) (*goph.ListResponse, error) {
	t.Helper()

	m := newUseCasesMock()
	m.Secrets.(*usecase.SecretsUseCaseMock).On(
		"List",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
	).
		Return(mockRV, mockErr)

	conn := createTestServerWithFakeAuth(t, m)
	req := &goph.ListRequest{}

	client := goph.NewSecretsServiceClient(conn)
	rv, err := client.List(context.Background(), req)

	m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)

	return rv, err
}

func doGetSecret(
	t *testing.T,
	mockRV *entity.Secret,
	mockErr error,
) (*goph.GetResponse, error) {
	t.Helper()

	m := newUseCasesMock()
	m.Secrets.(*usecase.SecretsUseCaseMock).On(
		"Get",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("uuid.UUID"),
	).
		Return(mockRV, mockErr)

	conn := createTestServerWithFakeAuth(t, m)
	req := &goph.GetRequest{Id: uuid.NewV4().String()}

	client := goph.NewSecretsServiceClient(conn)
	rv, err := client.Get(context.Background(), req)

	m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)

	return rv, err
}

func doDeleteSecret(
	t *testing.T,
	mockErr error,
) (*goph.DeleteResponse, error) {
	t.Helper()

	m := newUseCasesMock()
	m.Secrets.(*usecase.SecretsUseCaseMock).On(
		"Delete",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("uuid.UUID"),
	).
		Return(mockErr)

	conn := createTestServerWithFakeAuth(t, m)
	req := &goph.DeleteRequest{Id: uuid.NewV4().String()}

	client := goph.NewSecretsServiceClient(conn)
	rv, err := client.Delete(context.Background(), req)

	m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)

	return rv, err
}

func TestCreateSecret(t *testing.T) {
	tt := []struct {
		name       string
		secretName string
		metadata   []byte
		data       []byte
	}{
		{
			name:       "Create secret",
			secretName: gophtest.SecretName,
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret without metadata",
			secretName: gophtest.Username,
			metadata:   nil,
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret of maximum size",
			secretName: strings.Repeat("#", v1.DefaultMaxSecretNameLength),
			metadata:   []byte(strings.Repeat("#", v1.DefaultMetadataLimit)),
			data:       []byte(strings.Repeat("#", v1.DefaultDataLimit)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			expected := uuid.NewV4()

			m := newUseCasesMock()
			m.Secrets.(*usecase.SecretsUseCaseMock).On(
				"Create",
				mock.Anything,
				mock.AnythingOfType("uuid.UUID"),
				tc.secretName,
				goph.DataKind_DATA_KIND_BINARY,
				tc.metadata,
				tc.data,
			).
				Return(expected, nil)

			conn := createTestServerWithFakeAuth(t, m)

			req := &goph.CreateRequest{
				Name:     tc.secretName,
				Kind:     goph.DataKind_DATA_KIND_BINARY,
				Metadata: tc.metadata,
				Data:     tc.data,
			}

			client := goph.NewSecretsServiceClient(conn)
			resp, err := client.Create(context.Background(), req)

			require.NoError(t, err)
			require.Equal(t, expected.String(), resp.GetId())
			m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)
		})
	}
}

func TestCreateSecretWithBadRequest(t *testing.T) {
	tt := []struct {
		name       string
		secretName string
		metadata   []byte
		data       []byte
	}{
		{
			name:       "Create secret fails if secret name is empty",
			secretName: "",
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret fails if secret name is too long",
			secretName: strings.Repeat("#", v1.DefaultMaxSecretNameLength+1),
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret fails if metadata is too long",
			secretName: gophtest.Username,
			metadata:   []byte(strings.Repeat("#", v1.DefaultMetadataLimit+1)),
			data:       make([]byte, 0),
		},
		{
			name:       "Create secret fails if data is empty",
			secretName: gophtest.Username,
			metadata:   []byte(gophtest.Metadata),
			data:       make([]byte, 0),
		},
		{
			name:       "Create secret fails if data is too long",
			secretName: gophtest.Username,
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(strings.Repeat("#", v1.DefaultDataLimit+1)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conn := createTestServerWithFakeAuth(t, newUseCasesMock())

			req := &goph.CreateRequest{
				Name:     tc.secretName,
				Kind:     goph.DataKind_DATA_KIND_BINARY,
				Metadata: tc.metadata,
				Data:     tc.data,
			}

			client := goph.NewSecretsServiceClient(conn)
			_, err := client.Create(context.Background(), req)

			requireEqualCode(t, codes.InvalidArgument, err)
		})
	}
}

func TestCreateSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newUseCasesMock())

	client := goph.NewSecretsServiceClient(conn)
	_, err := client.Create(context.Background(), &goph.CreateRequest{})

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestCreateServerOnUseCaseFailure(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected codes.Code
	}{
		{
			name:     "Create secret fails if secret already exists",
			err:      entity.ErrSecretExists,
			expected: codes.AlreadyExists,
		},
		{
			name:     "Create secret fails if use case fails unexpectedly",
			err:      gophtest.ErrUnexpected,
			expected: codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newUseCasesMock()
			m.Secrets.(*usecase.SecretsUseCaseMock).On(
				"Create",
				mock.Anything,
				mock.AnythingOfType("uuid.UUID"),
				gophtest.SecretName,
				goph.DataKind_DATA_KIND_BINARY,
				[]byte(gophtest.Metadata),
				[]byte(gophtest.TextData),
			).
				Return(uuid.UUID{}, tc.err)

			conn := createTestServerWithFakeAuth(t, m)

			req := &goph.CreateRequest{
				Name:     gophtest.SecretName,
				Kind:     goph.DataKind_DATA_KIND_BINARY,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			}

			client := goph.NewSecretsServiceClient(conn)
			_, err := client.Create(context.Background(), req)

			requireEqualCode(t, tc.expected, err)
			m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)
		})
	}
}

func TestListSecrets(t *testing.T) {
	tt := []struct {
		name    string
		secrets []entity.Secret
	}{
		{
			name: "List secrets of a user",
			secrets: []entity.Secret{
				{
					ID:   gophtest.CreateUUID(t, "7728154c-9400-4f1b-a2a3-01deb83ece05"),
					Name: gophtest.SecretName,
					Kind: goph.DataKind_DATA_KIND_BINARY,
				},
				{
					ID:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45"),
					Name:     gophtest.SecretName + "ex",
					Kind:     goph.DataKind_DATA_KIND_TEXT,
					Metadata: []byte(gophtest.Metadata),
				},
			},
		},
		{
			name:    "List secrets when user has no secrets",
			secrets: []entity.Secret{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rv, err := doListSecrets(t, tc.secrets, nil)

			require.NoError(t, err)
			snaps.MatchSnapshot(t, rv.GetSecrets())
		})
	}
}

func TestListSecretsFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newUseCasesMock())

	client := goph.NewSecretsServiceClient(conn)
	_, err := client.List(context.Background(), &goph.ListRequest{})

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestListSecretsFailsOnUseCaseFailure(t *testing.T) {
	_, err := doListSecrets(t, nil, gophtest.ErrUnexpected)

	requireEqualCode(t, codes.Internal, err)
}

func TestGetSecret(t *testing.T) {
	tt := []struct {
		name   string
		secret *entity.Secret
	}{
		{
			name: "Get secret",
			secret: &entity.Secret{
				ID:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45"),
				Name:     gophtest.SecretName,
				Kind:     goph.DataKind_DATA_KIND_TEXT,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			},
		},
		{
			name: "Get secret without metadata",
			secret: &entity.Secret{
				ID:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45"),
				Name:     gophtest.SecretName,
				Kind:     goph.DataKind_DATA_KIND_TEXT,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := doGetSecret(t, tc.secret, nil)

			require.NoError(t, err)
			snaps.MatchSnapshot(t, resp.GetSecret())
			require.Equal(t, tc.secret.Data, resp.GetData())
		})
	}
}

func TestGetSecretOnBadRequest(t *testing.T) {
	conn := createTestServerWithFakeAuth(t, newUseCasesMock())

	req := &goph.GetRequest{Id: "xxx"}

	client := goph.NewSecretsServiceClient(conn)
	_, err := client.Get(context.Background(), req)

	requireEqualCode(t, codes.InvalidArgument, err)
}

func TestGetSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newUseCasesMock())

	client := goph.NewSecretsServiceClient(conn)
	_, err := client.Get(context.Background(), &goph.GetRequest{})

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestGetSecretOnUsecaseFailure(t *testing.T) {
	tt := []struct {
		name     string
		ucErr    error
		expected codes.Code
	}{
		{
			name:     "Get secret fails if secret not found",
			ucErr:    entity.ErrSecretNotFound,
			expected: codes.NotFound,
		},
		{
			name:     "Get secret fails on expected error",
			ucErr:    gophtest.ErrUnexpected,
			expected: codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := doGetSecret(t, nil, tc.ucErr)

			requireEqualCode(t, tc.expected, err)
		})
	}
}

func TestUpdateSecret(t *testing.T) {
	tt := []struct {
		name    string
		req     *goph.UpdateRequest
		changed []string
	}{
		{
			name: "Update all fields of a secret",
			req: &goph.UpdateRequest{
				Name:     gophtest.SecretName,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			},
			changed: []string{"data", "metadata", "name"},
		},
		{
			name: "Update secret's name",
			req: &goph.UpdateRequest{
				Name: gophtest.SecretName,
			},
			changed: []string{"name"},
		},
		{
			name: "Update secret's metadata",
			req: &goph.UpdateRequest{
				Metadata: []byte(gophtest.Metadata),
			},
			changed: []string{"metadata"},
		},
		{
			name: "Reset secret's metadata",
			req: &goph.UpdateRequest{
				Metadata: []byte(nil),
			},
			changed: []string{"metadata"},
		},
		{
			name: "Update secret's data",
			req: &goph.UpdateRequest{
				Data: []byte(gophtest.TextData),
			},
			changed: []string{"data"},
		},
		{
			name: "Update secret with maximum fields limits",
			req: &goph.UpdateRequest{
				Name:     strings.Repeat("#", v1.DefaultMaxSecretNameLength),
				Metadata: []byte(strings.Repeat("#", v1.DefaultMetadataLimit)),
				Data:     []byte(strings.Repeat("#", v1.DefaultDataLimit)),
			},
			changed: []string{"data", "metadata", "name"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.NewV4()

			mask, err := fieldmaskpb.New(tc.req, tc.changed...)
			require.NoError(t, err)

			tc.req.Id = id.String()
			tc.req.UpdateMask = mask

			m := newUseCasesMock()
			m.Secrets.(*usecase.SecretsUseCaseMock).On(
				"Update",
				mock.Anything,
				mock.AnythingOfType("uuid.UUID"),
				id,
				tc.changed,
				tc.req.Name,
				tc.req.Metadata,
				tc.req.Data,
			).
				Return(nil)

			conn := createTestServerWithFakeAuth(t, m)

			client := goph.NewSecretsServiceClient(conn)
			_, err = client.Update(context.Background(), tc.req)

			m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)

			require.NoError(t, err)
		})
	}
}

func TestUpdateSecretOnBadRequest(t *testing.T) {
	tt := []struct {
		name    string
		req     *goph.UpdateRequest
		changed []string
	}{
		{
			name: "Update fails if no mask specified",
			req: &goph.UpdateRequest{
				Id:   uuid.NewV4().String(),
				Name: gophtest.SecretName,
			},
			changed: nil,
		},
		{
			name: "Update fails if bad secret id provided",
			req: &goph.UpdateRequest{
				Id:   "xxx",
				Name: gophtest.SecretName,
			},
			changed: []string{"name"},
		},
		{
			name: "Update fails if empty name provided",
			req: &goph.UpdateRequest{
				Id:   uuid.NewV4().String(),
				Name: "",
			},
			changed: []string{"name"},
		},
		{
			name: "Update fails if too long name provided",
			req: &goph.UpdateRequest{
				Id:   uuid.NewV4().String(),
				Name: strings.Repeat("#", v1.DefaultMaxSecretNameLength+1),
			},
			changed: []string{"name"},
		},
		{
			name: "Update fails if too long metadata provided",
			req: &goph.UpdateRequest{
				Id:       uuid.NewV4().String(),
				Metadata: []byte(strings.Repeat("#", v1.DefaultMetadataLimit+1)),
			},
			changed: []string{"metadata"},
		},
		{
			name: "Update fails if empty data provided",
			req: &goph.UpdateRequest{
				Id: uuid.NewV4().String(),
			},
			changed: []string{"data"},
		},
		{
			name: "Update fails if too long data provided",
			req: &goph.UpdateRequest{
				Id:   uuid.NewV4().String(),
				Data: []byte(strings.Repeat("#", v1.DefaultDataLimit+1)),
			},
			changed: []string{"data"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conn := createTestServerWithFakeAuth(t, newUseCasesMock())

			mask, err := fieldmaskpb.New(tc.req, tc.changed...)
			require.NoError(t, err)

			tc.req.UpdateMask = mask

			client := goph.NewSecretsServiceClient(conn)
			_, err = client.Update(context.Background(), tc.req)

			requireEqualCode(t, codes.InvalidArgument, err)
		})
	}
}

func TestUpdateSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newUseCasesMock())

	client := goph.NewSecretsServiceClient(conn)
	_, err := client.Update(context.Background(), &goph.UpdateRequest{})

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestUpdateSecretOnUsecaseFailure(t *testing.T) {
	tt := []struct {
		name     string
		ucErr    error
		expected codes.Code
	}{
		{
			name:     "Update secret fails if secret not found",
			ucErr:    entity.ErrSecretNotFound,
			expected: codes.NotFound,
		},
		{
			name:     "Update secret fails if secret with provided name already exists",
			ucErr:    entity.ErrSecretNameConflict,
			expected: codes.AlreadyExists,
		},
		{
			name:     "Update secret fails on expected error",
			ucErr:    gophtest.ErrUnexpected,
			expected: codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.NewV4()

			m := newUseCasesMock()
			m.Secrets.(*usecase.SecretsUseCaseMock).On(
				"Update",
				mock.Anything,
				mock.AnythingOfType("uuid.UUID"),
				id,
				[]string{"name"},
				gophtest.SecretName,
				[]byte(nil),
				[]byte(nil),
			).
				Return(tc.ucErr)

			conn := createTestServerWithFakeAuth(t, m)
			req := &goph.UpdateRequest{
				Id:   id.String(),
				Name: gophtest.SecretName,
			}

			mask, err := fieldmaskpb.New(req, "name")
			require.NoError(t, err)

			req.UpdateMask = mask

			client := goph.NewSecretsServiceClient(conn)
			_, err = client.Update(context.Background(), req)

			m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)
			requireEqualCode(t, tc.expected, err)
		})
	}
}

func TestDeleteSecret(t *testing.T) {
	_, err := doDeleteSecret(t, nil)

	require.NoError(t, err)
}

func TestDeleteSecretOnBadRequest(t *testing.T) {
	conn := createTestServerWithFakeAuth(t, newUseCasesMock())

	req := &goph.DeleteRequest{Id: "xxx"}

	client := goph.NewSecretsServiceClient(conn)
	_, err := client.Delete(context.Background(), req)

	requireEqualCode(t, codes.InvalidArgument, err)
}

func TestDeleteSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newUseCasesMock())

	client := goph.NewSecretsServiceClient(conn)
	_, err := client.Delete(context.Background(), &goph.DeleteRequest{})

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestDeleteSecretOnUsecaseFailure(t *testing.T) {
	tt := []struct {
		name     string
		ucErr    error
		expected codes.Code
	}{
		{
			name:     "Delete secret fails if secret not found",
			ucErr:    entity.ErrSecretNotFound,
			expected: codes.NotFound,
		},
		{
			name:     "Delete secret fails on expected error",
			ucErr:    gophtest.ErrUnexpected,
			expected: codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := doDeleteSecret(t, tc.ucErr)

			requireEqualCode(t, tc.expected, err)
		})
	}
}
