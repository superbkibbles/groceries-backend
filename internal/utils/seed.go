package utils

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/superbkibbles/ecommerce/internal/adapters/repository/mongodb"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedData populates the database with sample data for testing
func SeedData(db *mongo.Database) error {
	// Initialize repositories
	productRepo := mongodb.NewProductRepository(db)
	categoryRepo := mongodb.NewCategoryRepository(db, productRepo)
	userRepo := mongodb.NewUserRepository(db)
	orderRepo := mongodb.NewOrderRepository(db)
	reviewRepo := mongodb.NewReviewRepository(db, orderRepo)

	ctx := context.Background()

	// Seed categories
	categories, err := seedCategories(ctx, categoryRepo)
	if err != nil {
		return err
	}

	// Seed products
	products, err := seedProducts(ctx, productRepo, categories)
	if err != nil {
		return err
	}

	// Seed users
	users, err := seedUsers(ctx, userRepo)
	if err != nil {
		return err
	}

	// Seed orders
	orders, err := seedOrders(ctx, orderRepo, users, products)
	if err != nil {
		return err
	}

	// Seed reviews
	_, err = seedReviews(ctx, reviewRepo, orders, users, products)
	if err != nil {
		return err
	}

	log.Println("Database seeded successfully!")
	return nil
}

// seedCategories creates sample categories
func seedCategories(ctx context.Context, repo *mongodb.CategoryRepository) (map[string]*entities.Category, error) {
	categories := make(map[string]*entities.Category)

	// Create main categories
	electronics := entities.NewCategory("Electronics", "Electronic devices and gadgets", "electronics", "")
	clothing := entities.NewCategory("Clothing", "Apparel and fashion items", "clothing", "")
	home := entities.NewCategory("Home & Kitchen", "Home goods and kitchen appliances", "home-kitchen", "")

	// Create subcategories for Electronics
	smartphones := entities.NewCategory("Smartphones", "Mobile phones and accessories", "smartphones", electronics.ID)
	laptops := entities.NewCategory("Laptops", "Notebook computers and accessories", "laptops", electronics.ID)
	audio := entities.NewCategory("Audio", "Headphones, speakers, and audio equipment", "audio", electronics.ID)

	// Create subcategories for Clothing
	mens := entities.NewCategory("Men's Clothing", "Clothing for men", "mens-clothing", clothing.ID)
	womens := entities.NewCategory("Women's Clothing", "Clothing for women", "womens-clothing", clothing.ID)
	kids := entities.NewCategory("Kids' Clothing", "Clothing for children", "kids-clothing", clothing.ID)

	// Create subcategories for Home & Kitchen
	furniture := entities.NewCategory("Furniture", "Home furniture and decor", "furniture", home.ID)
	kitchen := entities.NewCategory("Kitchen Appliances", "Appliances for cooking and food preparation", "kitchen-appliances", home.ID)
	bedding := entities.NewCategory("Bedding", "Sheets, pillows, and bedding accessories", "bedding", home.ID)

	// Save main categories first
	for _, category := range []*entities.Category{electronics, clothing, home} {
		if err := repo.Create(ctx, category); err != nil {
			return nil, err
		}
		categories[category.Name] = category
	}

	// Save subcategories
	subcategories := []*entities.Category{
		smartphones, laptops, audio, // Electronics subcategories
		mens, womens, kids, // Clothing subcategories
		furniture, kitchen, bedding, // Home & Kitchen subcategories
	}

	for _, category := range subcategories {
		if err := repo.Create(ctx, category); err != nil {
			return nil, err
		}
		categories[category.Name] = category
	}

	return categories, nil
}

