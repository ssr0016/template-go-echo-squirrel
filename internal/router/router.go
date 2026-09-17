package router

import (
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"

	"github.com/ssr0016/template/internal/handler"
	ourmiddleware "github.com/ssr0016/template/internal/middleware"
)

func Setup(e *echo.Echo, sm *scs.SessionManager, ah *handler.AuthHandler, uh *handler.UserHandler) {
	api := e.Group("/api/v1")

	auth := api.Group("/auth")

	// Rate limit: 5 requests per minute, burst 3
	authLimit := ourmiddleware.RateLimit(rate.Limit(5.0/60.0), 3)

	auth.POST("/register", ah.Register, authLimit)
	auth.POST("/login", ah.Login, authLimit)

	authProtected := auth.Group("", ourmiddleware.RequireAuth(sm))
	authProtected.POST("/logout", ah.Logout)
	authProtected.GET("/me", ah.Me)

	users := api.Group("/users", ourmiddleware.RequireAuth(sm))
	users.GET("", uh.ListUsers)
	users.GET("/:id", uh.GetUser)
}
