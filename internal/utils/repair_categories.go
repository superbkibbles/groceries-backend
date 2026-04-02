package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/superbkibbles/ecommerce/internal/adapters/repository/mongodb"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// seedProductCategorySlugsBySKU mirrors internal/utils/seed.go product→category slug links.
func seedProductCategorySlugsBySKU() map[string][]string {
	return map[string][]string{
		"IP13-GRA-128":     {"electronics", "smartphones"},
		"MBP16-PRO-16-512": {"electronics", "laptops"},
		"TS-BL-M":          {"clothing", "mens-clothing"},
		"BL-BL-600":        {"home-kitchen", "kitchen-appliances"},
		"ORG-BAN-1KG":      {"groceries", "fresh-produce"},
		"MILK-WH-1L":       {"groceries", "dairy-eggs"},
		"BRD-SOUR-1":       {"groceries", "bakery"},
		"RICE-BAS-5KG":     {"groceries", "pantry"},
		"OJ-FRESH-1L":      {"groceries", "beverages"},
	}
}

// RepairCategoryBSONKeysAndRefs fixes categories stored with legacy camel-ish BSON keys from
// struct fields without bson tags, then repairs product.categories and home_sections.category_ids
// when they contain only zero ObjectIDs (bad seed in-memory IDs).
func RepairCategoryBSONKeysAndRefs(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(mongodb.CategoryCollection)

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("find categories: %w", err)
	}
	defer cursor.Close(ctx)

	var normalized int
	for cursor.Next(ctx) {
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			return fmt.Errorf("decode category: %w", err)
		}
		idVal, ok := raw["_id"]
		if !ok {
			log.Printf("repair: skip category document without _id: %+v", raw)
			continue
		}

		clean := bson.M{"_id": idVal}

		if v, ok := raw["name"]; ok {
			clean["name"] = v
		}
		if v, ok := raw["description"]; ok {
			clean["description"] = v
		}
		if v, ok := raw["slug"]; ok {
			clean["slug"] = v
		}
		if v, ok := raw["level"]; ok {
			clean["level"] = v
		}
		if v, ok := raw["path"]; ok {
			clean["path"] = v
		}
		if v, ok := raw["translations"]; ok {
			clean["translations"] = v
		}

		switch {
		case raw["parent_id"] != nil:
			clean["parent_id"] = raw["parent_id"]
		case raw["parentid"] != nil:
			clean["parent_id"] = raw["parentid"]
		default:
			clean["parent_id"] = primitive.NilObjectID
		}

		switch {
		case raw["created_at"] != nil:
			clean["created_at"] = raw["created_at"]
		case raw["createdat"] != nil:
			clean["created_at"] = raw["createdat"]
		}
		switch {
		case raw["updated_at"] != nil:
			clean["updated_at"] = raw["updated_at"]
		case raw["updatedat"] != nil:
			clean["updated_at"] = raw["updatedat"]
		}

		if _, err := coll.ReplaceOne(ctx, bson.M{"_id": idVal}, clean); err != nil {
			return fmt.Errorf("replace category %v: %w", idVal, err)
		}
		normalized++
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	log.Printf("repair: normalized %d category documents", normalized)

	slugToID := make(map[string]primitive.ObjectID)
	cur2, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cur2.Close(ctx)
	for cur2.Next(ctx) {
		var c entities.Category
		if err := cur2.Decode(&c); err != nil {
			return fmt.Errorf("decode category after repair: %w", err)
		}
		if c.Slug != "" {
			slugToID[c.Slug] = c.ID
		}
	}
	if err := cur2.Err(); err != nil {
		return err
	}

	prodColl := db.Collection(mongodb.ProductCollection)
	skuMap := seedProductCategorySlugsBySKU()
	pcur, err := prodColl.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer pcur.Close(ctx)
	var productsUpdated int
	for pcur.Next(ctx) {
		var p entities.Product
		if err := pcur.Decode(&p); err != nil {
			return fmt.Errorf("decode product: %w", err)
		}
		if !productCategoryRefsNeedRepair(p.Categories) {
			continue
		}
		slugs, ok := skuMap[p.SKU]
		if !ok {
			continue
		}
		var ids []primitive.ObjectID
		for _, slug := range slugs {
			oid, ok := slugToID[slug]
			if !ok {
				return fmt.Errorf("unknown category slug %q for product SKU %s", slug, p.SKU)
			}
			ids = append(ids, oid)
		}
		if _, err := prodColl.UpdateOne(ctx, bson.M{"_id": p.ID}, bson.M{"$set": bson.M{"categories": ids}}); err != nil {
			return fmt.Errorf("update product %s: %w", p.SKU, err)
		}
		productsUpdated++
	}
	if err := pcur.Err(); err != nil {
		return err
	}
	log.Printf("repair: updated categories on %d products", productsUpdated)

	homeOrder := []string{"fresh-produce", "dairy-eggs", "bakery", "pantry", "beverages"}
	var homeCatIDs []primitive.ObjectID
	for _, slug := range homeOrder {
		oid, ok := slugToID[slug]
		if !ok {
			return fmt.Errorf("missing slug %q for home section repair", slug)
		}
		homeCatIDs = append(homeCatIDs, oid)
	}

	hsColl := db.Collection(mongodb.HomeSectionCollection)
	hcur, err := hsColl.Find(ctx, bson.M{"type": string(entities.HomeSectionTypeCategories)})
	if err != nil {
		return err
	}
	defer hcur.Close(ctx)
	var sectionsUpdated int
	for hcur.Next(ctx) {
		var s entities.HomeSection
		if err := hcur.Decode(&s); err != nil {
			return fmt.Errorf("decode home section: %w", err)
		}
		if !homeSectionCategoryRefsNeedRepair(s.CategoryIDs) {
			continue
		}
		if _, err := hsColl.UpdateOne(ctx, bson.M{"_id": s.ID}, bson.M{"$set": bson.M{"category_ids": homeCatIDs}}); err != nil {
			return fmt.Errorf("update home section %s: %w", s.ID.Hex(), err)
		}
		sectionsUpdated++
	}
	if err := hcur.Err(); err != nil {
		return err
	}
	log.Printf("repair: updated category_ids on %d home_sections", sectionsUpdated)

	return nil
}

func productCategoryRefsNeedRepair(cats []primitive.ObjectID) bool {
	if len(cats) == 0 {
		return true
	}
	for _, id := range cats {
		if id.IsZero() {
			return true
		}
	}
	return false
}

func homeSectionCategoryRefsNeedRepair(ids []primitive.ObjectID) bool {
	return productCategoryRefsNeedRepair(ids)
}
