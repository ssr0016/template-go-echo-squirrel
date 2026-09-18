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
	"github.com/ssr0016/template/internal/database"
	"github.com/ssr0016/template/internal/handler"
	"github.com/ssr0016/template/internal/logger"
	ourmiddleware "github.com/ssr0016/template/internal/middleware"
	"github.com/ssr0016/template/internal/repository"
	"github.com/ssr0016/template/internal/service"
	"github.com/ssr0016/template/internal/session"
	"github.com/ssr0016/template/internal/validator"
)

func main() {
	_ = godotenv.Load()

	// Initialize logger
	logLevel := getEnv("LOG_LEVEL", "info")
	logFormat := getEnv("LOG_FORMAT", "text")
	log := logger.New(logger.Config{
		Level:  logLevel,
		Format: logFormat,
	})
	slog.SetDefault(log)

	log.Info("starting server",
		"env", getEnv("APP_ENV", "local"),
		"log_level", logLevel,
		"log_format", logFormat,
	)

	ctx := context.Background()
	db, err := database.New(ctx)
	if err != nil {
		log.Error("db init failed", "error", err)
		os.Exit(1)
	}
	defer db.Pool.Close()

	if os.Getenv("APP_ENV") != "production" {
		if err := database.RunMigrations(db.Pool); err != nil {
			log.Error("migrations failed", "error", err)
			os.Exit(1)
		}
	}

	sm := session.New(db.Pool)
	userRepo := repository.NewUserRepo(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, sm)
	userHandler := handler.NewUserHandler(userRepo)

	e := echo.New()
	e.Validator = validator.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.RequestID())
	e.Use(ourmiddleware.SlogLogger(log))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("CORS_ALLOWED_ORIGINS")},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	api := e.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/me", authHandler.Me)

	users := api.Group("/users", ourmiddleware.RequireAuth(sm))
	users.GET("", userHandler.ListUsers)
	users.GET("/:id", userHandler.GetUser)

	// Wrap whole app with scs.LoadAndSave for session handling
	scsHandler := sm.LoadAndSave(e)

	go func() {
		port := getEnv("APP_PORT", "8080")
		server := &http.Server{
			Addr:         ":" + port,
			Handler:      scsHandler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		}
		log.Info("server starting",
			"port", port,
			"api_url", "http://localhost:"+port,
			"swagger_url", "http://localhost:"+port+"/swagger/index.html",
			"health_url", "http://localhost:"+port+"/health",
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", "error", err)
	}

	log.Info("server stopped")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
