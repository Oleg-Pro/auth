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
	"github.com/jackc/pgx/v4/pgxpool"
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
	pool *pgxpool.Pool
}

// NewRepository create UserRepository
func NewRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &repo{pool: pool}
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
	err = r.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return 0, err
	}

	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builderSelectOne := sq.Select(userColumnID, userColumnName, userColumnEmail, userColumnRoleID, userColumnCreatedAt, userColumnUpdateAt).
		From(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fmt.Sprintf(`"%s"`, userColumnID): id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Printf("Failed to build get query: %v", err)
		return nil, err
	}

	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Args: %v\n", args)

	var user modelRepo.User
	err = r.pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Info.Name, &user.Info.Email, &user.Info.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) Update(ctx context.Context, id int64, name *string, email *string, role model.Role) (int64, error) {
	builderUpdate := sq.Update(userTable).
		PlaceholderFormat(sq.Dollar).
		Set(userColumnUpdateAt, time.Now()).
		Set(userColumnRoleID, role).
		Where(sq.Eq{fmt.Sprintf(`"%s"`, userColumnID): id})

	if name != nil {
		builderUpdate = builderUpdate.Set(userColumnName, *name)
	}

	if email != nil {
		log.Printf("Email: %v", email)

		builderUpdate = builderUpdate.Set(userColumnEmail, email)
	}

	query, args, err := builderUpdate.ToSql()

	if err != nil {
		log.Printf("Failed to build update query: %v", err)
		return 0, err
	}

	res, err := r.pool.Exec(ctx, query, args...)
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

	log.Printf("DELETE SQL query: %s", query)

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to delete user with id %d: %v", id, err)
		return 0, err
	}

	return res.RowsAffected(), nil
}