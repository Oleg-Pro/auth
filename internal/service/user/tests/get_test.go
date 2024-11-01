package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"
	repoMocks "github.com/Oleg-Pro/auth/internal/repository/mocks"
	userService "github.com/Oleg-Pro/auth/internal/service/user"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type userCacheRepositoryMockFunc func(mc *minimock.Controller) repository.UserCacheRepository	

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
		info = model.UserInfo{
			Name:        name,
			Email:       email,
			Role:        role,
			PaswordHash: passwordHash,
		}

		userEntity = &model.User{
			ID: id,
			Info: info,
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true,},
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *model.User
		err                error
		userRepositoryMock userRepositoryMockFunc
		userCacheRepositoryMock userCacheRepositoryMockFunc		
	}{	
		{
			name: "get from database",
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
				mock.GetMock.Expect(ctx, id).Return(userEntity, nil)
				return mock
			},
			userCacheRepositoryMock: func(mc *minimock.Controller) repository.UserCacheRepository {
				mock := repoMocks.NewUserCacheRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, model.ErrorNoteNotFound)
				mock.CreateMock.Expect(ctx, id, &info).Return(0, model.ErrorNoteNotFound)				
				return mock
			},			
		},			
		{
			name: "get from cache",
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
				/*mock.GetMock.Expect(ctx, id).Return(&model.User{
					ID: id,
					Info: model.UserInfo{
						Name:        name,
						Email:       email,
						Role:        role,
						PaswordHash: passwordHash,
					},
					CreatedAt: createdAt,
					UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
				}, nil)*/
				return mock
			},

			userCacheRepositoryMock: func(mc *minimock.Controller) repository.UserCacheRepository {
				mock := repoMocks.NewUserCacheRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(userEntity, nil)
				return mock
			},			
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepoMock := tt.userRepositoryMock(mc)
			userCacheRepoMock := tt.userCacheRepositoryMock(mc)
			api := userService.New(userRepoMock, userCacheRepoMock)
			resonse, err := api.Get(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resonse)
		})
	}
}
