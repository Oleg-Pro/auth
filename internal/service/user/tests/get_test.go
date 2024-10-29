package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"
	repoMocks "github.com/Oleg-Pro/auth/internal/repository/mocks"
	"github.com/Oleg-Pro/auth/internal/service/user"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id           = gofakeit.Int64()
		name         = gofakeit.Name()
		email        = gofakeit.Email()
		passwordHash = "123456"
		role         = model.RoleADMIN
		createdAt    = gofakeit.Date()
		updatedAt    = gofakeit.Date()

		/*		ID        int64
				Info      UserInfo
				CreatedAt time.Time
				UpdatedAt sql.NullTime		*/

		/*		req = &model.User
				ID:           id,
				Info:         model.User{

				},
				PaswordHash:         passwordHash,
				Role:           role,
			}*/

	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *model.User
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: &model.User{
				ID: id,
				Info: model.UserInfo{
					Name:        name,
					Email:       email,
					Role:        role,
					PaswordHash: passwordHash,
				},
				CreatedAt: createdAt,
				UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(&model.User{
					ID: id,
					Info: model.UserInfo{
						Name:        name,
						Email:       email,
						Role:        role,
						PaswordHash: passwordHash,
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
			userRepoMock := tt.userRepositoryMock(mc)
			api := user.New(userRepoMock)
			resonse, err := api.Get(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resonse)
		})
	}
}
