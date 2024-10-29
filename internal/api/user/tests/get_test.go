package tests

import (
	"context"
	"database/sql"
	"testing"

	userAPI "github.com/Oleg-Pro/auth/internal/api/user"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/service"
	serviceMocks "github.com/Oleg-Pro/auth/internal/service/mocks"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id              = gofakeit.Int64()
		name            = gofakeit.Name()
		email           = gofakeit.Email()
		role            = desc.Role_ADMIN
		passwordHash = "$2a$10$krN.Ht8n2kfg12nPcHYMoeHB/dBB7Tvpj40b9U55VP6G.l.inwayO"
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()


		req = &desc.GetRequest{
			Id: id,
		}

		res = &desc.GetResponse{
			Id: id,
			Name: name,
			Email: email,
			Role: role,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),

		}	
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetResponse
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

				mock.GetMock.Expect(ctx, id).Return(&model.User{
					ID: id,
					Info: model.UserInfo{
						Name: name,
						Email: email, 
						PaswordHash: passwordHash,
						Role: model.RoleADMIN,					
					},	
					CreatedAt: createdAt,
					UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},

				}, nil)

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
			response, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}