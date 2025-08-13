// MongoDB initialization script
// This script creates the database and initial collections

// Switch to the groceries database
db = db.getSiblingDB("groceries_db");

// Create a regular user for the application
db.createUser({
  user: "groceries_user",
  pwd: "groceries_password",
  roles: [
    {
      role: "readWrite",
      db: "groceries_db",
    },
  ],
});

// Create collections with validation rules
db.createCollection("products", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["_id", "name", "sku", "price"],
      properties: {
        _id: {
          bsonType: "string",
          description: "must be a string and is required",
        },
        name: {
          bsonType: "string",
          description: "must be a string and is required",
        },
        description: {
          bsonType: "string",
          description: "must be a string",
        },
        sku: {
          bsonType: "string",
          description: "must be a string and is required",
        },
        price: {
          bsonType: "double",
          minimum: 0,
          description: "must be a positive number and is required",
        },
        stock_quantity: {
          bsonType: "int",
          minimum: 0,
          description: "must be a non-negative integer",
        },
        categories: {
          bsonType: "array",
          description: "must be an array of category IDs",
        },
        attributes: {
          bsonType: "object",
          description: "must be an object containing product attributes",
        },
        images: {
          bsonType: "array",
          description: "must be an array of image URLs",
        },
      },
    },
  },
});

db.createCollection("categories", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["_id", "name", "slug"],
      properties: {
        _id: {
          bsonType: "string",
          description: "must be a string and is required",
        },
        name: {
          bsonType: "string",
          description: "must be a string and is required",
        },
        slug: {
          bsonType: "string",
          description: "must be a string and is required",
        },
      },
    },
  },
});

db.createCollection("users", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["_id", "email"],
      properties: {
        _id: {
          bsonType: "string",
          description: "must be a string and is required",
        },
        email: {
          bsonType: "string",
          pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
          description: "must be a valid email and is required",
        },
      },
    },
  },
});

db.createCollection("orders");
db.createCollection("carts");
db.createCollection("reviews");
db.createCollection("settings");

// Create indexes for better performance
db.products.createIndex({ sku: 1 }, { unique: true });
db.products.createIndex({ name: "text", description: "text" });
db.products.createIndex({ categories: 1 });
db.products.createIndex({ price: 1 });

db.categories.createIndex({ slug: 1 }, { unique: true });
db.categories.createIndex({ parent_id: 1 });

db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ phone_number: 1 }, { unique: true, sparse: true });

db.orders.createIndex({ customer_id: 1 });
db.orders.createIndex({ status: 1 });
db.orders.createIndex({ created_at: 1 });

db.carts.createIndex({ user_id: 1 }, { unique: true });

db.reviews.createIndex({ product_id: 1 });
db.reviews.createIndex({ user_id: 1 });
db.reviews.createIndex({ order_id: 1 });

print("MongoDB initialization completed successfully!");
print("Database: groceries_db");
print("Collections created with validation rules and indexes");
print("Application user created: groceries_user");
