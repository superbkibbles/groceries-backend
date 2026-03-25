package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerPaymentMethodRepository struct {
	collection *mongo.Collection
}

func NewCustomerPaymentMethodRepository(db *mongo.Database) *CustomerPaymentMethodRepository {
	return &CustomerPaymentMethodRepository{
		collection: db.Collection("customer_payment_methods"),
	}
}

func (r *CustomerPaymentMethodRepository) Create(ctx context.Context, method *entities.CustomerPaymentMethod) error {
	_, err := r.collection.InsertOne(ctx, method)
	return err
}

func (r *CustomerPaymentMethodRepository) GetByID(ctx context.Context, id string) (*entities.CustomerPaymentMethod, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid customer payment method ID")
	}

	var method entities.CustomerPaymentMethod
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&method)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("customer payment method not found")
		}
		return nil, err
	}
	return &method, nil
}

func (r *CustomerPaymentMethodRepository) Update(ctx context.Context, method *entities.CustomerPaymentMethod) error {
	method.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": method.ID}, method)
	return err
}

func (r *CustomerPaymentMethodRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid customer payment method ID")
	}
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("customer payment method not found")
	}
	return nil
}

func (r *CustomerPaymentMethodRepository) GetByCustomer(ctx context.Context, customerID string) ([]*entities.CustomerPaymentMethod, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"customer_id": customerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var methods []*entities.CustomerPaymentMethod
	if err := cursor.All(ctx, &methods); err != nil {
		return nil, err
	}
	return methods, nil
}

func (r *CustomerPaymentMethodRepository) GetDefaultByCustomer(ctx context.Context, customerID string) (*entities.CustomerPaymentMethod, error) {
	var method entities.CustomerPaymentMethod
	err := r.collection.FindOne(ctx, bson.M{"customer_id": customerID, "is_default": true}).Decode(&method)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no default payment method found")
		}
		return nil, err
	}
	return &method, nil
}

func (r *CustomerPaymentMethodRepository) SetDefault(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid customer payment method ID")
	}
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"is_default": true, "updated_at": time.Now()}},
	)
	return err
}

func (r *CustomerPaymentMethodRepository) ClearDefault(ctx context.Context, customerID string) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"customer_id": customerID, "is_default": true},
		bson.M{"$set": bson.M{"is_default": false, "updated_at": time.Now()}},
	)
	return err
}
