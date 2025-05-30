package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CartItem represents a product variation in a user's cart
type CartItem struct {
	ID          string    `json:"id" bson:"id"`
	ProductID   string    `json:"product_id" bson:"product_id"`
	VariationID string    `json:"variation_id" bson:"variation_id"`
	SKU         string    `json:"sku" bson:"sku"`
	Name        string    `json:"name" bson:"name"`
	Price       float64   `json:"price" bson:"price"`
	Quantity    int       `json:"quantity" bson:"quantity"`
	Subtotal    float64   `json:"subtotal" bson:"subtotal"`
	AddedAt     time.Time `json:"added_at" bson:"added_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Cart represents a user's shopping cart in the e-commerce system
type Cart struct {
	ID          string      `json:"id" bson:"_id"`
	UserID      string      `json:"user_id" bson:"user_id"`
	Items       []*CartItem `json:"items" bson:"items"`
	TotalAmount float64     `json:"total_amount" bson:"total_amount"`
	CreatedAt   time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" bson:"updated_at"`
}

// NewCart creates a new cart with a unique ID
func NewCart(userID string) *Cart {
	now := time.Now()
	return &Cart{
		ID:          uuid.New().String(),
		UserID:      userID,
		Items:       []*CartItem{},
		TotalAmount: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddItem adds a product variation to the cart
func (c *Cart) AddItem(productID, variationID, sku, name string, price float64, quantity int) (*CartItem, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	// Check if the item already exists in the cart
	for i, item := range c.Items {
		if item.ProductID == productID && item.VariationID == variationID {
			// Update quantity instead of adding a new item
			c.Items[i].Quantity += quantity
			c.Items[i].Subtotal = float64(c.Items[i].Quantity) * c.Items[i].Price
			c.Items[i].UpdatedAt = time.Now()
			c.recalculateTotal()
			c.UpdatedAt = time.Now()
			return c.Items[i], nil
		}
	}

	// Add new item
	now := time.Now()
	item := &CartItem{
		ID:          uuid.New().String(),
		ProductID:   productID,
		VariationID: variationID,
		SKU:         sku,
		Name:        name,
		Price:       price,
		Quantity:    quantity,
		Subtotal:    price * float64(quantity),
		AddedAt:     now,
		UpdatedAt:   now,
	}

	c.Items = append(c.Items, item)
	c.recalculateTotal()
	c.UpdatedAt = time.Now()

	return item, nil
}

// UpdateItemQuantity updates the quantity of an item in the cart
func (c *Cart) UpdateItemQuantity(itemID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	for i, item := range c.Items {
		if item.ID == itemID {
			c.Items[i].Quantity = quantity
			c.Items[i].Subtotal = float64(quantity) * c.Items[i].Price
			c.Items[i].UpdatedAt = time.Now()
			c.recalculateTotal()
			c.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("item not found in cart")
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(itemID string) error {
	for i, item := range c.Items {
		if item.ID == itemID {
			// Remove item by replacing it with the last item and truncating the slice
			lastIndex := len(c.Items) - 1
			c.Items[i] = c.Items[lastIndex]
			c.Items = c.Items[:lastIndex]
			c.recalculateTotal()
			c.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("item not found in cart")
}

// ClearCart removes all items from the cart
func (c *Cart) ClearCart() {
	c.Items = []*CartItem{}
	c.TotalAmount = 0
	c.UpdatedAt = time.Now()
}

// ClearItems removes all items from the cart (alias for ClearCart)
func (c *Cart) ClearItems() {
	c.ClearCart()
}

// recalculateTotal recalculates the total amount of the cart
func (c *Cart) recalculateTotal() {
	total := 0.0
	for _, item := range c.Items {
		total += item.Subtotal
	}
	c.TotalAmount = total
}

// ConvertToOrder converts the cart to an order
func (c *Cart) ConvertToOrder(shippingInfo ShippingInfo) *Order {
	order := NewOrder(c.UserID, shippingInfo)

	// Copy items from cart to order
	for _, cartItem := range c.Items {
		order.AddItem(
			cartItem.ProductID,
			cartItem.VariationID,
			cartItem.SKU,
			cartItem.Name,
			cartItem.Price,
			cartItem.Quantity,
		)
	}

	return order
}
