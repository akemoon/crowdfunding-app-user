package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
	"github.com/akemoon/golib/pglib"
	"github.com/google/uuid"
)

const (
	usersUsernameUnique = "users_username_unique"
	credentialsEmailUnique = "credentials_email_unique"

	followsPK          = "follows_pk"
	followsNoSelfFollow = "follows_no_self_follow"
	followsFollowerFK  = "follows_follower_fk"
	followsFolloweeFK  = "follows_followee_fk"
)

var createUserConstraints = map[string]error{
	usersUsernameUnique:    domain.ErrUsernameExists,
	credentialsEmailUnique: domain.ErrEmailExists,
}

var followConstraints = map[string]error{
	followsPK:           domain.ErrAlreadyFollowing,
	followsNoSelfFollow: domain.ErrFollowToYourself,
	followsFollowerFK:   domain.ErrFollowerNotFound,
	followsFolloweeFK:   domain.ErrFolloweeNotFound,
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

//go:embed sql/create_user.sql
var createUserSQL string

//go:embed sql/create_credentials.sql
var createCredentialsSQL string

func (r *UserRepo) CreateUser(ctx context.Context, req domain.CreateUserReq, passwordHash string) (uuid.UUID, error) {
	t, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.UUID{}, err
	}

	defer func() {
		if err != nil {
			_ = t.Rollback()
		}
	}()

	var userID uuid.UUID

	err = t.QueryRowContext(ctx, createUserSQL, req.Username).Scan(&userID)
	if err != nil {
		return uuid.UUID{}, pglib.MapConstraintErr(err, createUserConstraints, err)
	}

	_, err = t.ExecContext(ctx, createCredentialsSQL, userID, req.Email, passwordHash)
	if err != nil {
		return uuid.UUID{}, pglib.MapConstraintErr(err, createUserConstraints, err)
	}

	err = t.Commit()
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

//go:embed sql/get_user_by_id.sql
var getUserByIDSQL string

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var u domain.User

	err := r.db.QueryRowContext(ctx, getUserByIDSQL, id).Scan(
		&u.ID,
		&u.Username,
		&u.DisplayName,
		&u.Description,
		&u.AvatarKey,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%w: %s", domain.ErrNotFound, err)
		}
		return domain.User{}, err
	}

	return u, nil
}

//go:embed sql/list_users.sql
var listUsersSQL string

func (r *UserRepo) ListUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, listUsersSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		err = rows.Scan(
			&u.ID,
			&u.Username,
			&u.DisplayName,
			&u.Description,
			&u.AvatarKey,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

//go:embed sql/update_profile.sql
var updateProfileSQL string

func (r *UserRepo) UpdateProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileReq) (domain.User, error) {
	var u domain.User

	err := r.db.QueryRowContext(ctx, updateProfileSQL, userID, req.DisplayName, req.Description).Scan(
		&u.ID,
		&u.Username,
		&u.DisplayName,
		&u.Description,
		&u.AvatarKey,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
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
	return err
}

//go:embed sql/list_follows.sql
var listFollowsSQL string

func (r *UserRepo) GetFollows(ctx context.Context, userID uuid.UUID) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, listFollowsSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		err = rows.Scan(
			&u.ID,
			&u.Username,
			&u.DisplayName,
			&u.Description,
			&u.AvatarKey,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}
