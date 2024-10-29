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
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id              = gofakeit.Int64()
		name            = gofakeit.Name()
		email           = gofakeit.Email()
		password        = "123456"
		passwordConfirm = "123456"
		role            = desc.Role_ADMIN

		req = &desc.CreateRequest{
			Name:            name,
			Email:           email,
			Pasword:         password,
			PasswordConfirm: passwordConfirm,
			Role:            role,
		}

		passwordConfirmNotEqual   = "1234567"
		reqWithDifferentPasswords = &desc.CreateRequest{
			Name:            name,
			Email:           email,
			Pasword:         password,
			PasswordConfirm: passwordConfirmNotEqual,
			Role:            role,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
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
				mock.CreateMock.ExpectCtxParam1(ctx).Inspect(func(_ context.Context, info *model.UserInfo) {
					require.Equal(t, email, info.Email)
					require.Equal(t, name, info.Name)
					require.Equal(t, model.RoleADMIN, info.Role)

				}).Return(id, nil)

				return mock
			},
		},
		{
			name: "passwords are not equal",
			args: args{
				ctx: ctx,
				req: reqWithDifferentPasswords,
			},
			want: nil,
			err:  userAPI.ErrPasswordsAreNotEqual,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewImplementation(userServiceMock)
			resonse, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resonse)
		})
	}
}
