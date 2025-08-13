package services

import (
	"context"
	"errors"
	"strings"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CategoryService implements the category service interface
type CategoryService struct {
	categoryRepo ports.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo ports.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(ctx context.Context, name, description, slug string, parentID string) (*entities.Category, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}

	if slug == "" {
		// Generate slug from name if not provided
		slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	}

	// Convert parentID to ObjectID
	var parentObjectID primitive.ObjectID
	if parentID != "" {
		var err error
		parentObjectID, err = primitive.ObjectIDFromHex(parentID)
		if err != nil {
			return nil, errors.New("invalid parent ID")
		}

		// Check if parent exists
		_, err = s.categoryRepo.GetByID(ctx, parentObjectID)
		if err != nil {
			return nil, errors.New("parent category not found")
		}
	} else {
		parentObjectID = primitive.NilObjectID
	}

	// Check if slug is unique
	existing, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err == nil && existing != nil {
		return nil, errors.New("category with this slug already exists")
	}

	category := entities.NewCategory(name, description, slug, parentObjectID)
	err = s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategory retrieves a category by ID
func (s *CategoryService) GetCategory(ctx context.Context, id string) (*entities.Category, error) {
	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}
	return s.categoryRepo.GetByID(ctx, categoryID)
}

// GetCategoryBySlug retrieves a category by slug
func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*entities.Category, error) {
	return s.categoryRepo.GetBySlug(ctx, slug)
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(ctx context.Context, category *entities.Category) error {
	// Verify category exists
	existing, err := s.categoryRepo.GetByID(ctx, category.ID)
	if err != nil {
		return err
	}

	// Check if slug is unique (if changed)
	if existing.Slug != category.Slug {
		existingWithSlug, err := s.categoryRepo.GetBySlug(ctx, category.Slug)
		if err == nil && existingWithSlug != nil && existingWithSlug.ID != category.ID {
			return errors.New("category with this slug already exists")
		}
	}

	return s.categoryRepo.Update(ctx, category)
}

// DeleteCategory removes a category
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	// Convert ID to ObjectID
	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid category ID")
	}

	// Verify category exists
	_, err = s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(ctx, categoryID)
}

// ListCategories retrieves categories with filtering and pagination
func (s *CategoryService) ListCategories(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Category, int, error) {
	return s.categoryRepo.List(ctx, filter, page, limit)
}

// GetRootCategories retrieves all top-level categories
func (s *CategoryService) GetRootCategories(ctx context.Context) ([]*entities.Category, error) {
	return s.categoryRepo.GetRootCategories(ctx)
}

// GetChildCategories retrieves all direct child categories of a parent
func (s *CategoryService) GetChildCategories(ctx context.Context, parentID string) ([]*entities.Category, error) {
	// Convert parentID to ObjectID
	parentObjectID, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, errors.New("invalid parent ID")
	}

	// Verify parent category exists
	_, err = s.categoryRepo.GetByID(ctx, parentObjectID)
	if err != nil {
		return nil, err
	}

	return s.categoryRepo.GetChildCategories(ctx, parentObjectID)
}

// GetCategoryTree builds a complete category tree starting from a root category
func (s *CategoryService) GetCategoryTree(ctx context.Context, rootID string) (*entities.Category, error) {
	rootObjectID, err := primitive.ObjectIDFromHex(rootID)
	if err != nil {
		return nil, errors.New("invalid root ID")
	}
	return s.categoryRepo.GetCategoryTree(ctx, rootObjectID)
}

// GetProductsByCategory retrieves products in a category and optionally its subcategories
func (s *CategoryService) GetProductsByCategory(ctx context.Context, categoryID string, includeSubcategories bool, page, limit int) ([]*entities.Product, int, error) {
	// Convert categoryID to ObjectID
	categoryObjectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, 0, errors.New("invalid category ID")
	}

	// Verify category exists
	_, err = s.categoryRepo.GetByID(ctx, categoryObjectID)
	if err != nil {
		return nil, 0, err
	}

	return s.categoryRepo.GetProductsByCategory(ctx, categoryObjectID, includeSubcategories, page, limit)
}
