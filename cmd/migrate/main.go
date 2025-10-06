package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
)

func main() {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("durra") // Replace with your database name

	fmt.Println("Starting migration to add translation support...")

	// Migrate products
	if err := migrateProducts(ctx, db); err != nil {
		log.Fatal("Failed to migrate products:", err)
	}

	// Migrate categories
	if err := migrateCategories(ctx, db); err != nil {
		log.Fatal("Failed to migrate categories:", err)
	}

	// Create indexes
	if err := createIndexes(ctx, db); err != nil {
		log.Fatal("Failed to create indexes:", err)
	}

	fmt.Println("Migration completed successfully!")
}

func migrateProducts(ctx context.Context, db *mongo.Database) error {
	fmt.Println("Migrating products...")

	collection := db.Collection("products")

	// Find all products that don't have translations field
	cursor, err := collection.Find(ctx, bson.M{"translations": bson.M{"$exists": false}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var productCount int
	for cursor.Next(ctx) {
		var product bson.M
		if err := cursor.Decode(&product); err != nil {
			fmt.Printf("Error decoding product: %v\n", err)
			continue
		}

		// Create translations from existing name and description
		translations := map[string]entities.Translation{
			"en": {
				Name:        getString(product, "name"),
				Description: getString(product, "description"),
			},
		}

		// Update the product
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": product["_id"]},
			bson.M{
				"$set": bson.M{
					"translations": translations,
					"updated_at":   time.Now(),
				},
			},
		)
		if err != nil {
			fmt.Printf("Error updating product %v: %v\n", product["_id"], err)
			continue
		}

		productCount++
		fmt.Printf("Migrated product: %s (ID: %v)\n", getString(product, "name"), product["_id"])
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	fmt.Printf("Products migration completed. Migrated %d products.\n", productCount)
	return nil
}

func migrateCategories(ctx context.Context, db *mongo.Database) error {
	fmt.Println("Migrating categories...")

	collection := db.Collection("categories")

	// Find all categories that don't have translations field
	cursor, err := collection.Find(ctx, bson.M{"translations": bson.M{"$exists": false}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var categoryCount int
	for cursor.Next(ctx) {
		var category bson.M
		if err := cursor.Decode(&category); err != nil {
			fmt.Printf("Error decoding category: %v\n", err)
			continue
		}

		// Create translations from existing name and description
		translations := map[string]entities.Translation{
			"en": {
				Name:        getString(category, "name"),
				Description: getString(category, "description"),
			},
		}

		// Update the category
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": category["_id"]},
			bson.M{
				"$set": bson.M{
					"translations": translations,
					"updated_at":   time.Now(),
				},
			},
		)
		if err != nil {
			fmt.Printf("Error updating category %v: %v\n", category["_id"], err)
			continue
		}

		categoryCount++
		fmt.Printf("Migrated category: %s (ID: %v)\n", getString(category, "name"), category["_id"])
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	fmt.Printf("Categories migration completed. Migrated %d categories.\n", categoryCount)
	return nil
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	fmt.Println("Creating indexes...")

	// Create index on products.translations
	productsCollection := db.Collection("products")
	_, err := productsCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{{Key: "translations", Value: 1}},
		},
	)
	if err != nil {
		fmt.Printf("Error creating index on products.translations: %v\n", err)
	} else {
		fmt.Println("Created index on products.translations")
	}

	// Create index on categories.translations
	categoriesCollection := db.Collection("categories")
	_, err = categoriesCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{{Key: "translations", Value: 1}},
		},
	)
	if err != nil {
		fmt.Printf("Error creating index on categories.translations: %v\n", err)
	} else {
		fmt.Println("Created index on categories.translations")
	}

	return nil
}

func getString(doc bson.M, key string) string {
	if val, ok := doc[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
