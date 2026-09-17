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
	"log"
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
	"github.com/ssr0016/template/internal/repository"
	"github.com/ssr0016/template/internal/service"
	"github.com/ssr0016/template/internal/session"
	"github.com/ssr0016/template/internal/validator"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()
	db, err := database.New(ctx)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer db.Pool.Close()

	if os.Getenv("APP_ENV") != "production" {
		if err := database.RunMigrations(db.Pool); err != nil {
			log.Fatalf("migrations: %v", err)
		}
	}

	sm := session.New(db.Pool)
	userRepo := repository.NewUserRepo(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, sm)
	userHandler := handler.NewUserHandler(userRepo)

	e := echo.New()
	e.Validator = validator.New()
	e.HideBanner = false

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("CORS_ALLOWED_ORIGINS")},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API routes
	api := e.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/me", authHandler.Me)

	users := api.Group("/users")
	users.GET("", userHandler.ListUsers)
	users.GET("/:id", userHandler.GetUser)

	// Wrap whole app with scs.LoadAndSave for session handling
	scsHandler := sm.LoadAndSave(e)

	go func() {
		port := os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}
		server := &http.Server{
			Addr:    ":" + port,
			Handler: scsHandler,
		}
		log.Printf("🚀 API Server:  http://localhost:%s", port)
		log.Printf("📖 Swagger UI:  http://localhost:%s/swagger/index.html", port)
		log.Printf("❤️  Health:      http://localhost:%s/health", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = e.Shutdown(shutdownCtx)
}
