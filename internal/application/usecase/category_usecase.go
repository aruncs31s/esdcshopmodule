package usecase

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
)

// CategoryUseCase handles category-related business logic.
type CategoryUseCase struct {
	categoryRepo repository.CategoryRepository
	idGenerator  IDGenerator
}

// NewCategoryUseCase creates a new CategoryUseCase.
func NewCategoryUseCase(
	categoryRepo repository.CategoryRepository,
	idGenerator IDGenerator,
) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
		idGenerator:  idGenerator,
	}
}

// CreateCategory creates a new category.
func (uc *CategoryUseCase) CreateCategory(ctx context.Context, req dto.CategoryRequest) (*dto.CategoryResponse, error) {
	// Validate parent category if provided
	if req.ParentID != "" {
		_, err := uc.categoryRepo.GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, entity.ErrCategoryNotFound
		}
	}

	// Generate ID
	id := uc.idGenerator.Generate()

	// Create category entity
	category, err := entity.NewCategory(id, req.Name, req.Description, req.ParentID)
	if err != nil {
		return nil, err
	}

	// Save category
	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return uc.toCategoryResponse(category), nil
}

// GetCategory retrieves a category by ID.
func (uc *CategoryUseCase) GetCategory(ctx context.Context, id string) (*dto.CategoryResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toCategoryResponse(category), nil
}

// GetCategories retrieves all categories.
func (uc *CategoryUseCase) GetCategories(ctx context.Context) (*dto.CategoryListResponse, error) {
	categories, err := uc.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	categoryResponses := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		categoryResponses[i] = *uc.toCategoryResponse(c)
	}

	return &dto.CategoryListResponse{
		Categories: categoryResponses,
		Total:      len(categories),
	}, nil
}

// GetRootCategories retrieves all root categories.
func (uc *CategoryUseCase) GetRootCategories(ctx context.Context) (*dto.CategoryListResponse, error) {
	categories, err := uc.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		return nil, err
	}

	categoryResponses := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		categoryResponses[i] = *uc.toCategoryResponse(c)
	}

	return &dto.CategoryListResponse{
		Categories: categoryResponses,
		Total:      len(categories),
	}, nil
}

// GetSubcategories retrieves subcategories of a category.
func (uc *CategoryUseCase) GetSubcategories(ctx context.Context, parentID string) (*dto.CategoryListResponse, error) {
	categories, err := uc.categoryRepo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	categoryResponses := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		categoryResponses[i] = *uc.toCategoryResponse(c)
	}

	return &dto.CategoryListResponse{
		Categories: categoryResponses,
		Total:      len(categories),
	}, nil
}

// UpdateCategory updates an existing category.
func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, id string, req dto.CategoryRequest) (*dto.CategoryResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate parent category if provided
	if req.ParentID != "" && req.ParentID != category.ParentID {
		_, err := uc.categoryRepo.GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, entity.ErrCategoryNotFound
		}
	}

	// Update category
	if err := category.Update(req.Name, req.Description); err != nil {
		return nil, err
	}
	category.ParentID = req.ParentID

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return uc.toCategoryResponse(category), nil
}

// DeleteCategory deletes a category.
func (uc *CategoryUseCase) DeleteCategory(ctx context.Context, id string) error {
	_, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.categoryRepo.Delete(ctx, id)
}

// toCategoryResponse converts a category entity to a response DTO.
func (uc *CategoryUseCase) toCategoryResponse(category *entity.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ParentID:    category.ParentID,
		Active:      category.Active,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
