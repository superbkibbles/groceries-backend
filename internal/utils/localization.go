package utils

import (
	"errors"

	"github.com/superbkibbles/ecommerce/internal/domain/entities"
)

// ValidateTranslations validates that translations are provided and English is included
func ValidateTranslations(translations map[string]entities.Translation) error {
	if len(translations) == 0 {
		return errors.New("at least one translation is required")
	}

	if _, hasEnglish := translations["en"]; !hasEnglish {
		return errors.New("English translation is required")
	}

	return nil
}

// ApplyLocalizationToProducts applies localization to a slice of products
func ApplyLocalizationToProducts(products []*entities.Product, language string) {
	for _, product := range products {
		product.ApplyLocalization(language)
	}
}

// ApplyLocalizationToCategories applies localization to a slice of categories
func ApplyLocalizationToCategories(categories []*entities.Category, language string) {
	for _, category := range categories {
		category.ApplyLocalization(language)
	}
}
