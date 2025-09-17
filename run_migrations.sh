#!/bin/bash

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Database connection parameters
DB_HOST=${BLUEPRINT_DB_HOST:-localhost}
DB_PORT=${BLUEPRINT_DB_PORT:-5432}
DB_NAME=${BLUEPRINT_DB_DATABASE:-ai_platform}
DB_USER=${BLUEPRINT_DB_USERNAME:-postgres}
DB_SCHEMA=${BLUEPRINT_DB_SCHEMA:-public}

# Check if password is provided
if [ -z "$BLUEPRINT_DB_PASSWORD" ]; then
    echo "Error: BLUEPRINT_DB_PASSWORD environment variable is not set"
    echo "Please set it in your .env file or export it directly"
    exit 1
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Running database migrations...${NC}"
echo "Host: $DB_HOST"
echo "Port: $DB_PORT"
echo "Database: $DB_NAME"
echo "User: $DB_USER"
echo "Schema: $DB_SCHEMA"
echo ""

# Check if migrations directory exists
if [ ! -d "migrations" ]; then
    echo -e "${RED}Error: migrations directory not found${NC}"
    exit 1
fi

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    echo -e "${RED}Error: psql command not found. Please install PostgreSQL client tools.${NC}"
    exit 1
fi

# Test database connection
echo -e "${YELLOW}Testing database connection...${NC}"
PGPASSWORD=$BLUEPRINT_DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1

if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Cannot connect to database. Please check your connection parameters.${NC}"
    exit 1
fi

echo -e "${GREEN}Database connection successful!${NC}"

# Create migrations table if it doesn't exist
echo -e "${YELLOW}Creating migrations tracking table...${NC}"
PGPASSWORD=$BLUEPRINT_DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME << EOF
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
EOF

if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to create migrations table${NC}"
    exit 1
fi

# Run migrations
for migration_file in migrations/*.sql; do
    if [ -f "$migration_file" ]; then
        # Extract version from filename (e.g., 001_create_users_table.sql -> 001)
        version=$(basename "$migration_file" | cut -d'_' -f1)

        # Check if migration has already been applied
        applied=$(PGPASSWORD=$BLUEPRINT_DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM schema_migrations WHERE version = '$version';" | xargs)

        if [ "$applied" -eq "0" ]; then
            echo -e "${YELLOW}Applying migration: $migration_file${NC}"

            # Run the migration
            PGPASSWORD=$BLUEPRINT_DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration_file"

            if [ $? -eq 0 ]; then
                # Record successful migration
                PGPASSWORD=$BLUEPRINT_DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "INSERT INTO schema_migrations (version) VALUES ('$version');"
                echo -e "${GREEN}Migration $version applied successfully!${NC}"
            else
                echo -e "${RED}Error: Migration $version failed!${NC}"
                exit 1
            fi
        else
            echo -e "${GREEN}Migration $version already applied, skipping...${NC}"
        fi
    fi
done

echo -e "${GREEN}All migrations completed successfully!${NC}"

# Show applied migrations
echo ""
echo -e "${YELLOW}Applied migrations:${NC}"
PGPASSWORD=$BLUEPRINT_DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT version, applied_at FROM schema_migrations ORDER BY version;"