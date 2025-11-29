package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
	"github.com/aruncs31s/esdcshopmodule/internal/application/usecase"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
)

// ProductHandler handles HTTP requests for products.
type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

// CreateProduct handles POST /products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.productUseCase.CreateProduct(r.Context(), req)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, product)
}

// GetProduct handles GET /products/{id}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	product, err := h.productUseCase.GetProduct(r.Context(), id)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, product)
}

// GetProducts handles GET /products
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPagination(r)

	products, err := h.productUseCase.GetProducts(r.Context(), limit, offset)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

// GetProductsByCategory handles GET /categories/{id}/products
func (h *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := getPathParam(r, "id")
	if categoryID == "" {
		respondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	limit, offset := getPagination(r)

	products, err := h.productUseCase.GetProductsByCategory(r.Context(), categoryID, limit, offset)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

// UpdateProduct handles PUT /products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.productUseCase.UpdateProduct(r.Context(), id, req)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, product)
}

// DeleteProduct handles DELETE /products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	if err := h.productUseCase.DeleteProduct(r.Context(), id); err != nil {
		handleDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Product deleted successfully"})
}

// UpdateStock handles PATCH /products/{id}/stock
func (h *ProductHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	var req dto.UpdateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.productUseCase.UpdateStock(r.Context(), id, req.Stock)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, product)
}

// SearchProducts handles GET /products/search
func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	limit, offset := getPagination(r)

	products, err := h.productUseCase.SearchProducts(r.Context(), query, limit, offset)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

// GetAvailableProducts handles GET /products/available
func (h *ProductHandler) GetAvailableProducts(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPagination(r)

	products, err := h.productUseCase.GetAvailableProducts(r.Context(), limit, offset)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

// Helper functions

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, dto.ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Code:    code,
	})
}

func handleDomainError(w http.ResponseWriter, err error) {
	switch err {
	case entity.ErrProductNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case entity.ErrCategoryNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case entity.ErrInvalidProductID, entity.ErrInvalidProductName, entity.ErrInvalidStock, entity.ErrInvalidQuantity:
		respondWithError(w, http.StatusBadRequest, err.Error())
	case entity.ErrInsufficientStock:
		respondWithError(w, http.StatusConflict, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func getPathParam(r *http.Request, name string) string {
	// This is a simple implementation. In production, use a proper router.
	return r.PathValue(name)
}

func getPagination(r *http.Request) (limit, offset int) {
	limit = 20
	offset = 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}
