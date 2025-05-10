package ports

import (
	"context"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
)

// WishlistRepository defines the interface for wishlist data access
type WishlistRepository interface {
	Create(ctx context.Context, wishlist *entities.Wishlist) error
	GetByID(ctx context.Context, id string) (*entities.Wishlist, error)
	GetByUserID(ctx context.Context, userID string) (*entities.Wishlist, error)
	Update(ctx context.Context, wishlist *entities.Wishlist) error
	Delete(ctx context.Context, id string) error
	AddItem(ctx context.Context, wishlistID string, item *entities.WishlistItem) error
	RemoveItem(ctx context.Context, wishlistID string, itemID string) error
	ClearWishlist(ctx context.Context, wishlistID string) error
}

// CartRepository defines the interface for cart data access
type CartRepository interface {
	Create(ctx context.Context, cart *entities.Cart) error
	GetByID(ctx context.Context, id string) (*entities.Cart, error)
	GetByUserID(ctx context.Context, userID string) (*entities.Cart, error)
	Update(ctx context.Context, cart *entities.Cart) error
	Delete(ctx context.Context, id string) error
	AddItem(ctx context.Context, cartID string, item *entities.CartItem) error
	UpdateItemQuantity(ctx context.Context, cartID string, itemID string, quantity int) error
	RemoveItem(ctx context.Context, cartID string, itemID string) error
	ClearCart(ctx context.Context, cartID string) error
}
