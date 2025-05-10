package entities

import (
	"time"

	"github.com/google/uuid"
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
	ID          string                 `json:"id" bson:"_id"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Type        PaymentMethodType      `json:"type" bson:"type"`
	Active      bool                   `json:"active" bson:"active"`
	Config      map[string]interface{} `json:"config,omitempty" bson:"config,omitempty"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

// PaymentGateway represents a payment gateway configuration
type PaymentGateway struct {
	ID        string                 `json:"id" bson:"_id"`
	Name      string                 `json:"name" bson:"name"`
	Provider  string                 `json:"provider" bson:"provider"`
	Active    bool                   `json:"active" bson:"active"`
	Config    map[string]interface{} `json:"config,omitempty" bson:"config,omitempty"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" bson:"updated_at"`
}

// CustomerPaymentMethod represents a payment method saved by a customer
type CustomerPaymentMethod struct {
	ID              string    `json:"id" bson:"_id"`
	CustomerID      string    `json:"customer_id" bson:"customer_id"`
	PaymentMethodID string    `json:"payment_method_id" bson:"payment_method_id"`
	Token           string    `json:"token" bson:"token"`
	Last4           string    `json:"last4,omitempty" bson:"last4,omitempty"`
	ExpiryMonth     int       `json:"expiry_month,omitempty" bson:"expiry_month,omitempty"`
	ExpiryYear      int       `json:"expiry_year,omitempty" bson:"expiry_year,omitempty"`
	IsDefault       bool      `json:"is_default" bson:"is_default"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// NewPaymentMethod creates a new payment method
func NewPaymentMethod(name, description string, methodType PaymentMethodType, config map[string]interface{}) *PaymentMethod {
	now := time.Now()
	return &PaymentMethod{
		ID:          uuid.New().String(),
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
		ID:        uuid.New().String(),
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
		ID:              uuid.New().String(),
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
