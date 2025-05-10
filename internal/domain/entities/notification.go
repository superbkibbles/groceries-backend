package entities

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeOrderStatus NotificationType = "order_status"
	NotificationTypePayment     NotificationType = "payment"
	NotificationTypeShipping    NotificationType = "shipping"
	NotificationTypeSystem      NotificationType = "system"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusUnread   NotificationStatus = "unread"
	NotificationStatusRead     NotificationStatus = "read"
	NotificationStatusArchived NotificationStatus = "archived"
)

// Notification represents a notification in the system
type Notification struct {
	ID        string                 `json:"id" bson:"_id"`
	UserID    string                 `json:"user_id" bson:"user_id"`
	Type      NotificationType       `json:"type" bson:"type"`
	Title     string                 `json:"title" bson:"title"`
	Message   string                 `json:"message" bson:"message"`
	Status    NotificationStatus     `json:"status" bson:"status"`
	Data      map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" bson:"updated_at"`
	ReadAt    *time.Time             `json:"read_at,omitempty" bson:"read_at,omitempty"`
}

// NotificationTemplate represents a template for generating notifications
type NotificationTemplate struct {
	ID              string           `json:"id" bson:"_id"`
	Name            string           `json:"name" bson:"name"`
	Type            NotificationType `json:"type" bson:"type"`
	TitleTemplate   string           `json:"title_template" bson:"title_template"`
	MessageTemplate string           `json:"message_template" bson:"message_template"`
	Active          bool             `json:"active" bson:"active"`
	CreatedAt       time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" bson:"updated_at"`
}

// NewNotification creates a new notification
func NewNotification(userID string, notificationType NotificationType, title, message string, data map[string]interface{}) *Notification {
	now := time.Now()
	return &Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      notificationType,
		Title:     title,
		Message:   message,
		Status:    NotificationStatusUnread,
		Data:      data,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// MarkAsRead marks a notification as read
func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.Status = NotificationStatusRead
	n.ReadAt = &now
	n.UpdatedAt = now
}

// MarkAsArchived marks a notification as archived
func (n *Notification) MarkAsArchived() {
	n.Status = NotificationStatusArchived
	n.UpdatedAt = time.Now()
}

// NewNotificationTemplate creates a new notification template
func NewNotificationTemplate(name string, notificationType NotificationType, titleTemplate, messageTemplate string) *NotificationTemplate {
	now := time.Now()
	return &NotificationTemplate{
		ID:              uuid.New().String(),
		Name:            name,
		Type:            notificationType,
		TitleTemplate:   titleTemplate,
		MessageTemplate: messageTemplate,
		Active:          true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
