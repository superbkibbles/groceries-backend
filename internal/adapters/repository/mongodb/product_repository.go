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

// ProductRepository implements the product repository interface using MongoDB
type ProductRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		db:         db,
		collection: db.Collection(ProductCollection),
	}
}

// Create adds a new product to the database
func (r *ProductRepository) Create(ctx context.Context, product *entities.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error) {
	var product entities.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, product *entities.Product) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": product.ID}, product)
	return err
}

// Delete removes a product from the database
func (r *ProductRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List retrieves products based on filters with pagination
func (r *ProductRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Product, int, error) {
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

	// Find products with pagination
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []*entities.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}

// GetByCategory retrieves products by category with pagination
func (r *ProductRepository) GetByCategory(ctx context.Context, category primitive.ObjectID, page, limit int) ([]*entities.Product, int, error) {
	filter := bson.M{"categories": category}
	return r.List(ctx, filter, page, limit)
}
