package repository

import (
	"context"
	"time"

	"github.com/AliSleiman0/Lacpa/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository defines the interface for user database operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, id string, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
	ActivateUser(ctx context.Context, id string) error
	DeactivateUser(ctx context.Context, id string) error
}

// MongoUserRepository implements UserRepository interface for MongoDB
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new MongoDB user repository
func NewUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

// CreateUser inserts a new user into the database
func (r *MongoUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetUser retrieves a single user by ID
func (r *MongoUserRepository) GetUser(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email address
func (r *MongoUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *MongoUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database
func (r *MongoUserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser updates an existing user
func (r *MongoUserRepository) UpdateUser(ctx context.Context, id string, user *models.User) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"is_active":  user.IsActive,
			"updated_at": user.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// DeleteUser removes a user from the database
func (r *MongoUserRepository) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// ActivateUser sets a user's is_active field to true
func (r *MongoUserRepository) ActivateUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":  true,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// DeactivateUser sets a user's is_active field to false
func (r *MongoUserRepository) DeactivateUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}
