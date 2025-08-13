package mongodb

import (
	"context"
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SettingRepository implements the repository interface for settings
type SettingRepository struct {
	collection *mongo.Collection
}

// NewSettingRepository creates a new setting repository
func NewSettingRepository(db *mongo.Database) *SettingRepository {
	collection := db.Collection(SettingCollection)

	// Create indexes for faster lookups
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key", Value: 1}, {Key: "scope", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "key", Value: 1}, {Key: "user_id", Value: 1}, {Key: "scope", Value: 1}},
			Options: options.Index().SetUnique(true).SetPartialFilterExpression(
				bson.M{"scope": entities.SettingScopeUser, "user_id": bson.M{"$exists": true}},
			),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexModels)
	if err != nil {
		// Log the error but continue
		// log.Printf("Error creating indexes for settings collection: %v", err)
	}

	return &SettingRepository{
		collection: collection,
	}
}

// Create creates a new setting
func (r *SettingRepository) Create(ctx context.Context, setting *entities.Setting) error {
	_, err := r.collection.InsertOne(ctx, setting)
	return err
}

// GetByID retrieves a setting by its ID
func (r *SettingRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Setting, error) {
	var setting entities.Setting
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&setting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("setting not found")
		}
		return nil, err
	}
	return &setting, nil
}

// GetByKey retrieves a system setting by its key
func (r *SettingRepository) GetByKey(ctx context.Context, key string) (*entities.Setting, error) {
	var setting entities.Setting
	err := r.collection.FindOne(ctx, bson.M{
		"key":   key,
		"scope": entities.SettingScopeSystem,
	}).Decode(&setting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("setting not found")
		}
		return nil, err
	}
	return &setting, nil
}

// GetUserSettingByKey retrieves a user setting by its key and user ID
func (r *SettingRepository) GetUserSettingByKey(ctx context.Context, key string, userID primitive.ObjectID) (*entities.Setting, error) {
	var setting entities.Setting
	err := r.collection.FindOne(ctx, bson.M{
		"key":     key,
		"user_id": userID,
		"scope":   entities.SettingScopeUser,
	}).Decode(&setting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("setting not found")
		}
		return nil, err
	}
	return &setting, nil
}

// Update updates an existing setting
func (r *SettingRepository) Update(ctx context.Context, setting *entities.Setting) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": setting.ID}, setting)
	return err
}

// Delete deletes a setting by its ID
func (r *SettingRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ListSystemSettings retrieves all system settings with optional filtering
func (r *SettingRepository) ListSystemSettings(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entities.Setting, int, error) {
	// Add scope filter
	filterBson := bson.M{"scope": entities.SettingScopeSystem}

	// Add additional filters if provided
	for k, v := range filter {
		filterBson[k] = v
	}

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Execute query
	cursor, err := r.collection.Find(ctx, filterBson, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var settings []*entities.Setting
	if err = cursor.All(ctx, &settings); err != nil {
		return nil, 0, err
	}

	// Count total documents for pagination
	total, err := r.collection.CountDocuments(ctx, filterBson)
	if err != nil {
		return nil, 0, err
	}

	return settings, int(total), nil
}

// ListUserSettings retrieves all settings for a specific user
func (r *SettingRepository) ListUserSettings(ctx context.Context, userID primitive.ObjectID, filter map[string]interface{}, page, limit int) ([]*entities.Setting, int, error) {
	// Add scope and user ID filters
	filterBson := bson.M{
		"scope":   entities.SettingScopeUser,
		"user_id": userID,
	}

	// Add additional filters if provided
	for k, v := range filter {
		filterBson[k] = v
	}

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Execute query
	cursor, err := r.collection.Find(ctx, filterBson, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var settings []*entities.Setting
	if err = cursor.All(ctx, &settings); err != nil {
		return nil, 0, err
	}

	// Count total documents for pagination
	total, err := r.collection.CountDocuments(ctx, filterBson)
	if err != nil {
		return nil, 0, err
	}

	return settings, int(total), nil
}
