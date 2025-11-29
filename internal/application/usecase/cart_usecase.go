package usecase

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

// CartUseCase handles cart-related business logic.
type CartUseCase struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	idGenerator IDGenerator
}

// NewCartUseCase creates a new CartUseCase.
func NewCartUseCase(
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	idGenerator IDGenerator,
) *CartUseCase {
	return &CartUseCase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		idGenerator: idGenerator,
	}
}

// GetCart retrieves a cart by customer ID.
func (uc *CartUseCase) GetCart(ctx context.Context, customerID string) (*dto.CartResponse, error) {
	cart, err := uc.cartRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		// If cart doesn't exist, create a new one
		cartID := uc.idGenerator.Generate()
		cart, err = entity.NewCart(cartID, customerID, valueobject.USD)
		if err != nil {
			return nil, err
		}
		if err := uc.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
	}
	return uc.toCartResponse(ctx, cart), nil
}

// AddItem adds an item to the cart.
func (uc *CartUseCase) AddItem(ctx context.Context, customerID string, req dto.CartItemRequest) (*dto.CartResponse, error) {
	// Get or create cart
	cart, err := uc.cartRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		cartID := uc.idGenerator.Generate()
		cart, err = entity.NewCart(cartID, customerID, valueobject.USD)
		if err != nil {
			return nil, err
		}
		if err := uc.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
	}

	// Get product
	product, err := uc.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, entity.ErrProductNotFound
	}

	// Check if product is available
	if !product.IsAvailable() {
		return nil, entity.ErrProductNotFound
	}

	// Check stock
	if product.Stock < req.Quantity {
		return nil, entity.ErrInsufficientStock
	}

	// Add item to cart
	if err := cart.AddItem(req.ProductID, req.Quantity, product.Price); err != nil {
		return nil, err
	}

	// Save cart
	if err := uc.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	return uc.toCartResponse(ctx, cart), nil
}

// UpdateItemQuantity updates the quantity of an item in the cart.
func (uc *CartUseCase) UpdateItemQuantity(ctx context.Context, customerID string, productID string, quantity int) (*dto.CartResponse, error) {
	cart, err := uc.cartRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, entity.ErrCartNotFound
	}

	// Check stock if increasing quantity
	if quantity > 0 {
		product, err := uc.productRepo.GetByID(ctx, productID)
		if err != nil {
			return nil, entity.ErrProductNotFound
		}
		if product.Stock < quantity {
			return nil, entity.ErrInsufficientStock
		}
	}

	// Update quantity
	if err := cart.UpdateItemQuantity(productID, quantity); err != nil {
		return nil, err
	}

	// Save cart
	if err := uc.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	return uc.toCartResponse(ctx, cart), nil
}

// RemoveItem removes an item from the cart.
func (uc *CartUseCase) RemoveItem(ctx context.Context, customerID string, productID string) (*dto.CartResponse, error) {
	cart, err := uc.cartRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, entity.ErrCartNotFound
	}

	if err := cart.RemoveItem(productID); err != nil {
		return nil, err
	}

	if err := uc.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	return uc.toCartResponse(ctx, cart), nil
}

// ClearCart clears all items from the cart.
func (uc *CartUseCase) ClearCart(ctx context.Context, customerID string) (*dto.CartResponse, error) {
	cart, err := uc.cartRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, entity.ErrCartNotFound
	}

	cart.Clear()

	if err := uc.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	return uc.toCartResponse(ctx, cart), nil
}

// toCartResponse converts a cart entity to a response DTO.
func (uc *CartUseCase) toCartResponse(ctx context.Context, cart *entity.Cart) *dto.CartResponse {
	items := make([]dto.CartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		productName := ""
		product, err := uc.productRepo.GetByID(ctx, item.ProductID)
		if err == nil {
			productName = product.Name
		}

		items[i] = dto.CartItemResponse{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice.Amount,
			Currency:    string(item.UnitPrice.Currency),
			TotalPrice:  item.TotalPrice.Amount,
			AddedAt:     item.AddedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &dto.CartResponse{
		ID:         cart.ID,
		CustomerID: cart.CustomerID,
		Items:      items,
		ItemCount:  cart.GetItemCount(),
		Total:      cart.Total.Amount,
		Currency:   string(cart.Currency),
		CreatedAt:  cart.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  cart.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
