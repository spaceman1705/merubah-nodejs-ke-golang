package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"backend-go/internal/dto"
	"backend-go/internal/service"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body format",
		})
	}

	if err := h.authSvc.Register(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Account registered successfully",
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body format",
		})
	}

	user, err := h.authSvc.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	// sementara kita kembalikan data user sukses login. tidak menggunakkan jwt
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data":    user,
	})
}
