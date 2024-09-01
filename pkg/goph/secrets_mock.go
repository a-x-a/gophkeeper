package goph

import (
	context "context"

	"github.com/stretchr/testify/mock"
	grpc "google.golang.org/grpc"
)

var _ SecretsServiceClient = (*SecretsClientMock)(nil)

type SecretsClientMock struct {
	mock.Mock
}

func (m *SecretsClientMock) Create(
	ctx context.Context,
	in *CreateRequest,
	opts ...grpc.CallOption,
) (*CreateResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateResponse), args.Error(1)
}

func (m *SecretsClientMock) List(
	ctx context.Context,
	in *ListRequest,
	opts ...grpc.CallOption,
) (*ListResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListResponse), args.Error(1)
}

func (m *SecretsClientMock) Get(
	ctx context.Context,
	in *GetRequest,
	opts ...grpc.CallOption,
) (*GetResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetResponse), args.Error(1)
}

func (m *SecretsClientMock) Update(
	ctx context.Context,
	in *UpdateRequest,
	opts ...grpc.CallOption,
) (*UpdateResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateResponse), args.Error(1)
}

func (m *SecretsClientMock) Delete(
	ctx context.Context,
	in *DeleteRequest,
	opts ...grpc.CallOption,
) (*DeleteResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteResponse), args.Error(1)
}
