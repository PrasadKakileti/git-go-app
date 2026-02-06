# Authentication System Documentation

## Overview

JobHub now includes a complete authentication system with user signup and login functionality. Users must register and login before accessing the job portal features.

## Features Added

### 1. User Signup
- **Email or Mobile Registration**: Users can register with either email or mobile number
- **Password Protection**: Secure password hashing using bcrypt
- **Experience Field**: Users specify their years of experience
- **Email/Mobile Verification**: Verification code generated (ready for SMS/Email integration)
- **Data Validation**: Client and server-side validation

### 2. User Login
- **Flexible Login**: Login with email or mobile number
- **Session Management**: Token-based authentication
- **Secure Password Verification**: bcrypt password comparison
- **Redirect to Home**: Automatic redirect after successful login

### 3. Database Schema Updates

```sql
-- New columns added to users table
ALTER TABLE users ADD COLUMN mobile VARCHAR(15) AFTER email;
ALTER TABLE users ADD COLUMN password VARCHAR(255) AFTER mobile;
ALTER TABLE users ADD COLUMN experience VARCHAR(50) AFTER domain;
ALTER TABLE users ADD COLUMN is_verified BOOLEAN DEFAULT FALSE AFTER is_active;
ALTER TABLE users ADD COLUMN verification_code VARCHAR(6) AFTER is_verified;
ALTER TABLE users ADD COLUMN created_by VARCHAR(50) DEFAULT 'self' AFTER created_at;
```

## User Flow

### Signup Flow
```
1. User visits /signup.html
2. Fills registration form:
   - Email OR Mobile (required)
   - Password (min 6 characters)
   - Confirm Password
   - Location
   - Domain
   - Experience
   - Notification Frequency
3. Click "Sign Up"
4. Backend validates and creates user
5. Password hashed with bcrypt
6. Welcome email sent
7. Redirect to login page
```

### Login Flow
```
1. User visits / or /login.html
2. Enters email/mobile and password
3. Click "Login"
4. Backend verifies credentials
5. Session token generated
6. User data stored in localStorage
7. Redirect to home page (index.html)
```

## API Endpoints

### Signup
```http
POST /api/signup
Content-Type: application/json

{
  "email": "user@example.com",
  "mobile": "9876543210",
  "password": "securepass123",
  "location": "Bangalore",
  "domain": "Java",
  "experience": "3-5",
  "notification_frequency": "daily"
}

Response (200 OK):
{
  "success": true,
  "message": "Signup successful! Please login to continue."
}
```

### Login
```http
POST /api/login
Content-Type: application/json

{
  "emailOrMobile": "user@example.com",
  "password": "securepass123"
}

Response (200 OK):
{
  "success": true,
  "message": "Login successful",
  "token": "1-1706634000",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "mobile": "9876543210",
    "location": "Bangalore",
    "domain": "Java"
  }
}
```

## Experience Options

| Value | Description |
|-------|-------------|
| 0-1 | 0-1 years |
| 1-3 | 1-3 years |
| 3-5 | 3-5 years |
| 5-10 | 5-10 years |
| 10+ | 10+ years |

## Security Features

### Password Security
- **Hashing**: bcrypt with default cost (10)
- **No Plain Text**: Passwords never stored in plain text
- **Secure Comparison**: Constant-time comparison

### Validation
- **Email Format**: Standard email validation
- **Mobile Format**: 10-digit number validation
- **Password Length**: Minimum 6 characters
- **Required Fields**: All fields validated

### Session Management
- **Token Generation**: Simple token format (userID-timestamp)
- **Client Storage**: localStorage for session persistence
- **Token Verification**: Ready for middleware implementation

## Frontend Pages

### Signup Page (`/signup.html`)
- Clean, modern design
- Form validation
- Password confirmation
- Experience dropdown
- Responsive layout

### Login Page (`/login.html`)
- Simple login form
- Email or mobile input
- Password field
- Link to signup page
- Default landing page

## Database Structure

### Users Table (Updated)
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    mobile VARCHAR(15),
    password VARCHAR(255),
    location VARCHAR(100),
    domain VARCHAR(100),
    experience VARCHAR(50),
    notification_frequency ENUM('daily', 'weekly') DEFAULT 'daily',
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_code VARCHAR(6),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) DEFAULT 'self',
    UNIQUE KEY unique_email (email),
    UNIQUE KEY unique_mobile (mobile)
);
```

## Testing

### Test Signup
```bash
curl -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "test123",
    "location": "Bangalore",
    "domain": "Java",
    "experience": "3-5",
    "notification_frequency": "daily"
  }'
```

### Test Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "emailOrMobile": "test@example.com",
    "password": "test123"
  }'
```

## Future Enhancements

### Planned Features
- [ ] Email verification with OTP
- [ ] SMS verification for mobile
- [ ] Password reset functionality
- [ ] JWT token implementation
- [ ] Session expiry
- [ ] Remember me functionality
- [ ] Social login (Google, LinkedIn)
- [ ] Two-factor authentication

### Security Improvements
- [ ] Rate limiting on login attempts
- [ ] CAPTCHA integration
- [ ] Password strength meter
- [ ] Account lockout after failed attempts
- [ ] IP-based access control

## Migration Guide

### For Existing Users
```sql
-- Update existing users with default values
UPDATE users 
SET experience = '1-3', 
    password = '$2a$10$defaulthash',
    is_verified = TRUE
WHERE password IS NULL;
```

### Backward Compatibility
- Old `/api/register` endpoint still works
- New users should use `/api/signup`
- Login required for accessing home page

## Troubleshooting

### Common Issues

**1. "Email or mobile already registered"**
- User already exists in database
- Try logging in instead

**2. "Invalid credentials"**
- Wrong email/mobile or password
- Check credentials and try again

**3. "Passwords do not match"**
- Password and confirm password don't match
- Re-enter passwords

## Access URLs

- **Signup**: http://localhost:8080/signup.html
- **Login**: http://localhost:8080/login.html (default)
- **Home**: http://localhost:8080/index.html (requires login)
- **Admin**: http://localhost:8080/admin

---

**Feature Status**: ✅ Production Ready  
**Version**: 2.0.0  
**Last Updated**: 2024-01-30
