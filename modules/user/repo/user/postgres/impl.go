package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/golib/pglib"
	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
	"github.com/google/uuid"
)

const (
	usersUsernameUnique = "users_username_unique"

	credentialsEmailUnique = "credentials_email_unique"

	followsPK           = "follows_pk"
	followsNoToYourself = "follows_no_to_yourself"
	followsFollowerFK   = "follows_follower_fk"
	followsFolloweeFK   = "follows_followee_fk"
)

var createUserConstraints = map[string]error{
	credentialsEmailUnique: lib.ErrEmailExists,
	usersUsernameUnique:    lib.ErrUsernameExists,
}

var updateProfileConstraints = map[string]error{
	usersUsernameUnique: lib.ErrUsernameExists,
}

var followConstraints = map[string]error{
	followsPK:           domain.ErrAlreadyFollowing,
	followsNoToYourself: domain.ErrFollowToYourself,
	followsFollowerFK:   domain.ErrFollowerNotFound,
	followsFolloweeFK:   domain.ErrFolloweeNotFound,
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

//go:embed sql/create_credentials.sql
var createCredentialsSQL string

//go:embed sql/create_user.sql
var createUserSQL string

func (r *UserRepo) CreateUser(ctx context.Context, req lib.CreateUserReq) (uuid.UUID, error) {
	t, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.UUID{}, err
	}

	defer func() {
		if err != nil {
			// TODO: handle err
			_ = t.Rollback()
		}
	}()

	var userID uuid.UUID

	err = t.QueryRowContext(ctx, createCredentialsSQL, req.Email, req.PasswordHash).Scan(&userID)
	if err != nil {
		return uuid.UUID{}, pglib.MapConstraintErr(err, createUserConstraints, err)
	}

	_, err = t.ExecContext(ctx, createUserSQL, userID, req.Username)
	if err != nil {
		return uuid.UUID{}, pglib.MapConstraintErr(err, createUserConstraints, err)
	}

	err = t.Commit()
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

//go:embed sql/get_credentials_by_email.sql
var getCredentialsByEmailSQL string

func (r *UserRepo) GetCredentialsByEmail(ctx context.Context, email string) (lib.UserCredentials, error) {
	var c lib.UserCredentials

	err := r.db.QueryRowContext(ctx, getCredentialsByEmailSQL, email).Scan(&c.UserID, &c.PasswordHash, &c.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.UserCredentials{}, lib.ErrInvalidCredentials
		}
		return lib.UserCredentials{}, err
	}

	return c, nil
}

//go:embed sql/get_credentials_by_id.sql
var getCredentialsByIDSQL string

func (r *UserRepo) GetCredentialsByID(ctx context.Context, id uuid.UUID) (lib.UserCredentials, error) {
	var c lib.UserCredentials

	err := r.db.QueryRowContext(ctx, getCredentialsByIDSQL, id).Scan(&c.UserID, &c.PasswordHash, &c.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.UserCredentials{}, lib.ErrNotFound
		}
		return lib.UserCredentials{}, err
	}

	return c, nil
}

//go:embed sql/get_user_by_id.sql
var getUserByIDSQL string

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var u domain.User

	err := r.db.QueryRowContext(
		ctx, getUserByIDSQL, id,
	).Scan(
		&u.ID,
		&u.Username,
		&u.DisplayName,
		&u.Description,
		&u.AvatarUrl,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%w: %s", lib.ErrNotFound, err)
		}
		return domain.User{}, err
	}

	return u, nil
}

//go:embed sql/update_profile.sql
var updateProfileSQL string

func (r *UserRepo) UpdateProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileReq) (domain.User, error) {
	var u domain.User

	err := r.db.QueryRowContext(
		ctx, updateProfileSQL, userID, req.Username, req.DisplayName, req.Description,
	).Scan(
		&u.ID,
		&u.Username,
		&u.DisplayName,
		&u.Description,
		&u.AvatarUrl,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, lib.ErrNotFound
		}
		// TODO: maybe use default error like internal
		return domain.User{}, pglib.MapConstraintErr(err, updateProfileConstraints, err)
	}

	return u, nil
}

//go:embed sql/follow.sql
var followSQL string

func (r *UserRepo) Follow(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, followSQL, followerID, followeeID)
	if err != nil {
		return pglib.MapConstraintErr(err, followConstraints, err)
	}
	return nil
}

//go:embed sql/unfollow.sql
var unfollowSQL string

func (r *UserRepo) Unfollow(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, unfollowSQL, followerID, followeeID)
	if err != nil {
		return err
	}
	return nil
}

//go:embed sql/search_users.sql
var searchUsersSQL string

func (r *UserRepo) SearchUsers(ctx context.Context, req domain.SearchUsersReq) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, searchUsersSQL, req.Query, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.User

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.DisplayName,
			&u.Description,
			&u.AvatarUrl,
		); err != nil {
			return nil, err
		}
		out = append(out, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

//go:embed sql/update_role.sql
var updateRoleSQL string

func (r *UserRepo) UpdateRole(ctx context.Context, userID uuid.UUID, role string) error {
	roleID, err := MapRoleToDB(role)
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, updateRoleSQL, userID, roleID)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return lib.ErrNotFound
	}

	return nil
}
