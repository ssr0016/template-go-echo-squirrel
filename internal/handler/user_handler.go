package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ssr0016/template/internal/apperror"
	"github.com/ssr0016/template/internal/repository"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

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
		return apperror.Internal("failed to list users").WithError(err)
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

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return apperror.BadRequest("invalid user id")
	}

	user, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return apperror.Internal("failed to get user").WithError(err)
	}
	if user == nil {
		return apperror.NotFound("user not found")
	}
	return c.JSON(http.StatusOK, user.ToResponse())
}
