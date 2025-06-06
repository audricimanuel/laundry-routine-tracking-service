package controller

import (
	"fmt"
	"github.com/audricimanuel/errorutils"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth/service"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/constants"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/httputils"
	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type (
	AuthController interface {
		GetLoginPage(ctx *gin.Context)
		Login(ctx *gin.Context)
		SignUp(ctx *gin.Context)
		ForgotPassword(ctx *gin.Context)
		VerifyEmail(ctx *gin.Context)
	}

	AuthControllerImpl struct {
		cfg         config.Config
		authService service.AuthService
	}
)

func NewAuthController(cfg config.Config, a service.AuthService) AuthController {
	return &AuthControllerImpl{
		cfg:         cfg,
		authService: a,
	}
}

func (a *AuthControllerImpl) GetLoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{})
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

	token, err := ctx.Cookie(constants.COOKIE_AUTH_TOKEN)
	if err == nil && token != "" {
		if isValidToken := func(jwtString string) bool {
			token, err := jwt2.ParseWithClaims(jwtString, &model.UserClaims{}, func(token *jwt2.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt2.SigningMethodHMAC); !ok {
					return nil, errorutils.ErrorInvalidToken.CustomMessage(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
				}
				return []byte(a.cfg.JWTSecret), nil
			})
			if err != nil {
				return false
			}

			claims, ok := token.Claims.(*model.UserClaims)
			if !ok || !token.Valid {
				return false
			}

			if claims.IsExpired() {
				return false
			}

			return true
		}(token); isValidToken {
			httputils.SetHttpResponse(ctx, map[string]interface{}{"token": token}, nil, nil)
			return
		}
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     constants.COOKIE_AUTH_TOKEN,
		Value:    *jwt,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  utils.TimeNow().Add(2 * time.Hour),
	})

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

	jwt, err := a.authService.SignUpUser(ctx, request)
	if err != nil {
		httputils.SetHttpResponse(ctx, nil, err, nil)
		return
	}

	httputils.SetHttpResponse(ctx, map[string]interface{}{"token": jwt}, nil, nil)
}

func (a *AuthControllerImpl) ForgotPassword(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthControllerImpl) VerifyEmail(ctx *gin.Context) {
	userClaims, _ := ctx.Get(constants.USER_DATA)
	userData := userClaims.(model.UserClaims)

	var request model.UserVerifyEmailRequest

	if err := errorutils.ValidatePayload(ctx.Request, &request); err != nil {
		httputils.SetHttpResponse(ctx, nil, err, nil)
		return
	}

	if err := a.authService.VerifyEmail(ctx, userData.UserId, request.Token); err != nil {
		httputils.SetHttpResponse(ctx, nil, err, nil)
		return
	}

	httputils.SetHttpResponse(ctx, "Verification success. Please login using your email.", nil, nil)
}
