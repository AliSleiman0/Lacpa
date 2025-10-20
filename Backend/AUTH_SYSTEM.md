# Authentication System Documentation

## Overview

Complete JWT-based authentication system for LACPA with the following features:

- User registration with email verification
- Login with LACPA ID and password
- Password reset via OTP
- JWT token-based authentication
- Role-based access control middleware
- Secure password hashing with bcrypt

## API Endpoints

### Public Endpoints (No Authentication Required)

#### 1. User Signup
```http
POST /api/auth/signup
Content-Type: application/json

{
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Registration successful. Please verify your email with the OTP sent.",
  "data": {
    "email": "john@example.com",
    "lacpa_id": "LACPA-2025-12345"
  }
}
```

#### 2. User Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "lacpa_id": "LACPA-2025-12345",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "507f1f77bcf86cd799439011",
      "lacpa_id": "LACPA-2025-12345",
      "full_name": "John Doe",
      "email": "john@example.com",
      "role": "member",
      "is_verified": true,
      "is_active": true,
      "created_at": "2025-10-20T10:00:00Z"
    }
  }
}
```

#### 3. Forgot Password
```http
POST /api/auth/forgot-password
Content-Type: application/json

{
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "If the email exists, a reset OTP has been sent."
}
```

#### 4. Verify OTP
```http
POST /api/auth/verify-otp
Content-Type: application/json

{
  "email": "john@example.com",
  "otp": "123456"
}
```

**Response:**
```json
{
  "success": true,
  "message": "OTP verified successfully",
  "data": {
    "reset_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"
  }
}
```

#### 5. Resend OTP
```http
POST /api/auth/resend-otp
Content-Type: application/json

{
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "If the email exists, a new OTP has been sent."
}
```

#### 6. Reset Password
```http
POST /api/auth/reset-password
Content-Type: application/json

{
  "token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6",
  "new_password": "NewSecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password reset successful. You can now login with your new password."
}
```

### Protected Endpoints (Authentication Required)

All protected endpoints require an `Authorization` header with a valid JWT token:

```http
Authorization: Bearer <your_jwt_token>
```

#### 7. Get Profile
```http
GET /api/auth/profile
Authorization: Bearer <your_jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "lacpa_id": "LACPA-2025-12345",
    "full_name": "John Doe",
    "email": "john@example.com",
    "role": "member",
    "is_verified": true,
    "is_active": true,
    "last_login": "2025-10-20T10:30:00Z",
    "created_at": "2025-10-20T10:00:00Z"
  }
}
```

#### 8. Logout
```http
POST /api/auth/logout
Authorization: Bearer <your_jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Logout successful"
}
```

## Middleware

### AuthMiddleware

Validates JWT tokens and extracts user information:

```go
import "lacpa/middleware"

// Protect a route
app.Get("/protected", middleware.AuthMiddleware, handler)

// Access user info in handler
userID := c.Locals("userID").(string)
lacpaID := c.Locals("lacpaID").(string)
email := c.Locals("email").(string)
role := c.Locals("role").(string)
```

### RoleMiddleware

Restricts access based on user roles:

```go
// Admin only route
app.Get("/admin", middleware.AuthMiddleware, middleware.RoleMiddleware("admin"), adminHandler)

// Multiple roles allowed
app.Get("/members", middleware.AuthMiddleware, middleware.RoleMiddleware("admin", "member"), membersHandler)
```

### OptionalAuthMiddleware

Validates token if present, but doesn't require it:

```go
app.Get("/public-with-context", middleware.OptionalAuthMiddleware, handler)
```

## User Roles

- `admin` - Full access to all resources
- `member` - Standard member access
- `guest` - Limited access

## Security Features

1. **Password Hashing**: bcrypt with default cost (10)
2. **JWT Tokens**: 24-hour expiration
3. **OTP**: 6-digit codes with 10-minute expiry
4. **Reset Tokens**: Secure 32-byte hex tokens with 15-minute expiry
5. **Email Verification**: Required before login
6. **Account Status**: Active/Inactive flag

## Environment Variables

Add to `.env` file:

```env
JWT_SECRET=your-super-secret-jwt-key-change-this
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=lacpa
PORT=3000
APP_ENV=development
```

## Database Schema

### Users Collection

```json
{
  "_id": "ObjectId",
  "lacpa_id": "LACPA-2025-12345",
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "$2a$10$...",
  "role": "member",
  "is_verified": true,
  "is_active": true,
  "otp": "123456",
  "otp_expiry": "2025-10-20T10:10:00Z",
  "reset_token": "a1b2c3...",
  "reset_token_expiry": "2025-10-20T10:15:00Z",
  "last_login": "2025-10-20T10:00:00Z",
  "created_at": "2025-10-20T09:00:00Z",
  "updated_at": "2025-10-20T10:00:00Z"
}
```

## Error Responses

All errors follow this format:

```json
{
  "success": false,
  "error": "Error message here"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid/missing token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `409` - Conflict (duplicate email)
- `500` - Internal Server Error

## Integration with Frontend

The authentication pages in `LACPA_Web/src/` are ready to connect:

1. **login.html** → `POST /api/auth/login`
2. **signup.html** → `POST /api/auth/signup`
3. **forgot-password.html** → `POST /api/auth/forgot-password`
4. **verify-account.html** → `POST /api/auth/verify-otp` and `POST /api/auth/resend-otp`
5. **reset-password.html** → `POST /api/auth/reset-password`

### Frontend Token Storage

```javascript
// Store token after login
localStorage.setItem('token', response.data.token);

// Add token to requests
fetch('/api/auth/profile', {
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('token')}`
  }
});

// Clear token on logout
localStorage.removeItem('token');
```

## Testing

Use the provided test script:

```bash
node Backend/scripts/test_auth.js
```

Or test manually with curl:

```bash
# Signup
curl -X POST http://localhost:3000/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Test User","email":"test@example.com","password":"Test123!"}'

# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"lacpa_id":"LACPA-2025-12345","password":"Test123!"}'

# Get Profile (with token)
curl -X GET http://localhost:3000/api/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## TODO: Email Service Integration

Currently, OTPs are printed to console. To integrate email service:

1. Install email package (e.g., SendGrid, AWS SES)
2. Add email templates
3. Update `auth_handler.go`:
   - Replace `fmt.Printf` in `Signup()` with email send
   - Replace `fmt.Printf` in `ForgotPassword()` with email send
   - Replace `fmt.Printf` in `ResendOTP()` with email send

## Notes

- LACPA IDs are auto-generated in format: `LACPA-YYYY-XXXXX`
- Default role for new users is `member`
- Users must verify email before login
- OTP is valid for 10 minutes
- Reset token is valid for 15 minutes
- JWT token is valid for 24 hours
