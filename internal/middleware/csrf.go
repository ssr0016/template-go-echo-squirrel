package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	csrfTokenKey   = "csrf_token"
	csrfHeaderName = "X-CSRF-Token"
	csrfCookieName = "csrf_token"
)

// CSRFProtection implements the double-submit cookie pattern.
func CSRFProtection() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method

			cookie, err := c.Cookie(csrfCookieName)
			var token string
			if err == nil && cookie.Value != "" {
				token = cookie.Value
			}

			if token == "" {
				token = generateCSRFToken()
				c.SetCookie(&http.Cookie{
					Name:     csrfCookieName,
					Value:    token,
					Path:     "/",
					HttpOnly: false,
					SameSite: http.SameSiteStrictMode,
					Secure:   false,
					MaxAge:   86400,
				})
			}

			c.Set(csrfTokenKey, token)

			if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
				return next(c)
			}

			headerToken := c.Request().Header.Get(csrfHeaderName)
			if headerToken == "" || headerToken != token {
				return echo.NewHTTPError(http.StatusForbidden, "CSRF token invalid or missing")
			}

			return next(c)
		}
	}
}

func GetCSRFToken(c echo.Context) string {
	if token, ok := c.Get(csrfTokenKey).(string); ok {
		return token
	}
	return ""
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
