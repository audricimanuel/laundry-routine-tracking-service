package service

import (
	"context"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth/repository"
)

type (
	AuthService interface {
		SignUpUser(ctx context.Context, request model.UserSignUpRequest) error
		LoginUser(ctx context.Context, request model.UserLoginRequest) (*string, error)
	}

	AuthServiceImpl struct {
		cfg            config.Config
		authRepository repository.AuthRepository
	}
)

func NewAuthService(cfg config.Config, a repository.AuthRepository) AuthService {
	return &AuthServiceImpl{
		cfg:            cfg,
		authRepository: a,
	}
}

func (a *AuthServiceImpl) SignUpUser(ctx context.Context, request model.UserSignUpRequest) error {
	return a.authRepository.AddUser(ctx, request.Email, request.FullName, request.Password)
}

func (a *AuthServiceImpl) LoginUser(ctx context.Context, request model.UserLoginRequest) (*string, error) {
	userData, err := a.authRepository.LoginUser(ctx, request.Email, request.Password)
	if err != nil {
		return nil, err
	}

	jwt, err := a.generateJWT(*userData)
	if err != nil {
		return nil, err
	}

	return &jwt, nil
}

func (a *AuthServiceImpl) generateJWT(userData model.UserInfoResponse) (string, error) {
	jwtSecret := []byte(a.cfg.JWTSecret)
	token := userData.ToJWT()

	return token.SignedString(jwtSecret)
}
