package services

import (
	"context"
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Check if user already has a cart
	existingCart, err := s.cartRepo.GetByUserID(ctx, userObjectID)
	if err == nil {
		return existingCart, nil
	}

	// Create new cart
	cart := entities.NewCart(userObjectID)
	err = s.cartRepo.Create(ctx, cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

// GetCart retrieves a cart by its ID
func (s *CartService) GetCart(ctx context.Context, id string) (*entities.Cart, error) {
	cartID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid cart ID")
	}
	return s.cartRepo.GetByID(ctx, cartID)
}

// GetUserCart retrieves a user's cart
func (s *CartService) GetUserCart(ctx context.Context, userID string) (*entities.Cart, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	return s.cartRepo.GetByUserID(ctx, userObjectID)
}

// AddItem adds a product to a user's cart
func (s *CartService) AddItem(ctx context.Context, userID string, productID string, quantity int) (*entities.CartItem, error) {
	// Convert IDs to ObjectID
	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Verify product exists
	product, err := s.productRepo.GetByID(ctx, productObjectID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check if product is in stock
	if product.StockQuantity < quantity {
		return nil, errors.New("insufficient stock")
	}

	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	// Add item to cart
	item, err := cart.AddItem(productObjectID, product.SKU, product.Name, product.Price, quantity)
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

	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	itemObjectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return errors.New("invalid item ID")
	}

	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userObjectID)
	if err != nil {
		return err
	}

	// Find the item
	var item *entities.CartItem
	for _, i := range cart.Items {
		if i.ID == itemObjectID {
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

	// Check if product is in stock
	if product.StockQuantity < quantity {
		return errors.New("insufficient stock")
	}

	// Update item quantity
	err = cart.UpdateItemQuantity(itemObjectID, quantity)
	if err != nil {
		return err
	}

	// Save to repository
	return s.cartRepo.UpdateItemQuantity(ctx, cart.ID, itemObjectID, quantity)
}

// RemoveItem removes an item from a user's cart
func (s *CartService) RemoveItem(ctx context.Context, userID string, itemID string) error {
	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	itemObjectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return errors.New("invalid item ID")
	}

	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userObjectID)
	if err != nil {
		return err
	}

	// Remove item from cart
	err = cart.RemoveItem(itemObjectID)
	if err != nil {
		return err
	}

	// Save to repository
	return s.cartRepo.RemoveItem(ctx, cart.ID, itemObjectID)
}

// ClearCart removes all items from a user's cart
func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userObjectID)
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
	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	// Check if cart is empty
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Create order from cart
	order := entities.NewOrder(userObjectID, shippingInfo)

	// Add items from cart to order
	for _, cartItem := range cart.Items {
		order.AddItem(
			cartItem.ProductID,
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
