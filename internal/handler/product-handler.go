package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"backend-go/internal/dto"
	"backend-go/internal/service"
)

type ProductHandler struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	var req dto.ProductQueryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid query parameters",
		})
	}

	products, total, err := h.svc.GetProducts(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"products": products,
			"total":    total,
		},
	})
}

func (h *ProductHandler) GetProductDetail(c echo.Context) error {
	productID := c.Param("id")

	product, err := h.svc.GetProductDetail(c.Request().Context(), productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    product,
	})
}
