package http

import (
	"fmt"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/middleware"
	authController "github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth/controller"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/laundry/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"html/template"
	"net/http"
	"os"
)

func RegisterRouter(
	cfg config.Config,
	// additional middlewares
	authMiddleware middleware.AuthMiddleware,
	// register new controllers here
	authController authController.AuthController,
	laundryController controller.LaundryController,
) *gin.Engine {
	r := gin.Default()

	setHTMLTemplate(r)

	mid := middleware.InitMiddleware(cfg)

	setMiddlewareGlobal(cfg, mid, r)

	// Swagger
	r.Handle("GET", "/docs/*any", mid.BasicAuth(cfg.SwaggerUsername, cfg.SwaggerPassword), ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Handle("GET", "/docs", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})

	r.GET("/ping", func(ctx *gin.Context) {
		staticText := fmt.Sprintf("hello world: %s", cfg.Env)
		ctx.JSON(http.StatusOK, gin.H{"message": staticText})
	})

	// route of FE
	viewApi := r.Group("")
	{
		// /login
		viewApi.GET("/login", authMiddleware.ValidateGetLoginPage(), authController.GetLoginPage)

		viewApi.GET("/", authMiddleware.ValidateJWTFromCookie(), laundryController.GetLaundryList)
	}

	api := r.Group("/api")
	{
		// /api/v1/auth
		authApi := api.Group("/v1/auth")
		{
			// /api/v1/auth/login
			authApi.POST("/login", authController.Login)
			// /api/v1/auth/signup
			authApi.POST("/signup", authController.SignUp)
			// /api/v1/auth/verify-email
			authApi.POST("/verify-email", authMiddleware.ValidateJWT(), authController.VerifyEmail)
			// /api/v1/auth/forgot-password
			authApi.POST("/forgot-password", authController.ForgotPassword)
			// /api/v1/auth/refresh
			authApi.POST("/refresh", authMiddleware.RefreshJWT())
		}

		// /api/v1/laundry
		laundryApi := api.Group("/v1/laundry")
		{
			laundryApi.POST("/", authMiddleware.ValidateJWT(), laundryController.AddLaundry)
		}
	}

	return r
}

func setMiddlewareGlobal(cfg config.Config, mid middleware.GoMiddleware, r *gin.Engine) {
	// Logger
	r.Use(mid.LogRequest())

	// Cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Host.FEBaseUrl},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowCredentials: false,
		AllowHeaders:     []string{"*"},
		MaxAge:           300,
	}))

	// Recovery
	r.Use(mid.RecoverPanic())
}

func setHTMLTemplate(r *gin.Engine) {
	var templates *template.Template

	r.Static("/static/", "./view/static")

	dir, _ := os.Getwd()
	fmt.Println("Current working directory:", dir)
	templates = template.Must(template.ParseGlob(dir + "/view/templates/*.html"))

	r.SetHTMLTemplate(templates)
}
