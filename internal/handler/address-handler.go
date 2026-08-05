package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"backend-go/internal/dto"
	"backend-go/internal/service"
)

type AddressHandler struct {
	svc service.AddressService
}

func NewAddressHandler(svc service.AddressService) *AddressHandler {
	return &AddressHandler{svc: svc}
}

func (h *AddressHandler) GetAddresses(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	addresses, err := h.svc.GetUserAddresses(c.Request().Context(), userID.(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    addresses,
	})
}

func (h *AddressHandler) AddAddress(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	var req dto.CreateAddressRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body format",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid required fields",
			"error":   err.Error(),
		})
	}

	if err := h.svc.AddAddress(c.Request().Context(), userID.(string), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Address added successfully",
	})
}

func (h *AddressHandler) SetPrimary(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	addressID := c.Param("id")
	if err := h.svc.SetPrimaryAddress(c.Request().Context(), userID.(string), addressID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Primary address updated successfully",
	})
}

func (h *AddressHandler) DeleteAddress(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	addressID := c.Param("id")
	if err := h.svc.DeleteAddress(c.Request().Context(), userID.(string), addressID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Address deleted successfully",
	})
}
