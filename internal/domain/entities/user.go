package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleCustomer UserRole = "customer"
	UserRoleAdmin    UserRole = "admin"
)

// User represents a user in the e-commerce system
type User struct {
	ID           string    `json:"id" bson:"_id"`
	Email        string    `json:"email" bson:"email"`
	PasswordHash string    `json:"password_hash,omitempty" bson:"password_hash"`
	FirstName    string    `json:"first_name" bson:"first_name"`
	LastName     string    `json:"last_name" bson:"last_name"`
	PhoneNumber  string    `json:"phone_number" bson:"phone_number"`
	Role         UserRole  `json:"role" bson:"role"`
	Active       bool      `json:"active" bson:"active"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

// Address represents a user's address
type Address struct {
	ID           string    `json:"id" bson:"id"`
	UserID       string    `json:"user_id" bson:"user_id"`
	Name         string    `json:"name" bson:"name"`
	AddressLine1 string    `json:"address_line1" bson:"address_line1"`
	AddressLine2 string    `json:"address_line2,omitempty" bson:"address_line2,omitempty"`
	City         string    `json:"city" bson:"city"`
	State        string    `json:"state" bson:"state"`
	Country      string    `json:"country" bson:"country"`
	PostalCode   string    `json:"postal_code" bson:"postal_code"`
	Phone        string    `json:"phone" bson:"phone"`
	IsDefault    bool      `json:"is_default" bson:"is_default"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

// NewUser creates a new user with a unique ID
func NewUser(email, password, firstName, lastName string, role UserRole) (*User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		Role:         role,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// ValidatePassword checks if the provided password matches the stored hash
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(currentPassword, newPassword string) error {
	if !u.ValidatePassword(currentPassword) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

// NewAddress creates a new address for a user
func NewAddress(userID, name, addressLine1, addressLine2, city, state, country, postalCode, phone string, isDefault bool) *Address {
	now := time.Now()
	return &Address{
		ID:           uuid.New().String(),
		UserID:       userID,
		Name:         name,
		AddressLine1: addressLine1,
		AddressLine2: addressLine2,
		City:         city,
		State:        state,
		Country:      country,
		PostalCode:   postalCode,
		Phone:        phone,
		IsDefault:    isDefault,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
