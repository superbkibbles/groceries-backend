package utils

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParseObjectID safely parses a string to ObjectID
func ParseObjectID(id string) (primitive.ObjectID, error) {
	if id == "" {
		return primitive.NilObjectID, errors.New("empty ID")
	}
	return primitive.ObjectIDFromHex(id)
}

// ObjectIDToString safely converts ObjectID to string
func ObjectIDToString(id primitive.ObjectID) string {
	if id.IsZero() {
		return ""
	}
	return id.Hex()
}

// ParseObjectIDSlice converts a slice of string IDs to ObjectIDs
func ParseObjectIDSlice(ids []string) ([]primitive.ObjectID, error) {
	result := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objID, err := ParseObjectID(id)
		if err != nil {
			return nil, err
		}
		result[i] = objID
	}
	return result, nil
}

// ObjectIDSliceToString converts a slice of ObjectIDs to strings
func ObjectIDSliceToString(ids []primitive.ObjectID) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = ObjectIDToString(id)
	}
	return result
}

// IsValidObjectID checks if a string is a valid ObjectID
func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

// NewObjectID generates a new ObjectID
func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// NilObjectID returns a zero ObjectID
func NilObjectID() primitive.ObjectID {
	return primitive.NilObjectID
}
