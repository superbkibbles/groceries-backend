package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// WishlistItem represents a product in a user's wishlist
type WishlistItem struct {
	ID          string    `json:"id" bson:"id"`
	ProductID   string    `json:"product_id" bson:"product_id"`
	VariationID string    `json:"variation_id,omitempty" bson:"variation_id,omitempty"`
	AddedAt     time.Time `json:"added_at" bson:"added_at"`
}

// Wishlist represents a user's wishlist in the e-commerce system
type Wishlist struct {
	ID        string          `json:"id" bson:"_id"`
	UserID    string          `json:"user_id" bson:"user_id"`
	Items     []*WishlistItem `json:"items" bson:"items"`
	CreatedAt time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" bson:"updated_at"`
}

// NewWishlist creates a new wishlist with a unique ID
func NewWishlist(userID string) *Wishlist {
	now := time.Now()
	return &Wishlist{
		ID:        uuid.New().String(),
		UserID:    userID,
		Items:     []*WishlistItem{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddItem adds a product to the wishlist
func (w *Wishlist) AddItem(productID string, variationID string) (*WishlistItem, error) {
	// Check if the item already exists in the wishlist
	for _, item := range w.Items {
		if item.ProductID == productID && item.VariationID == variationID {
			return nil, errors.New("item already exists in wishlist")
		}
	}

	// Add new item
	item := &WishlistItem{
		ID:          uuid.New().String(),
		ProductID:   productID,
		VariationID: variationID,
		AddedAt:     time.Now(),
	}

	w.Items = append(w.Items, item)
	w.UpdatedAt = time.Now()

	return item, nil
}

// RemoveItem removes a product from the wishlist
func (w *Wishlist) RemoveItem(itemID string) error {
	for i, item := range w.Items {
		if item.ID == itemID {
			// Remove item by replacing it with the last item and truncating the slice
			lastIndex := len(w.Items) - 1
			w.Items[i] = w.Items[lastIndex]
			w.Items = w.Items[:lastIndex]
			w.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("item not found in wishlist")
}

// RemoveItemByProduct removes a product from the wishlist by product ID and variation ID
func (w *Wishlist) RemoveItemByProduct(productID string, variationID string) error {
	for i, item := range w.Items {
		if item.ProductID == productID && item.VariationID == variationID {
			// Remove item by replacing it with the last item and truncating the slice
			lastIndex := len(w.Items) - 1
			w.Items[i] = w.Items[lastIndex]
			w.Items = w.Items[:lastIndex]
			w.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("item not found in wishlist")
}

// ClearWishlist removes all items from the wishlist
func (w *Wishlist) ClearWishlist() {
	w.Items = []*WishlistItem{}
	w.UpdatedAt = time.Now()
}

// ClearItems removes all items from the wishlist (alias for ClearWishlist)
func (w *Wishlist) ClearItems() {
	w.ClearWishlist()
}

// ContainsProduct checks if a product is in the wishlist
func (w *Wishlist) ContainsProduct(productID string, variationID string) bool {
	for _, item := range w.Items {
		if item.ProductID == productID && item.VariationID == variationID {
			return true
		}
	}

	return false
}
