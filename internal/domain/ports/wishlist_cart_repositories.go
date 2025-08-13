package ports

import (
	"context"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WishlistRepository defines the interface for wishlist data access
type WishlistRepository interface {
	Create(ctx context.Context, wishlist *entities.Wishlist) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Wishlist, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.Wishlist, error)
	Update(ctx context.Context, wishlist *entities.Wishlist) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	AddItem(ctx context.Context, wishlistID primitive.ObjectID, item *entities.WishlistItem) error
	RemoveItem(ctx context.Context, wishlistID primitive.ObjectID, itemID primitive.ObjectID) error
	ClearWishlist(ctx context.Context, wishlistID primitive.ObjectID) error
}

// CartRepository defines the interface for cart data access
type CartRepository interface {
	Create(ctx context.Context, cart *entities.Cart) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Cart, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.Cart, error)
	Update(ctx context.Context, cart *entities.Cart) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	AddItem(ctx context.Context, cartID primitive.ObjectID, item *entities.CartItem) error
	UpdateItemQuantity(ctx context.Context, cartID primitive.ObjectID, itemID primitive.ObjectID, quantity int) error
	RemoveItem(ctx context.Context, cartID primitive.ObjectID, itemID primitive.ObjectID) error
	ClearCart(ctx context.Context, cartID primitive.ObjectID) error
}
