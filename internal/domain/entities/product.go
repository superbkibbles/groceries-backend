package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the e-commerce system
type Product struct {
	ID          string       `json:"id" bson:"_id"`
	Name        string       `json:"name" bson:"name"`
	Description string       `json:"description" bson:"description"`
	BasePrice   float64      `json:"base_price" bson:"base_price"`
	Categories  []string     `json:"categories" bson:"categories"` // Category IDs the product belongs to
	Variations  []*Variation `json:"variations" bson:"variations"`
	CreatedAt   time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" bson:"updated_at"`
}

// Variation represents a specific variation of a product
type Variation struct {
	ID            string                 `json:"id" bson:"id"`
	Attributes    map[string]interface{} `json:"attributes" bson:"attributes"`
	SKU           string                 `json:"sku" bson:"sku"`
	Price         float64                `json:"price" bson:"price"`
	StockQuantity int                    `json:"stock_quantity" bson:"stock_quantity"`
	Images        []string               `json:"images" bson:"images"`
}

// NewProduct creates a new product with a unique ID
func NewProduct(name, description string, basePrice float64, categories []string) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		BasePrice:   basePrice,
		Categories:  categories,
		Variations:  []*Variation{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddVariation adds a new variation to the product
func (p *Product) AddVariation(attributes map[string]interface{}, sku string, price float64, stockQuantity int, images []string) (*Variation, error) {
	// Validate SKU uniqueness
	for _, v := range p.Variations {
		if v.SKU == sku {
			return nil, errors.New("variation with this SKU already exists")
		}
	}

	variation := &Variation{
		ID:            uuid.New().String(),
		Attributes:    attributes,
		SKU:           sku,
		Price:         price,
		StockQuantity: stockQuantity,
		Images:        images,
	}

	p.Variations = append(p.Variations, variation)
	p.UpdatedAt = time.Now()

	return variation, nil
}

// UpdateVariation updates an existing variation
func (p *Product) UpdateVariation(variationID string, attributes map[string]interface{}, sku string, price float64, stockQuantity int, images []string) error {
	for i, v := range p.Variations {
		if v.ID == variationID {
			// Check if the new SKU conflicts with another variation
			if v.SKU != sku {
				for _, otherV := range p.Variations {
					if otherV.ID != variationID && otherV.SKU == sku {
						return errors.New("another variation with this SKU already exists")
					}
				}
			}

			p.Variations[i].Attributes = attributes
			p.Variations[i].SKU = sku
			p.Variations[i].Price = price
			p.Variations[i].StockQuantity = stockQuantity
			p.Variations[i].Images = images
			p.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("variation not found")
}

// RemoveVariation removes a variation from the product
func (p *Product) RemoveVariation(variationID string) error {
	for i, v := range p.Variations {
		if v.ID == variationID {
			p.Variations = append(p.Variations[:i], p.Variations[i+1:]...)
			p.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("variation not found")
}

// GetVariation returns a variation by ID
func (p *Product) GetVariation(variationID string) (*Variation, error) {
	for _, v := range p.Variations {
		if v.ID == variationID {
			return v, nil
		}
	}

	return nil, errors.New("variation not found")
}

// UpdateStock updates the stock quantity for a specific variation
func (p *Product) UpdateStock(variationID string, quantity int) error {
	for i, v := range p.Variations {
		if v.ID == variationID {
			p.Variations[i].StockQuantity = quantity
			p.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("variation not found")
}
