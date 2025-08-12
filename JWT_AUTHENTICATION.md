# JWT Authentication for Hydration APIs

This document explains how JWT authentication has been implemented for all hydration APIs using CNumber and UserName.

## Overview

The JWT authentication system has been updated to use:
- **CNumber**: User's contact number (unique identifier)
- **UserName**: User's display name

Instead of the previous UserID-based system.

## JWT Token Structure

### Claims
```json
{
  "c_number": "1234567890",
  "username": "John Doe",
  "exp": 1735689600,
  "iat": 1733097600,
  "nbf": 1733097600
}
```

### Token Expiration
- **Duration**: 30 days from generation
- **Format**: JWT with HS256 signing algorithm

## Authentication Flow

### 1. Login
```http
POST /Services/innovologin
Content-Type: application/json

{
  "cnumber": "1234567890",
  "userpin": "user_password"
}
```

**Response:**
```json
{
  "code": 0,
  "message": "OK",
  "userid": 123,
  "userdetails": {...},
  "jwt_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. Using Protected APIs
Include the JWT token in the Authorization header:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Protected Hydration Endpoints

All hydration APIs now require JWT authentication and automatically use the authenticated user's information:

### Core Hydration APIs
- `POST /Services/protected/innovoHyderation` - Record basic hydration data
- `POST /Services/protected/newinnovoHyderation` - Record enhanced hydration data
- `POST /Services/protected/updateHyderationValue` - Update hydration data
- `POST /Services/protected/updateSweatData` - Update sweat data

### Report APIs
- `POST /Services/protected/getSummary` - Get summary data
- `POST /Services/protected/getUserDetailedSummary` - Get detailed summary
- `POST /Services/protected/getClientHistory` - Get client history
- `POST /Services/protected/getHyderartionHistory` - Get hydration history
- `POST /Services/protected/getElectrolyteHistory` - Get electrolyte history

## Security Features

### Automatic User Identification
- User ID is automatically extracted from JWT claims
- No need to pass user ID in request body for most endpoints
- Prevents users from accessing other users' data

### Token Validation
- JWT tokens are validated on every protected request
- Expired tokens are automatically rejected
- Invalid tokens return 401 Unauthorized

### Context Injection
The middleware automatically injects user information into the request context:
- `jwt_claims`: Full JWT claims object
- `user_cnumber`: User's contact number
- `username`: User's display name

## Implementation Details

### JWT Service (`services/jwt_service.go`)
- `GenerateToken(cNumber, userName)`: Creates JWT token
- `ValidateToken(tokenString)`: Validates and parses JWT token

### Middleware (`middleware/jwt_auth.go`)
- `JWTAuthMiddleware()`: Validates JWT tokens
- `GetJWTClaimsFromContext()`: Retrieves claims from context
- `GetUserCNumberFromJWTContext()`: Gets CNumber from context
- `GetUserNameFromJWTContext()`: Gets UserName from context

### User Service (`services/user_service.go`)
- `GetUserIDByCNumber(cnumber)`: Converts CNumber to UserID for database operations

## Example Usage

### 1. Login and Get Token
```bash
curl -X POST http://localhost:8080/Services/innovologin \
  -H "Content-Type: application/json" \
  -d '{"cnumber": "1234567890", "userpin": "password123"}'
```

### 2. Use Token for Hydration API
```bash
curl -X POST http://localhost:8080/Services/protected/innovoHyderation \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "weight": 70.5,
    "height": 175.0,
    "sweat_position": 25.0,
    "time_taken": 45.0,
    "device_type": 1,
    "image_path": "path/to/image.jpg"
  }'
```

## Error Handling

### Authentication Errors
- **401 Unauthorized**: Missing or invalid JWT token
- **Token Expired**: JWT token has expired
- **Invalid Format**: Malformed Authorization header

### User Errors
- **User Not Found**: CNumber in JWT doesn't exist in database
- **Account Deleted**: User account has been deactivated

## Migration Notes

### For Existing Users
- Existing login functionality remains the same
- JWT tokens now contain CNumber and UserName instead of UserID
- All protected endpoints automatically use authenticated user's data

### For API Consumers
- No changes required in request payloads for most endpoints
- User ID is automatically determined from JWT token
- Enhanced security with automatic user isolation

## Environment Variables

Set the following environment variable for JWT secret:
```bash
export JWT_SECRET="your-secret-key-here"
```

If not set, a default development key will be used (not recommended for production).
