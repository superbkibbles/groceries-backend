package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Category represents a product category in the e-commerce system
type Category struct {
	ID          string     `json:"id" bson:"_id"`
	Name        string     `json:"name" bson:"name"`
	Description string     `json:"description" bson:"description"`
	Slug        string     `json:"slug" bson:"slug"`
	ParentID    string     `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	Level       int        `json:"level" bson:"level"`
	Path        []string   `json:"path" bson:"path"`
	Children    []Category `json:"children,omitempty" bson:"-"`
	CreatedAt   time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" bson:"updated_at"`
}

// NewCategory creates a new category with a unique ID
func NewCategory(name, description, slug string, parentID string) *Category {
	now := time.Now()
	level := 1
	path := []string{}

	// If this is a subcategory, set level to 2 or more
	if parentID != "" {
		level = 2 // This will be updated when saved to reflect the actual level
	}

	return &Category{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Slug:        slug,
		ParentID:    parentID,
		Level:       level,
		Path:        path,
		Children:    []Category{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddChild adds a child category to this category
func (c *Category) AddChild(child *Category) error {
	if child == nil {
		return errors.New("child category cannot be nil")
	}

	// Set the parent ID and update level
	child.ParentID = c.ID
	child.Level = c.Level + 1

	// Update the path to include all ancestors
	child.Path = append(append([]string{}, c.Path...), c.ID)

	// Add to children collection
	c.Children = append(c.Children, *child)
	c.UpdatedAt = time.Now()

	return nil
}

// Update updates the category details
func (c *Category) Update(name, description, slug string) {
	c.Name = name
	c.Description = description
	c.Slug = slug
	c.UpdatedAt = time.Now()
}

// IsRoot checks if this category is a root category (no parent)
func (c *Category) IsRoot() bool {
	return c.ParentID == ""
}

// IsLeaf checks if this category is a leaf category (no children)
func (c *Category) IsLeaf() bool {
	return len(c.Children) == 0
}

// GetAncestorIDs returns all ancestor IDs including this category's ID
func (c *Category) GetAncestorIDs() []string {
	result := append(append([]string{}, c.Path...), c.ID)
	return result
}
