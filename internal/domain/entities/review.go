package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Review represents a product review in the e-commerce system
type Review struct {
	ID        string    `json:"id" bson:"_id"`
	ProductID string    `json:"product_id" bson:"product_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	OrderID   string    `json:"order_id" bson:"order_id"` // Links review to a specific order
	Rating    int       `json:"rating" bson:"rating"`     // Rating from 1-5
	Comment   string    `json:"comment" bson:"comment"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// NewReview creates a new review with validation
func NewReview(productID, userID, orderID string, rating int, comment string) (*Review, error) {
	// Validate rating (1-5 stars)
	if rating < 1 || rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	// Validate required fields
	if productID == "" || userID == "" || orderID == "" {
		return nil, errors.New("product ID, user ID, and order ID are required")
	}

	now := time.Now()
	return &Review{
		ID:        uuid.New().String(),
		ProductID: productID,
		UserID:    userID,
		OrderID:   orderID,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update updates the review content
func (r *Review) Update(rating int, comment string) error {
	// Validate rating
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	r.Rating = rating
	r.Comment = comment
	r.UpdatedAt = time.Now()
	return nil
}
