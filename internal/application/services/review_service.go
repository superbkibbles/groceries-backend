package services

import (
	"context"
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// ReviewService implements the review service interface
type ReviewService struct {
	reviewRepo  ports.ReviewRepository
	orderRepo   ports.OrderRepository
	productRepo ports.ProductRepository
}

// NewReviewService creates a new review service
func NewReviewService(reviewRepo ports.ReviewRepository, orderRepo ports.OrderRepository, productRepo ports.ProductRepository) *ReviewService {
	return &ReviewService{
		reviewRepo:  reviewRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// CreateReview creates a new product review
// Only users who have purchased and received the product can leave a review
func (s *ReviewService) CreateReview(ctx context.Context, productID, userID, orderID string, rating int, comment string) (*entities.Review, error) {
	// Verify the product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check if the user is eligible to review this product
	eligible, err := s.reviewRepo.CheckUserReviewEligibility(ctx, userID, productID)
	if err != nil {
		return nil, err
	}
	if !eligible {
		return nil, errors.New("you can only review products you have purchased and received")
	}

	// Create the review
	review, err := entities.NewReview(productID, userID, orderID, rating, comment)
	if err != nil {
		return nil, err
	}

	// Save the review
	err = s.reviewRepo.Create(ctx, review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// GetReview retrieves a review by ID
func (s *ReviewService) GetReview(ctx context.Context, id string) (*entities.Review, error) {
	return s.reviewRepo.GetByID(ctx, id)
}

// UpdateReview updates an existing review
// Only the user who created the review can update it
func (s *ReviewService) UpdateReview(ctx context.Context, id, userID string, rating int, comment string) (*entities.Review, error) {
	// Get the existing review
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify the user owns this review
	if review.UserID != userID {
		return nil, errors.New("you can only update your own reviews")
	}

	// Update the review
	err = review.Update(rating, comment)
	if err != nil {
		return nil, err
	}

	// Save the updated review
	err = s.reviewRepo.Update(ctx, review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// DeleteReview removes a review
// Only the user who created the review or an admin can delete it
func (s *ReviewService) DeleteReview(ctx context.Context, id, userID string, isAdmin bool) error {
	// Get the existing review
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify the user owns this review or is an admin
	if review.UserID != userID && !isAdmin {
		return errors.New("you can only delete your own reviews")
	}

	// Delete the review
	return s.reviewRepo.Delete(ctx, id)
}

// ListReviews retrieves reviews with filtering and pagination
func (s *ReviewService) ListReviews(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Review, int, error) {
	return s.reviewRepo.List(ctx, filter, page, limit)
}

// GetProductReviews retrieves reviews for a specific product
func (s *ReviewService) GetProductReviews(ctx context.Context, productID string, page, limit int) ([]*entities.Review, int, error) {
	return s.reviewRepo.GetByProductID(ctx, productID, page, limit)
}

// GetUserReviews retrieves reviews by a specific user
func (s *ReviewService) GetUserReviews(ctx context.Context, userID string, page, limit int) ([]*entities.Review, int, error) {
	return s.reviewRepo.GetByUserID(ctx, userID, page, limit)
}

// CheckUserReviewEligibility checks if a user is eligible to review a product
func (s *ReviewService) CheckUserReviewEligibility(ctx context.Context, userID, productID string) (bool, error) {
	return s.reviewRepo.CheckUserReviewEligibility(ctx, userID, productID)
}
