package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/akemoon/crowdfunding-app-user/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	constraintUsernameUnique = "users_username_unique"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

//go:embed sql/create_user.sql
var createUserSQL string

func (r *UserRepo) CreateUser(ctx context.Context, req domain.CreateUserReq) (domain.CreateUserResp, error) {
	var resp domain.CreateUserResp

	err := r.db.QueryRowContext(ctx, createUserSQL, req.UserID, req.Username).Scan(&resp.UserID)
	if err != nil {
		pgErr := asPostgresError(err)
		if pgErr != nil {
			mappedErr := mapPostgresError(pgErr)
			return domain.CreateUserResp{}, fmt.Errorf("%w: %s", mappedErr, pgErr.Detail)
		}

		return domain.CreateUserResp{}, fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	return resp, nil
}

//go:embed sql/get_user_by_id.sql
var getUserByIDSQL string

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var u domain.User

	err := r.db.QueryRowContext(ctx, getUserByIDSQL, id).Scan(&u.ID, &u.Username, &u.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%w: %s", domain.ErrUserNotFound, err)
		}

		return domain.User{}, fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	return u, nil
}

//go:embed sql/get_user_by_username.sql
var getUserByUsernameSQL string

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	var u domain.User

	err := r.db.QueryRowContext(ctx, getUserByUsernameSQL, username).Scan(&u.ID, &u.Username, &u.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%w: %s", domain.ErrUserNotFound, err)
		}

		return domain.User{}, fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	return u, nil
}

func asPostgresError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}

func mapPostgresError(err *pgconn.PgError) error {
	if err.Code == "23505" {
		switch err.ConstraintName {
		case constraintUsernameUnique:
			return domain.ErrUsernameExists
		default:
			return domain.ErrUnknownConflict
		}
	}

	return domain.ErrInternal
}
