package repository

import (
	"context"

	"github.com/AliSleiman0/Lacpa/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ItemRepository defines the interface for item database operations
type ItemRepository interface {
	CreateItem(ctx context.Context, item *models.Item) error
	GetItem(ctx context.Context, id string) (*models.Item, error)
	GetAllItems(ctx context.Context) ([]*models.Item, error)
	UpdateItem(ctx context.Context, id string, item *models.Item) error
	DeleteItem(ctx context.Context, id string) error
}

// MongoItemRepository implements ItemRepository interface for MongoDB
type MongoItemRepository struct {
	collection *mongo.Collection
}

// NewItemRepository creates a new MongoDB item repository
func NewItemRepository(db *mongo.Database) ItemRepository {
	return &MongoItemRepository{
		collection: db.Collection("items"),
	}
}

// CreateItem inserts a new item into the database
func (r *MongoItemRepository) CreateItem(ctx context.Context, item *models.Item) error {
	result, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	item.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetItem retrieves a single item by ID
func (r *MongoItemRepository) GetItem(ctx context.Context, id string) (*models.Item, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var item models.Item
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetAllItems retrieves all items from the database
func (r *MongoItemRepository) GetAllItems(ctx context.Context) ([]*models.Item, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*models.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateItem updates an existing item
func (r *MongoItemRepository) UpdateItem(ctx context.Context, id string, item *models.Item) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":        item.Name,
			"description": item.Description,
			"updated_at":  item.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// DeleteItem removes an item from the database
func (r *MongoItemRepository) DeleteItem(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
