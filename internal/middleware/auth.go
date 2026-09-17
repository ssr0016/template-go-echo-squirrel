package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
)

const UserIDKey = "user_id"

func RequireAuth(sm *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := sm.GetInt64(c.Request().Context(), UserIDKey)
			if userID == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
			}
			return next(c)
		}
	}
}

func GetUserID(c echo.Context, sm *scs.SessionManager) int64 {
	return sm.GetInt64(c.Request().Context(), UserIDKey)
}
