package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the current status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents a customer order in the e-commerce system
type Order struct {
	ID           string       `json:"id" bson:"_id"`
	CustomerID   string       `json:"customer_id" bson:"customer_id"`
	Items        []*OrderItem `json:"items" bson:"items"`
	TotalAmount  float64      `json:"total_amount" bson:"total_amount"`
	Status       OrderStatus  `json:"status" bson:"status"`
	ShippingInfo ShippingInfo `json:"shipping_info" bson:"shipping_info"`
	PaymentInfo  PaymentInfo  `json:"payment_info" bson:"payment_info"`
	CreatedAt    time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" bson:"updated_at"`
}

// OrderItem represents a product variation in an order
type OrderItem struct {
	ProductID   string  `json:"product_id" bson:"product_id"`
	VariationID string  `json:"variation_id" bson:"variation_id"`
	SKU         string  `json:"sku" bson:"sku"`
	Name        string  `json:"name" bson:"name"`
	Price       float64 `json:"price" bson:"price"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Subtotal    float64 `json:"subtotal" bson:"subtotal"`
}

// ShippingInfo contains shipping details for an order
type ShippingInfo struct {
	Address     string `json:"address" bson:"address"`
	City        string `json:"city" bson:"city"`
	State       string `json:"state" bson:"state"`
	Country     string `json:"country" bson:"country"`
	PostalCode  string `json:"postal_code" bson:"postal_code"`
	Carrier     string `json:"carrier" bson:"carrier"`
	TrackingNum string `json:"tracking_num" bson:"tracking_num"`
}

// PaymentInfo contains payment details for an order
type PaymentInfo struct {
	Method        string    `json:"method" bson:"method"`
	TransactionID string    `json:"transaction_id" bson:"transaction_id"`
	PaidAt        time.Time `json:"paid_at,omitempty" bson:"paid_at,omitempty"`
	Status        string    `json:"status" bson:"status"` // pending, paid, failed, refunded, etc.
	Amount        float64   `json:"amount" bson:"amount"`
	Timestamp     time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

// NewOrder creates a new order with a unique ID
func NewOrder(customerID string, shippingInfo ShippingInfo) *Order {
	now := time.Now()
	return &Order{
		ID:           uuid.New().String(),
		CustomerID:   customerID,
		Items:        []*OrderItem{},
		TotalAmount:  0,
		Status:       OrderStatusPending,
		ShippingInfo: shippingInfo,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddItem adds a product variation to the order
func (o *Order) AddItem(productID, variationID, sku, name string, price float64, quantity int) error {
	if o.Status != OrderStatusPending {
		return errors.New("cannot modify a non-pending order")
	}

	// Check if the item already exists in the order
	for i, item := range o.Items {
		if item.ProductID == productID && item.VariationID == variationID {
			// Update quantity instead of adding a new item
			o.Items[i].Quantity += quantity
			o.Items[i].Subtotal = float64(o.Items[i].Quantity) * o.Items[i].Price
			o.recalculateTotal()
			o.UpdatedAt = time.Now()
			return nil
		}
	}

	// Add new item
	item := &OrderItem{
		ProductID:   productID,
		VariationID: variationID,
		SKU:         sku,
		Name:        name,
		Price:       price,
		Quantity:    quantity,
		Subtotal:    price * float64(quantity),
	}

	o.Items = append(o.Items, item)
	o.recalculateTotal()
	o.UpdatedAt = time.Now()

	return nil
}

// UpdateItemQuantity updates the quantity of an item in the order
func (o *Order) UpdateItemQuantity(productID, variationID string, quantity int) error {
	if o.Status != OrderStatusPending {
		return errors.New("cannot modify a non-pending order")
	}

	if quantity <= 0 {
		return o.RemoveItem(productID, variationID)
	}

	for i, item := range o.Items {
		if item.ProductID == productID && item.VariationID == variationID {
			o.Items[i].Quantity = quantity
			o.Items[i].Subtotal = float64(quantity) * item.Price
			o.recalculateTotal()
			o.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("item not found in order")
}

// RemoveItem removes an item from the order
func (o *Order) RemoveItem(productID, variationID string) error {
	if o.Status != OrderStatusPending {
		return errors.New("cannot modify a non-pending order")
	}

	for i, item := range o.Items {
		if item.ProductID == productID && item.VariationID == variationID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			o.recalculateTotal()
			o.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("item not found in order")
}

// UpdateStatus updates the status of the order
func (o *Order) UpdateStatus(status OrderStatus) error {
	// Validate status transition
	switch o.Status {
	case OrderStatusPending:
		if status != OrderStatusPaid && status != OrderStatusCancelled {
			return errors.New("invalid status transition")
		}
	case OrderStatusPaid:
		if status != OrderStatusShipped && status != OrderStatusCancelled {
			return errors.New("invalid status transition")
		}
	case OrderStatusShipped:
		if status != OrderStatusDelivered {
			return errors.New("invalid status transition")
		}
	case OrderStatusDelivered, OrderStatusCancelled:
		return errors.New("cannot change status of a delivered or cancelled order")
	}

	o.Status = status
	o.UpdatedAt = time.Now()

	return nil
}

// SetPaymentInfo sets the payment information for the order
func (o *Order) SetPaymentInfo(method, transactionID string, amount float64) error {
	if o.Status != OrderStatusPending {
		return errors.New("cannot set payment info for a non-pending order")
	}

	o.PaymentInfo = PaymentInfo{
		Method:        method,
		TransactionID: transactionID,
		PaidAt:        time.Now(),
		Amount:        amount,
	}

	o.UpdatedAt = time.Now()

	return nil
}

// SetTrackingInfo sets the shipping tracking information
func (o *Order) SetTrackingInfo(carrier, trackingNum string) error {
	if o.Status != OrderStatusPaid && o.Status != OrderStatusShipped {
		return errors.New("cannot set tracking info for an order that is not paid or shipped")
	}

	o.ShippingInfo.Carrier = carrier
	o.ShippingInfo.TrackingNum = trackingNum
	o.UpdatedAt = time.Now()

	return nil
}

// recalculateTotal recalculates the total amount of the order
func (o *Order) recalculateTotal() {
	total := 0.0
	for _, item := range o.Items {
		total += item.Subtotal
	}
	o.TotalAmount = total
}
