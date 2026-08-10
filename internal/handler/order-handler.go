package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"backend-go/internal/service"
)

type OrderHandler struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) GetUserOrders(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	statusQuery := c.QueryParam("status")
	var status *string
	if statusQuery != "" {
		status = &statusQuery
	}

	orders, err := h.svc.GetUserOrders(c.Request().Context(), userID.(string), status)
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

func (h *OrderHandler) GetOrderDetail(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized, missing user context",
		})
	}

	orderID := c.Param("id")
	order, err := h.svc.GetOrderDetail(c.Request().Context(), orderID, userID.(string))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    order,
	})
}
