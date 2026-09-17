package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ssr0016/template/internal/repository"
)

type UserHandler struct{ repo *repository.UserRepo }

func NewUserHandler(repo *repository.UserRepo) *UserHandler {
	return &UserHandler{repo: repo}
}

// ListUsers godoc
// @Summary      List users
// @Tags         users
// @Produce      json
// @Security     CookieAuth
// @Param        email query string false "Filter by email"
// @Param        limit query int false "Limit (default 20, max 100)"
// @Success      200 {array} map[string]interface{}
// @Router       /users [get]
func (h *UserHandler) ListUsers(c echo.Context) error {
	emailFilter := c.QueryParam("email")
	limit := 20
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	users, err := h.repo.List(c.Request().Context(), emailFilter, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	out := make([]map[string]interface{}, 0, len(users))
	for i := range users {
		out = append(out, map[string]interface{}{
			"id": users[i].ID, "email": users[i].Email,
			"name": users[i].Name, "created_at": users[i].CreatedAt,
		})
	}
	return c.JSON(http.StatusOK, out)
}

// GetUser godoc
// @Summary      Get user by ID
// @Tags         users
// @Produce      json
// @Security     CookieAuth
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	user, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, user.ToResponse())
}
