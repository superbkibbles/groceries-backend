package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository implements the user repository interface using MongoDB
type UserRepository struct {
	db             *mongo.Database
	userCollection *mongo.Collection
	addrCollection *mongo.Collection
	redisClient    *redis.Client
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *mongo.Database, redisClient *redis.Client) *UserRepository {
	return &UserRepository{
		db:             db,
		userCollection: db.Collection(UserCollection),
		addrCollection: db.Collection(AddressCollection),
		redisClient:    redisClient,
	}
}

// Create adds a new user to the database
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	// Check if email already exists
	count, err := r.userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email already exists")
	}

	_, err = r.userCollection.InsertOne(ctx, user)
	return err
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	err := r.userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	_, err := r.userCollection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	return err
}

// Delete removes a user from the database
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// Delete user
	_, err := r.userCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	// Delete all addresses for this user
	_, err = r.addrCollection.DeleteMany(ctx, bson.M{"user_id": id})
	return err
}

// List retrieves users based on filters with pagination
func (r *UserRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	// Convert map to bson.M
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}

	// Get total count
	total, err := r.userCollection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	// Find users with pagination
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := r.userCollection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*entities.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

// AddAddress adds a new address for a user
func (r *UserRepository) AddAddress(ctx context.Context, address *entities.Address) error {
	// If this is the default address, unset any existing default
	if address.IsDefault {
		_, err := r.addrCollection.UpdateMany(
			ctx,
			bson.M{"user_id": address.UserID, "is_default": true},
			bson.M{"$set": bson.M{"is_default": false}},
		)
		if err != nil {
			return err
		}
	}

	_, err := r.addrCollection.InsertOne(ctx, address)
	return err
}

// UpdateAddress updates an existing address
func (r *UserRepository) UpdateAddress(ctx context.Context, address *entities.Address) error {
	// If this is being set as default, unset any existing default
	if address.IsDefault {
		_, err := r.addrCollection.UpdateMany(
			ctx,
			bson.M{"user_id": address.UserID, "_id": bson.M{"$ne": address.ID}, "is_default": true},
			bson.M{"$set": bson.M{"is_default": false}},
		)
		if err != nil {
			return err
		}
	}

	_, err := r.addrCollection.ReplaceOne(ctx, bson.M{"id": address.ID}, address)
	return err
}

// DeleteAddress removes an address
func (r *UserRepository) DeleteAddress(ctx context.Context, id string) error {
	// Check if this is a default address
	var address entities.Address
	err := r.addrCollection.FindOne(ctx, bson.M{"id": id}).Decode(&address)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("address not found")
		}
		return err
	}

	// Delete the address
	_, err = r.addrCollection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

// GetAddressesByUserID retrieves all addresses for a user
func (r *UserRepository) GetAddressesByUserID(ctx context.Context, userID string) ([]*entities.Address, error) {
	cursor, err := r.addrCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var addresses []*entities.Address
	if err = cursor.All(ctx, &addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *UserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error) {
	var user entities.User
	err := r.userCollection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SaveOTP(ctx context.Context, phoneNumber string, otp string) error {
	// Save OTP to Redis with a 5-minute expiration
	err := r.redisClient.Set(ctx, phoneNumber, otp, 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
