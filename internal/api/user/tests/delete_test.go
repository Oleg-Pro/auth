package tests

import (
	"context"
	"testing"

	userAPI "github.com/Oleg-Pro/auth/internal/api/user"
	"github.com/Oleg-Pro/auth/internal/service"
	serviceMocks "github.com/Oleg-Pro/auth/internal/service/mocks"
	userSaverProducer "github.com/Oleg-Pro/auth/internal/service/producer/user_saver"
	userSaverProducerMocks "github.com/Oleg-Pro/auth/internal/service/producer/user_saver/mocks"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService
	type userSaverProducerMockFunc func(mc *minimock.Controller) userSaverProducer.UserSaverProducer

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id           = gofakeit.Int64()
		numberOfRows = 1
		req          = &desc.DeleteRequest{
			Id: id,
		}

		res = &empty.Empty{}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name                  string
		args                  args
		want                  *empty.Empty
		err                   error
		userServiceMock       userServiceMockFunc
		userSaverProducerMock userSaverProducerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(int64(numberOfRows), nil)
				return mock
			},
			userSaverProducerMock: func(mc *minimock.Controller) userSaverProducer.UserSaverProducer {
				mock := userSaverProducerMocks.NewUserSaverProducerMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userServiceMock := tt.userServiceMock(mc)
			userSaverProducerMock := tt.userSaverProducerMock(mc)
			api := userAPI.NewImplementation(userServiceMock, userSaverProducerMock)
			response, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}
