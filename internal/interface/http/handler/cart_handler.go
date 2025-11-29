package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
	"github.com/aruncs31s/esdcshopmodule/internal/application/usecase"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
)

// CartHandler handles HTTP requests for carts.
type CartHandler struct {
	cartUseCase *usecase.CartUseCase
}

// NewCartHandler creates a new CartHandler.
func NewCartHandler(cartUseCase *usecase.CartUseCase) *CartHandler {
	return &CartHandler{
		cartUseCase: cartUseCase,
	}
}

// GetCart handles GET /customers/{id}/cart
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	customerID := getPathParam(r, "id")
	if customerID == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	cart, err := h.cartUseCase.GetCart(r.Context(), customerID)
	if err != nil {
		handleCartError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, cart)
}

// AddItem handles POST /customers/{id}/cart/items
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	customerID := getPathParam(r, "id")
	if customerID == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	var req dto.CartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cart, err := h.cartUseCase.AddItem(r.Context(), customerID, req)
	if err != nil {
		handleCartError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, cart)
}

// UpdateItemQuantity handles PUT /customers/{id}/cart/items/{productId}
func (h *CartHandler) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
	customerID := getPathParam(r, "id")
	productID := getPathParam(r, "productId")

	if customerID == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}
	if productID == "" {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	var req dto.CartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cart, err := h.cartUseCase.UpdateItemQuantity(r.Context(), customerID, productID, req.Quantity)
	if err != nil {
		handleCartError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, cart)
}

// RemoveItem handles DELETE /customers/{id}/cart/items/{productId}
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	customerID := getPathParam(r, "id")
	productID := getPathParam(r, "productId")

	if customerID == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}
	if productID == "" {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	cart, err := h.cartUseCase.RemoveItem(r.Context(), customerID, productID)
	if err != nil {
		handleCartError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, cart)
}

// ClearCart handles DELETE /customers/{id}/cart
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	customerID := getPathParam(r, "id")
	if customerID == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	cart, err := h.cartUseCase.ClearCart(r.Context(), customerID)
	if err != nil {
		handleCartError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, cart)
}

func handleCartError(w http.ResponseWriter, err error) {
	switch err {
	case entity.ErrCartNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case entity.ErrProductNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case entity.ErrInvalidCartID, entity.ErrInvalidCustomerID, entity.ErrInvalidQuantity:
		respondWithError(w, http.StatusBadRequest, err.Error())
	case entity.ErrInsufficientStock:
		respondWithError(w, http.StatusConflict, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
