package services

import (
	"context"
	"errors"
	"time"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
)

// SettingService implements the setting service interface
type SettingService struct {
	settingRepo ports.SettingRepository
}

// NewSettingService creates a new setting service
func NewSettingService(settingRepo ports.SettingRepository) *SettingService {
	return &SettingService{
		settingRepo: settingRepo,
	}
}

// CreateSystemSetting creates a new system-wide setting
func (s *SettingService) CreateSystemSetting(ctx context.Context, key string, value interface{}, settingType entities.SettingType, description string) (*entities.Setting, error) {
	// Check if setting already exists
	existing, err := s.settingRepo.GetByKey(ctx, key)
	if err == nil && existing != nil {
		return nil, errors.New("setting already exists")
	}

	// Create new setting
	setting := entities.NewSystemSetting(key, value, settingType, description)

	// Validate setting
	if err := setting.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.settingRepo.Create(ctx, setting); err != nil {
		return nil, err
	}

	return setting, nil
}

// CreateUserSetting creates a new user-specific setting
func (s *SettingService) CreateUserSetting(ctx context.Context, key string, value interface{}, settingType entities.SettingType, description string, userID string) (*entities.Setting, error) {
	// Check if setting already exists for this user
	existing, err := s.settingRepo.GetUserSettingByKey(ctx, key, userID)
	if err == nil && existing != nil {
		return nil, errors.New("user setting already exists")
	}

	// Create new setting
	setting := entities.NewUserSetting(key, value, settingType, description, userID)

	// Validate setting
	if err := setting.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.settingRepo.Create(ctx, setting); err != nil {
		return nil, err
	}

	return setting, nil
}

// GetSystemSetting retrieves a system setting by its key
func (s *SettingService) GetSystemSetting(ctx context.Context, key string) (*entities.Setting, error) {
	return s.settingRepo.GetByKey(ctx, key)
}

// GetUserSetting retrieves a user setting by its key and user ID
func (s *SettingService) GetUserSetting(ctx context.Context, key string, userID string) (*entities.Setting, error) {
	return s.settingRepo.GetUserSettingByKey(ctx, key, userID)
}

// UpdateSetting updates an existing setting
func (s *SettingService) UpdateSetting(ctx context.Context, setting *entities.Setting) error {
	// Get existing setting
	existing, err := s.settingRepo.GetByID(ctx, setting.ID)
	if err != nil {
		return err
	}

	// Update fields
	existing.Value = setting.Value
	existing.Description = setting.Description
	existing.UpdatedAt = time.Now()

	// Validate setting
	if err := existing.Validate(); err != nil {
		return err
	}

	// Save to repository
	return s.settingRepo.Update(ctx, existing)
}

// DeleteSetting deletes a setting by its ID
func (s *SettingService) DeleteSetting(ctx context.Context, id string) error {
	return s.settingRepo.Delete(ctx, id)
}

// ListSystemSettings retrieves all system settings with optional filtering
func (s *SettingService) ListSystemSettings(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Setting, int, error) {
	return s.settingRepo.ListSystemSettings(ctx, filter, page, limit)
}

// ListUserSettings retrieves all settings for a specific user
func (s *SettingService) ListUserSettings(ctx context.Context, userID string, filter map[string]interface{}, page, limit int) ([]*entities.Setting, int, error) {
	return s.settingRepo.ListUserSettings(ctx, userID, filter, page, limit)
}
