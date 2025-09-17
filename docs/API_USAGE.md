# API Usage Guide

## Login API

### Endpoint: `POST /api/login`

Authenticates a user with email and password.

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
