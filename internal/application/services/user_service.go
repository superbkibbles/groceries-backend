package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"github.com/superbkibbles/ecommerce/internal/domain/ports"
	"github.com/superbkibbles/ecommerce/internal/utils"
)

// JWT secret key - in production, this should be loaded from environment variables
const jwtSecret = "SJSDH#$!!^&#dsds9%^!sajh"

// UserService implements the user service interface
type UserService struct {
	userRepo ports.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo ports.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Register creates a new user account
func (s *UserService) Register(ctx context.Context, email, password, firstName, lastName string) (*entities.User, error) {
	// Validate input
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Check if email already exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	// Create new user
	user, err := entities.NewUser(email, password, firstName, lastName, entities.UserRoleCustomer)
	if err != nil {
		return nil, err
	}

	// Save user to database
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(ctx context.Context, phoneNumber string) (*entities.User, string, error) {
	// validate input and check if phone number is valid
	if phoneNumber == "" {
		return nil, "", errors.New("phone number is required")
	}
	if len(phoneNumber) < 10 {
		return nil, "", errors.New("phone number is invalid")
	}

	// generate otp of 6 digits
	otp, err := utils.GenerateOTP(6)
	if err != nil {
		return nil, "", err
	}

	s.userRepo.SaveOTP(ctx, phoneNumber, otp) // save otp to redis
	// send otp to the phone number

	// validate otp

	// Get user by phone number
	user, err := s.userRepo.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// check if user does not exists
	if user == nil {
		// create new user with phone number only
	}

	// send OTP to user

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, user *entities.User) error {
	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}

	return s.userRepo.Update(ctx, user)
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Update password
	err = user.UpdatePassword(currentPassword, newPassword)
	if err != nil {
		return err
	}

	// Save user to database
	return s.userRepo.Update(ctx, user)
}

// AddAddress adds a new address for a user
func (s *UserService) AddAddress(ctx context.Context, userID, name, addressLine1, addressLine2, city, state, country, postalCode, phone string, isDefault bool) (*entities.Address, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create new address
	address := entities.NewAddress(userID, name, addressLine1, addressLine2, city, state, country, postalCode, phone, isDefault)

	// Save address to database
	err = s.userRepo.AddAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

// UpdateAddress updates an existing address
func (s *UserService) UpdateAddress(ctx context.Context, address *entities.Address) error {
	// Verify address exists
	addresses, err := s.userRepo.GetAddressesByUserID(ctx, address.UserID)
	if err != nil {
		return err
	}

	found := false
	for _, a := range addresses {
		if a.ID == address.ID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("address not found")
	}

	return s.userRepo.UpdateAddress(ctx, address)
}

// DeleteAddress removes an address
func (s *UserService) DeleteAddress(ctx context.Context, addressID string) error {
	return s.userRepo.DeleteAddress(ctx, addressID)
}

// GetAddresses retrieves all addresses for a user
func (s *UserService) GetAddresses(ctx context.Context, userID string) ([]*entities.Address, error) {
	return s.userRepo.GetAddressesByUserID(ctx, userID)
}

// SetDefaultAddress sets an address as the default for a user
func (s *UserService) SetDefaultAddress(ctx context.Context, userID, addressID string) error {
	// Get all addresses for the user
	addresses, err := s.userRepo.GetAddressesByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Find the address to set as default
	var targetAddress *entities.Address
	for _, a := range addresses {
		if a.ID == addressID {
			targetAddress = a
			break
		}
	}

	if targetAddress == nil {
		return errors.New("address not found")
	}

	// Set as default
	targetAddress.IsDefault = true
	targetAddress.UpdatedAt = time.Now()

	return s.userRepo.UpdateAddress(ctx, targetAddress)
}

// ListUsers retrieves users based on filters with pagination
func (s *UserService) ListUsers(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.User, int, error) {
	return s.userRepo.List(ctx, filter, page, limit)
}

// Helper function to generate JWT token
func generateJWT(user *entities.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
