package user

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"
	"github.com/Oleg-Pro/auth/internal/repository/user/converter"
	modelRepo "github.com/Oleg-Pro/auth/internal/repository/user/model"
	"github.com/Oleg-Pro/platform-common/pkg/db"
)

const (
	userTable = "users"

	userColumnID           = "id"
	userColumnName         = "name"
	userColumnEmail        = "email"
	userColumnRoleID       = "role_id"
	userColumnCreatedAt    = "created_at"
	userColumnUpdateAt     = "updated_at"
	userColumnPasswordHash = "password_hash"
)

type repo struct {
	db db.Client
}

// NewRepository create UserRepository
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.UserInfo) (int64, error) {

	builderInsert := sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns(userColumnName, userColumnEmail, userColumnPasswordHash, userColumnRoleID).
		Values(info.Name, info.Email, info.PaswordHash, info.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("Failed to build insert query: %v", err)
		return 0, err
	}

	var userID int64

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return 0, err
	}

	return userID, nil
}

func (r *repo) Get(ctx context.Context, filter repository.UserFilter) (*model.User, error) {
	builderSelectOne := sq.Select(userColumnID, userColumnName, userColumnEmail, userColumnRoleID, userColumnPasswordHash, userColumnCreatedAt, userColumnUpdateAt).
		From(userTable).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	if filter.ID != nil {
		builderSelectOne = builderSelectOne.Where(sq.Eq{fmt.Sprintf(`"%s"`, userColumnID): filter.ID})
	}

	if filter.Email != nil {
		builderSelectOne = builderSelectOne.Where(sq.Eq{fmt.Sprintf(`"%s"`, userColumnEmail): filter.Email})
	}

	query, args, err := builderSelectOne.ToSql()
	log.Println("Get sql Log")
	log.Printf("Get sql : %s\n", query)
	if err != nil {
		log.Printf("Failed to build get query: %v", err)
		return nil, err
	}

	var user modelRepo.User

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Info.Name, &user.Info.Email, &user.Info.Role, &user.Info.PaswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) Update(ctx context.Context, id int64, info *model.UserUpdateInfo) (int64, error) {
	builderUpdate := sq.Update(userTable).
		PlaceholderFormat(sq.Dollar).
		Set(userColumnUpdateAt, time.Now()).
		Set(userColumnRoleID, info.Role).
		Where(sq.Eq{fmt.Sprintf(`"%s"`, userColumnID): id})

	if info.Name != nil {
		builderUpdate = builderUpdate.Set(userColumnName, info.Name)
	}

	if info.Email != nil {
		builderUpdate = builderUpdate.Set(userColumnEmail, info.Email)
	}

	query, args, err := builderUpdate.ToSql()

	if err != nil {
		log.Printf("Failed to build update query: %v", err)
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("Failed to update user with id %d: %v", id, err)
		return 0, err
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return res.RowsAffected(), err
}

func (r *repo) Delete(ctx context.Context, id int64) (int64, error) {

	builderDelete := sq.Delete(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fmt.Sprintf(`"%s"`, userColumnID): id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("Failed to delete user with id %d: %v", id, err)
		return 0, err
	}

	return res.RowsAffected(), nil
}
