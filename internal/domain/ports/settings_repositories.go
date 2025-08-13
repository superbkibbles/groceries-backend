package ports

import (
	"context"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SettingRepository defines the interface for settings data access
type SettingRepository interface {
	// Create creates a new setting
	Create(ctx context.Context, setting *entities.Setting) error

	// GetByID retrieves a setting by its ID
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Setting, error)

	// GetByKey retrieves a system setting by its key
	GetByKey(ctx context.Context, key string) (*entities.Setting, error)

	// GetUserSettingByKey retrieves a user setting by its key and user ID
	GetUserSettingByKey(ctx context.Context, key string, userID primitive.ObjectID) (*entities.Setting, error)

	// Update updates an existing setting
	Update(ctx context.Context, setting *entities.Setting) error

	// Delete deletes a setting by its ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// ListSystemSettings retrieves all system settings with optional filtering
	ListSystemSettings(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Setting, int, error)

	// ListUserSettings retrieves all settings for a specific user
	ListUserSettings(ctx context.Context, userID primitive.ObjectID, filter map[string]interface{}, page, limit int) ([]*entities.Setting, int, error)
}
