package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ssr0016/template/internal/apperror"
	"github.com/ssr0016/template/internal/model"
	"github.com/ssr0016/template/internal/repository"
	"github.com/ssr0016/template/pkg/pagination"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	params := pagination.FromContext(c)
	emailFilter := c.QueryParam("email")

	users, total, err := h.repo.ListWithPagination(
		c.Request().Context(),
		emailFilter,
		params.Page,
		params.Limit,
	)
	if err != nil {
		return apperror.Internal("failed to list users").WithError(err)
	}

	// Convert to response DTOs
	out := make([]model.UserResponse, 0, len(users))
	for i := range users {
		out = append(out, users[i].ToResponse())
	}

	return c.JSON(http.StatusOK, pagination.NewResponse(out, params, total))
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
