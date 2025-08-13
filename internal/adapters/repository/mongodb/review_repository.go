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

// ReviewCollection is the name of the reviews collection in MongoDB
const ReviewCollection = "reviews"

// ReviewRepository implements the review repository interface using MongoDB
type ReviewRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
	orderRepo  *OrderRepository
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(db *mongo.Database, orderRepo *OrderRepository) *ReviewRepository {
	return &ReviewRepository{
		db:         db,
		collection: db.Collection(ReviewCollection),
		orderRepo:  orderRepo,
	}
}

// Create adds a new review to the database
func (r *ReviewRepository) Create(ctx context.Context, review *entities.Review) error {
	// Check if the user is eligible to review this product
	eligible, err := r.CheckUserReviewEligibility(ctx, review.UserID, review.ProductID)
	if err != nil {
		return err
	}
	if !eligible {
		return errors.New("user is not eligible to review this product")
	}

	// Check if the user has already reviewed this product for this order
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"user_id":    review.UserID,
		"product_id": review.ProductID,
		"order_id":   review.OrderID,
	})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user has already reviewed this product for this order")
	}

	_, err = r.collection.InsertOne(ctx, review)
	return err
}

// GetByID retrieves a review by its ID
func (r *ReviewRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Review, error) {
	var review entities.Review
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&review)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("review not found")
		}
		return nil, err
	}
	return &review, nil
}

// Update updates an existing review
func (r *ReviewRepository) Update(ctx context.Context, review *entities.Review) error {
	// Verify the review exists and belongs to the user
	existing, err := r.GetByID(ctx, review.ID)
	if err != nil {
		return err
	}

	// Only allow updating the rating and comment
	existing.Rating = review.Rating
	existing.Comment = review.Comment
	existing.UpdatedAt = review.UpdatedAt

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": review.ID}, existing)
	return err
}

// Delete removes a review from the database
func (r *ReviewRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List retrieves reviews based on filters with pagination
func (r *ReviewRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Review, int, error) {
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

	// Find reviews with pagination
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})
	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reviews []*entities.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, 0, err
	}

	return reviews, int(total), nil
}

// GetByProductID retrieves reviews for a specific product with pagination
func (r *ReviewRepository) GetByProductID(ctx context.Context, productID string, page, limit int) ([]*entities.Review, int, error) {
	filter := bson.M{"product_id": productID}
	return r.List(ctx, filter, page, limit)
}

// GetByUserID retrieves reviews by a specific user with pagination
func (r *ReviewRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Review, int, error) {
	filter := bson.M{"user_id": userID}
	return r.List(ctx, filter, page, limit)
}

// GetByOrderID retrieves reviews for a specific order
func (r *ReviewRepository) GetByOrderID(ctx context.Context, orderID string) ([]*entities.Review, error) {
	filter := bson.M{"order_id": orderID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []*entities.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

// CheckUserReviewEligibility checks if a user is eligible to review a product
// A user is eligible if they have purchased the product in a completed order
func (r *ReviewRepository) CheckUserReviewEligibility(ctx context.Context, userID primitive.ObjectID, productID primitive.ObjectID) (bool, error) {
	// Find orders by this user that contain this product and are delivered
	filter := bson.M{
		"customer_id":      userID,
		"status":           entities.OrderStatusDelivered,
		"items.product_id": productID,
	}

	count, err := r.orderRepo.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
