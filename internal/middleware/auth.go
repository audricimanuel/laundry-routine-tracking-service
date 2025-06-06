package middleware

import (
	"fmt"
	"github.com/audricimanuel/errorutils"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth/repository"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/constants"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/httputils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type (
	AuthMiddleware interface {
		ValidateJWT() gin.HandlerFunc
		RefreshJWT() gin.HandlerFunc
		ValidateJWTFromCookie() gin.HandlerFunc
		ValidateGetLoginPage() gin.HandlerFunc
	}

	AuthMiddlewareImpl struct {
		cfg      config.Config
		authRepo repository.AuthRepository
	}
)

func NewAuthMiddleware(cfg config.Config, a repository.AuthRepository) AuthMiddleware {
	return &AuthMiddlewareImpl{
		cfg:      cfg,
		authRepo: a,
	}
}

func (a *AuthMiddlewareImpl) ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader(constants.HEADER_AUTHORIZATION)
		if bearerToken == "" {
			httputils.SetHttpResponse(c, nil, errorutils.ErrorInvalidToken, nil)
			c.Abort()
			return
		}

		splitBearer := strings.SplitN(bearerToken, "Bearer ", 2)
		if len(splitBearer) < 2 {
			httputils.SetHttpResponse(c, nil, errorutils.ErrorInvalidToken, nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(splitBearer[1])
		claims, err := a.validateToken(tokenString)
		if err != nil {
			httputils.SetHttpResponse(c, nil, err, nil)
			c.Abort()
			return
		}

		if claims.IsExpired() {
			httputils.SetHttpResponse(c, nil, errorutils.ErrorInvalidToken.CustomMessage("expired token"), nil)
			c.Abort()
			return
		}

		c.Set(constants.USER_DATA, *claims)
		c.Set(constants.USER_TOKEN, tokenString)

		c.Next()
	}
}

func (a *AuthMiddlewareImpl) ValidateJWTFromCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(constants.COOKIE_AUTH_TOKEN)
		if err != nil || cookie.Value == "" {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		tokenString := cookie.Value

		claims, err := a.validateToken(tokenString)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		if claims.IsExpired() {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		c.Set(constants.USER_DATA, *claims)
		c.Set(constants.USER_TOKEN, tokenString)

		c.Next()
	}
}

func (a *AuthMiddlewareImpl) RefreshJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader(constants.HEADER_AUTHORIZATION)
		if bearerToken == "" {
			httputils.SetHttpResponse(c, nil, errorutils.ErrorInvalidToken, nil)
			c.Abort()
			return
		}

		splitBearer := strings.SplitN(bearerToken, "Bearer ", 2)
		if len(splitBearer) < 2 {
			httputils.SetHttpResponse(c, nil, errorutils.ErrorInvalidToken, nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(splitBearer[1])
		claims, err := a.validateToken(tokenString)
		if err != nil {
			httputils.SetHttpResponse(c, nil, err, nil)
			c.Abort()
			return
		}

		data := map[string]string{
			"token": tokenString,
		}

		if claims.IsExpired() {
			claims.ExpiresAt.Time = utils.TimeNow().Add(time.Duration(a.cfg.JWTExpirationDuration) * time.Hour)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			data["token"], _ = token.SignedString([]byte(a.cfg.JWTSecret))
		}

		httputils.SetHttpResponse(c, data, nil, nil)
	}
}

func (a *AuthMiddlewareImpl) ValidateGetLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(constants.COOKIE_AUTH_TOKEN)
		if err != nil || cookie.Value == "" {
			c.Next()
			return
		}

		tokenString := cookie.Value

		claims, err := a.validateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		if claims.IsExpired() {
			c.Next()
			return
		}

		//c.Set(constants.USER_DATA, *claims)
		//c.Set(constants.USER_TOKEN, tokenString)

		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

func (a *AuthMiddlewareImpl) validateToken(jwtString string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(jwtString, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorutils.ErrorInvalidToken.CustomMessage(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(a.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, errorutils.ErrorInvalidToken
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok || !token.Valid {
		return nil, errorutils.ErrorInvalidToken.CustomMessage("invalid claims")
	}

	return claims, nil
}
