package entities

import (
	"time"

	"github.com/google/uuid"
)

// ShippingMethod represents a shipping method available in the system
type ShippingMethod struct {
	ID                    string    `json:"id" bson:"_id"`
	Name                  string    `json:"name" bson:"name"`
	Description           string    `json:"description" bson:"description"`
	BasePrice             float64   `json:"base_price" bson:"base_price"`
	EstimatedDeliveryDays int       `json:"estimated_delivery_days" bson:"estimated_delivery_days"`
	Active                bool      `json:"active" bson:"active"`
	CreatedAt             time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" bson:"updated_at"`
}

// ShippingZone represents a geographical zone with specific shipping rates
type ShippingZone struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Countries []string  `json:"countries" bson:"countries"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// ShippingRate represents the cost of a shipping method for a specific zone
type ShippingRate struct {
	ID               string    `json:"id" bson:"_id"`
	ShippingZoneID   string    `json:"shipping_zone_id" bson:"shipping_zone_id"`
	ShippingMethodID string    `json:"shipping_method_id" bson:"shipping_method_id"`
	Price            float64   `json:"price" bson:"price"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}

// NewShippingMethod creates a new shipping method
func NewShippingMethod(name, description string, basePrice float64, estimatedDeliveryDays int) *ShippingMethod {
	now := time.Now()
	return &ShippingMethod{
		ID:                    uuid.New().String(),
		Name:                  name,
		Description:           description,
		BasePrice:             basePrice,
		EstimatedDeliveryDays: estimatedDeliveryDays,
		Active:                true,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}

// NewShippingZone creates a new shipping zone
func NewShippingZone(name string, countries []string) *ShippingZone {
	now := time.Now()
	return &ShippingZone{
		ID:        uuid.New().String(),
		Name:      name,
		Countries: countries,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewShippingRate creates a new shipping rate
func NewShippingRate(shippingZoneID, shippingMethodID string, price float64) *ShippingRate {
	now := time.Now()
	return &ShippingRate{
		ID:               uuid.New().String(),
		ShippingZoneID:   shippingZoneID,
		ShippingMethodID: shippingMethodID,
		Price:            price,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}
