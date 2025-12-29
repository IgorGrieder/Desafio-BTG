package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/domain"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/ports"
)

type OrderHandler struct {
	orderService ports.OrderService
}

func NewOrderHandler(service ports.OrderService) *OrderHandler {
	return &OrderHandler{orderService: service}
}

// GetOrderTotal godoc
// @Summary Get total value of an order
// @Description Get the total value of an order by its code
// @Tags orders
// @Accept json
// @Produce json
// @Param code path int true "Order Code" minimum(1)
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/orders/{code}/total [get]
func (h *OrderHandler) GetOrderTotal(w http.ResponseWriter, r *http.Request) {
	codeStr := r.PathValue("code")

	code, err := strconv.Atoi(codeStr)
	if err != nil || code < 1 {
		RespondError(w, http.StatusBadRequest, "Invalid order code", map[string]string{
			"code": "Order code must be a positive integer",
		})
		return
	}

	// TODO: Call service to get order total
	// For now, return mock data
	RespondJSON(w, http.StatusOK, map[string]any{
		"order_code":  code,
		"total_value": 120.50,
	})
}

// CountCustomerOrders godoc
// @Summary Count orders by customer
// @Description Get the total number of orders for a specific customer
// @Tags customers
// @Accept json
// @Produce json
// @Param code path int true "Customer Code" minimum(1)
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/customers/{code}/orders/count [get]
func (h *OrderHandler) CountCustomerOrders(w http.ResponseWriter, r *http.Request) {
	codeStr := r.PathValue("code")

	code, err := strconv.Atoi(codeStr)
	if err != nil || code < 1 {
		RespondError(w, http.StatusBadRequest, "Invalid customer code", map[string]string{
			"code": "Customer code must be a positive integer",
		})
		return
	}

	// TODO: Call service to count orders
	// For now, return mock data
	RespondJSON(w, http.StatusOK, map[string]any{
		"customer_code": code,
		"order_count":   5,
	})
}

// ListCustomerOrders godoc
// @Summary List customer orders
// @Description Get list of all orders for a specific customer
// @Tags customers
// @Accept json
// @Produce json
// @Param code path int true "Customer Code" minimum(1)
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/customers/{code}/orders [get]
func (h *OrderHandler) ListCustomerOrders(w http.ResponseWriter, r *http.Request) {
	codeStr := r.PathValue("code")

	code, err := strconv.Atoi(codeStr)
	if err != nil || code < 1 {
		RespondError(w, http.StatusBadRequest, "Invalid customer code", map[string]string{
			"code": "Customer code must be a positive integer",
		})
		return
	}

	// TODO: Call service to get customer orders
	// For now, return mock data
	RespondJSON(w, http.StatusOK, map[string]any{
		"customer_code": code,
		"orders": []map[string]any{
			{
				"code":        1001,
				"total_value": 120.50,
				"created_at":  "2024-12-08T15:30:00Z",
			},
			{
				"code":        1002,
				"total_value": 50.00,
				"created_at":  "2024-12-07T10:15:00Z",
			},
		},
	})
}

type CreateOrderRequest struct {
	Code         int64                    `json:"codigoPedido" validate:"required,gt=0" example:"1001"`
	CustomerCode int                      `json:"codigoCliente" validate:"required,gt=0" example:"1"`
	Items        []CreateOrderItemRequest `json:"itens" validate:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	Product  string  `json:"produto" validate:"required,min=1" example:"lápis"`
	Quantity int     `json:"quantidade" validate:"required,gt=0" example:"100"`
	Price    float64 `json:"preco" validate:"required,gt=0" example:"1.10"`
}

func (r *CreateOrderRequest) ToDomain() *domain.Order {
	orderItems := make([]domain.OrderItem, len(r.Items))

	for _, item := range r.Items {
		newItem := domain.OrderItem{
			Product:  item.Product,
			Quantity: item.Quantity,
			Price:    item.Price,
		}

		orderItems = append(orderItems, newItem)
	}

	return &domain.Order{
		CustomerCode: r.CustomerCode,
		OrderCode:    r.Code,
		Items:        orderItems,
		CreatedAt:    time.Now().UTC(),
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with items
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Order data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body", nil)
		return

	}

	if err := ValidateStruct(req); err != nil {
		RespondValidationError(w, err)
		return
	}

	err := h.orderService.CreateOrder(r.Context(), req.ToDomain())

	// TODO: Call service to create order
	RespondJSON(w, http.StatusCreated, map[string]any{
		"order_code":    req.Code,
		"customer_code": req.CustomerCode,
		"message":       "Order created successfully",
	})
}
