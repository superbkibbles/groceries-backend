package ports

import (
	"context"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
)

// SettingService defines the interface for settings business logic
type SettingService interface {
	// CreateSystemSetting creates a new system-wide setting
	CreateSystemSetting(ctx context.Context, key string, value interface{}, settingType entities.SettingType, description string) (*entities.Setting, error)

	// CreateUserSetting creates a new user-specific setting
	CreateUserSetting(ctx context.Context, key string, value interface{}, settingType entities.SettingType, description string, userID string) (*entities.Setting, error)

	// GetSystemSetting retrieves a system setting by its key
	GetSystemSetting(ctx context.Context, key string) (*entities.Setting, error)

	// GetUserSetting retrieves a user setting by its key and user ID
	GetUserSetting(ctx context.Context, key string, userID string) (*entities.Setting, error)

	// UpdateSetting updates an existing setting
	UpdateSetting(ctx context.Context, setting *entities.Setting) error

	// DeleteSetting deletes a setting by its ID
	DeleteSetting(ctx context.Context, id string) error

	// ListSystemSettings retrieves all system settings with optional filtering
	ListSystemSettings(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Setting, int, error)

	// ListUserSettings retrieves all settings for a specific user
	ListUserSettings(ctx context.Context, userID string, filter map[string]interface{}, page, limit int) ([]*entities.Setting, int, error)
}
