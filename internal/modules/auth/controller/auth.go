package controller

import (
	"github.com/audricimanuel/errorutils"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth/service"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/httputils"
	"github.com/gin-gonic/gin"
)

type (
	AuthController interface {
		Login(ctx *gin.Context)
		SignUp(ctx *gin.Context)
		ForgotPassword(ctx *gin.Context)
		VerifyEmail(ctx *gin.Context)
	}

	AuthControllerImpl struct {
		authService service.AuthService
	}
)

func NewAuthController(a service.AuthService) AuthController {
	return &AuthControllerImpl{
		authService: a,
	}
}

func (a *AuthControllerImpl) Login(ctx *gin.Context) {
	var request model.UserLoginRequest

	if err := errorutils.ValidatePayload(ctx.Request, &request); err != nil {
		httputils.SetHttpResponse(ctx, nil, err, nil)
		return
	}

	jwt, err := a.authService.LoginUser(ctx, request)
	if err != nil {
		httputils.SetHttpResponse(ctx, nil, err, nil)
		return
	}

	httputils.SetHttpResponse(ctx, map[string]interface{}{"token": jwt}, nil, nil)
}

func (a *AuthControllerImpl) SignUp(ctx *gin.Context) {
	var request model.UserSignUpRequest

	// validate payload
	if err := errorutils.ValidatePayload(ctx.Request, &request); err != nil {
		httputils.SetHttpResponse(ctx, nil, err, nil)
		return
	}

	// validate password confirmation
	if request.Password != request.ConfirmPassword {
		httputils.SetHttpResponse(ctx, nil, errorutils.ErrorBadRequest.CustomMessage("mismatched password confirmation"), nil)
		return
	}

	if err := a.authService.SignUpUser(ctx, request); err != nil {
		httputils.SetHttpResponse(ctx, nil, err, nil)
		return
	}
}

func (a *AuthControllerImpl) ForgotPassword(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthControllerImpl) VerifyEmail(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
