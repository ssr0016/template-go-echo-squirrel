package main

// @title           Template Go Echo API
// @version         1.0
// @description     Go Echo + Squirrel + pgx template with session auth
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name app_session

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/ssr0016/template/docs"
	"github.com/ssr0016/template/internal/apperror"
	"github.com/ssr0016/template/internal/config"
	"github.com/ssr0016/template/internal/database"
	"github.com/ssr0016/template/internal/handler"
	"github.com/ssr0016/template/internal/logger"
	ourmiddleware "github.com/ssr0016/template/internal/middleware"
	"github.com/ssr0016/template/internal/repository"
	"github.com/ssr0016/template/internal/router"
	"github.com/ssr0016/template/internal/service"
	"github.com/ssr0016/template/internal/session"
	"github.com/ssr0016/template/internal/validator"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Initialize logger
	log := logger.New(logger.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
	})
	slog.SetDefault(log)

	log.Info("starting server",
		"env", cfg.App.Env,
		"stage", cfg.Stage.String(),
		"log_level", cfg.Logging.Level,
		"log_format", cfg.Logging.Format,
	)

	// Database
	ctx := context.Background()
	db, err := database.NewWithURL(ctx, cfg.Database.URL)
	if err != nil {
		return err
	}
	defer db.Pool.Close()

	if !cfg.IsProduction() {
		if err := database.RunMigrations(db.Pool); err != nil {
			return err
		}
	}

	// Session
	sm := session.New(db.Pool)

	// Repos
	userRepo := repository.NewUserRepo(db)

	// Services
	authService := service.NewAuthService(userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, sm)
	userHandler := handler.NewUserHandler(userRepo)

	// Echo
	e := echo.New()
	e.Validator = validator.New()
	e.HTTPErrorHandler = apperror.ErrorHandler
	e.HideBanner = true
	e.HidePort = true

	// Global middleware
	e.Use(middleware.RequestID())
	e.Use(ourmiddleware.SlogLogger(log))
	e.Use(ourmiddleware.CSRFProtection())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Public routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API routes
	router.Setup(e, sm, authHandler, userHandler)

	// Wrap whole app with scs.LoadAndSave
	scsHandler := sm.LoadAndSave(e)

	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      scsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server starting",
			"port", cfg.App.Port,
			"api_url", cfg.App.URL,
			"swagger_url", cfg.App.URL+"/swagger/index.html",
			"health_url", cfg.App.URL+"/health",
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			quit <- syscall.SIGTERM
		}
	}()

	<-quit

	log.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", "error", err)
		return err
	}

	log.Info("server stopped")
	return nil
}
