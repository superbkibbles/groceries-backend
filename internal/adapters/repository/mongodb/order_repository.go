package mongodb

import (
	"context"
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OrderRepository implements the order repository interface using MongoDB
type OrderRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		db:         db,
		collection: db.Collection(OrderCollection),
	}
}

// Create adds a new order to the database
func (r *OrderRepository) Create(ctx context.Context, order *entities.Order) error {
	_, err := r.collection.InsertOne(ctx, order)
	return err
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Order, error) {
	var order entities.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

// Update updates an existing order
func (r *OrderRepository) Update(ctx context.Context, order *entities.Order) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": order.ID}, order)
	return err
}

// Delete removes an order from the database
func (r *OrderRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List retrieves orders based on filters with pagination
func (r *OrderRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	// Convert map to bson.M
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	// Find orders with pagination
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})
	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var orders []*entities.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

// GetByCustomerID retrieves orders for a specific customer with pagination
func (r *OrderRepository) GetByCustomerID(ctx context.Context, customerID primitive.ObjectID, page, limit int) ([]*entities.Order, int, error) {
	filter := bson.M{"customer_id": customerID}
	return r.List(ctx, filter, page, limit)
}
