package repository

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/audricimanuel/errorutils"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/database"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/logging"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/constants"
	"time"
)

type (
	AuthRepository interface {
		AddUser(ctx context.Context, email, fullName, password string) (*model.UserInfoResponse, error)
		LoginUser(ctx context.Context, email, password string) (*model.UserInfoResponse, error)
		SaveOTP(ctx context.Context, userId, otp string, action auth.OTPAction) error
		IsValidOTP(ctx context.Context, userId, otp string, action auth.OTPAction) bool
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

func (a *AuthRepositoryImpl) AddUser(ctx context.Context, email, fullName, password string) (*model.UserInfoResponse, error) {
	log := logging.WithContext(ctx)

	hashedPassword, _ := utils.HashPassword(password)

	query, args := squirrel.Insert("users").
		Columns(
			"id", "full_name", "email", "password", "role",
		).
		Values(
			utils.GenerateCleanUUID(), fullName, email, hashedPassword, constants.ROLE_USER,
		).
		Suffix("RETURNING id, full_name, email, password, role, is_verified, is_active, created_at, updated_at, deleted_at, last_login").
		PlaceholderFormat(squirrel.Dollar).MustSql()

	var result model.UserInfoResponse
	if err := a.db.PostgresDBSqlx.QueryRowxContext(ctx, query, args...).StructScan(&result); err != nil {
		log.Error("error when add user:", err.Error())
		errDb := errorutils.DefineSQLError(err)
		if errors.Is(errDb, errorutils.ErrorDuplicateData) {
			return nil, errorutils.ErrorDuplicateData.CustomMessage("this email has been used, please login using your email")
		}
		return nil, errDb
	}

	return &result, nil
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

	queryUpdate, args := squirrel.Update("users").
		Set("last_login", time.Now()).
		Where(squirrel.Eq{"id": user.Id}).
		PlaceholderFormat(squirrel.Dollar).MustSql()

	a.db.PostgresDBSqlx.ExecContext(ctx, queryUpdate, args...)

	return &user, nil
}

func (a *AuthRepositoryImpl) SaveOTP(ctx context.Context, userId, otp string, action auth.OTPAction) error {
	currentTime := utils.TimeNow()
	query, args := squirrel.Insert("otps").
		Columns("user_id", "otp_code", "created_at", "expired_at", "action").
		Values(userId, otp, currentTime, currentTime.Add(5*time.Minute), action).
		PlaceholderFormat(squirrel.Dollar).
		MustSql()

	if _, err := a.db.PostgresDBSqlx.ExecContext(ctx, query, args...); err != nil {
		logging.WithContext(ctx).Error("error when add otp:", err.Error())
		return errorutils.DefineSQLError(err)
	}

	return nil
}

func (a *AuthRepositoryImpl) IsValidOTP(ctx context.Context, userId, otp string, action auth.OTPAction) bool {
	log := logging.WithContext(ctx)

	tx, err := a.db.PostgresDBSqlx.BeginTxx(ctx, nil)
	if err != nil {
		log.Error("error when begin transaction:", err.Error())
		return false
	}

	defer tx.Rollback()

	query, args := squirrel.Select("id").
		From("otps").
		Where(squirrel.Eq{"user_id": userId, "otp_code": otp, "action": action, "is_active": true}).
		Where(squirrel.Expr("expired_at > ?", utils.TimeNow().Format(constants.FORMAT_DATETIME_DEFAULT))).
		OrderBy("id DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).MustSql()

	var id int
	if err := tx.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		log.Error("error when check otp:", err.Error())
		return false
	}

	queryUpdate, args := squirrel.Update("otps").
		Set("is_active", false).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).MustSql()

	if _, err := tx.ExecContext(ctx, queryUpdate, args...); err != nil {
		log.Error("error when update otp:", err.Error())
		return false
	}

	return id != 0
}
