package http

import (
	"fmt"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/middleware"
	authController "github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func RegisterRouter(
	cfg config.Config,
	// register new controllers here
	authController authController.AuthController,
) *gin.Engine {
	r := gin.Default()

	mid := middleware.InitMiddleware(cfg)

	setMiddlewareGlobal(mid, r)

	// Swagger
	r.Handle("GET", "/docs/*any", mid.BasicAuth(cfg.SwaggerUsername, cfg.SwaggerPassword), ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Handle("GET", "/docs", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})

	r.GET("/ping", func(ctx *gin.Context) {
		staticText := fmt.Sprintf("hello world: %s", cfg.Env)
		ctx.JSON(http.StatusOK, gin.H{"message": staticText})
	})

	api := r.Group("/api")
	{
		// api/v1/auth
		authApi := api.Group("/v1/auth")
		{
			// api/v1/auth/login
			authApi.POST("/login", authController.Login)
			// api/v1/auth/signup
			authApi.POST("/signup", authController.SignUp)
			// api/v1/auth/verify-email
			authApi.POST("/verify-email", authController.VerifyEmail)
			// api/v1/auth/forgot-password
			authApi.POST("/forgot-password", authController.ForgotPassword)
		}
	}

	return r
}

func setMiddlewareGlobal(mid middleware.GoMiddleware, r *gin.Engine) {
	// Logger
	r.Use(mid.LogRequest())

	// Cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowCredentials: false,
		AllowHeaders:     []string{"*"},
		MaxAge:           300,
	}))

	// Recovery
	r.Use(mid.RecoverPanic())
}
