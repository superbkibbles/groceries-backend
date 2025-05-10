package ports

import (
	"context"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
)

// ShippingMethodRepository defines the interface for shipping method data access
type ShippingMethodRepository interface {
	Create(ctx context.Context, method *entities.ShippingMethod) error
	GetByID(ctx context.Context, id string) (*entities.ShippingMethod, error)
	Update(ctx context.Context, method *entities.ShippingMethod) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, active bool) ([]*entities.ShippingMethod, error)
}

// ShippingZoneRepository defines the interface for shipping zone data access
type ShippingZoneRepository interface {
	Create(ctx context.Context, zone *entities.ShippingZone) error
	GetByID(ctx context.Context, id string) (*entities.ShippingZone, error)
	Update(ctx context.Context, zone *entities.ShippingZone) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*entities.ShippingZone, error)
	GetByCountry(ctx context.Context, countryCode string) (*entities.ShippingZone, error)
}

// ShippingRateRepository defines the interface for shipping rate data access
type ShippingRateRepository interface {
	Create(ctx context.Context, rate *entities.ShippingRate) error
	GetByID(ctx context.Context, id string) (*entities.ShippingRate, error)
	Update(ctx context.Context, rate *entities.ShippingRate) error
	Delete(ctx context.Context, id string) error
	GetByZone(ctx context.Context, zoneID string) ([]*entities.ShippingRate, error)
	GetByMethod(ctx context.Context, methodID string) ([]*entities.ShippingRate, error)
	GetByZoneAndMethod(ctx context.Context, zoneID, methodID string) (*entities.ShippingRate, error)
}

// PaymentMethodRepository defines the interface for payment method data access
type PaymentMethodRepository interface {
	Create(ctx context.Context, method *entities.PaymentMethod) error
	GetByID(ctx context.Context, id string) (*entities.PaymentMethod, error)
	Update(ctx context.Context, method *entities.PaymentMethod) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, active bool) ([]*entities.PaymentMethod, error)
	GetByType(ctx context.Context, methodType entities.PaymentMethodType) ([]*entities.PaymentMethod, error)
}

// PaymentGatewayRepository defines the interface for payment gateway data access
type PaymentGatewayRepository interface {
	Create(ctx context.Context, gateway *entities.PaymentGateway) error
	GetByID(ctx context.Context, id string) (*entities.PaymentGateway, error)
	Update(ctx context.Context, gateway *entities.PaymentGateway) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, active bool) ([]*entities.PaymentGateway, error)
	GetByProvider(ctx context.Context, provider string) (*entities.PaymentGateway, error)
}

// CustomerPaymentMethodRepository defines the interface for customer payment method data access
type CustomerPaymentMethodRepository interface {
	Create(ctx context.Context, method *entities.CustomerPaymentMethod) error
	GetByID(ctx context.Context, id string) (*entities.CustomerPaymentMethod, error)
	Update(ctx context.Context, method *entities.CustomerPaymentMethod) error
	Delete(ctx context.Context, id string) error
	GetByCustomer(ctx context.Context, customerID string) ([]*entities.CustomerPaymentMethod, error)
	GetDefaultByCustomer(ctx context.Context, customerID string) (*entities.CustomerPaymentMethod, error)
	SetDefault(ctx context.Context, id string) error
	ClearDefault(ctx context.Context, customerID string) error
}

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	Create(ctx context.Context, notification *entities.Notification) error
	GetByID(ctx context.Context, id string) (*entities.Notification, error)
	Update(ctx context.Context, notification *entities.Notification) error
	Delete(ctx context.Context, id string) error
	GetByUser(ctx context.Context, userID string, status entities.NotificationStatus, page, limit int) ([]*entities.Notification, int, error)
	CountUnread(ctx context.Context, userID string) (int, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAsArchived(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

// NotificationTemplateRepository defines the interface for notification template data access
type NotificationTemplateRepository interface {
	Create(ctx context.Context, template *entities.NotificationTemplate) error
	GetByID(ctx context.Context, id string) (*entities.NotificationTemplate, error)
	Update(ctx context.Context, template *entities.NotificationTemplate) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, active bool) ([]*entities.NotificationTemplate, error)
	GetByType(ctx context.Context, notificationType entities.NotificationType) ([]*entities.NotificationTemplate, error)
	GetByName(ctx context.Context, name string) (*entities.NotificationTemplate, error)
}