// seedProducts creates sample products
func seedProducts(ctx context.Context, repo *mongodb.ProductRepository, categories map[string]*entities.Category) (map[string]*entities.Product, error) {
	products := make(map[string]*entities.Product)

	// Smartphone products
	smartphoneCategories := []string{categories["Electronics"].ID, categories["Smartphones"].ID}
	iphone := entities.NewProduct(
		"iPhone 13 Pro",
		"Apple's flagship smartphone with A15 Bionic chip and Pro camera system",
		999.99,
		smartphoneCategories,
	)

	// Add iPhone variations
	iPhoneColors := []string{"Graphite", "Gold", "Silver", "Sierra Blue"}
	iPhoneStorage := []int{128, 256, 512, 1024}

	for _, color := range iPhoneColors {
		for _, storage := range iPhoneStorage {
			attributes := map[string]interface{}{
				"color":   color,
				"storage": storage,
			}

			// Price increases with storage size
			price := 999.99
			switch storage {
			case 256:
				price = 1099.99
			case 512:
				price = 1299.99
			case 1024:
				price = 1499.99
			}

			sku := "IP13-" + color[:3] + "-" + string(rune(storage/128+64))
			images := []string{"iphone13-" + color + ".jpg"}

			_, err := iphone.AddVariation(attributes, sku, price, 50, images)
			if err != nil {
				return nil, err
			}
		}
	}

	// Laptop product
	laptopCategories := []string{categories["Electronics"].ID, categories["Laptops"].ID}
	macbook := entities.NewProduct(
		"MacBook Pro 16",
		"Powerful laptop for professionals with M1 Pro or M1 Max chip",
		2499.99,
		laptopCategories,
	)

	// Add MacBook variations
	macbookChips := []string{"M1 Pro", "M1 Max"}
	macbookRAM := []int{16, 32, 64}
	macbookStorage := []int{512, 1024, 2048, 4096}

	for _, chip := range macbookChips {
		for _, ram := range macbookRAM {
			// Skip invalid combinations
			if chip == "M1 Pro" && ram == 64 {
				continue
			}

			for _, storage := range macbookStorage {
				attributes := map[string]interface{}{
					"chip":    chip,
					"ram":     ram,
					"storage": storage,
				}

				// Base price for M1 Pro with 16GB RAM and 512GB storage
				price := 2499.99

				// Add for chip upgrade
				if chip == "M1 Max" {
					price += 200
				}

				// Add for RAM upgrade
				if ram == 32 {
					price += 400
				} else if ram == 64 {
					price += 800
				}

				// Add for storage upgrade
				switch storage {
				case 1024:
					price += 200
				case 2048:
					price += 600
				case 4096:
					price += 1200
				}

				sku := "MBP16-" + chip[3:] + "-" + string(rune(ram/16+64)) + "-" + string(rune(storage/512+64))
				images := []string{"macbook-pro-16.jpg"}

				_, err := macbook.AddVariation(attributes, sku, price, 20, images)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// Clothing product
	shirtCategories := []string{categories["Clothing"].ID, categories["Men's Clothing"].ID}
	tshirt := entities.NewProduct(
		"Premium Cotton T-Shirt",
		"Soft, comfortable 100% cotton t-shirt",
		29.99,
		shirtCategories,
	)

	// Add T-shirt variations
	tshirtColors := []string{"Black", "White", "Navy", "Gray", "Red"}
	tshirtSizes := []string{"S", "M", "L", "XL", "XXL"}

	for _, color := range tshirtColors {
		for _, size := range tshirtSizes {
			attributes := map[string]interface{}{
				"color": color,
				"size":  size,
			}

			// Larger sizes cost slightly more
			price := 29.99
			if size == "XL" {
				price = 32.99
			} else if size == "XXL" {
				price = 34.99
			}

			sku := "TS-" + color[:2] + "-" + size
			images := []string{"tshirt-" + color + ".jpg"}

			_, err := tshirt.AddVariation(attributes, sku, price, 100, images)
			if err != nil {
				return nil, err
			}
		}
	}

	// Kitchen product
	kitchenCategories := []string{categories["Home & Kitchen"].ID, categories["Kitchen Appliances"].ID}
	blender := entities.NewProduct(
		"High-Performance Blender",
		"Powerful blender for smoothies, soups, and more",
		149.99,
		kitchenCategories,
	)

	// Add blender variations
	blenderColors := []string{"Black", "Silver", "Red"}
	blenderWattages := []int{600, 900, 1200}

	for _, color := range blenderColors {
		for _, wattage := range blenderWattages {
			attributes := map[string]interface{}{
				"color":   color,
				"wattage": wattage,
			}

			// Higher wattage costs more
			price := 149.99
			if wattage == 900 {
				price = 179.99
			} else if wattage == 1200 {
				price = 199.99
			}

			sku := "BL-" + color[:2] + "-" + string(rune(wattage/300+64))
			images := []string{"blender-" + color + ".jpg"}

			_, err := blender.AddVariation(attributes, sku, price, 30, images)
			if err != nil {
				return nil, err
			}
		}
	}

	// Save all products
	for _, product := range []*entities.Product{iphone, macbook, tshirt, blender} {
		if err := repo.Create(ctx, product); err != nil {
			return nil, err
		}
		products[product.Name] = product
	}

	return products, nil
}

// seedUsers creates sample users
func seedUsers(ctx context.Context, repo *mongodb.UserRepository) (map[string]*entities.User, error) {
	users := make(map[string]*entities.User)

	// Create admin user
	admin, err := entities.NewUser("admin@example.com", "Admin123!", "Admin", "User", entities.UserRoleAdmin)
	if err != nil {
		return nil, err
	}

	// Create customer users
	john, err := entities.NewUser("john@example.com", "John123!", "John", "Doe", entities.UserRoleCustomer)
	if err != nil {
		return nil, err
	}

	jane, err := entities.NewUser("jane@example.com", "Jane123!", "Jane", "Smith", entities.UserRoleCustomer)
	if err != nil {
		return nil, err
	}

	// Save users
	for _, user := range []*entities.User{admin, john, jane} {
		if err := repo.Create(ctx, user); err != nil {
			return nil, err
		}
		users[user.Email] = user
	}

	// Add addresses for customers
	johnAddress := &entities.Address{
		UserID:       john.ID,
		Name:         "Home",
		AddressLine1: "123 Main St",
		City:         "New York",
		State:        "NY",
		Country:      "USA",
		PostalCode:   "10001",
		Phone:        "555-123-4567",
		IsDefault:    true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	janeAddress := &entities.Address{
		UserID:       jane.ID,
		Name:         "Home",
		AddressLine1: "456 Oak Ave",
		City:         "Los Angeles",
		State:        "CA",
		Country:      "USA",
		PostalCode:   "90001",
		Phone:        "555-987-6543",
		IsDefault:    true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save addresses
	for _, address := range []*entities.Address{johnAddress, janeAddress} {
		if err := repo.AddAddress(ctx, address); err != nil {
			return nil, err
		}
	}

	return users, nil
}

// seedOrders creates sample orders
func seedOrders(ctx context.Context, repo *mongodb.OrderRepository, users map[string]*entities.User, products map[string]*entities.Product) (map[string]*entities.Order, error) {
	orders := make(map[string]*entities.Order)

	// Get John's user ID
	john := users["john@example.com"]

	// Create an order for John
	johnOrder := &entities.Order{
		CustomerID: john.ID,
		Items:      []*entities.OrderItem{},
		Status:     entities.OrderStatusDelivered,
		ShippingInfo: entities.ShippingInfo{
			Address:    "123 Main St",
			City:       "New York",
			State:      "NY",
			Country:    "USA",
			PostalCode: "10001",
		},
		PaymentInfo: entities.PaymentInfo{
			Method:    "Credit Card",
			Status:    "Paid",
			Timestamp: time.Now().Add(-7 * 24 * time.Hour), // 1 week ago
		},
		CreatedAt: time.Now().Add(-10 * 24 * time.Hour), // 10 days ago
		UpdatedAt: time.Now().Add(-7 * 24 * time.Hour),  // 7 days ago
	}

	// Add iPhone to John's order
	iphone := products["iPhone 13 Pro"]
	iphoneVariation := iphone.Variations[0] // First variation
	johnOrder.Items = append(johnOrder.Items, &entities.OrderItem{
		ProductID:   iphone.ID,
		VariationID: iphoneVariation.ID,
		SKU:         iphoneVariation.SKU,
		Name:        iphone.Name,
		Price:       iphoneVariation.Price,
		Quantity:    1,
		Subtotal:    iphoneVariation.Price * 1,
	})

	// Add T-shirt to John's order
	tshirt := products["Premium Cotton T-Shirt"]
	tshirtVariation := tshirt.Variations[0] // First variation
	johnOrder.Items = append(johnOrder.Items, &entities.OrderItem{
		ProductID:   tshirt.ID,
		VariationID: tshirtVariation.ID,
		SKU:         tshirtVariation.SKU,
		Name:        tshirt.Name,
		Price:       tshirtVariation.Price,
		Quantity:    2,
		Subtotal:    tshirtVariation.Price * 2,
	})

	// Calculate total amount
	totalAmount := 0.0
	for _, item := range johnOrder.Items {
		totalAmount += item.Subtotal
	}
	johnOrder.TotalAmount = totalAmount

	// Get Jane's user ID
	jane := users["jane@example.com"]

	// Create an order for Jane
	janeOrder := &entities.Order{
		CustomerID: jane.ID,
		Items:      []*entities.OrderItem{},
		Status:     entities.OrderStatusPaid,
		ShippingInfo: entities.ShippingInfo{
			Address:    "456 Oak Ave",
			City:       "Los Angeles",
			State:      "CA",
			Country:    "USA",
			PostalCode: "90001",
		},
		PaymentInfo: entities.PaymentInfo{
			Method:    "PayPal",
			Status:    "Paid",
			Timestamp: time.Now().Add(-2 * 24 * time.Hour), // 2 days ago
		},
		CreatedAt: time.Now().Add(-3 * 24 * time.Hour), // 3 days ago
		UpdatedAt: time.Now().Add(-2 * 24 * time.Hour), // 2 days ago
	}

	// Add MacBook to Jane's order
	macbook := products["MacBook Pro 16"]
	macbookVariation := macbook.Variations[0] // First variation
	janeOrder.Items = append(janeOrder.Items, &entities.OrderItem{
		ProductID:   macbook.ID,
		VariationID: macbookVariation.ID,
		SKU:         macbookVariation.SKU,
		Name:        macbook.Name,
		Price:       macbookVariation.Price,
		Quantity:    1,
		Subtotal:    macbookVariation.Price * 1,
	})

	// Calculate total amount
	totalAmount = 0.0
	for _, item := range janeOrder.Items {
		totalAmount += item.Subtotal
	}
	janeOrder.TotalAmount = totalAmount

	// Save orders
	for _, order := range []*entities.Order{johnOrder, janeOrder} {
		if err := repo.Create(ctx, order); err != nil {
			return nil, err
		}
		orders[order.ID] = order
	}

	return orders, nil
}

// seedReviews creates sample reviews for products in orders
func seedReviews(ctx context.Context, repo *mongodb.ReviewRepository, orders map[string]*entities.Order, users map[string]*entities.User, products map[string]*entities.Product) (map[string]*entities.Review, error) {
	reviews := make(map[string]*entities.Review)

	// Get users
	john := users["john@example.com"]

	// Find John's order
	var johnOrder *entities.Order
	for _, order := range orders {
		if order.CustomerID == john.ID && order.Status == entities.OrderStatusDelivered {
			johnOrder = order
			break
		}
	}

	if johnOrder == nil {
		return nil, errors.New("could not find John's delivered order")
	}

	// Create reviews for John's order items
	for _, item := range johnOrder.Items {
		var rating int
		var comment string

		// Different ratings and comments based on product
		if item.ProductID == products["iPhone 13 Pro"].ID {
			rating = 5
			comment = "Excellent phone! The camera is amazing and battery life is great. Highly recommended!"
		} else if item.ProductID == products["Premium Cotton T-Shirt"].ID {
			rating = 4
			comment = "Very comfortable t-shirt, good quality cotton. Fits well but slightly larger than expected."
		} else {
			// Default rating and comment for other products
			rating = 4
			comment = "Good product, satisfied with my purchase."
		}

		// Create the review
		review, err := entities.NewReview(item.ProductID, john.ID, johnOrder.ID, rating, comment)
		if err != nil {
			return nil, err
		}

		// Set creation date to a few days after delivery
		review.CreatedAt = johnOrder.UpdatedAt.Add(3 * 24 * time.Hour)
		review.UpdatedAt = review.CreatedAt

		// Save the review
		if err := repo.Create(ctx, review); err != nil {
			return nil, err
		}

		reviews[review.ID] = review
	}

	return reviews, nil
}
