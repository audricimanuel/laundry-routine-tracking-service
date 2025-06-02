package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/audricimanuel/errorutils"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/database"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/logging"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/constants"
)

type (
	AuthRepository interface {
		AddUser(ctx context.Context, email, fullName, password string) error
		LoginUser(ctx context.Context, email, password string) (*model.UserInfoResponse, error)
	}

	AuthRepositoryImpl struct {
		cfg config.Config
		db  database.DBCollection
	}
)

func NewAuthRepository(cfg config.Config, db database.DBCollection) AuthRepository {
	return &AuthRepositoryImpl{
		cfg: cfg,
		db:  db,
	}
}

func (a *AuthRepositoryImpl) AddUser(ctx context.Context, email, fullName, password string) error {
	log := logging.WithContext(ctx)

	hashedPassword, _ := utils.HashPassword(password)

	query, args := squirrel.Insert("users").
		Columns(
			"id", "full_name", "email", "password", "role",
		).
		Values(
			utils.GenerateCleanUUID(), fullName, email, hashedPassword, constants.ROLE_USER,
		).
		PlaceholderFormat(squirrel.Dollar).MustSql()

	if _, err := a.db.PostgresDBSqlx.ExecContext(ctx, query, args...); err != nil {
		log.Error("error when add user:", err.Error())
		return errorutils.DefineSQLError(err)
	}

	return nil
}

func (a *AuthRepositoryImpl) LoginUser(ctx context.Context, email, password string) (*model.UserInfoResponse, error) {
	log := logging.WithContext(ctx)

	query, args := squirrel.Select("id", "full_name", "email", "password", "role", "is_verified",
		"is_active", "created_at", "updated_at", "deleted_at", "last_login").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).
		MustSql()

	var user model.UserInfoResponse
	if err := a.db.PostgresDBSqlx.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		log.Error("error when login user:", err.Error())
		return nil, errorutils.DefineSQLError(err)
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errorutils.ErrorUnauthorized.CustomMessage("invalid email or password")
	}

	return &user, nil
}
