package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestItem(t *testing.T) {
	t.Run("Create Item with valid data", func(t *testing.T) {
		item := Item{
			ID:          primitive.NewObjectID(),
			Name:        "Test Item",
			Description: "This is a test item",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if item.Name != "Test Item" {
			t.Errorf("Expected name to be 'Test Item', got '%s'", item.Name)
		}

		if item.Description != "This is a test item" {
			t.Errorf("Expected description to be 'This is a test item', got '%s'", item.Description)
		}
	})
}
