package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// WishlistRepository is a MongoDB implementation of the WishlistRepository interface
type WishlistRepository struct {
	collection *mongo.Collection
}

// NewWishlistRepository creates a new WishlistRepository
func NewWishlistRepository(db *mongo.Database) *WishlistRepository {
	return &WishlistRepository{
		collection: db.Collection("wishlists"),
	}
}

// Create adds a new wishlist to the database
func (r *WishlistRepository) Create(ctx context.Context, wishlist *entities.Wishlist) error {
	_, err := r.collection.InsertOne(ctx, wishlist)
	return err
}

// GetByID retrieves a wishlist by its ID
func (r *WishlistRepository) GetByID(ctx context.Context, id string) (*entities.Wishlist, error) {
	var wishlist entities.Wishlist
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&wishlist)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("wishlist not found")
		}
		return nil, err
	}
	return &wishlist, nil
}

// GetByUserID retrieves a wishlist by user ID
func (r *WishlistRepository) GetByUserID(ctx context.Context, userID string) (*entities.Wishlist, error) {
	var wishlist entities.Wishlist
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&wishlist)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If no wishlist exists for the user, create a new one
			wishlist := entities.NewWishlist(userID)
			err := r.Create(ctx, wishlist)
			if err != nil {
				return nil, err
			}
			return wishlist, nil
		}
		return nil, err
	}
	return &wishlist, nil
}

// Update updates a wishlist in the database
func (r *WishlistRepository) Update(ctx context.Context, wishlist *entities.Wishlist) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": wishlist.ID}, wishlist)
	return err
}

// Delete removes a wishlist from the database
func (r *WishlistRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// AddItem adds an item to a wishlist
func (r *WishlistRepository) AddItem(ctx context.Context, wishlistID string, item *entities.WishlistItem) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": wishlistID},
		bson.M{"$push": bson.M{"items": item}, "$set": bson.M{"updated_at": item.AddedAt}},
	)
	return err
}

// RemoveItem removes an item from a wishlist
func (r *WishlistRepository) RemoveItem(ctx context.Context, wishlistID string, itemID string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": wishlistID},
		bson.M{
			"$pull": bson.M{"items": bson.M{"id": itemID}},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// ClearWishlist removes all items from a wishlist
func (r *WishlistRepository) ClearWishlist(ctx context.Context, wishlistID string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": wishlistID},
		bson.M{"$set": bson.M{"items": []entities.WishlistItem{}, "updated_at": time.Now()}},
	)
	return err
}
