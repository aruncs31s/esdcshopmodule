package usecase

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

// ProductUseCase handles product-related business logic.
// Following the Single Responsibility Principle (SRP) from SOLID.
type ProductUseCase struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	idGenerator  IDGenerator
}

// IDGenerator defines the interface for generating unique IDs.
// Following the Dependency Inversion Principle (DIP) from SOLID.
type IDGenerator interface {
	Generate() string
}

// NewProductUseCase creates a new ProductUseCase.
func NewProductUseCase(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	idGenerator IDGenerator,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		idGenerator:  idGenerator,
	}
}

// CreateProduct creates a new product.
func (uc *ProductUseCase) CreateProduct(ctx context.Context, req dto.ProductRequest) (*dto.ProductResponse, error) {
	// Validate category if provided
	if req.CategoryID != "" {
		_, err := uc.categoryRepo.GetByID(ctx, req.CategoryID)
		if err != nil {
			return nil, entity.ErrCategoryNotFound
		}
	}

	// Create money value object
	price, err := valueobject.NewMoney(req.Price, valueobject.Currency(req.Currency))
	if err != nil {
		return nil, err
	}

	// Generate ID
	id := uc.idGenerator.Generate()

	// Create product entity
	product, err := entity.NewProduct(id, req.Name, req.Description, price, req.Stock, req.CategoryID)
	if err != nil {
		return nil, err
	}

	// Add images
	for _, img := range req.Images {
		product.AddImage(img)
	}

	// Set active status
	if !req.Active {
		product.Deactivate()
	}

	// Save product
	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return uc.toProductResponse(product), nil
}

// GetProduct retrieves a product by ID.
func (uc *ProductUseCase) GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toProductResponse(product), nil
}

// GetProducts retrieves all products with pagination.
func (uc *ProductUseCase) GetProducts(ctx context.Context, limit, offset int) (*dto.ProductListResponse, error) {
	products, err := uc.productRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.productRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	productResponses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		productResponses[i] = *uc.toProductResponse(p)
	}

	return &dto.ProductListResponse{
		Products: productResponses,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
		HasMore:  int64(offset+len(products)) < total,
	}, nil
}

// GetProductsByCategory retrieves products by category.
func (uc *ProductUseCase) GetProductsByCategory(ctx context.Context, categoryID string, limit, offset int) (*dto.ProductListResponse, error) {
	products, err := uc.productRepo.GetByCategory(ctx, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}

	productResponses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		productResponses[i] = *uc.toProductResponse(p)
	}

	return &dto.ProductListResponse{
		Products: productResponses,
		Total:    int64(len(products)),
		Limit:    limit,
		Offset:   offset,
		HasMore:  false,
	}, nil
}

// UpdateProduct updates an existing product.
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id string, req dto.ProductRequest) (*dto.ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate category if provided
	if req.CategoryID != "" {
		_, err := uc.categoryRepo.GetByID(ctx, req.CategoryID)
		if err != nil {
			return nil, entity.ErrCategoryNotFound
		}
	}

	// Create money value object
	price, err := valueobject.NewMoney(req.Price, valueobject.Currency(req.Currency))
	if err != nil {
		return nil, err
	}

	// Update product fields
	product.Name = req.Name
	product.Description = req.Description
	product.Price = price
	product.Stock = req.Stock
	product.CategoryID = req.CategoryID
	product.Images = req.Images

	if req.Active {
		product.Activate()
	} else {
		product.Deactivate()
	}

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return uc.toProductResponse(product), nil
}

// DeleteProduct deletes a product.
func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	_, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.productRepo.Delete(ctx, id)
}

// UpdateStock updates product stock.
func (uc *ProductUseCase) UpdateStock(ctx context.Context, id string, stock int) (*dto.ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := product.UpdateStock(stock); err != nil {
		return nil, err
	}

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return uc.toProductResponse(product), nil
}

// SearchProducts searches products by query.
func (uc *ProductUseCase) SearchProducts(ctx context.Context, query string, limit, offset int) (*dto.ProductListResponse, error) {
	products, err := uc.productRepo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	productResponses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		productResponses[i] = *uc.toProductResponse(p)
	}

	return &dto.ProductListResponse{
		Products: productResponses,
		Total:    int64(len(products)),
		Limit:    limit,
		Offset:   offset,
		HasMore:  false,
	}, nil
}

// GetAvailableProducts retrieves all available products.
func (uc *ProductUseCase) GetAvailableProducts(ctx context.Context, limit, offset int) (*dto.ProductListResponse, error) {
	products, err := uc.productRepo.GetActiveProducts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Filter available products
	available := make([]dto.ProductResponse, 0)
	for _, p := range products {
		if p.IsAvailable() {
			available = append(available, *uc.toProductResponse(p))
		}
	}

	return &dto.ProductListResponse{
		Products: available,
		Total:    int64(len(available)),
		Limit:    limit,
		Offset:   offset,
		HasMore:  false,
	}, nil
}

// toProductResponse converts a product entity to a response DTO.
func (uc *ProductUseCase) toProductResponse(product *entity.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price.Amount,
		Currency:    string(product.Price.Currency),
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		Images:      product.Images,
		Active:      product.Active,
		Available:   product.IsAvailable(),
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
