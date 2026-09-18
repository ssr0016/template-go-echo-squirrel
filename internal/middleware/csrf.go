package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ssr0016/template/internal/apperror"
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
					HttpOnly: true,                              // #nosec G124 - not sensitive, readable by JS needed for header
					SameSite: http.SameSiteStrictMode,           // #nosec G124 - CSRF protection
					Secure:   isSecure(c),                       // #nosec G124 - set based on environment
					MaxAge:   86400,
				})
			}

			c.Set(csrfTokenKey, token)

			if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
				return next(c)
			}

			headerToken := c.Request().Header.Get(csrfHeaderName)
			if headerToken == "" || headerToken != token {
				return apperror.Forbidden("CSRF token invalid or missing")
			}

			return next(c)
		}
	}
}

// isSecure returns true if the request is over HTTPS.
func isSecure(c echo.Context) bool {
	return c.Scheme() == "https"
}

// GetCSRFToken returns the current CSRF token.
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
