package handler

import (
	"errors"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/ssr0016/template/internal/model"
	"github.com/ssr0016/template/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	sm          *scs.SessionManager
}

func NewAuthHandler(authService *service.AuthService, sm *scs.SessionManager) *AuthHandler {
	return &AuthHandler{authService: authService, sm: sm}
}

// Register godoc
// @Summary      Register new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body model.RegisterRequest true "Registration payload"
// @Success      201  {object}  model.UserResponse
// @Failure      409  {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req model.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	user, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, user.ToResponse())
}

// Login godoc
// @Summary      Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body model.LoginRequest true "Login payload"
// @Success      200  {object}  model.UserResponse
// @Failure      401  {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	user, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := h.sm.RenewToken(c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "session error")
	}
	h.sm.Put(c.Request().Context(), "user_id", user.ID)
	return c.JSON(http.StatusOK, user.ToResponse())
}

// Logout godoc
// @Summary      Logout user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	if err := h.sm.Destroy(c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "logged out"})
}

// Me godoc
// @Summary      Get current user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]int64
// @Failure      401  {object}  map[string]string
// @Security     CookieAuth
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c echo.Context) error {
	userID := h.sm.GetInt64(c.Request().Context(), "user_id")
	if userID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "not authenticated")
	}
	return c.JSON(http.StatusOK, map[string]int64{"user_id": userID})
}
