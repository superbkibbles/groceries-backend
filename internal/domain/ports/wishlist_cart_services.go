package ports

import (
	"context"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
)

// WishlistService defines the interface for wishlist business logic
type WishlistService interface {
	CreateWishlist(ctx context.Context, userID string) (*entities.Wishlist, error)
	GetWishlist(ctx context.Context, id string) (*entities.Wishlist, error)
	GetUserWishlist(ctx context.Context, userID string) (*entities.Wishlist, error)
	AddItem(ctx context.Context, userID string, productID string, variationID string) (*entities.WishlistItem, error)
	RemoveItem(ctx context.Context, userID string, itemID string) error
	RemoveItemByProduct(ctx context.Context, userID string, productID string, variationID string) error
	ClearWishlist(ctx context.Context, userID string) error
	IsProductInWishlist(ctx context.Context, userID string, productID string, variationID string) (bool, error)
}

// CartService defines the interface for cart business logic
type CartService interface {
	CreateCart(ctx context.Context, userID string) (*entities.Cart, error)
	GetCart(ctx context.Context, id string) (*entities.Cart, error)
	GetUserCart(ctx context.Context, userID string) (*entities.Cart, error)
	AddItem(ctx context.Context, userID string, productID string, variationID string, quantity int) (*entities.CartItem, error)
	UpdateItemQuantity(ctx context.Context, userID string, itemID string, quantity int) error
	RemoveItem(ctx context.Context, userID string, itemID string) error
	ClearCart(ctx context.Context, userID string) error
	ConvertToOrder(ctx context.Context, userID string, shippingInfo entities.ShippingInfo) (*entities.Order, error)
}
