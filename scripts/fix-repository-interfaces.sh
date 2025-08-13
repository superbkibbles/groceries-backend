#!/bin/bash

# Fix Repository Interfaces Script
# This script updates all repository interfaces to use primitive.ObjectID

echo "🔧 Updating repository interfaces for ObjectID..."

REPO_FILE="internal/domain/ports/repositories.go"

# Function to update ID parameters in repository interfaces
update_repository_interfaces() {
    echo "📝 Updating repository interface methods in $REPO_FILE"
    
    # Update all ID-related parameters to use primitive.ObjectID
    sed -i '' 's/id string/id primitive.ObjectID/g' "$REPO_FILE"
    sed -i '' 's/userID string/userID primitive.ObjectID/g' "$REPO_FILE" 
    sed -i '' 's/customerID string/customerID primitive.ObjectID/g' "$REPO_FILE"
    sed -i '' 's/productID string/productID primitive.ObjectID/g' "$REPO_FILE"
    sed -i '' 's/orderID string/orderID primitive.ObjectID/g' "$REPO_FILE"
    sed -i '' 's/categoryID string/categoryID primitive.ObjectID/g' "$REPO_FILE"
    sed -i '' 's/parentID string/parentID primitive.ObjectID/g' "$REPO_FILE"
    sed -i '' 's/rootID string/rootID primitive.ObjectID/g' "$REPO_FILE"
    
    echo "✅ Repository interfaces updated!"
}

# Run the update
if [[ -f "$REPO_FILE" ]]; then
    update_repository_interfaces
    echo "🎉 Repository interface ObjectID migration complete!"
else
    echo "❌ Repository interface file not found: $REPO_FILE"
    exit 1
fi

echo "📋 Next: Update repository implementations to match the new interfaces"
