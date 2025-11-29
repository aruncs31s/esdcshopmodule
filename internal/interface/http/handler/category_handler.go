package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
	"github.com/aruncs31s/esdcshopmodule/internal/application/usecase"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
)

// CategoryHandler handles HTTP requests for categories.
type CategoryHandler struct {
	categoryUseCase *usecase.CategoryUseCase
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(categoryUseCase *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}

// CreateCategory handles POST /categories
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category, err := h.categoryUseCase.CreateCategory(r.Context(), req)
	if err != nil {
		handleCategoryError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, category)
}

// GetCategory handles GET /categories/{id}
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	category, err := h.categoryUseCase.GetCategory(r.Context(), id)
	if err != nil {
		handleCategoryError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, category)
}

// GetCategories handles GET /categories
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryUseCase.GetCategories(r.Context())
	if err != nil {
		handleCategoryError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

// GetRootCategories handles GET /categories/root
func (h *CategoryHandler) GetRootCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryUseCase.GetRootCategories(r.Context())
	if err != nil {
		handleCategoryError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

// GetSubcategories handles GET /categories/{id}/subcategories
func (h *CategoryHandler) GetSubcategories(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	categories, err := h.categoryUseCase.GetSubcategories(r.Context(), id)
	if err != nil {
		handleCategoryError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

// UpdateCategory handles PUT /categories/{id}
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	var req dto.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category, err := h.categoryUseCase.UpdateCategory(r.Context(), id, req)
	if err != nil {
		handleCategoryError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, category)
}

// DeleteCategory handles DELETE /categories/{id}
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	if err := h.categoryUseCase.DeleteCategory(r.Context(), id); err != nil {
		handleCategoryError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Category deleted successfully"})
}

func handleCategoryError(w http.ResponseWriter, err error) {
	switch err {
	case entity.ErrCategoryNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case entity.ErrInvalidCategoryID, entity.ErrInvalidCategoryName:
		respondWithError(w, http.StatusBadRequest, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
