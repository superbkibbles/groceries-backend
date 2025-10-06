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

// CreateCategory creates a new category with translations
func (s *CategoryService) CreateCategory(ctx context.Context, slug string, parentID string, translations map[string]entities.Translation) (*entities.Category, error) {
	if len(translations) == 0 {
		return nil, errors.New("at least one translation is required")
	}

	// Validate that we have at least English translation
	if _, hasEnglish := translations["en"]; !hasEnglish {
		return nil, errors.New("English translation is required")
	}

	if slug == "" {
		// Generate slug from English name if not provided
		if enTranslation, exists := translations["en"]; exists {
			slug = strings.ToLower(strings.ReplaceAll(enTranslation.Name, " ", "-"))
		} else {
			return nil, errors.New("slug is required when no English translation provided")
		}
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

	category := entities.NewCategoryWithTranslations(slug, parentObjectID, translations)
	err = s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategory retrieves a category by ID and applies localization
func (s *CategoryService) GetCategory(ctx context.Context, id string, language string) (*entities.Category, error) {
	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Apply localization
	category.ApplyLocalization(language)
	return category, nil
}

// GetCategoryBySlug retrieves a category by slug and applies localization
func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string, language string) (*entities.Category, error) {
	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Apply localization
	category.ApplyLocalization(language)
	return category, nil
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

// ListCategories retrieves categories with filtering and pagination, applying localization
func (s *CategoryService) ListCategories(ctx context.Context, filter map[string]interface{}, page, limit int, language string) ([]*entities.Category, int, error) {
	categories, total, err := s.categoryRepo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Apply localization to each category
	for _, category := range categories {
		category.ApplyLocalization(language)
	}

	return categories, total, nil
}

// GetRootCategories retrieves all top-level categories, applying localization
func (s *CategoryService) GetRootCategories(ctx context.Context, language string) ([]*entities.Category, error) {
	categories, err := s.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		return nil, err
	}

	// Apply localization to each category
	for _, category := range categories {
		category.ApplyLocalization(language)
	}

	return categories, nil
}

// GetChildCategories retrieves all direct child categories of a parent, applying localization
func (s *CategoryService) GetChildCategories(ctx context.Context, parentID string, language string) ([]*entities.Category, error) {
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

	categories, err := s.categoryRepo.GetChildCategories(ctx, parentObjectID)
	if err != nil {
		return nil, err
	}

	// Apply localization to each category
	for _, category := range categories {
		category.ApplyLocalization(language)
	}

	return categories, nil
}

// GetCategoryTree builds a complete category tree starting from a root category, applying localization
func (s *CategoryService) GetCategoryTree(ctx context.Context, rootID string, language string) (*entities.Category, error) {
	rootObjectID, err := primitive.ObjectIDFromHex(rootID)
	if err != nil {
		return nil, errors.New("invalid root ID")
	}

	category, err := s.categoryRepo.GetCategoryTree(ctx, rootObjectID)
	if err != nil {
		return nil, err
	}

	// Apply localization recursively to the tree
	applyLocalizationToTree(category, language)
	return category, nil
}

// applyLocalizationToTree recursively applies localization to a category tree
func applyLocalizationToTree(category *entities.Category, language string) {
	if category == nil {
		return
	}

	// Apply localization to current category
	category.ApplyLocalization(language)

	// Apply localization to children
	for i := range category.Children {
		applyLocalizationToTree(&category.Children[i], language)
	}
}

// GetProductsByCategory retrieves products in a category and optionally its subcategories, applying localization
func (s *CategoryService) GetProductsByCategory(ctx context.Context, categoryID string, includeSubcategories bool, page, limit int, language string) ([]*entities.Product, int, error) {
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

	products, total, err := s.categoryRepo.GetProductsByCategory(ctx, categoryObjectID, includeSubcategories, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Apply localization to each product
	for _, product := range products {
		product.ApplyLocalization(language)
	}

	return products, total, nil
}

// AddCategoryTranslation adds a translation for a specific language
func (s *CategoryService) AddCategoryTranslation(ctx context.Context, categoryID string, language string, translation entities.Translation) error {
	categoryObjectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return errors.New("invalid category ID")
	}

	// Verify category exists
	_, err = s.categoryRepo.GetByID(ctx, categoryObjectID)
	if err != nil {
		return err
	}

	return s.categoryRepo.AddTranslation(ctx, categoryObjectID, language, translation)
}

// UpdateCategoryTranslation updates a translation for a specific language
func (s *CategoryService) UpdateCategoryTranslation(ctx context.Context, categoryID string, language string, translation entities.Translation) error {
	categoryObjectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return errors.New("invalid category ID")
	}

	// Verify category exists
	_, err = s.categoryRepo.GetByID(ctx, categoryObjectID)
	if err != nil {
		return err
	}

	return s.categoryRepo.UpdateTranslation(ctx, categoryObjectID, language, translation)
}

// DeleteCategoryTranslation deletes a translation for a specific language
func (s *CategoryService) DeleteCategoryTranslation(ctx context.Context, categoryID string, language string) error {
	categoryObjectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return errors.New("invalid category ID")
	}

	// Verify category exists
	_, err = s.categoryRepo.GetByID(ctx, categoryObjectID)
	if err != nil {
		return err
	}

	return s.categoryRepo.DeleteTranslation(ctx, categoryObjectID, language)
}

// GetCategoryTranslations retrieves all translations for a category
func (s *CategoryService) GetCategoryTranslations(ctx context.Context, categoryID string) (map[string]entities.Translation, error) {
	categoryObjectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	// Verify category exists
	_, err = s.categoryRepo.GetByID(ctx, categoryObjectID)
	if err != nil {
		return nil, err
	}

	return s.categoryRepo.GetTranslations(ctx, categoryObjectID)
}
