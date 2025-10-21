# Authentication System Testing Checklist

## 🧪 Pre-Test Setup

- [ ] Server is running (`cd Backend && air`)
- [ ] MongoDB is running (port 27017)
- [ ] Database name is "LACPA" (check .env)
- [ ] JWT_SECRET is set in .env
- [ ] Console is visible (for OTP codes)

## ✅ Test Suite

### 1. User Registration (Signup)

#### Test 1.1: Successful Signup
- [ ] Navigate to http://localhost:3000/src/signup.html
- [ ] Fill in valid data:
  - Full Name: "Test User"
  - Email: "test@example.com"
  - Password: "Test123!@#"
  - Confirm Password: "Test123!@#"
- [ ] Password strength shows "Strong" (green)
- [ ] Click "Sign Up"
- [ ] **Expected:** Success message, redirected to login
- [ ] **Expected:** OTP printed in server console
- [ ] **Expected:** LACPA ID created (format: LACPA-2025-XXXXX)

#### Test 1.2: Duplicate Email
- [ ] Try to signup again with same email
- [ ] **Expected:** Error "Email already registered"

#### Test 1.3: Weak Password
- [ ] Try password: "weak"
- [ ] **Expected:** Password strength shows "Weak" (red)
- [ ] **Expected:** Button disabled or validation error

#### Test 1.4: Password Mismatch
- [ ] Password: "Test123!@#"
- [ ] Confirm: "Different123!@#"
- [ ] **Expected:** Error message "Passwords do not match"

#### Test 1.5: Invalid Email
- [ ] Email: "notanemail"
- [ ] **Expected:** Email validation error

---

### 2. Email Verification (OTP)

#### Test 2.1: Successful OTP Verification
- [ ] Copy OTP from server console (after signup)
- [ ] Navigate to verify-account page (or redirected)
- [ ] Enter the 6-digit OTP
- [ ] **Expected:** All 6 boxes filled
- [ ] **Expected:** "Verify" button enabled
- [ ] Click "Verify"
- [ ] **Expected:** Success message
- [ ] **Expected:** Account is now verified
- [ ] **Expected:** Redirected to login page

#### Test 2.2: Invalid OTP
- [ ] Enter wrong OTP: "999999"
- [ ] **Expected:** Error "Invalid OTP"
- [ ] **Expected:** Shake animation on inputs

#### Test 2.3: Expired OTP
- [ ] Wait 11+ minutes (OTP expires after 10 min)
- [ ] Try to verify
- [ ] **Expected:** Error "OTP expired or invalid"

#### Test 2.4: Resend OTP
- [ ] Click "Resend Code"
- [ ] **Expected:** Countdown timer starts (60s)
- [ ] **Expected:** New OTP printed in console
- [ ] **Expected:** Can verify with new OTP

#### Test 2.5: Auto-focus Navigation
- [ ] Enter digits one by one
- [ ] **Expected:** Auto-focus to next box
- [ ] Press Backspace
- [ ] **Expected:** Auto-focus to previous box

#### Test 2.6: Paste 6-digit Code
- [ ] Copy "123456" to clipboard
- [ ] Paste in first input box
- [ ] **Expected:** All 6 boxes filled automatically
- [ ] **Expected:** "Verify" button enabled

---

### 3. User Login

#### Test 3.1: Login Before Verification
- [ ] Try to login with unverified account
- [ ] **Expected:** Error "Please verify your email first"

#### Test 3.2: Successful Login
- [ ] Navigate to http://localhost:3000/src/login.html
- [ ] Enter LACPA ID: "LACPA-2025-12345" (from signup)
- [ ] Enter Password: "Test123!@#"
- [ ] Click "Login"
- [ ] **Expected:** Success message
- [ ] **Expected:** JWT token received
- [ ] **Expected:** Token stored in localStorage
- [ ] **Expected:** User info stored in localStorage
- [ ] **Expected:** Redirected to dashboard

#### Test 3.3: Invalid LACPA ID
- [ ] Enter LACPA ID: "INVALID-ID"
- [ ] Enter correct password
- [ ] **Expected:** Error "Invalid credentials"

#### Test 3.4: Invalid Password
- [ ] Enter correct LACPA ID
- [ ] Enter password: "WrongPassword"
- [ ] **Expected:** Error "Invalid credentials"

#### Test 3.5: Empty Fields
- [ ] Leave fields empty
- [ ] Try to submit
- [ ] **Expected:** Validation errors

---

### 4. Protected Routes

#### Test 4.1: Access Profile Without Token
- [ ] Clear localStorage (delete token)
- [ ] Try to access: GET http://localhost:3000/api/auth/profile
- [ ] **Expected:** 401 Unauthorized
- [ ] **Expected:** Error "Missing authorization header"

