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

type PaymentMethodRepository struct {
	collection *mongo.Collection
}

func NewPaymentMethodRepository(db *mongo.Database) *PaymentMethodRepository {
	return &PaymentMethodRepository{
		collection: db.Collection("payment_methods"),
	}
}

func (r *PaymentMethodRepository) Create(ctx context.Context, method *entities.PaymentMethod) error {
	_, err := r.collection.InsertOne(ctx, method)
	return err
}

func (r *PaymentMethodRepository) GetByID(ctx context.Context, id string) (*entities.PaymentMethod, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid payment method ID")
	}

	var method entities.PaymentMethod
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&method)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("payment method not found")
		}
		return nil, err
	}
	return &method, nil
}

func (r *PaymentMethodRepository) Update(ctx context.Context, method *entities.PaymentMethod) error {
	method.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": method.ID}, method)
	return err
}

func (r *PaymentMethodRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid payment method ID")
	}
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("payment method not found")
	}
	return nil
}

func (r *PaymentMethodRepository) List(ctx context.Context, active bool) ([]*entities.PaymentMethod, error) {
	filter := bson.M{}
	if active {
		filter["active"] = true
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var methods []*entities.PaymentMethod
	if err := cursor.All(ctx, &methods); err != nil {
		return nil, err
	}
	return methods, nil
}

func (r *PaymentMethodRepository) GetByType(ctx context.Context, methodType entities.PaymentMethodType) ([]*entities.PaymentMethod, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"type": methodType})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var methods []*entities.PaymentMethod
	if err := cursor.All(ctx, &methods); err != nil {
		return nil, err
	}
	return methods, nil
}
