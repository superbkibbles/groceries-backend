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

	// Create main categories with translations
	electronics := entities.NewCategoryWithTranslations("electronics", primitive.NilObjectID, map[string]entities.Translation{
		"en": {Name: "Electronics", Description: "Electronic devices and gadgets"},
		"ar": {Name: "الإلكترونيات", Description: "الأجهزة الإلكترونية والأدوات"},
	})
	clothing := entities.NewCategoryWithTranslations("clothing", primitive.NilObjectID, map[string]entities.Translation{
		"en": {Name: "Clothing", Description: "Apparel and fashion items"},
		"ar": {Name: "الملابس", Description: "الملابس وعناصر الموضة"},
	})
	home := entities.NewCategoryWithTranslations("home-kitchen", primitive.NilObjectID, map[string]entities.Translation{
		"en": {Name: "Home & Kitchen", Description: "Home goods and kitchen appliances"},
		"ar": {Name: "المنزل والمطبخ", Description: "السلع المنزلية وأجهزة المطبخ"},
	})

	// Save main categories first to get their IDs
	for _, category := range []*entities.Category{electronics, clothing, home} {
		if err := repo.Create(ctx, category); err != nil {
			return nil, err
		}
		categories[category.Slug] = category
	}

	// Create subcategories for Electronics
	smartphones := entities.NewCategoryWithTranslations("smartphones", electronics.ID, map[string]entities.Translation{
		"en": {Name: "Smartphones", Description: "Mobile phones and accessories"},
		"ar": {Name: "الهواتف الذكية", Description: "الهواتف المحمولة والإكسسوارات"},
	})
	laptops := entities.NewCategoryWithTranslations("laptops", electronics.ID, map[string]entities.Translation{
		"en": {Name: "Laptops", Description: "Notebook computers and accessories"},
		"ar": {Name: "أجهزة الكمبيوتر المحمولة", Description: "أجهزة الكمبيوتر المحمول والإكسسوارات"},
	})
	audio := entities.NewCategoryWithTranslations("audio", electronics.ID, map[string]entities.Translation{
		"en": {Name: "Audio", Description: "Headphones, speakers, and audio equipment"},
		"ar": {Name: "الصوتيات", Description: "سماعات الرأس والمكبرات والمعدات الصوتية"},
	})

	// Create subcategories for Clothing
	mens := entities.NewCategoryWithTranslations("mens-clothing", clothing.ID, map[string]entities.Translation{
		"en": {Name: "Men's Clothing", Description: "Clothing for men"},
		"ar": {Name: "ملابس رجالية", Description: "ملابس للرجال"},
	})
	womens := entities.NewCategoryWithTranslations("womens-clothing", clothing.ID, map[string]entities.Translation{
		"en": {Name: "Women's Clothing", Description: "Clothing for women"},
		"ar": {Name: "ملابس نسائية", Description: "ملابس للنساء"},
	})
	kids := entities.NewCategoryWithTranslations("kids-clothing", clothing.ID, map[string]entities.Translation{
		"en": {Name: "Kids' Clothing", Description: "Clothing for children"},
		"ar": {Name: "ملابس أطفال", Description: "ملابس للأطفال"},
	})

	// Create subcategories for Home & Kitchen
	furniture := entities.NewCategoryWithTranslations("furniture", home.ID, map[string]entities.Translation{
		"en": {Name: "Furniture", Description: "Home furniture and decor"},
		"ar": {Name: "الأثاث", Description: "أثاث المنزل والديكور"},
	})
	kitchen := entities.NewCategoryWithTranslations("kitchen-appliances", home.ID, map[string]entities.Translation{
		"en": {Name: "Kitchen Appliances", Description: "Appliances for cooking and food preparation"},
		"ar": {Name: "أجهزة المطبخ", Description: "أجهزة للطبخ وتحضير الطعام"},
	})
	bedding := entities.NewCategoryWithTranslations("bedding", home.ID, map[string]entities.Translation{
		"en": {Name: "Bedding", Description: "Sheets, pillows, and bedding accessories"},
		"ar": {Name: "الفراش", Description: "ملاءات ووسائد وإكسسوارات الفراش"},
	})

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
		categories[category.Slug] = category
	}

	return categories, nil
}

