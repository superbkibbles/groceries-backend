package services

import (
	"context"
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductService implements the product service interface
type ProductService struct {
	productRepo ports.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo ports.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, name, description string, categories []string, attributes map[string]interface{}, sku string, price float64, stockQuantity int, images []string) (*entities.Product, error) {
	if name == "" {
		return nil, errors.New("product name is required")
	}
	if sku == "" {
		return nil, errors.New("product SKU is required")
	}
	if price <= 0 {
		return nil, errors.New("product price must be greater than zero")
	}

	// Convert categories from strings to ObjectIDs
	categoryIDs := make([]primitive.ObjectID, len(categories))
	for i, catStr := range categories {
		catID, err := primitive.ObjectIDFromHex(catStr)
		if err != nil {
			return nil, errors.New("invalid category ID: " + catStr)
		}
		categoryIDs[i] = catID
	}

	product := entities.NewProduct(name, description, categoryIDs, attributes, sku, price, stockQuantity, images)
	err := s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id string) (*entities.Product, error) {
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}
	return s.productRepo.GetByID(ctx, productID)
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, product *entities.Product) error {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return err
	}

	return s.productRepo.Update(ctx, product)
}

// DeleteProduct removes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	// Convert ID to ObjectID
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	// Verify product exists
	_, err = s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(ctx, productID)
}

// ListProducts retrieves products with filtering and pagination
func (s *ProductService) ListProducts(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Product, int, error) {
	return s.productRepo.List(ctx, filter, page, limit)
}

// GetProductsByCategory retrieves products by category
func (s *ProductService) GetProductsByCategory(ctx context.Context, category string, page, limit int) ([]*entities.Product, int, error) {
	categoryID, err := primitive.ObjectIDFromHex(category)
	if err != nil {
		return nil, 0, errors.New("invalid category ID")
	}
	return s.productRepo.GetByCategory(ctx, categoryID, page, limit)
}

// UpdateStock updates the stock quantity for a product
func (s *ProductService) UpdateStock(ctx context.Context, productID string, quantity int) error {
	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New("invalid product ID")
	}

	product, err := s.productRepo.GetByID(ctx, productObjectID)
	if err != nil {
		return err
	}

	err = product.UpdateStock(quantity)
	if err != nil {
		return err
	}

	return s.productRepo.Update(ctx, product)
}
