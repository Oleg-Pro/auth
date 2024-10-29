package tests

import (
	"context"
	"testing"

	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"	
	"github.com/Oleg-Pro/auth/internal/service/user"	
	repoMocks "github.com/Oleg-Pro/auth/internal/repository/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *model.UserInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id              = gofakeit.Int64()
		name            = gofakeit.Name()
		email           = gofakeit.Email()
		passwordHash        = "123456"
		role            = model.RoleADMIN

		req = &model.UserInfo{
			Name:            name,
			Email:           email,
			PaswordHash:         passwordHash,
			Role:            role,
		}

	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            int64
		err             error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepoMock := tt.userRepositoryMock(mc)
			api := user.New(userRepoMock)
			resonse, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resonse)
		})
	}
}
