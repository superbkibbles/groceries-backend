package ports

import (
	"context"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, id string) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Product, int, error)
	GetByCategory(ctx context.Context, category string, page, limit int) ([]*entities.Product, int, error)
	UpdateStock(ctx context.Context, variationID string, quantity int) error
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(ctx context.Context, category *entities.Category) error
	GetByID(ctx context.Context, id string) (*entities.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entities.Category, error)
	Update(ctx context.Context, category *entities.Category) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Category, int, error)
	GetRootCategories(ctx context.Context) ([]*entities.Category, error)
	GetChildCategories(ctx context.Context, parentID string) ([]*entities.Category, error)
	GetCategoryTree(ctx context.Context, rootID string) (*entities.Category, error)
	GetProductsByCategory(ctx context.Context, categoryID string, includeSubcategories bool, page, limit int) ([]*entities.Product, int, error)
}

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	GetByID(ctx context.Context, id string) (*entities.Order, error)
	Update(ctx context.Context, order *entities.Order) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Order, int, error)
	GetByCustomerID(ctx context.Context, customerID string, page, limit int) ([]*entities.Order, int, error)
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.User, int, error)
	AddAddress(ctx context.Context, address *entities.Address) error
	UpdateAddress(ctx context.Context, address *entities.Address) error
	DeleteAddress(ctx context.Context, id string) error
	GetAddressesByUserID(ctx context.Context, userID string) ([]*entities.Address, error)
}

// ReviewRepository defines the interface for review data access
type ReviewRepository interface {
	Create(ctx context.Context, review *entities.Review) error
	GetByID(ctx context.Context, id string) (*entities.Review, error)
	Update(ctx context.Context, review *entities.Review) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Review, int, error)
	GetByProductID(ctx context.Context, productID string, page, limit int) ([]*entities.Review, int, error)
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Review, int, error)
	GetByOrderID(ctx context.Context, orderID string) ([]*entities.Review, error)
	CheckUserReviewEligibility(ctx context.Context, userID string, productID string) (bool, error)
}
