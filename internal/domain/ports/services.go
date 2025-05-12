package ports

import (
	"context"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
)

// ProductService defines the interface for product business logic
type ProductService interface {
	CreateProduct(ctx context.Context, name, description string, basePrice float64, categories []string) (*entities.Product, error)
	GetProduct(ctx context.Context, id string) (*entities.Product, error)
	UpdateProduct(ctx context.Context, product *entities.Product) error
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Product, int, error)
	GetProductsByCategory(ctx context.Context, category string, page, limit int) ([]*entities.Product, int, error)
	AddVariation(ctx context.Context, productID string, attributes map[string]interface{}, sku string, price float64, stockQuantity int, images []string) (*entities.Variation, error)
	UpdateVariation(ctx context.Context, productID, variationID string, attributes map[string]interface{}, sku string, price float64, stockQuantity int, images []string) error
	RemoveVariation(ctx context.Context, productID, variationID string) error
	UpdateStock(ctx context.Context, productID, variationID string, quantity int) error
}

// OrderService defines the interface for order business logic
type OrderService interface {
	CreateOrder(ctx context.Context, customerID string, shippingInfo entities.ShippingInfo) (*entities.Order, error)
	GetOrder(ctx context.Context, id string) (*entities.Order, error)
	AddItem(ctx context.Context, orderID, productID, variationID string, quantity int) error
	UpdateItemQuantity(ctx context.Context, orderID, productID, variationID string, quantity int) error
	RemoveItem(ctx context.Context, orderID, productID, variationID string) error
	UpdateOrderStatus(ctx context.Context, orderID string, status entities.OrderStatus) error
	SetPaymentInfo(ctx context.Context, orderID, method, transactionID string, amount float64) error
	SetTrackingInfo(ctx context.Context, orderID, carrier, trackingNum string) error
	ListOrders(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Order, int, error)
	GetCustomerOrders(ctx context.Context, customerID string, page, limit int) ([]*entities.Order, int, error)
}

// UserService defines the interface for user business logic
type UserService interface {
	Register(ctx context.Context, email, password, firstName, lastName string) (*entities.User, error)
	Login(ctx context.Context, phoneNumber string) (*entities.User, string, error) // Returns user, token, error
	SendOTP(ctx context.Context, phoneNumber string) error
	GetUser(ctx context.Context, id string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
	AddAddress(ctx context.Context, userID, name, addressLine1, addressLine2, city, state, country, postalCode, phone string, isDefault bool) (*entities.Address, error)
	UpdateAddress(ctx context.Context, address *entities.Address) error
	DeleteAddress(ctx context.Context, addressID string) error
	GetAddresses(ctx context.Context, userID string) ([]*entities.Address, error)
	SetDefaultAddress(ctx context.Context, userID, addressID string) error
	ListUsers(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.User, int, error)
}

// CategoryService defines the interface for category business logic
type CategoryService interface {
	CreateCategory(ctx context.Context, name, description, slug string, parentID string) (*entities.Category, error)
	GetCategory(ctx context.Context, id string) (*entities.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*entities.Category, error)
	UpdateCategory(ctx context.Context, category *entities.Category) error
	DeleteCategory(ctx context.Context, id string) error
	ListCategories(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Category, int, error)
	GetRootCategories(ctx context.Context) ([]*entities.Category, error)
	GetChildCategories(ctx context.Context, parentID string) ([]*entities.Category, error)
	GetCategoryTree(ctx context.Context, rootID string) (*entities.Category, error)
	GetProductsByCategory(ctx context.Context, categoryID string, includeSubcategories bool, page, limit int) ([]*entities.Product, int, error)
}
