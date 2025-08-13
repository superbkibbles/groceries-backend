package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CartRepository is a MongoDB implementation of the CartRepository interface
type CartRepository struct {
	collection *mongo.Collection
}

// NewCartRepository creates a new CartRepository instance
func NewCartRepository(db *mongo.Database) *CartRepository {
	return &CartRepository{
		collection: db.Collection("carts"),
	}
}

// Create creates a new cart in the database
func (r *CartRepository) Create(ctx context.Context, cart *entities.Cart) error {
	_, err := r.collection.InsertOne(ctx, cart)
	return err
}

// GetByID retrieves a cart by its ID
func (r *CartRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Cart, error) {
	var cart entities.Cart
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cart)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// GetByUserID retrieves a cart by user ID
func (r *CartRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.Cart, error) {
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
	cart.UpdatedAt = time.Now()
	filter := bson.M{"_id": cart.ID}
	update := bson.M{"$set": cart}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete removes a cart from the database
func (r *CartRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// AddItem adds an item to a cart
func (r *CartRepository) AddItem(ctx context.Context, cartID primitive.ObjectID, item *entities.CartItem) error {
	filter := bson.M{"_id": cartID}
	update := bson.M{
		"$push": bson.M{"items": item},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateItem updates an item in a cart
func (r *CartRepository) UpdateItem(ctx context.Context, cartID primitive.ObjectID, itemID primitive.ObjectID, quantity int, price float64) error {
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

	// Update the item
	filter := bson.M{
		"_id":      cartID,
		"items.id": itemID,
	}
	update := bson.M{
		"$set": bson.M{
			"items.$.quantity":   quantity,
			"items.$.price":      price,
			"items.$.subtotal":   float64(quantity) * price,
			"updated_at":         time.Now(),
		},
	}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// RemoveItem removes an item from a cart
func (r *CartRepository) RemoveItem(ctx context.Context, cartID primitive.ObjectID, itemID primitive.ObjectID) error {
	// First get the cart to calculate the new total amount
	cart, err := r.GetByID(ctx, cartID)
	if err != nil {
		return err
	}

	// Find the item to remove
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

	// Remove the item
	filter := bson.M{"_id": cartID}
	update := bson.M{
		"$pull": bson.M{"items": bson.M{"id": itemID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// ClearCart removes all items from a cart
func (r *CartRepository) ClearCart(ctx context.Context, cartID primitive.ObjectID) error {
	filter := bson.M{"_id": cartID}
	update := bson.M{
		"$set": bson.M{
			"items":       []*entities.CartItem{},
			"total_amount": 0,
			"updated_at":   time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// List retrieves a paginated list of carts
func (r *CartRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Cart, int, error) {
	// Convert filter to bson.M
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip value
	skip := (page - 1) * limit

	// Find documents with pagination
	skip64 := int64(skip)
	limit64 := int64(limit)
	cursor, err := r.collection.Find(ctx, bsonFilter, &options.FindOptions{
		Skip:  &skip64,
		Limit: &limit64,
	})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var carts []*entities.Cart
	for cursor.Next(ctx) {
		var cart entities.Cart
		if err := cursor.Decode(&cart); err != nil {
			return nil, 0, err
		}
		carts = append(carts, &cart)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return carts, int(total), nil
}
