package services

import (
	"context"
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// WishlistService implements the WishlistService interface
type WishlistService struct {
	wishlistRepo ports.WishlistRepository
	productRepo  ports.ProductRepository
}

// NewWishlistService creates a new WishlistService
func NewWishlistService(wishlistRepo ports.WishlistRepository, productRepo ports.ProductRepository) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
		productRepo:  productRepo,
	}
}

// CreateWishlist creates a new wishlist for a user
func (s *WishlistService) CreateWishlist(ctx context.Context, userID string) (*entities.Wishlist, error) {
	// Check if user already has a wishlist
	existingWishlist, err := s.wishlistRepo.GetByUserID(ctx, userID)
	if err == nil {
		return existingWishlist, nil
	}

	// Create new wishlist
	wishlist := entities.NewWishlist(userID)
	err = s.wishlistRepo.Create(ctx, wishlist)
	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

// GetWishlist retrieves a wishlist by its ID
func (s *WishlistService) GetWishlist(ctx context.Context, id string) (*entities.Wishlist, error) {
	return s.wishlistRepo.GetByID(ctx, id)
}

// GetUserWishlist retrieves a user's wishlist
func (s *WishlistService) GetUserWishlist(ctx context.Context, userID string) (*entities.Wishlist, error) {
	return s.wishlistRepo.GetByUserID(ctx, userID)
}

// AddItem adds a product to a user's wishlist
func (s *WishlistService) AddItem(ctx context.Context, userID string, productID string, variationID string) (*entities.WishlistItem, error) {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Get user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Add item to wishlist
	item, err := wishlist.AddItem(productID, variationID)
	if err != nil {
		return nil, err
	}

	// Save to repository
	err = s.wishlistRepo.AddItem(ctx, wishlist.ID, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// RemoveItem removes an item from a user's wishlist
func (s *WishlistService) RemoveItem(ctx context.Context, userID string, itemID string) error {
	// Get user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Remove item from wishlist
	err = wishlist.RemoveItem(itemID)
	if err != nil {
		return err
	}

	// Save to repository
	return s.wishlistRepo.RemoveItem(ctx, wishlist.ID, itemID)
}

// RemoveItemByProduct removes an item from a user's wishlist by product ID
func (s *WishlistService) RemoveItemByProduct(ctx context.Context, userID string, productID string, variationID string) error {
	// Get user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Find the item with the given product ID
	var itemID string
	for _, item := range wishlist.Items {
		if item.ProductID == productID && item.VariationID == variationID {
			itemID = item.ID
			break
		}
	}

	if itemID == "" {
		return errors.New("item not found in wishlist")
	}

	// Remove the item
	return s.RemoveItem(ctx, userID, itemID)
}

// ClearWishlist removes all items from a user's wishlist
func (s *WishlistService) ClearWishlist(ctx context.Context, userID string) error {
	// Get user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Clear wishlist
	wishlist.ClearItems()

	// Save to repository
	return s.wishlistRepo.ClearWishlist(ctx, wishlist.ID)
}

// IsProductInWishlist checks if a product is in a user's wishlist
func (s *WishlistService) IsProductInWishlist(ctx context.Context, userID string, productID string, variationID string) (bool, error) {
	// Get user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Check if product is in wishlist
	for _, item := range wishlist.Items {
		if item.ProductID == productID && item.VariationID == variationID {
			return true, nil
		}
	}

	return false, nil
}
