package services

import (
	"context"
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// CartService implements the CartService interface
type CartService struct {
	cartRepo    ports.CartRepository
	productRepo ports.ProductRepository
	orderRepo   ports.OrderRepository
}

// NewCartService creates a new CartService
func NewCartService(
	cartRepo ports.CartRepository,
	productRepo ports.ProductRepository,
	orderRepo ports.OrderRepository,
) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

// CreateCart creates a new cart for a user
func (s *CartService) CreateCart(ctx context.Context, userID string) (*entities.Cart, error) {
	// Check if user already has a cart
	existingCart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err == nil {
		return existingCart, nil
	}

	// Create new cart
	cart := entities.NewCart(userID)
	err = s.cartRepo.Create(ctx, cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

// GetCart retrieves a cart by its ID
func (s *CartService) GetCart(ctx context.Context, id string) (*entities.Cart, error) {
	return s.cartRepo.GetByID(ctx, id)
}

// GetUserCart retrieves a user's cart
func (s *CartService) GetUserCart(ctx context.Context, userID string) (*entities.Cart, error) {
	return s.cartRepo.GetByUserID(ctx, userID)
}

// AddItem adds a product to a user's cart
func (s *CartService) AddItem(ctx context.Context, userID string, productID string, variationID string, quantity int) (*entities.CartItem, error) {
	// Verify product exists
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Find the variation
	var variation *entities.Variation
	for _, v := range product.Variations {
		if v.ID == variationID {
			variation = v
			break
		}
	}

	if variation == nil {
		return nil, errors.New("product variation not found")
	}

	// Check if variation is in stock
	if variation.StockQuantity < quantity {
		return nil, errors.New("insufficient stock")
	}

	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Add item to cart
	item, err := cart.AddItem(productID, variationID, variation.SKU, product.Name, variation.Price, quantity)
	if err != nil {
		return nil, err
	}

	// Save to repository
	err = s.cartRepo.AddItem(ctx, cart.ID, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateItemQuantity updates the quantity of an item in a user's cart
func (s *CartService) UpdateItemQuantity(ctx context.Context, userID string, itemID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Find the item
	var item *entities.CartItem
	for _, i := range cart.Items {
		if i.ID == itemID {
			item = i
			break
		}
	}

	if item == nil {
		return errors.New("item not found in cart")
	}

	// Verify product exists and has sufficient stock
	product, err := s.productRepo.GetByID(ctx, item.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	// Find the variation
	var variation *entities.Variation
	for _, v := range product.Variations {
		if v.ID == item.VariationID {
			variation = v
			break
		}
	}

	if variation == nil {
		return errors.New("product variation not found")
	}

	// Check if variation is in stock
	if variation.StockQuantity < quantity {
		return errors.New("insufficient stock")
	}

	// Update item quantity
	err = cart.UpdateItemQuantity(itemID, quantity)
	if err != nil {
		return err
	}

	// Save to repository
	return s.cartRepo.UpdateItemQuantity(ctx, cart.ID, itemID, quantity)
}

// RemoveItem removes an item from a user's cart
func (s *CartService) RemoveItem(ctx context.Context, userID string, itemID string) error {
	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Remove item from cart
	err = cart.RemoveItem(itemID)
	if err != nil {
		return err
	}

	// Save to repository
	return s.cartRepo.RemoveItem(ctx, cart.ID, itemID)
}

// ClearCart removes all items from a user's cart
func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Clear cart
	cart.ClearItems()

	// Save to repository
	return s.cartRepo.ClearCart(ctx, cart.ID)
}

// ConvertToOrder converts a cart to an order
func (s *CartService) ConvertToOrder(ctx context.Context, userID string, shippingInfo entities.ShippingInfo) (*entities.Order, error) {
	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if cart is empty
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Create order from cart
	order := entities.NewOrder(userID, shippingInfo)

	// Add items from cart to order
	for _, cartItem := range cart.Items {
		order.AddItem(
			cartItem.ProductID,
			cartItem.VariationID,
			cartItem.SKU,
			cartItem.Name,
			cartItem.Price,
			cartItem.Quantity,
		)
	}

	// Save order
	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	// Clear cart after successful order creation
	err = s.ClearCart(ctx, userID)
	if err != nil {
		// Log error but don't fail the order creation
		// TODO: Add proper logging
	}

	return order, nil
}
