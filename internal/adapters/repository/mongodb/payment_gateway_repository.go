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

type PaymentGatewayRepository struct {
	collection *mongo.Collection
}

func NewPaymentGatewayRepository(db *mongo.Database) *PaymentGatewayRepository {
	return &PaymentGatewayRepository{
		collection: db.Collection("payment_gateways"),
	}
}

func (r *PaymentGatewayRepository) Create(ctx context.Context, gateway *entities.PaymentGateway) error {
	_, err := r.collection.InsertOne(ctx, gateway)
	return err
}

func (r *PaymentGatewayRepository) GetByID(ctx context.Context, id string) (*entities.PaymentGateway, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid payment gateway ID")
	}

	var gateway entities.PaymentGateway
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&gateway)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("payment gateway not found")
		}
		return nil, err
	}
	return &gateway, nil
}

func (r *PaymentGatewayRepository) Update(ctx context.Context, gateway *entities.PaymentGateway) error {
	gateway.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": gateway.ID}, gateway)
	return err
}

func (r *PaymentGatewayRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid payment gateway ID")
	}
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("payment gateway not found")
	}
	return nil
}

func (r *PaymentGatewayRepository) List(ctx context.Context, active bool) ([]*entities.PaymentGateway, error) {
	filter := bson.M{}
	if active {
		filter["active"] = true
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var gateways []*entities.PaymentGateway
	if err := cursor.All(ctx, &gateways); err != nil {
		return nil, err
	}
	return gateways, nil
}

func (r *PaymentGatewayRepository) GetByProvider(ctx context.Context, provider string) (*entities.PaymentGateway, error) {
	var gateway entities.PaymentGateway
	err := r.collection.FindOne(ctx, bson.M{"provider": provider}).Decode(&gateway)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("payment gateway not found")
		}
		return nil, err
	}
	return &gateway, nil
}
