package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"backend-go/internal/dto"
	"backend-go/internal/service"
)

type CartHandler struct {
	svc service.CartService
}

func NewCartHandler(svc service.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

func (h *CartHandler) GetCart(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	carts, err := h.svc.GetCartByUserID(c.Request().Context(), userID.(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    carts,
	})
}

func (h *CartHandler) AddToCart(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	var req dto.AddToCartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body format",
		})
	}

	if err := h.svc.AddToCart(c.Request().Context(), userID.(string), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Item added to cart successfully",
	})
}

func (h *CartHandler) UpdateCart(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	cartID := c.Param("id")
	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body format",
		})
	}

	if err := h.svc.UpdateCartQuantity(c.Request().Context(), cartID, userID.(string), req.Quantity); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Cart updated successfully",
	})
}

func (h *CartHandler) RemoveCartItem(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	cartID := c.Param("id")

	if err := h.svc.RemoveCartItem(c.Request().Context(), cartID, userID.(string)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Item removed from cart successfully",
	})
}
