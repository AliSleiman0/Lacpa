# Authentication System - Quick Start Guide

## 🎯 What Was Implemented

A complete JWT-based authentication system with:

✅ **User Registration** - Signup with email verification via OTP
✅ **User Login** - Authenticate with LACPA ID and password  
✅ **Password Reset** - Full forgot password flow with OTP verification
✅ **JWT Tokens** - 24-hour expiring tokens for session management
✅ **Auth Middleware** - Protect routes and check user roles
✅ **Password Security** - bcrypt hashing with strength validation
✅ **Email Verification** - OTP-based account verification

## 📁 Files Created

### Models
- `Backend/models/user.go` - User struct and auth request/response models

### Utils
- `Backend/utils/jwt.go` - JWT token generation and validation
- `Backend/utils/password.go` - Password hashing and verification
- `Backend/utils/otp.go` - OTP and reset token generation
- `Backend/utils/validation.go` - Enhanced with ValidateStruct function

### Middleware
- `Backend/middleware/auth.go` - Auth, Role, and Optional Auth middleware

### Repository
- `Backend/repository/auth_repository.go` - Database operations for users

### Handler
- `Backend/handler/auth_handler.go` - All authentication endpoints

### Routes
- `Backend/routes/auth_routes.go` - Auth route registration

### Configuration
- `Backend/.env` - Updated with JWT_SECRET
- `Backend/.env.example` - Template for environment variables

### Documentation
- `Backend/AUTH_SYSTEM.md` - Complete API documentation
- `Backend/scripts/test_auth.js` - Comprehensive test suite
- `Backend/QUICKSTART_AUTH.md` - This file!

## 🚀 Quick Start

### 1. Install Dependencies

```powershell
cd Backend
go mod tidy
```

### 2. Set JWT Secret

The `.env` file has been updated with a JWT secret. For production, change it:

```env
JWT_SECRET=your-super-secret-production-key-here
```

### 3. Start the Server

```powershell
# In Backend directory
air
```

The server will run on `http://localhost:3000`

### 4. Test the Authentication

Run the automated test suite:

```powershell
cd Backend
node scripts/test_auth.js
```

Or test manually with the frontend pages:
- `http://localhost:3000/src/login.html`
- `http://localhost:3000/src/signup.html`
- `http://localhost:3000/src/forgot-password.html`
- `http://localhost:3000/src/verify-account.html`
- `http://localhost:3000/src/reset-password.html`

## 📡 API Endpoints

### Public Routes (No Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/signup` | Register new user |
| POST | `/api/auth/login` | Login with LACPA ID |
| POST | `/api/auth/forgot-password` | Request password reset |
| POST | `/api/auth/verify-otp` | Verify OTP code |
| POST | `/api/auth/resend-otp` | Resend OTP |
| POST | `/api/auth/reset-password` | Reset password with token |

### Protected Routes (Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/profile` | Get current user profile |
| POST | `/api/auth/logout` | Logout (client-side) |

## 🔐 Using the Middleware

### Protect a Single Route

```go
app.Get("/protected", middleware.AuthMiddleware, yourHandler)

// Access user info in handler
func yourHandler(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)
    role := c.Locals("role").(string)
    // ... your logic
}
```

### Protect Multiple Routes with Group

```go
api := app.Group("/api")
api.Use(middleware.AuthMiddleware)

api.Get("/dashboard", dashboardHandler)
api.Get("/profile", profileHandler)
```

### Role-Based Access Control

```go
// Admin only
app.Get("/admin", 
    middleware.AuthMiddleware, 
    middleware.RoleMiddleware("admin"), 
    adminHandler)

// Multiple roles
app.Get("/members", 
    middleware.AuthMiddleware, 
    middleware.RoleMiddleware("admin", "member"), 
    membersHandler)
```

### Optional Authentication

```go
// Available to everyone, but extracts user if token present
app.Get("/public", middleware.OptionalAuthMiddleware, publicHandler)
```

## 🧪 Example: Test with cURL

### 1. Signup

```bash
curl -X POST http://localhost:3000/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

Response:
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

**Important:** Check the server console for the OTP (until email service is integrated)

### 2. Verify OTP

```bash
curl -X POST http://localhost:3000/api/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "otp": "123456"
  }'
```

### 3. Login

```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "lacpa_id": "LACPA-2025-12345",
    "password": "SecurePass123!"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "...",
      "lacpa_id": "LACPA-2025-12345",
      "full_name": "John Doe",
      "email": "john@example.com",
      "role": "member",
      "is_verified": true,
      "is_active": true
    }
  }
}
```

### 4. Get Profile (with token)

```bash
curl -X GET http://localhost:3000/api/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 🔧 Frontend Integration

### Store Token After Login

