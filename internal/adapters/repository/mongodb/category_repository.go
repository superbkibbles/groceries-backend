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

// CategoryCollection is the name of the categories collection in MongoDB
const CategoryCollection = "categories"

// CategoryRepository implements the category repository interface using MongoDB
type CategoryRepository struct {
	db          *mongo.Database
	collection  *mongo.Collection
	productRepo *ProductRepository
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *mongo.Database, productRepo *ProductRepository) *CategoryRepository {
	return &CategoryRepository{
		db:          db,
		collection:  db.Collection(CategoryCollection),
		productRepo: productRepo,
	}
}

// Create adds a new category to the database
func (r *CategoryRepository) Create(ctx context.Context, category *entities.Category) error {
	// If this is a subcategory, update its path and level based on parent
	if !category.ParentID.IsZero() {
		parent, err := r.GetByID(ctx, category.ParentID)
		if err != nil {
			return err
		}

		// Set the level based on parent's level
		category.Level = parent.Level + 1

		// Set the path to include all ancestors
		category.Path = append(append([]primitive.ObjectID{}, parent.Path...), parent.ID)
	}

	_, err := r.collection.InsertOne(ctx, category)
	return err
}

// GetByID retrieves a category by its ID
func (r *CategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Category, error) {
	var category entities.Category
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// GetBySlug retrieves a category by its slug
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*entities.Category, error) {
	var category entities.Category
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// Update updates an existing category
func (r *CategoryRepository) Update(ctx context.Context, category *entities.Category) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": category.ID}, category)
	return err
}

// Delete removes a category from the database
func (r *CategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	// First check if this category has children
	count, err := r.collection.CountDocuments(ctx, bson.M{"parent_id": id})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete category with subcategories")
	}

	// Check if there are products in this category
	products, _, err := r.productRepo.GetByCategory(ctx, id, 1, 1)
	if err != nil {
		return err
	}
	if len(products) > 0 {
		return errors.New("cannot delete category with products")
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List retrieves categories based on filters with pagination
func (r *CategoryRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Category, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	// Convert map to bson.M
	// bsonFilter := bson.M{}
	// for k, v := range filter {
	// 	bsonFilter[k] = v
	// }

	// Get total count
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Find categories with pagination
	// opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"level": 1, "name": 1})
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var categories []*entities.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, 0, err
	}

	return categories, int(total), nil
}

// GetRootCategories retrieves all top-level categories (no parent)
func (r *CategoryRepository) GetRootCategories(ctx context.Context) ([]*entities.Category, error) {
	filter := bson.M{"parent_id": ""}
	opts := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*entities.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetChildCategories retrieves all direct child categories of a parent
func (r *CategoryRepository) GetChildCategories(ctx context.Context, parentID primitive.ObjectID) ([]*entities.Category, error) {
	filter := bson.M{"parent_id": parentID}
	opts := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*entities.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryTree builds a complete category tree starting from a root category
func (r *CategoryRepository) GetCategoryTree(ctx context.Context, rootID primitive.ObjectID) (*entities.Category, error) {
	// Get the root category
	root, err := r.GetByID(ctx, rootID)
	if err != nil {
		return nil, err
	}

	// Recursively build the tree
	err = r.buildCategoryTree(ctx, root)
	if err != nil {
		return nil, err
	}

	return root, nil
}

// buildCategoryTree is a helper function to recursively build the category tree
func (r *CategoryRepository) buildCategoryTree(ctx context.Context, category *entities.Category) error {
	// Get all children of this category
	children, err := r.GetChildCategories(ctx, category.ID)
	if err != nil {
		return err
	}

	// If no children, return
	if len(children) == 0 {
		category.Children = []entities.Category{}
		return nil
	}

	// Process each child
	for _, child := range children {
		// Recursively build the tree for this child
		err = r.buildCategoryTree(ctx, child)
		if err != nil {
			return err
		}

		// Add the child to the parent's children
		category.Children = append(category.Children, *child)
	}

	return nil
}

// GetProductsByCategory retrieves products in a category and optionally its subcategories
func (r *CategoryRepository) GetProductsByCategory(ctx context.Context, categoryID primitive.ObjectID, includeSubcategories bool, page, limit int) ([]*entities.Product, int, error) {
	// If we don't need to include subcategories, just use the product repo directly
	if !includeSubcategories {
		return r.productRepo.GetByCategory(ctx, categoryID, page, limit)
	}

	// Check if the category exists
	_, err := r.GetByID(ctx, categoryID)
	if err != nil {
		return nil, 0, err
	}

	// Get all category IDs to include in the search
	categoryIDs := []primitive.ObjectID{categoryID}

	// If we need to include subcategories, build the full tree and collect all IDs
	if includeSubcategories {
		// Build the full tree
		categoryTree, err := r.GetCategoryTree(ctx, categoryID)
		if err != nil {
			return nil, 0, err
		}

		// Collect all category IDs in the tree
		categoryIDs = r.collectCategoryIDs(categoryTree)
	}

	// Create a filter to find products in any of these categories
	filter := bson.M{"categories": bson.M{"$in": categoryIDs}}

	// Use the product repo to get the products
	return r.productRepo.List(ctx, filter, page, limit)
}

// collectCategoryIDs is a helper function to collect all category IDs in a tree
func (r *CategoryRepository) collectCategoryIDs(category *entities.Category) []primitive.ObjectID {
	ids := []primitive.ObjectID{category.ID}

	for _, child := range category.Children {
		childIDs := r.collectCategoryIDs(&child)
		ids = append(ids, childIDs...)
	}

	return ids
}

// AddTranslation adds a translation for a specific language
func (r *CategoryRepository) AddTranslation(ctx context.Context, categoryID primitive.ObjectID, language string, translation entities.Translation) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": categoryID},
		bson.M{"$set": bson.M{"translations." + language: translation}},
	)
	return err
}

// UpdateTranslation updates a translation for a specific language
func (r *CategoryRepository) UpdateTranslation(ctx context.Context, categoryID primitive.ObjectID, language string, translation entities.Translation) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": categoryID},
		bson.M{"$set": bson.M{"translations." + language: translation}},
	)
	return err
}

// DeleteTranslation deletes a translation for a specific language
func (r *CategoryRepository) DeleteTranslation(ctx context.Context, categoryID primitive.ObjectID, language string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": categoryID},
		bson.M{"$unset": bson.M{"translations." + language: ""}},
	)
	return err
}

// GetTranslations retrieves all translations for a category
func (r *CategoryRepository) GetTranslations(ctx context.Context, categoryID primitive.ObjectID) (map[string]entities.Translation, error) {
	var category entities.Category
	err := r.collection.FindOne(ctx, bson.M{"_id": categoryID}).Decode(&category)
	if err != nil {
		return nil, err
	}
	return category.Translations, nil
}
