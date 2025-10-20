package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	LACPAID           string             `json:"lacpa_id" bson:"lacpa_id"`
	FullName          string             `json:"full_name" bson:"full_name"`
	Email             string             `json:"email" bson:"email"`
	Password          string             `json:"-" bson:"password"` // Never expose in JSON
	Role              string             `json:"role" bson:"role"`  // "admin", "member", "guest"
	IsVerified        bool               `json:"is_verified" bson:"is_verified"`
	IsActive          bool               `json:"is_active" bson:"is_active"`
	VerificationToken string             `json:"-" bson:"verification_token,omitempty"`
	ResetToken        string             `json:"-" bson:"reset_token,omitempty"`
	ResetTokenExpiry  time.Time          `json:"-" bson:"reset_token_expiry,omitempty"`
	OTP               string             `json:"-" bson:"otp,omitempty"`
	OTPExpiry         time.Time          `json:"-" bson:"otp_expiry,omitempty"`
	LastLogin         time.Time          `json:"last_login" bson:"last_login"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserResponse is the sanitized user response (without sensitive data)
type UserResponse struct {
	ID         primitive.ObjectID `json:"id"`
	LACPAID    string             `json:"lacpa_id"`
	FullName   string             `json:"full_name"`
	Email      string             `json:"email"`
	Role       string             `json:"role"`
	IsVerified bool               `json:"is_verified"`
	IsActive   bool               `json:"is_active"`
	LastLogin  time.Time          `json:"last_login"`
	CreatedAt  time.Time          `json:"created_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		LACPAID:    u.LACPAID,
		FullName:   u.FullName,
		Email:      u.Email,
		Role:       u.Role,
		IsVerified: u.IsVerified,
		IsActive:   u.IsActive,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
	}
}

// LoginRequest represents login credentials
type LoginRequest struct {
	LACPAID  string `json:"lacpa_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// SignupRequest represents signup data
type SignupRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// ForgotPasswordRequest represents forgot password data
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyOTPRequest represents OTP verification data
type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

// ResetPasswordRequest represents password reset data
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ResendOTPRequest represents resend OTP data
type ResendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
