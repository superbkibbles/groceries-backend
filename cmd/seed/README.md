# E-Commerce Seed Data Utility

This utility populates your MongoDB database with dummy data for testing and development purposes.

## What Data is Added

### Categories
- Main categories: Electronics, Clothing, Home & Kitchen
- Subcategories for each main category (e.g., Smartphones, Laptops under Electronics)

### Products
- iPhone 13 Pro with variations (different colors and storage options)
- MacBook Pro 16 with variations (different chip, RAM, and storage configurations)
- Premium Cotton T-Shirt with variations (different colors and sizes)
- High-Performance Blender with variations (different colors and wattages)

### Users
- Admin user: admin@example.com (password: Admin123!)
- Customer users:
  - John Doe: john@example.com (password: John123!)
  - Jane Smith: jane@example.com (password: Jane123!)
- Addresses for each customer

### Orders
- Completed order for John (iPhone and T-shirts)
- Paid order for Jane (MacBook)

## How to Use

### Run the Seed Utility

```bash
# Navigate to the seed directory
cd cmd/seed

# Build the seed utility
go build

# Run the seed utility
./seed
```

Alternatively, you can run it directly with:

```bash
go run cmd/seed/main.go
```

## Notes

- The seed utility will connect to MongoDB using the connection details from your environment variables or default values
- If you've already populated the database, running the seed utility again may result in duplicate data
- All passwords for test users are hashed securely