// seedProducts creates sample products
func seedProducts(ctx context.Context, repo *mongodb.ProductRepository, categories map[string]*entities.Category) (map[string]*entities.Product, error) {
	products := make(map[string]*entities.Product)

	// Smartphone products
	smartphoneCategories := []primitive.ObjectID{categories["electronics"].ID, categories["smartphones"].ID}
	iphone := entities.NewProductWithTranslations(
		smartphoneCategories,
		map[string]interface{}{
			"color":   "Graphite",
			"storage": 128,
		},
		"IP13-GRA-128",
		999.99,
		50,
		[]string{"iphone13-graphite.jpg"},
		map[string]entities.Translation{
			"en": {
				Name:        "iPhone 13 Pro",
				Description: "Apple's flagship smartphone with A15 Bionic chip and Pro camera system",
			},
			"ar": {
				Name:        "آيفون 13 برو",
				Description: "الهاتف الذكي الرائد من Apple مع شريحة A15 Bionic ونظام كاميرا Pro",
			},
		},
	)

	// Laptop product
	laptopCategories := []primitive.ObjectID{categories["electronics"].ID, categories["laptops"].ID}
	macbook := entities.NewProductWithTranslations(
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
		map[string]entities.Translation{
			"en": {
				Name:        "MacBook Pro 16",
				Description: "Powerful laptop for professionals with M1 Pro or M1 Max chip",
			},
			"ar": {
				Name:        "ماك بوك برو 16",
				Description: "كمبيوتر محمول قوي للمحترفين مع شريحة M1 Pro أو M1 Max",
			},
		},
	)

	// Clothing product
	shirtCategories := []primitive.ObjectID{categories["clothing"].ID, categories["mens-clothing"].ID}
	tshirt := entities.NewProductWithTranslations(
		shirtCategories,
		map[string]interface{}{
			"color": "Black",
			"size":  "M",
		},
		"TS-BL-M",
		29.99,
		100,
		[]string{"tshirt-black.jpg"},
		map[string]entities.Translation{
			"en": {
				Name:        "Premium Cotton T-Shirt",
				Description: "Soft, comfortable 100% cotton t-shirt",
			},
			"ar": {
				Name:        "قميص قطني فاخر",
				Description: "قميص قطني ناعم ومريح 100٪",
			},
		},
	)

	// Kitchen product
	kitchenCategories := []primitive.ObjectID{categories["home-kitchen"].ID, categories["kitchen-appliances"].ID}
	blender := entities.NewProductWithTranslations(
		kitchenCategories,
		map[string]interface{}{
			"color":   "Black",
			"wattage": 600,
		},
		"BL-BL-600",
		149.99,
		30,
		[]string{"blender-black.jpg"},
		map[string]entities.Translation{
			"en": {
				Name:        "High-Performance Blender",
				Description: "Powerful blender for smoothies, soups, and more",
			},
			"ar": {
				Name:        "خلاط عالي الأداء",
				Description: "خلاط قوي للعصائر والحساء والمزيد",
			},
		},
	)

	// Save all products
	for _, product := range []*entities.Product{iphone, macbook, tshirt, blender} {
		if err := repo.Create(ctx, product); err != nil {
			return nil, err
		}
		products[product.SKU] = product
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
	johnOrder := entities.NewOrder(john.ID, entities.ShippingInfo{
		Address:    "123 Main St",
		City:       "New York",
		State:      "NY",
		Country:    "USA",
		PostalCode: "10001",
	})

	// Add iPhone to John's order
	iphone := products["IP13-GRA-128"]
	if err := johnOrder.AddItem(iphone.ID, iphone.SKU, iphone.Name, iphone.Price, 1); err != nil {
		return nil, err
	}

	// Add T-shirt to John's order
	tshirt := products["TS-BL-M"]
	if err := johnOrder.AddItem(tshirt.ID, tshirt.SKU, tshirt.Name, tshirt.Price, 2); err != nil {
		return nil, err
	}

	// Set order status and payment info after adding items
	johnOrder.Status = entities.OrderStatusDelivered
	johnOrder.PaymentInfo = entities.PaymentInfo{
		Method:    "Credit Card",
		Status:    "Paid",
		Timestamp: time.Now().Add(-7 * 24 * time.Hour), // 1 week ago
	}
	johnOrder.CreatedAt = time.Now().Add(-10 * 24 * time.Hour) // 10 days ago
	johnOrder.UpdatedAt = time.Now().Add(-7 * 24 * time.Hour)  // 7 days ago

	// Get Jane's user ID
	jane := users["jane@example.com"]

	// Create an order for Jane
	janeOrder := entities.NewOrder(jane.ID, entities.ShippingInfo{
		Address:    "456 Oak Ave",
		City:       "Los Angeles",
		State:      "CA",
		Country:    "USA",
		PostalCode: "90001",
	})

	// Add MacBook to Jane's order
	macbook := products["MBP16-PRO-16-512"]
	if err := janeOrder.AddItem(macbook.ID, macbook.SKU, macbook.Name, macbook.Price, 1); err != nil {
		return nil, err
	}

	// Set order status and payment info after adding items
	janeOrder.Status = entities.OrderStatusPaid
	janeOrder.PaymentInfo = entities.PaymentInfo{
		Method:    "PayPal",
		Status:    "Paid",
		Timestamp: time.Now().Add(-2 * 24 * time.Hour), // 2 days ago
	}
	janeOrder.CreatedAt = time.Now().Add(-3 * 24 * time.Hour) // 3 days ago
	janeOrder.UpdatedAt = time.Now().Add(-2 * 24 * time.Hour) // 2 days ago

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
		if item.ProductID == products["IP13-GRA-128"].ID {
			rating = 5
			comment = "Excellent phone! The camera is amazing and battery life is great. Highly recommended!"
		} else if item.ProductID == products["TS-BL-M"].ID {
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
