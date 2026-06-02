package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"backend-go/internal/dto"
	"backend-go/internal/service"
)

type CheckoutHandler struct {
	svc service.CheckoutService
}

func NewCheckoutHandler(svc service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{svc: svc}
}

func (h *CheckoutHandler) GetUserOrders(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	orders, err := h.svc.GetUserOrders(c.Request().Context(), userID.(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    orders,
	})
}

func (h *CheckoutHandler) CreateOrder(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	var req dto.CreateOrderRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body format",
		})
	}

	req.UserID = userID.(string)

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "missing or invalid required fields",
			"error":   err.Error(),
		})
	}

	if req.PaymentMethod != "TRANSFER" && req.PaymentMethod != "COD" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid payment method",
		})
	}

	order, err := h.svc.CreateOrder(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Order created successfully",
		"data":    order,
	})
}

func (h *CheckoutHandler) CancelOrder(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	orderID := c.Param("id")

	var reqBody struct {
		Reason string `json:"reason"`
	}

	_ = c.Bind(&reqBody)

	order, err := h.svc.CancelOrder(c.Request().Context(), orderID, userID.(string), reqBody.Reason)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Order cancelled successfully, stock has been restored",
		"data":    order,
	})
}
