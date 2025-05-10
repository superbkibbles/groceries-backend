package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CartRepository is a MongoDB implementation of the CartRepository interface
type CartRepository struct {
	collection *mongo.Collection
}

// NewCartRepository creates a new CartRepository
func NewCartRepository(db *mongo.Database) *CartRepository {
	return &CartRepository{
		collection: db.Collection("carts"),
	}
}

// Create adds a new cart to the database
func (r *CartRepository) Create(ctx context.Context, cart *entities.Cart) error {
	_, err := r.collection.InsertOne(ctx, cart)
	return err
}

// GetByID retrieves a cart by its ID
func (r *CartRepository) GetByID(ctx context.Context, id string) (*entities.Cart, error) {
	var cart entities.Cart
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("cart not found")
		}
		return nil, err
	}
	return &cart, nil
}

// GetByUserID retrieves a cart by user ID
func (r *CartRepository) GetByUserID(ctx context.Context, userID string) (*entities.Cart, error) {
	var cart entities.Cart
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If no cart exists for the user, create a new one
			cart := entities.NewCart(userID)
			err := r.Create(ctx, cart)
			if err != nil {
				return nil, err
			}
			return cart, nil
		}
		return nil, err
	}
	return &cart, nil
}

// Update updates a cart in the database
func (r *CartRepository) Update(ctx context.Context, cart *entities.Cart) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": cart.ID}, cart)
	return err
}

// Delete removes a cart from the database
func (r *CartRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// AddItem adds an item to a cart
func (r *CartRepository) AddItem(ctx context.Context, cartID string, item *entities.CartItem) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": cartID},
		bson.M{
			"$push": bson.M{"items": item},
			"$set": bson.M{
				"updated_at":   time.Now(),
				"total_amount": bson.M{"$add": []interface{}{"$total_amount", item.Subtotal}},
			},
		},
	)
	return err
}

// UpdateItemQuantity updates the quantity of an item in a cart
func (r *CartRepository) UpdateItemQuantity(ctx context.Context, cartID string, itemID string, quantity int) error {
	// First get the cart to calculate the new total amount
	cart, err := r.GetByID(ctx, cartID)
	if err != nil {
		return err
	}

	// Find the item and calculate price difference
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

	// Calculate price difference
	oldSubtotal := item.Subtotal
	newSubtotal := item.Price * float64(quantity)
	priceDifference := newSubtotal - oldSubtotal

	// Update the item quantity and subtotal
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": cartID, "items.id": itemID},
		bson.M{
			"$set": bson.M{
				"items.$.quantity":   quantity,
				"items.$.subtotal":   newSubtotal,
				"items.$.updated_at": time.Now(),
				"updated_at":         time.Now(),
				"total_amount":       bson.M{"$add": []interface{}{"$total_amount", priceDifference}},
			},
		},
	)

	return err
}

// RemoveItem removes an item from a cart
func (r *CartRepository) RemoveItem(ctx context.Context, cartID string, itemID string) error {
	// First get the cart to calculate the new total amount
	cart, err := r.GetByID(ctx, cartID)
	if err != nil {
		return err
	}

	// Find the item to get its subtotal
	var itemSubtotal float64
	for _, item := range cart.Items {
		if item.ID == itemID {
			itemSubtotal = item.Subtotal
			break
		}
	}

	// Remove the item and update the total amount
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": cartID},
		bson.M{
			"$pull": bson.M{"items": bson.M{"id": itemID}},
			"$set": bson.M{
				"updated_at":   time.Now(),
				"total_amount": bson.M{"$subtract": []interface{}{"$total_amount", itemSubtotal}},
			},
		},
	)

	return err
}

// ClearCart removes all items from a cart
func (r *CartRepository) ClearCart(ctx context.Context, cartID string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": cartID},
		bson.M{
			"$set": bson.M{
				"items":        []entities.CartItem{},
				"total_amount": 0,
				"updated_at":   time.Now(),
			},
		},
	)
	return err
}
