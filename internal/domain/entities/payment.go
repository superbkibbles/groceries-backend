package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentMethodType represents the type of payment method
type PaymentMethodType string

const (
	PaymentMethodCreditCard     PaymentMethodType = "credit_card"
	PaymentMethodPayPal         PaymentMethodType = "paypal"
	PaymentMethodBankTransfer   PaymentMethodType = "bank_transfer"
	PaymentMethodCryptocurrency PaymentMethodType = "cryptocurrency"
)

// PaymentMethod represents a payment method available in the system
type PaymentMethod struct {
	ID primitive.ObjectID `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        PaymentMethodType      `json:"type"`
	Active      bool                   `json:"active"`
	Config      map[string]interface{} `json:"config,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// PaymentGateway represents a payment gateway configuration
type PaymentGateway struct {
	ID primitive.ObjectID `json:"id,omitempty"`
	Name      string                 `json:"name"`
	Provider  string                 `json:"provider"`
	Active    bool                   `json:"active"`
	Config    map[string]interface{} `json:"config,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// CustomerPaymentMethod represents a payment method saved by a customer
type CustomerPaymentMethod struct {
	ID primitive.ObjectID `json:"id,omitempty"`
	CustomerID      string    `json:"customer_id"`
	PaymentMethodID string    `json:"payment_method_id"`
	Token           string    `json:"token"`
	Last4           string    `json:"last4,omitempty"`
	ExpiryMonth     int       `json:"expiry_month,omitempty"`
	ExpiryYear      int       `json:"expiry_year,omitempty"`
	IsDefault       bool      `json:"is_default"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// NewPaymentMethod creates a new payment method
func NewPaymentMethod(name, description string, methodType PaymentMethodType, config map[string]interface{}) *PaymentMethod {
	now := time.Now()
	return &PaymentMethod{
		Name:        name,
		Description: description,
		Type:        methodType,
		Active:      true,
		Config:      config,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewPaymentGateway creates a new payment gateway
func NewPaymentGateway(name, provider string, config map[string]interface{}) *PaymentGateway {
	now := time.Now()
	return &PaymentGateway{
		Name:      name,
		Provider:  provider,
		Active:    true,
		Config:    config,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewCustomerPaymentMethod creates a new customer payment method
func NewCustomerPaymentMethod(customerID, paymentMethodID, token string, last4 string, expiryMonth, expiryYear int, isDefault bool) *CustomerPaymentMethod {
	now := time.Now()
	return &CustomerPaymentMethod{
		CustomerID:      customerID,
		PaymentMethodID: paymentMethodID,
		Token:           token,
		Last4:           last4,
		ExpiryMonth:     expiryMonth,
		ExpiryYear:      expiryYear,
		IsDefault:       isDefault,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
