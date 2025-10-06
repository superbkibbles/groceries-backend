package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SettingScope defines the scope of a setting
type SettingScope string

const (
	// SettingScopeSystem represents system-wide settings
	SettingScopeSystem SettingScope = "system"
	// SettingScopeUser represents user-specific settings
	SettingScopeUser SettingScope = "user"
)

// SettingType defines the data type of a setting value
type SettingType string

const (
	// SettingTypeString represents string values
	SettingTypeString SettingType = "string"
	// SettingTypeNumber represents numeric values
	SettingTypeNumber SettingType = "number"
	// SettingTypeBoolean represents boolean values
	SettingTypeBoolean SettingType = "boolean"
	// SettingTypeJSON represents JSON object values
	SettingTypeJSON SettingType = "json"
)

// Setting represents a configuration setting in the e-commerce system
type Setting struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	Key         string             `json:"key"`
	Value       interface{}        `json:"value"`
	Type        SettingType        `json:"type"`
	Scope       SettingScope       `json:"scope"`
	UserID      primitive.ObjectID `json:"user_id,omitempty"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// NewSetting creates a new setting with a unique ID
func NewSetting(key string, value interface{}, settingType SettingType, scope SettingScope, description string, userID primitive.ObjectID) *Setting {
	now := time.Now()
	return &Setting{
		Key:         key,
		Value:       value,
		Type:        settingType,
		Scope:       scope,
		UserID:      userID,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewSystemSetting creates a new system-wide setting
func NewSystemSetting(key string, value interface{}, settingType SettingType, description string) *Setting {
	return NewSetting(key, value, settingType, SettingScopeSystem, description, primitive.NilObjectID)
}

// NewUserSetting creates a new user-specific setting
func NewUserSetting(key string, value interface{}, settingType SettingType, description string, userID primitive.ObjectID) *Setting {
	return NewSetting(key, value, settingType, SettingScopeUser, description, userID)
}

// Validate checks if the setting is valid
func (s *Setting) Validate() error {
	// Add validation logic here
	return nil
}
