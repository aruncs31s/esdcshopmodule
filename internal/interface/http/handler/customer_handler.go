package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
	"github.com/aruncs31s/esdcshopmodule/internal/application/usecase"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
)

// CustomerHandler handles HTTP requests for customers.
type CustomerHandler struct {
	customerUseCase *usecase.CustomerUseCase
}

// NewCustomerHandler creates a new CustomerHandler.
func NewCustomerHandler(customerUseCase *usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{
		customerUseCase: customerUseCase,
	}
}

// CreateCustomer handles POST /customers
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req dto.CustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	customer, err := h.customerUseCase.CreateCustomer(r.Context(), req)
	if err != nil {
		handleCustomerError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, customer)
}

// GetCustomer handles GET /customers/{id}
func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	customer, err := h.customerUseCase.GetCustomer(r.Context(), id)
	if err != nil {
		handleCustomerError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, customer)
}

// GetCustomers handles GET /customers
func (h *CustomerHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPagination(r)

	customers, err := h.customerUseCase.GetCustomers(r.Context(), limit, offset)
	if err != nil {
		handleCustomerError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, customers)
}

// UpdateCustomer handles PUT /customers/{id}
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	var req dto.CustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	customer, err := h.customerUseCase.UpdateCustomer(r.Context(), id, req)
	if err != nil {
		handleCustomerError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, customer)
}

// DeleteCustomer handles DELETE /customers/{id}
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	if err := h.customerUseCase.DeleteCustomer(r.Context(), id); err != nil {
		handleCustomerError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Customer deleted successfully"})
}

// UpdateShippingAddress handles PUT /customers/{id}/shipping-address
func (h *CustomerHandler) UpdateShippingAddress(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	var req dto.AddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	customer, err := h.customerUseCase.UpdateShippingAddress(r.Context(), id, req)
	if err != nil {
		handleCustomerError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, customer)
}

// UpdateBillingAddress handles PUT /customers/{id}/billing-address
func (h *CustomerHandler) UpdateBillingAddress(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	var req dto.AddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	customer, err := h.customerUseCase.UpdateBillingAddress(r.Context(), id, req)
	if err != nil {
		handleCustomerError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, customer)
}

// SearchCustomers handles GET /customers/search
func (h *CustomerHandler) SearchCustomers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	limit, offset := getPagination(r)

	customers, err := h.customerUseCase.SearchCustomers(r.Context(), query, limit, offset)
	if err != nil {
		handleCustomerError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, customers)
}

func handleCustomerError(w http.ResponseWriter, err error) {
	switch err {
	case entity.ErrCustomerNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case entity.ErrInvalidCustomerID, entity.ErrInvalidCustomerEmail:
		respondWithError(w, http.StatusBadRequest, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
