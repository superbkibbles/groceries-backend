#!/bin/bash

# Quick ObjectID Migration Script
# This script applies systematic changes to migrate from string IDs to ObjectID

echo "🔄 Starting ObjectID migration..."

# Function to update entity imports
update_imports() {
    local file=$1
    echo "📝 Updating imports in $file"
    
    # Replace uuid import with primitive import
    sed -i '' 's|"github.com/google/uuid"|"go.mongodb.org/mongo-driver/bson/primitive"|g' "$file"
    
    # Remove uuid import if primitive already exists
    if grep -q "primitive" "$file"; then
        sed -i '' '/github\.com\/google\/uuid/d' "$file"
    fi
}

# Function to update ID fields in structs
update_id_fields() {
    local file=$1
    echo "📝 Updating ID fields in $file"
    
    # Update main ID field
    sed -i '' 's|ID.*string.*`json:"id".*bson:"_id"`|ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`|g' "$file"
    
    # Update reference ID fields
    sed -i '' 's|UserID.*string.*`json:"user_id"|UserID primitive.ObjectID `json:"user_id"|g' "$file"
    sed -i '' 's|ProductID.*string.*`json:"product_id"|ProductID primitive.ObjectID `json:"product_id"|g' "$file"
    sed -i '' 's|OrderID.*string.*`json:"order_id"|OrderID primitive.ObjectID `json:"order_id"|g' "$file"
    sed -i '' 's|CategoryID.*string.*`json:"category_id"|CategoryID primitive.ObjectID `json:"category_id"|g' "$file"
    sed -i '' 's|ParentID.*string.*`json:"parent_id"|ParentID primitive.ObjectID `json:"parent_id"|g' "$file"
    sed -i '' 's|CartID.*string.*`json:"cart_id"|CartID primitive.ObjectID `json:"cart_id"|g' "$file"
}

# Function to remove UUID generation from constructors
remove_uuid_generation() {
    local file=$1
    echo "📝 Removing UUID generation in $file"
    
    # Remove ID assignment lines in constructors
    sed -i '' '/ID:.*uuid\.New()\.String()/d' "$file"
}

# Update remaining entity files
ENTITY_DIR="internal/domain/entities"

echo "🔧 Updating remaining entity files..."

# Files that still need updates
entities=("review.go" "notification.go" "wishlist.go" "setting.go" "payment.go" "shipping.go")

for entity in "${entities[@]}"; do
    entity_file="$ENTITY_DIR/$entity"
    if [[ -f "$entity_file" ]]; then
        echo "⚡ Processing $entity_file"
        update_imports "$entity_file"
        update_id_fields "$entity_file"
        remove_uuid_generation "$entity_file"
    else
        echo "⚠️  File not found: $entity_file"
    fi
done

echo "✅ Entity files updated!"

# Note: After running this script, manual fixes will still be needed for:
# - Constructor parameter types (string -> primitive.ObjectID)
# - String comparison logic (== "" -> .IsZero())
# - Repository implementations
# - Service layers
# - HTTP handlers

echo "📋 Next steps:"
echo "1. Fix constructor parameter types manually"
echo "2. Update string comparison logic to use .IsZero()"
echo "3. Update repository implementations"
echo "4. Update service layers"
echo "5. Update HTTP handlers"
echo "6. Test the build: make build"

echo "🎉 Basic entity ObjectID migration complete!"
