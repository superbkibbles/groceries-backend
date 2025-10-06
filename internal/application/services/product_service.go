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

// CreateProduct creates a new product with translations
func (s *ProductService) CreateProduct(ctx context.Context, categories []string, attributes map[string]interface{}, sku string, price float64, stockQuantity int, images []string, translations map[string]entities.Translation) (*entities.Product, error) {
	if sku == "" {
		return nil, errors.New("product SKU is required")
	}
	if price <= 0 {
		return nil, errors.New("product price must be greater than zero")
	}
	if len(translations) == 0 {
		return nil, errors.New("at least one translation is required")
	}

	// Validate that we have at least English translation
	if _, hasEnglish := translations["en"]; !hasEnglish {
		return nil, errors.New("English translation is required")
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

	product := entities.NewProductWithTranslations(categoryIDs, attributes, sku, price, stockQuantity, images, translations)
	err := s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct retrieves a product by ID and applies localization
func (s *ProductService) GetProduct(ctx context.Context, id string, language string) (*entities.Product, error) {
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Apply localization
	product.ApplyLocalization(language)
	return product, nil
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

// ListProducts retrieves products with filtering and pagination, applying localization
func (s *ProductService) ListProducts(ctx context.Context, filter map[string]interface{}, page, limit int, language string) ([]*entities.Product, int, error) {
	products, total, err := s.productRepo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Apply localization to each product
	for _, product := range products {
		product.ApplyLocalization(language)
	}

	return products, total, nil
}

// GetProductsByCategory retrieves products by category, applying localization
func (s *ProductService) GetProductsByCategory(ctx context.Context, category string, page, limit int, language string) ([]*entities.Product, int, error) {
	categoryID, err := primitive.ObjectIDFromHex(category)
	if err != nil {
		return nil, 0, errors.New("invalid category ID")
	}

	products, total, err := s.productRepo.GetByCategory(ctx, categoryID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Apply localization to each product
	for _, product := range products {
		product.ApplyLocalization(language)
	}

	return products, total, nil
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

// AddProductTranslation adds a translation for a specific language
func (s *ProductService) AddProductTranslation(ctx context.Context, productID string, language string, translation entities.Translation) error {
	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New("invalid product ID")
	}

	// Verify product exists
	_, err = s.productRepo.GetByID(ctx, productObjectID)
	if err != nil {
		return err
	}

	return s.productRepo.AddTranslation(ctx, productObjectID, language, translation)
}

// UpdateProductTranslation updates a translation for a specific language
func (s *ProductService) UpdateProductTranslation(ctx context.Context, productID string, language string, translation entities.Translation) error {
	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New("invalid product ID")
	}

	// Verify product exists
	_, err = s.productRepo.GetByID(ctx, productObjectID)
	if err != nil {
		return err
	}

	return s.productRepo.UpdateTranslation(ctx, productObjectID, language, translation)
}

// DeleteProductTranslation deletes a translation for a specific language
func (s *ProductService) DeleteProductTranslation(ctx context.Context, productID string, language string) error {
	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New("invalid product ID")
	}

	// Verify product exists
	_, err = s.productRepo.GetByID(ctx, productObjectID)
	if err != nil {
		return err
	}

	return s.productRepo.DeleteTranslation(ctx, productObjectID, language)
}

// GetProductTranslations retrieves all translations for a product
func (s *ProductService) GetProductTranslations(ctx context.Context, productID string) (map[string]entities.Translation, error) {
	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	// Verify product exists
	_, err = s.productRepo.GetByID(ctx, productObjectID)
	if err != nil {
		return nil, err
	}

	return s.productRepo.GetTranslations(ctx, productObjectID)
}
