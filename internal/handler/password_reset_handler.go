package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ssr0016/template/internal/apperror"
	"github.com/ssr0016/template/internal/service"
)

type PasswordResetHandler struct {
	service *service.PasswordResetService
}

func NewPasswordResetHandler(svc *service.PasswordResetService) *PasswordResetHandler {
	return &PasswordResetHandler{service: svc}
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

func (h *PasswordResetHandler) ForgotPassword(c echo.Context) error {
	var req ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return apperror.BadRequest("invalid request body").WithError(err)
	}
	if err := c.Validate(&req); err != nil {
		return apperror.Validation(err.Error())
	}

	if err := h.service.RequestReset(c.Request().Context(), req.Email); err != nil {
		return apperror.Internal("failed to process request").WithError(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "if the email exists, a reset link has been sent",
	})
}

func (h *PasswordResetHandler) ResetPassword(c echo.Context) error {
	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return apperror.BadRequest("invalid request body").WithError(err)
	}
	if err := c.Validate(&req); err != nil {
		return apperror.Validation(err.Error())
	}

	if err := h.service.ResetPassword(c.Request().Context(), req.Token, req.NewPassword); err != nil {
		if err == service.ErrInvalidResetToken {
			return apperror.BadRequest("invalid or expired token")
		}
		return apperror.Internal("password reset failed").WithError(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "password reset successful",
	})
}
