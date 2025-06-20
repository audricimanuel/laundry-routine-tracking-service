package main

import (
	"context"
	"fmt"
	"github.com/audricimanuel/laundry-routine-tracking-service/docs"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/database"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/middleware"
	authController "github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth/controller"
	authRepository "github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth/repository"
	authService "github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/auth/service"
	laundryController "github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/laundry/controller"
	laundryRepository "github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/laundry/repository"
	laundryService "github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/laundry/service"
	httpServer "github.com/audricimanuel/laundry-routine-tracking-service/internal/server/http"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/tools"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func setSwaggerInfo() {
	docs.SwaggerInfo.Title = "Laundry Tracking API"
	docs.SwaggerInfo.Description = "Laundry Tracking API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
}

func main() {
	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	// initialize mongodb connection
	databaseCollection := database.NewDatabaseCollection(cfg)
	defer func() {
		databaseCollection.PostgresDBSqlx.Close()
	}()

	// tools
	smtpClient := tools.NewSMTPClient(cfg)

	// repositories
	authRepo := authRepository.NewAuthRepository(cfg, databaseCollection)
	laundryRepo := laundryRepository.NewLaundryRepository(databaseCollection)

	// services
	authServ := authService.NewAuthService(cfg, smtpClient, authRepo)
	laundrySvc := laundryService.NewLaundryService(laundryRepo)

	// controllers
	authCtrl := authController.NewAuthController(cfg, authServ)
	laundryCtrl := laundryController.NewLaundryController(laundrySvc)

	// set swagger info
	setSwaggerInfo()

	// registering router
	router := httpServer.RegisterRouter(
		cfg,

		// register additional middlewares here
		middleware.NewAuthMiddleware(cfg, authRepo),

		// register controllers in here
		authCtrl,
		laundryCtrl,
	)

	// running server
	logrus.Println("[INFO] Loading server")
	runServer(cfg, router)
}

func runServer(cfg config.Config, route http.Handler) {
	// The HTTP Server
	server := &http.Server{
		WriteTimeout: time.Second * time.Duration(cfg.Host.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(cfg.Host.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.Host.IdleTimeout),
		Handler:      route,
	}

	if cfg.Host.Port != "" {
		server.Addr = fmt.Sprintf("%s:%s", cfg.Host.Address, cfg.Host.Port)
	}

	// Run Server
	go func() {
		logrus.Printf("⇨ http server started on %s\n", server.Addr)
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	logrus.Println("received shutdown signal. Trying to shutdown gracefully", sig)

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop Server
	logrus.Println("Stopping HTTP Server")
	server.SetKeepAlivesEnabled(false)
	err := server.Shutdown(ctx)
	if err != nil {
		logrus.Fatal("Failure while shutting down gracefully, errApp: ", err)
	}

	logrus.Println("Shutdown gracefully completed")
}
