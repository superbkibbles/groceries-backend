package utils

import (
	"context"

	"github.com/superbkibbles/ecommerce/internal/adapters/repository/mongodb"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// seedHomeSections inserts default home page sections using seeded product and category IDs.
func seedHomeSections(
	ctx context.Context,
	repo *mongodb.HomeSectionRepository,
	categories map[string]*entities.Category,
	products map[string]*entities.Product,
) error {
	// 1) Featured — mix of electronics + groceries (highlights)
	feat, err := entities.NewHomeSection(
		entities.HomeSectionTypeProducts,
		entities.LocalizedText{
			"en": "Featured",
			"ar": "مختارات مميزة",
		},
		[]primitive.ObjectID{
			products["IP13-GRA-128"].ID,
			products["MILK-WH-1L"].ID,
			products["ORG-BAN-1KG"].ID,
		},
		nil,
		1,
		true,
	)
	if err != nil {
		return err
	}
	if err := repo.Create(ctx, feat); err != nil {
		return err
	}

	// 2) Fresh today — grocery picks
	fresh, err := entities.NewHomeSection(
		entities.HomeSectionTypeProducts,
		entities.LocalizedText{
			"en": "Fresh today",
			"ar": "طازج اليوم",
		},
		[]primitive.ObjectID{
			products["ORG-BAN-1KG"].ID,
			products["BRD-SOUR-1"].ID,
			products["OJ-FRESH-1L"].ID,
			products["RICE-BAS-5KG"].ID,
		},
		nil,
		2,
		true,
	)
	if err != nil {
		return err
	}
	if err := repo.Create(ctx, fresh); err != nil {
		return err
	}

	// 3) Shop by aisle — category shortcuts
	shop, err := entities.NewHomeSection(
		entities.HomeSectionTypeCategories,
		entities.LocalizedText{
			"en": "Shop by aisle",
			"ar": "تسوق حسب الممر",
		},
		nil,
		[]primitive.ObjectID{
			categories["fresh-produce"].ID,
			categories["dairy-eggs"].ID,
			categories["bakery"].ID,
			categories["pantry"].ID,
			categories["beverages"].ID,
		},
		3,
		true,
	)
	if err != nil {
		return err
	}
	if err := repo.Create(ctx, shop); err != nil {
		return err
	}

	return nil
}
