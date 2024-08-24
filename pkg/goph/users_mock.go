package goph

import (
	context "context"

	"github.com/stretchr/testify/mock"
	grpc "google.golang.org/grpc"
)

var _ UsersServiceClient = (*UsersClientMock)(nil)

type UsersClientMock struct {
	mock.Mock
}

func (m *UsersClientMock) Register(
	ctx context.Context,
	in *RegisterRequest,
	opts ...grpc.CallOption,
) (*RegisterResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RegisterResponse), args.Error(1)
}
