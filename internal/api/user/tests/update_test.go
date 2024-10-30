package tests

import (
	"context"
	"testing"

	userAPI "github.com/Oleg-Pro/auth/internal/api/user"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/service"
	serviceMocks "github.com/Oleg-Pro/auth/internal/service/mocks"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id           = gofakeit.Int64()
		numberOfRows = int64(1)
		name         = gofakeit.Name()
		email        = gofakeit.Email()
		role         = desc.Role_ADMIN

		req = &desc.UpdateRequest{
			Id:    id,
			Name:  &wrappers.StringValue{Value: name},
			Email: &wrappers.StringValue{Value: email},
			Role:  role,
		}

		res = &empty.Empty{}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *empty.Empty
		err             error
		userServiceMock userServiceMockFunc
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

				mock.UpdateMock.Expect(ctx, id, &model.UserUpdateInfo{
					Name:  &name,
					Email: &email,
					Role:  model.RoleADMIN,
				}).Return(numberOfRows, nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewImplementation(userServiceMock)
			resonse, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resonse)
		})
	}
}