```javascript
// In login.html
const response = await fetch('/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ lacpa_id, password })
});

const data = await response.json();
if (data.success) {
  localStorage.setItem('token', data.data.token);
  localStorage.setItem('user', JSON.stringify(data.data.user));
  window.location.href = '/dashboard.html';
}
```

### Add Token to Requests

```javascript
const token = localStorage.getItem('token');

fetch('/api/auth/profile', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
```

### Logout

```javascript
localStorage.removeItem('token');
localStorage.removeItem('user');
window.location.href = '/src/login.html';
```

## 📊 Database Schema

The `users` collection will be automatically created with this structure:

```javascript
{
  _id: ObjectId,
  lacpa_id: "LACPA-2025-12345",    // Auto-generated
  full_name: "John Doe",
  email: "john@example.com",
  password: "$2a$10$...",            // bcrypt hashed
  role: "member",                    // admin, member, guest
  is_verified: true,                 // Email verified
  is_active: true,                   // Account active
  otp: "123456",                     // Temporary OTP
  otp_expiry: ISODate,               // 10 minutes
  reset_token: "abc123...",          // Reset token
  reset_token_expiry: ISODate,       // 15 minutes
  last_login: ISODate,
  created_at: ISODate,
  updated_at: ISODate
}
```

## ⚙️ Configuration Details

### JWT Token
- **Expiration:** 24 hours
- **Algorithm:** HS256
- **Claims:** userID, lacpaID, email, role

### OTP
- **Length:** 6 digits
- **Expiration:** 10 minutes
- **Use cases:** Email verification, password reset

### Reset Token
- **Length:** 64 characters (32 bytes hex)
- **Expiration:** 15 minutes
- **Use case:** Password reset

### Password Requirements
- Minimum 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 digit
- At least 1 special character

## 📧 TODO: Email Integration

Currently, OTPs are printed to the console. To add email functionality:

### Option 1: SendGrid

```go
// Install: go get github.com/sendgrid/sendgrid-go
import "github.com/sendgrid/sendgrid-go"
import "github.com/sendgrid/sendgrid-go/helpers/mail"

func sendOTPEmail(toEmail, otp string) error {
    from := mail.NewEmail("LACPA", "noreply@lacpa.org")
    to := mail.NewEmail("", toEmail)
    subject := "Your LACPA Verification Code"
    content := mail.NewContent("text/html", 
        fmt.Sprintf("<h1>Your OTP is: %s</h1>", otp))
    
    m := mail.NewV3MailInit(from, subject, to, content)
    client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
    _, err := client.Send(m)
    return err
}
```

### Option 2: AWS SES

```go
// Install: go get github.com/aws/aws-sdk-go/service/ses
// Use AWS SES service to send emails
```

### Option 3: SMTP

```go
// Install: go get gopkg.in/gomail.v2
import "gopkg.in/gomail.v2"

func sendOTPEmail(toEmail, otp string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", "noreply@lacpa.org")
    m.SetHeader("To", toEmail)
    m.SetHeader("Subject", "Your LACPA Verification Code")
    m.SetBody("text/html", fmt.Sprintf("<h1>Your OTP is: %s</h1>", otp))
    
    d := gomail.NewDialer("smtp.gmail.com", 587, 
        os.Getenv("SMTP_USER"), 
        os.Getenv("SMTP_PASSWORD"))
    
    return d.DialAndSend(m)
}
```

Then replace `fmt.Printf` calls in `auth_handler.go` with your email function.

## 🐛 Troubleshooting

### "Invalid or expired token"
- Token expired (24 hours)
- JWT_SECRET was changed
- Token format is incorrect (must be `Bearer <token>`)

### "Please verify your email first"
- User hasn't completed OTP verification
- Check console for OTP during development

### "Invalid credentials"
- Wrong LACPA ID or password
- User doesn't exist

### OTP not working
- OTP expired (10 minutes)
- Use resend-otp endpoint to get new code

## 📚 Further Reading

For complete API documentation, see: `Backend/AUTH_SYSTEM.md`

For testing the system, see: `Backend/scripts/test_auth.js`

## ✨ Next Steps

1. **Integrate Email Service** - Replace console OTP with real emails
2. **Add Rate Limiting** - Prevent brute force attacks
3. **Add Session Management** - Track active sessions
4. **Add Token Blacklist** - Invalidate tokens on logout
5. **Add 2FA** - Two-factor authentication
6. **Add OAuth** - Google, GitHub, etc.
7. **Add Audit Logging** - Track auth events

## 🎉 You're Ready!

The authentication system is fully functional and ready to use. All frontend pages are connected and waiting for the backend routes we just created.

Start the server and test the complete flow:
1. Signup → Get OTP from console
2. Verify OTP
3. Login → Get JWT token
4. Access protected routes with token
5. Test password reset flow

Happy coding! 🚀
