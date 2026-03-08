package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	lib "github.com/akemoon/crowdfunding-app-user/lib/domain"
	"github.com/akemoon/crowdfunding-app-user/lib/pglib"
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

func (r *UserRepo) CreateUser(ctx context.Context, req lib.CreateUserReq) error {
	t, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
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
		return pglib.MapConstraintErr(err, createUserConstraints, err)
	}

	_, err = t.ExecContext(ctx, createUserSQL, userID, req.Username)
	if err != nil {
		return pglib.MapConstraintErr(err, createUserConstraints, err)
	}

	err = t.Commit()
	if err != nil {
		return err
	}

	return nil
}

//go:embed sql/get_credentials_by_email.sql
var getCredentialsByEmailSQL string

func (r *UserRepo) GetCredentialsByEmail(ctx context.Context, email string) (lib.UserCredentials, error) {
	var c lib.UserCredentials

	err := r.db.QueryRowContext(ctx, getCredentialsByEmailSQL, email).Scan(&c.UserID, &c.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lib.UserCredentials{}, lib.ErrInvalidCredentials
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
			return domain.User{}, fmt.Errorf("%w: %s", domain.ErrNotFound, err)
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
			return domain.User{}, domain.ErrNotFound
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
