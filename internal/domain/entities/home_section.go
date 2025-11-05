package entities

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HomeSectionType represents the type of a home section
type HomeSectionType string

const (
	HomeSectionTypeProducts   HomeSectionType = "products"
	HomeSectionTypeCategories HomeSectionType = "categories"
)

// LocalizedText stores a localized title/name per language
type LocalizedText map[string]string

// HomeSection represents a configurable section on the home page
type HomeSection struct {
	ID          primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Type        HomeSectionType      `json:"type" bson:"type"`
	Title       LocalizedText        `json:"title" bson:"title"` // e.g., {"en": "Best Sellers", "ar": "الأكثر مبيعًا"}
	ProductIDs  []primitive.ObjectID `json:"product_ids,omitempty" bson:"product_ids,omitempty"`
	CategoryIDs []primitive.ObjectID `json:"category_ids,omitempty" bson:"category_ids,omitempty"`
	Order       int                  `json:"order" bson:"order"`
	Active      bool                 `json:"active" bson:"active"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}

func NewHomeSection(sectionType HomeSectionType, title LocalizedText, productIDs, categoryIDs []primitive.ObjectID, order int, active bool) (*HomeSection, error) {
	if sectionType != HomeSectionTypeProducts && sectionType != HomeSectionTypeCategories {
		return nil, errors.New("invalid home section type")
	}
	if len(title) == 0 {
		return nil, errors.New("title is required")
	}
	now := time.Now()
	return &HomeSection{
		ID:          primitive.NewObjectID(),
		Type:        sectionType,
		Title:       title,
		ProductIDs:  productIDs,
		CategoryIDs: categoryIDs,
		Order:       order,
		Active:      active,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
