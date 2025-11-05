package mongodb

import (
    "context"

    "github.com/superbkibbles/ecommerce/internal/domain/entities"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const HomeSectionCollection = "home_sections"

type HomeSectionRepository struct {
    collection *mongo.Collection
}

func NewHomeSectionRepository(db *mongo.Database) *HomeSectionRepository {
    return &HomeSectionRepository{collection: db.Collection(HomeSectionCollection)}
}

func (r *HomeSectionRepository) Create(ctx context.Context, section *entities.HomeSection) error {
    _, err := r.collection.InsertOne(ctx, section)
    return err
}

func (r *HomeSectionRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.HomeSection, error) {
    var section entities.HomeSection
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&section)
    if err != nil {
        return nil, err
    }
    return &section, nil
}

func (r *HomeSectionRepository) Update(ctx context.Context, section *entities.HomeSection) error {
    _, err := r.collection.ReplaceOne(ctx, bson.M{"_id": section.ID}, section)
    return err
}

func (r *HomeSectionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
    _, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

func (r *HomeSectionRepository) List(ctx context.Context, onlyActive bool) ([]*entities.HomeSection, error) {
    filter := bson.M{}
    if onlyActive {
        filter["active"] = true
    }
    opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}})
    cursor, err := r.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var sections []*entities.HomeSection
    if err := cursor.All(ctx, &sections); err != nil {
        return nil, err
    }
    return sections, nil
}


