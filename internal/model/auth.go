package model

import (
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type (
	UserLoginRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	UserSignUpRequest struct {
		FullName        string `json:"full_name" validate:"required"`
		Email           string `json:"email" validate:"required"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}

	UserVerifyEmailRequest struct {
		Token string `json:"token" validate:"required"`
	}

	UserForgotPasswordRequest struct {
		Email string `json:"email" validate:"required"`
	}
)

type (
	UserInfoResponse struct {
		Id         string     `json:"id" db:"id"`
		FullName   string     `json:"full_name" db:"full_name"`
		Email      string     `json:"email" db:"email"`
		Password   string     `json:"-" db:"password"`
		Role       int        `json:"role" db:"role"`
		IsVerified bool       `json:"is_verified" db:"is_verified"`
		IsActive   bool       `json:"-" db:"is_active"`
		CreatedAt  time.Time  `json:"-" db:"created_at"`
		UpdatedAt  time.Time  `json:"-" db:"updated_at"`
		DeletedAt  *time.Time `json:"-" db:"deleted_at"`
		LastLogin  time.Time  `json:"-" db:"last_login"`
	}

	UserClaims struct {
		UserId string `json:"user_id"`
		Email  string `json:"email"`
		Role   int    `json:"role"`
		jwt.RegisteredClaims
	}
)

func (u *UserInfoResponse) ToJWT(cfg config.Config) *jwt.Token {
	expireDuration := time.Duration(cfg.JWTExpirationDuration) * time.Hour
	claims := UserClaims{
		UserId: u.Id,
		Email:  u.Email,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "laundry-routine-tracking-service",
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func (u *UserClaims) IsExpired() bool {
	return u.ExpiresAt.Time.Before(time.Now())
}