#### Test 4.2: Access Profile With Invalid Token
- [ ] Set Authorization header: "Bearer invalid-token"
- [ ] Try to access profile endpoint
- [ ] **Expected:** 401 Unauthorized
- [ ] **Expected:** Error "Invalid or expired token"

#### Test 4.3: Access Profile With Valid Token
- [ ] Login first (get valid token)
- [ ] GET http://localhost:3000/api/auth/profile
- [ ] Header: Authorization: Bearer {token}
- [ ] **Expected:** 200 OK
- [ ] **Expected:** User profile data returned

#### Test 4.4: Token Expiry
- [ ] Wait 24+ hours (or modify JWT expiry for testing)
- [ ] Try to access protected route
- [ ] **Expected:** 401 Unauthorized
- [ ] **Expected:** Error "Invalid or expired token"

---

### 5. Password Reset

#### Test 5.1: Request Password Reset
- [ ] Navigate to http://localhost:3000/src/forgot-password.html
- [ ] Enter email: "test@example.com"
- [ ] Click "Send Reset Link"
- [ ] **Expected:** Success message
- [ ] **Expected:** OTP printed in console
- [ ] **Expected:** Redirected to verify page

#### Test 5.2: Verify Reset OTP
- [ ] Copy OTP from console
- [ ] Enter in verify-account page
- [ ] Click "Verify"
- [ ] **Expected:** Success
- [ ] **Expected:** Reset token returned
- [ ] **Expected:** Redirected to reset-password page with token in URL

#### Test 5.3: Reset Password
- [ ] Navigate to reset-password page (with token)
- [ ] Enter new password: "NewPass123!@#"
- [ ] Confirm new password: "NewPass123!@#"
- [ ] **Expected:** All 4 requirements turn green:
  - ✓ At least 8 characters
  - ✓ Lowercase letter
  - ✓ Uppercase letter
  - ✓ Number or special character
- [ ] **Expected:** Password strength indicator shows "Strong"
- [ ] **Expected:** "Passwords match" message (green)
- [ ] Click "Reset Password"
- [ ] **Expected:** Success message
- [ ] **Expected:** Redirected to login page

#### Test 5.4: Login With New Password
- [ ] Login with LACPA ID
- [ ] Use new password: "NewPass123!@#"
- [ ] **Expected:** Successful login
- [ ] **Expected:** JWT token received

#### Test 5.5: Old Password No Longer Works
- [ ] Logout
- [ ] Try to login with old password: "Test123!@#"
- [ ] **Expected:** Error "Invalid credentials"

#### Test 5.6: Invalid Reset Token
- [ ] Manually edit token in URL to invalid value
- [ ] Try to reset password
- [ ] **Expected:** Error "Invalid or expired reset token"

#### Test 5.7: Expired Reset Token
- [ ] Wait 16+ minutes after getting token
- [ ] Try to reset password
- [ ] **Expected:** Error "Invalid or expired reset token"

