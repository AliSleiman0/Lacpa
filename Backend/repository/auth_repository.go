package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/AliSleiman0/Lacpa/models"
)

type AuthRepository struct {
	collection *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) *AuthRepository {
	return &AuthRepository{
		collection: db.Collection("users"),
	}
}

// CreateUser creates a new user
func (r *AuthRepository) CreateUser(user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true
	user.IsVerified = false

	_, err := r.collection.InsertOne(context.Background(), user)
	return err
}

// GetUserByEmail retrieves a user by email
func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByLACPAID retrieves a user by LACPA ID
func (r *AuthRepository) GetUserByLACPAID(lacpaID string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"lacpa_id": lacpaID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *AuthRepository) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information
func (r *AuthRepository) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

// SetOTP sets OTP and expiry for a user
func (r *AuthRepository) SetOTP(email, otp string, expiry time.Time) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"email": email},
		bson.M{
			"$set": bson.M{
				"otp":        otp,
				"otp_expiry": expiry,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// ClearOTP clears OTP fields for a user
func (r *AuthRepository) ClearOTP(email string) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"email": email},
		bson.M{
			"$unset": bson.M{
				"otp":        "",
				"otp_expiry": "",
			},
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// SetResetToken sets reset token and expiry for a user
func (r *AuthRepository) SetResetToken(email, token string, expiry time.Time) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"email": email},
		bson.M{
			"$set": bson.M{
				"reset_token":        token,
				"reset_token_expiry": expiry,
				"updated_at":         time.Now(),
			},
		},
	)
	return err
}

// GetUserByResetToken retrieves a user by reset token
func (r *AuthRepository) GetUserByResetToken(token string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(
		context.Background(),
		bson.M{
			"reset_token": token,
			"reset_token_expiry": bson.M{
				"$gt": time.Now(),
			},
		},
	).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ClearResetToken clears reset token fields for a user
func (r *AuthRepository) ClearResetToken(userID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{
			"$unset": bson.M{
				"reset_token":        "",
				"reset_token_expiry": "",
			},
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// UpdatePassword updates user password
func (r *AuthRepository) UpdatePassword(userID primitive.ObjectID, hashedPassword string) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"password":   hashedPassword,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// UpdateLastLogin updates the last login time
func (r *AuthRepository) UpdateLastLogin(userID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"last_login": time.Now(),
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// VerifyUser marks user as verified
func (r *AuthRepository) VerifyUser(email string) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"email": email},
		bson.M{
			"$set": bson.M{
				"is_verified": true,
				"updated_at":  time.Now(),
			},
		},
	)
	return err
}
