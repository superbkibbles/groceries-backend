package utils

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/superbkibbles/ecommerce/internal/adapters/repository/mongodb"
	"github.com/superbkibbles/ecommerce/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedData populates the database with sample data for testing
func SeedData(db *mongo.Database, redisClient *redis.Client) error {
	// Initialize repositories
	productRepo := mongodb.NewProductRepository(db)
	categoryRepo := mongodb.NewCategoryRepository(db, productRepo)
	userRepo := mongodb.NewUserRepository(db, redisClient)
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
	electronics := entities.NewCategory("Electronics", "Electronic devices and gadgets", "electronics", primitive.NilObjectID)
	clothing := entities.NewCategory("Clothing", "Apparel and fashion items", "clothing", primitive.NilObjectID)
	home := entities.NewCategory("Home & Kitchen", "Home goods and kitchen appliances", "home-kitchen", primitive.NilObjectID)

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
	smartphoneCategories := []primitive.ObjectID{categories["Electronics"].ID, categories["Smartphones"].ID}
	iphone := entities.NewProduct(
		"iPhone 13 Pro",
		"Apple's flagship smartphone with A15 Bionic chip and Pro camera system",
		smartphoneCategories,
		map[string]interface{}{
			"color":   "Graphite",
			"storage": 128,
		},
		"IP13-GRA-128",
		999.99,
		50,
		[]string{"iphone13-graphite.jpg"},
	)

	// Laptop product
	laptopCategories := []primitive.ObjectID{categories["Electronics"].ID, categories["Laptops"].ID}
	macbook := entities.NewProduct(
		"MacBook Pro 16",
		"Powerful laptop for professionals with M1 Pro or M1 Max chip",
		laptopCategories,
		map[string]interface{}{
			"chip":    "M1 Pro",
			"ram":     16,
			"storage": 512,
		},
		"MBP16-PRO-16-512",
		2499.99,
		20,
		[]string{"macbook-pro-16.jpg"},
	)

	// Clothing product
	shirtCategories := []primitive.ObjectID{categories["Clothing"].ID, categories["Men's Clothing"].ID}
	tshirt := entities.NewProduct(
		"Premium Cotton T-Shirt",
		"Soft, comfortable 100% cotton t-shirt",
		shirtCategories,
		map[string]interface{}{
			"color": "Black",
			"size":  "M",
		},
		"TS-BL-M",
		29.99,
		100,
		[]string{"tshirt-black.jpg"},
	)

	// Kitchen product
	kitchenCategories := []primitive.ObjectID{categories["Home & Kitchen"].ID, categories["Kitchen Appliances"].ID}
	blender := entities.NewProduct(
		"High-Performance Blender",
		"Powerful blender for smoothies, soups, and more",
		kitchenCategories,
		map[string]interface{}{
			"color":   "Black",
			"wattage": 600,
		},
		"BL-BL-600",
		149.99,
		30,
		[]string{"blender-black.jpg"},
	)

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
func seedOrders(ctx context.Context, repo *mongodb.OrderRepository, users map[string]*entities.User, products map[string]*entities.Product) (map[primitive.ObjectID]*entities.Order, error) {
	orders := make(map[primitive.ObjectID]*entities.Order)

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
	johnOrder.Items = append(johnOrder.Items, &entities.OrderItem{
		ProductID: iphone.ID,
		SKU:       iphone.SKU,
		Name:      iphone.Name,
		Price:     iphone.Price,
		Quantity:  1,
		Subtotal:  iphone.Price * 1,
	})

	// Add T-shirt to John's order
	tshirt := products["Premium Cotton T-Shirt"]
	johnOrder.Items = append(johnOrder.Items, &entities.OrderItem{
		ProductID: tshirt.ID,
		SKU:       tshirt.SKU,
		Name:      tshirt.Name,
		Price:     tshirt.Price,
		Quantity:  2,
		Subtotal:  tshirt.Price * 2,
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
	janeOrder.Items = append(janeOrder.Items, &entities.OrderItem{
		ProductID: macbook.ID,
		SKU:       macbook.SKU,
		Name:      macbook.Name,
		Price:     macbook.Price,
		Quantity:  1,
		Subtotal:  macbook.Price * 1,
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
func seedReviews(ctx context.Context, repo *mongodb.ReviewRepository, orders map[primitive.ObjectID]*entities.Order, users map[string]*entities.User, products map[string]*entities.Product) (map[primitive.ObjectID]*entities.Review, error) {
	reviews := make(map[primitive.ObjectID]*entities.Review)

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
