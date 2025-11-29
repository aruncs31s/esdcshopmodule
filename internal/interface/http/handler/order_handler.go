package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
	"github.com/aruncs31s/esdcshopmodule/internal/application/usecase"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
)

// OrderHandler handles HTTP requests for orders.
type OrderHandler struct {
	orderUseCase *usecase.OrderUseCase
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(orderUseCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	order, err := h.orderUseCase.CreateOrder(r.Context(), req)
	if err != nil {
		handleOrderError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, order)
}

// CreateOrderFromCart handles POST /customers/{id}/orders
func (h *OrderHandler) CreateOrderFromCart(w http.ResponseWriter, r *http.Request) {
	customerID := getPathParam(r, "id")
	if customerID == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	var req dto.CreateOrderFromCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	order, err := h.orderUseCase.CreateOrderFromCart(r.Context(), customerID, req)
	if err != nil {
		handleOrderError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, order)
}

// GetOrder handles GET /orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	order, err := h.orderUseCase.GetOrder(r.Context(), id)
	if err != nil {
		handleOrderError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

// GetOrders handles GET /orders
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPagination(r)

	orders, err := h.orderUseCase.GetOrders(r.Context(), limit, offset)
	if err != nil {
		handleOrderError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, orders)
}

// GetCustomerOrders handles GET /customers/{id}/orders
func (h *OrderHandler) GetCustomerOrders(w http.ResponseWriter, r *http.Request) {
	customerID := getPathParam(r, "id")
	if customerID == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	limit, offset := getPagination(r)

	orders, err := h.orderUseCase.GetCustomerOrders(r.Context(), customerID, limit, offset)
	if err != nil {
		handleOrderError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, orders)
}

// UpdateOrderStatus handles PATCH /orders/{id}/status
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	order, err := h.orderUseCase.UpdateOrderStatus(r.Context(), id, req.Status)
	if err != nil {
		handleOrderError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

// CancelOrder handles POST /orders/{id}/cancel
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	order, err := h.orderUseCase.CancelOrder(r.Context(), id)
	if err != nil {
		handleOrderError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

func handleOrderError(w http.ResponseWriter, err error) {
	switch err {
	case entity.ErrOrderNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case entity.ErrCustomerNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case entity.ErrProductNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case entity.ErrCartNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case entity.ErrInvalidOrderID, entity.ErrInvalidCustomerID, entity.ErrEmptyOrder:
		respondWithError(w, http.StatusBadRequest, err.Error())
	case entity.ErrInvalidOrderStatus:
		respondWithError(w, http.StatusConflict, err.Error())
	case entity.ErrInsufficientStock:
		respondWithError(w, http.StatusConflict, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