#### Test 5.8: Non-existent Email
- [ ] Request password reset for "nonexistent@example.com"
- [ ] **Expected:** Generic success message (don't reveal if email exists)

---

### 6. Role-Based Access Control

#### Test 6.1: Member Access
- [ ] Login as member role user
- [ ] Try to access member-only route
- [ ] **Expected:** 200 OK, access granted

#### Test 6.2: Admin Access Required
- [ ] Login as member role user
- [ ] Try to access admin-only route
- [ ] **Expected:** 403 Forbidden
- [ ] **Expected:** Error "Insufficient permissions"

#### Test 6.3: Multiple Roles Allowed
- [ ] Login as any allowed role
- [ ] Access route with multiple allowed roles
- [ ] **Expected:** 200 OK, access granted

---

### 7. Frontend Validation

#### Test 7.1: Password Strength Indicator
- [ ] Type "weak" → **Expected:** 1 bar, red, "Weak"
- [ ] Type "Better1" → **Expected:** 2 bars, orange, "Fair"
- [ ] Type "Better1Pass" → **Expected:** 3 bars, blue, "Good"
- [ ] Type "Better1Pass!" → **Expected:** 4 bars, green, "Strong"

#### Test 7.2: Real-time Email Validation
- [ ] Type "invalid" → **Expected:** Red border
- [ ] Type "valid@email.com" → **Expected:** Green border

#### Test 7.3: Real-time Password Match
- [ ] Password: "Test123!@#"
- [ ] Confirm: "Test123!@#"
- [ ] **Expected:** Green border, "✓ Passwords match"
- [ ] Change confirm to: "Different"
- [ ] **Expected:** Red border, "✗ Passwords do not match"

#### Test 7.4: Loading States
- [ ] Click submit button
- [ ] **Expected:** Button shows spinner
- [ ] **Expected:** Button disabled during request
- [ ] **Expected:** Returns to normal after response

#### Test 7.5: Error Messages
- [ ] Trigger various errors
- [ ] **Expected:** Error message shown in red box
- [ ] **Expected:** Can be dismissed

---

### 8. Security Tests

#### Test 8.1: XSS Prevention
- [ ] Try to signup with name: "<script>alert('xss')</script>"
- [ ] **Expected:** Script not executed
- [ ] **Expected:** Stored safely or sanitized

#### Test 8.2: SQL Injection (N/A for MongoDB)
- [ ] Try LACPA ID: "' OR '1'='1"
- [ ] **Expected:** Treated as literal string, no injection

#### Test 8.3: Token in URL
- [ ] Never send JWT in URL parameters
- [ ] **Expected:** Token only in Authorization header
- [ ] Exception: Reset token in query param (one-time use)

#### Test 8.4: CORS Policy
- [ ] Try to access API from different origin
- [ ] **Expected:** Proper CORS headers present
- [ ] **Expected:** Preflight requests handled

---

### 9. Edge Cases

#### Test 9.1: Very Long Input
- [ ] Try 1000-character email
- [ ] **Expected:** Validation error or max length enforced

#### Test 9.2: Special Characters
- [ ] Try email: "test+tag@example.com"
- [ ] **Expected:** Accepted (valid email format)

#### Test 9.3: Case Sensitivity
- [ ] Signup with "Test@Example.com"
- [ ] Login with "test@example.com"
- [ ] **Expected:** Works (emails lowercase in DB)

#### Test 9.4: Unicode in Name
- [ ] Name: "José García 李明"
- [ ] **Expected:** Accepted and stored correctly

#### Test 9.5: Multiple Simultaneous Requests
- [ ] Click submit button multiple times rapidly
- [ ] **Expected:** Only one request sent (button disabled)

---

### 10. Automated Test Suite

#### Test 10.1: Run Full Test Suite
- [ ] Server running
- [ ] Run: `node Backend/scripts/test_auth.js`
- [ ] **Expected:** All 12 tests pass
- [ ] **Expected:** No errors in console
- [ ] **Expected:** Success rate: 100%

---

## 🔍 Manual Testing Commands

### Using cURL

```bash
# 1. Signup
curl -X POST http://localhost:3000/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Test User","email":"test@test.com","password":"Test123!@#"}'

# 2. Verify OTP (use OTP from console)
curl -X POST http://localhost:3000/api/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","otp":"123456"}'

# 3. Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"lacpa_id":"LACPA-2025-12345","password":"Test123!@#"}'

# 4. Get Profile (replace TOKEN)
curl http://localhost:3000/api/auth/profile \
  -H "Authorization: Bearer TOKEN"

# 5. Forgot Password
curl -X POST http://localhost:3000/api/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com"}'

# 6. Reset Password (replace TOKEN)
curl -X POST http://localhost:3000/api/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"token":"TOKEN","new_password":"NewPass123!@#"}'
```

### Using PowerShell (Invoke-RestMethod)

```powershell
# Signup
$body = @{
    full_name = "Test User"
    email = "test@test.com"
    password = "Test123!@#"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:3000/api/auth/signup" `
  -Method POST -Body $body -ContentType "application/json"

# Login
$body = @{
    lacpa_id = "LACPA-2025-12345"
    password = "Test123!@#"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:3000/api/auth/login" `
  -Method POST -Body $body -ContentType "application/json"
```

---

## 📊 Test Results Template

```
Date: _______________
Tester: _______________
Environment: Development / Staging / Production

┌────────────────────────────┬──────┬──────────┐
│ Test Category              │ Pass │ Comments │
├────────────────────────────┼──────┼──────────┤
│ 1. User Registration       │ □    │          │
│ 2. Email Verification      │ □    │          │
│ 3. User Login              │ □    │          │
│ 4. Protected Routes        │ □    │          │
│ 5. Password Reset          │ □    │          │
│ 6. Role-Based Access       │ □    │          │
│ 7. Frontend Validation     │ □    │          │
│ 8. Security Tests          │ □    │          │
│ 9. Edge Cases              │ □    │          │
│ 10. Automated Tests        │ □    │          │
└────────────────────────────┴──────┴──────────┘

Overall Status: ⬜ Pass  ⬜ Fail  ⬜ Partial

Issues Found:
1. _________________________________
2. _________________________________
3. _________________________________

Notes:
_________________________________________
_________________________________________
_________________________________________
```

---

## 🐛 Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| OTP not appearing in console | Check server logs, verify OTP generation code |
| Token always invalid | Check JWT_SECRET, ensure same secret used for sign & verify |
| Can't access protected routes | Verify token in localStorage, check Authorization header format |
| Password reset not working | Check reset token in URL, ensure not expired |
| Email validation failing | Check regex pattern, try different email formats |
| MongoDB connection error | Verify MongoDB running on port 27017 |

---

**Testing Status:** Ready for Testing ✅

**Last Updated:** October 20, 2025

**Next Steps:** 
1. Start server
2. Run automated test suite
3. Manually test frontend flows
4. Report any issues
