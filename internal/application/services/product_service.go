package services

import (
	"context"
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
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
func (s *ProductService) CreateProduct(ctx context.Context, name, description string, basePrice float64, categories []string) (*entities.Product, error) {
	if name == "" {
		return nil, errors.New("product name is required")
	}

	product := entities.NewProduct(name, description, basePrice, categories)
	err := s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id string) (*entities.Product, error) {
	return s.productRepo.GetByID(ctx, id)
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
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(ctx, id)
}

// ListProducts retrieves products with filtering and pagination
func (s *ProductService) ListProducts(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Product, int, error) {
	return s.productRepo.List(ctx, filter, page, limit)
}

// GetProductsByCategory retrieves products by category
func (s *ProductService) GetProductsByCategory(ctx context.Context, category string, page, limit int) ([]*entities.Product, int, error) {
	return s.productRepo.GetByCategory(ctx, category, page, limit)
}

// AddVariation adds a new variation to a product
func (s *ProductService) AddVariation(ctx context.Context, productID string, attributes map[string]interface{}, sku string, price float64, stockQuantity int, images []string) (*entities.Variation, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	variation, err := product.AddVariation(attributes, sku, price, stockQuantity, images)
	if err != nil {
		return nil, err
	}

	err = s.productRepo.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return variation, nil
}

// UpdateVariation updates an existing product variation
func (s *ProductService) UpdateVariation(ctx context.Context, productID, variationID string, attributes map[string]interface{}, sku string, price float64, stockQuantity int, images []string) error {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	err = product.UpdateVariation(variationID, attributes, sku, price, stockQuantity, images)
	if err != nil {
		return err
	}

	return s.productRepo.Update(ctx, product)
}

// RemoveVariation removes a variation from a product
func (s *ProductService) RemoveVariation(ctx context.Context, productID, variationID string) error {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	err = product.RemoveVariation(variationID)
	if err != nil {
		return err
	}

	return s.productRepo.Update(ctx, product)
}

// UpdateStock updates the stock quantity for a product variation
func (s *ProductService) UpdateStock(ctx context.Context, productID, variationID string, quantity int) error {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// Verify variation exists
	_, err = product.GetVariation(variationID)
	if err != nil {
		return err
	}

	return s.productRepo.UpdateStock(ctx, variationID, quantity)
}
