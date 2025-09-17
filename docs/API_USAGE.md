# API Usage Guide

## Authentication

This API uses JWT (JSON Web Token) Bearer authentication. After logging in, you'll receive a token that must be included in the `Authorization` header for protected endpoints.

## Login API

### Endpoint: `POST /api/login`

Authenticates a user with email and password. Returns a JWT token for accessing protected endpoints.

#### Request Body

```json
{
    "email": "user@example.com",
    "password": "your_password"
}
```

### Setting up users

Since there's no register endpoint yet, you'll need to manually add users to the database.

1. Run the migrations to create and configure the users table:

```bash
./run_migrations.sh
```

2. Insert a test user with a hashed password:

```sql
-- Example: Insert a user with email 'test@example.com' and password 'password123'
-- The password hash below corresponds to 'password123'
-- Note: User IDs are UUIDs that are automatically generated
INSERT INTO users (email, password)
VALUES ('test@example.com', '$2a$12$k2WRsfc9868pKseoXaGAf.YdtXrp8uXumJiWoTxq1UxBWQ5m0df96');
```

**Note:** User IDs are now UUIDs that are automatically generated. The database will assign a unique UUID to each user when inserted.

You can generate password hashes using any bcrypt tool or by creating a simple Go script with the `golang.org/x/crypto/bcrypt` package.

## Protected Endpoints

All API endpoints (except `/api/login`) require a valid JWT token in the Authorization header.

### Header Format

```
Authorization: Bearer <your-jwt-token>
```

### Example Protected Endpoint: `GET /api/profile`

Returns the current user's profile information.

#### Request Headers

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Success Response (200 OK)

```json
{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "message": "This is a protected endpoint"
}
```

#### Error Responses

**401 Unauthorized** - Missing or invalid token

```json
{
    "error": "Authorization header is required"
}
```

```json
{
    "error": "Authorization header must start with 'Bearer '"
}
```

```json
{
    "error": "Invalid or expired token"
}
```

## Usage Example

```bash
# 1. Login to get a token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'

# Response:
# {
#   "user": {"id": "...", "email": "test@example.com"},
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "message": "Login successful"
# }

# 2. Use the token to access protected endpoints
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Token Details

- **Expiration**: 24 hours from issue time
- **Algorithm**: HS256
- **Claims**: Contains user ID and email
- **Issuer**: ai-platform
