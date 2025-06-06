package service

import (
	"context"
	"github.com/audricimanuel/errorutils"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth/repository"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/tools"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/constants"
)

type (
	AuthService interface {
		SignUpUser(ctx context.Context, request model.UserSignUpRequest) (*string, error)
		LoginUser(ctx context.Context, request model.UserLoginRequest) (*string, error)
		VerifyEmail(ctx context.Context, userId, token string) error
	}

	AuthServiceImpl struct {
		cfg            config.Config
		smtpClient     tools.SMTPClient
		authRepository repository.AuthRepository
	}
)

func NewAuthService(cfg config.Config, smtp tools.SMTPClient, a repository.AuthRepository) AuthService {
	return &AuthServiceImpl{
		cfg:            cfg,
		smtpClient:     smtp,
		authRepository: a,
	}
}

func (a *AuthServiceImpl) SignUpUser(ctx context.Context, request model.UserSignUpRequest) (*string, error) {
	userData, err := a.authRepository.AddUser(ctx, request.Email, request.FullName, request.Password)
	if err != nil {
		return nil, err
	}

	// send verification email
	a.SendVerificationEmail(ctx, *userData)

	jwt, err := a.generateJWT(*userData)
	if err != nil {
		return nil, err
	}

	return &jwt, nil
}

func (a *AuthServiceImpl) LoginUser(ctx context.Context, request model.UserLoginRequest) (*string, error) {
	userData, err := a.authRepository.LoginUser(ctx, request.Email, request.Password)
	if err != nil {
		return nil, errorutils.ErrorNotFound.CustomMessage("user not found")
	}

	if !userData.IsActive || userData.DeletedAt != nil {
		return nil, errorutils.ErrorNotFound.CustomMessage("user not found")
	}

	jwt, err := a.generateJWT(*userData)
	if err != nil {
		return nil, err
	}

	return &jwt, nil
}

func (a *AuthServiceImpl) generateJWT(userData model.UserInfoResponse) (string, error) {
	jwtSecret := []byte(a.cfg.JWTSecret)
	token := userData.ToJWT(a.cfg)

	return token.SignedString(jwtSecret)
}

func (a *AuthServiceImpl) SendVerificationEmail(ctx context.Context, userData model.UserInfoResponse) error {
	otpCode, _ := utils.GenerateOTP(6)
	message := "Here's your verification code:\n" + otpCode
	message += "\nDon't share this code to anyone else, including me lol."

	err := a.smtpClient.SendEmail("", tools.EMAIL_TYPE_OTP, constants.SUBJECT_OTP_SIGNUP, message, []string{userData.Email}, nil)
	if err != nil {
		return err
	}

	a.authRepository.SaveOTP(ctx, userData.Id, otpCode, auth.SIGNUP_ACTION)

	return nil
}

func (a *AuthServiceImpl) VerifyEmail(ctx context.Context, userId, token string) error {
	isValidOTP := a.authRepository.IsValidOTP(ctx, userId, token, auth.SIGNUP_ACTION)
	if !isValidOTP {
		return errorutils.ErrorBadRequest.CustomMessage("invalid OTP")
	}

	return nil
}
