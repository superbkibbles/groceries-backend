package services

import (
	"context"
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Check if user already has a wishlist
	existingWishlist, err := s.wishlistRepo.GetByUserID(ctx, userObjectID)
	if err == nil {
		return existingWishlist, nil
	}

	// Create new wishlist
	wishlist := entities.NewWishlist(userObjectID)
	err = s.wishlistRepo.Create(ctx, wishlist)
	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

// GetWishlist retrieves a wishlist by its ID
func (s *WishlistService) GetWishlist(ctx context.Context, id string) (*entities.Wishlist, error) {
	wishlistID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid wishlist ID")
	}
	return s.wishlistRepo.GetByID(ctx, wishlistID)
}

// GetUserWishlist retrieves a user's wishlist
func (s *WishlistService) GetUserWishlist(ctx context.Context, userID string) (*entities.Wishlist, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	return s.wishlistRepo.GetByUserID(ctx, userObjectID)
}

// AddItem adds a product to a user's wishlist
func (s *WishlistService) AddItem(ctx context.Context, userID string, productID string) (*entities.WishlistItem, error) {
	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	// Verify product exists
	_, err = s.productRepo.GetByID(ctx, productObjectID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Get user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	// Add item to wishlist
	item, err := wishlist.AddItem(productObjectID)
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
	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	itemObjectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return errors.New("invalid item ID")
	}

	// Get user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userObjectID)
	if err != nil {
		return err
	}

	// Remove item from wishlist
	err = wishlist.RemoveItem(itemObjectID)
	if err != nil {
		return err
	}

	// Save to repository
	return s.wishlistRepo.RemoveItem(ctx, wishlist.ID, itemObjectID)
}

// RemoveItemByProduct removes an item from a user's wishlist by product ID
func (s *WishlistService) RemoveItemByProduct(ctx context.Context, userID string, productID string) error {
	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New("invalid product ID")
	}

	// Get user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userObjectID)
	if err != nil {
		return err
	}

	// Find the item with the given product ID
	var itemID primitive.ObjectID
	for _, item := range wishlist.Items {
		if item.ProductID == productObjectID {
			itemID = item.ID
			break
		}
	}

	if itemID.IsZero() {
		return errors.New("item not found in wishlist")
	}

	// Remove the item
	return s.RemoveItem(ctx, userID, itemID.Hex())
}

// ClearWishlist removes all items from a user's wishlist
func (s *WishlistService) ClearWishlist(ctx context.Context, userID string) error {
	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Get user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userObjectID)
	if err != nil {
		return err
	}

	// Clear wishlist
	wishlist.ClearItems()

	// Save to repository
	return s.wishlistRepo.ClearWishlist(ctx, wishlist.ID)
}

// IsProductInWishlist checks if a product is in a user's wishlist
func (s *WishlistService) IsProductInWishlist(ctx context.Context, userID string, productID string) (bool, error) {
	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New("invalid user ID")
	}

	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return false, errors.New("invalid product ID")
	}
	// Get user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userObjectID)
	if err != nil {
		return false, err
	}

	// Check if product is in wishlist
	return wishlist.ContainsProduct(productObjectID), nil
}
