package ports

import (
	"context"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
)

// ShippingService defines the interface for shipping business logic
type ShippingService interface {
	// Shipping Methods
	CreateShippingMethod(ctx context.Context, name, description string, basePrice float64, estimatedDeliveryDays int) (*entities.ShippingMethod, error)
	GetShippingMethod(ctx context.Context, id string) (*entities.ShippingMethod, error)
	UpdateShippingMethod(ctx context.Context, method *entities.ShippingMethod) error
	DeleteShippingMethod(ctx context.Context, id string) error
	ListShippingMethods(ctx context.Context, active bool) ([]*entities.ShippingMethod, error)

	// Shipping Zones
	CreateShippingZone(ctx context.Context, name string, countries []string) (*entities.ShippingZone, error)
	GetShippingZone(ctx context.Context, id string) (*entities.ShippingZone, error)
	UpdateShippingZone(ctx context.Context, zone *entities.ShippingZone) error
	DeleteShippingZone(ctx context.Context, id string) error
	ListShippingZones(ctx context.Context) ([]*entities.ShippingZone, error)

	// Shipping Rates
	CreateShippingRate(ctx context.Context, zoneID, methodID string, price float64) (*entities.ShippingRate, error)
	GetShippingRate(ctx context.Context, id string) (*entities.ShippingRate, error)
	UpdateShippingRate(ctx context.Context, rate *entities.ShippingRate) error
	DeleteShippingRate(ctx context.Context, id string) error
	GetShippingRatesByZone(ctx context.Context, zoneID string) ([]*entities.ShippingRate, error)
	GetShippingRatesByMethod(ctx context.Context, methodID string) ([]*entities.ShippingRate, error)
	CalculateShippingCost(ctx context.Context, countryCode string, items []*entities.CartItem) ([]map[string]interface{}, error)
}

// PaymentService defines the interface for payment business logic
type PaymentService interface {
	// Payment Methods
	CreatePaymentMethod(ctx context.Context, name, description string, methodType entities.PaymentMethodType, config map[string]interface{}) (*entities.PaymentMethod, error)
	GetPaymentMethod(ctx context.Context, id string) (*entities.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, method *entities.PaymentMethod) error
	DeletePaymentMethod(ctx context.Context, id string) error
	ListPaymentMethods(ctx context.Context, active bool) ([]*entities.PaymentMethod, error)

	// Payment Gateways
	CreatePaymentGateway(ctx context.Context, name, provider string, config map[string]interface{}) (*entities.PaymentGateway, error)
	GetPaymentGateway(ctx context.Context, id string) (*entities.PaymentGateway, error)
	UpdatePaymentGateway(ctx context.Context, gateway *entities.PaymentGateway) error
	DeletePaymentGateway(ctx context.Context, id string) error
	ListPaymentGateways(ctx context.Context, active bool) ([]*entities.PaymentGateway, error)

	// Customer Payment Methods
	CreateCustomerPaymentMethod(ctx context.Context, customerID, paymentMethodID, token string, last4 string, expiryMonth, expiryYear int, isDefault bool) (*entities.CustomerPaymentMethod, error)
	GetCustomerPaymentMethod(ctx context.Context, id string) (*entities.CustomerPaymentMethod, error)
	UpdateCustomerPaymentMethod(ctx context.Context, method *entities.CustomerPaymentMethod) error
	DeleteCustomerPaymentMethod(ctx context.Context, id string) error
	ListCustomerPaymentMethods(ctx context.Context, customerID string) ([]*entities.CustomerPaymentMethod, error)
	SetDefaultCustomerPaymentMethod(ctx context.Context, customerID, methodID string) error

	// Payment Processing
	ProcessPayment(ctx context.Context, orderID string, paymentMethodID string, amount float64) (*entities.PaymentInfo, error)
	VerifyPayment(ctx context.Context, transactionID string) (bool, error)
	RefundPayment(ctx context.Context, orderID string, amount float64, reason string) error
}

// NotificationService defines the interface for notification business logic
type NotificationService interface {
	// Notifications
	CreateNotification(ctx context.Context, userID string, notificationType entities.NotificationType, title, message string, data map[string]interface{}) (*entities.Notification, error)
	GetNotification(ctx context.Context, id string) (*entities.Notification, error)
	MarkNotificationAsRead(ctx context.Context, id string) error
	MarkNotificationAsArchived(ctx context.Context, id string) error
	DeleteNotification(ctx context.Context, id string) error
	ListUserNotifications(ctx context.Context, userID string, status entities.NotificationStatus, page, limit int) ([]*entities.Notification, int, error)
	CountUnreadNotifications(ctx context.Context, userID string) (int, error)
	MarkAllNotificationsAsRead(ctx context.Context, userID string) error

	// Notification Templates
	CreateNotificationTemplate(ctx context.Context, name string, notificationType entities.NotificationType, titleTemplate, messageTemplate string) (*entities.NotificationTemplate, error)
	GetNotificationTemplate(ctx context.Context, id string) (*entities.NotificationTemplate, error)
	UpdateNotificationTemplate(ctx context.Context, template *entities.NotificationTemplate) error
	DeleteNotificationTemplate(ctx context.Context, id string) error
	ListNotificationTemplates(ctx context.Context, active bool) ([]*entities.NotificationTemplate, error)

	// Notification Sending
	SendOrderStatusNotification(ctx context.Context, orderID string, status entities.OrderStatus) error
	SendPaymentNotification(ctx context.Context, orderID string, paymentInfo entities.PaymentInfo) error
	SendShippingNotification(ctx context.Context, orderID string, trackingInfo entities.ShippingInfo) error
	SendSystemNotification(ctx context.Context, userID string, title, message string) error
}
