#!/bin/bash

# =============================================
# Kasir API - Database Migration Script
# =============================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Default values
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-kasir}"

# Check if DB_CON is set (connection string)
if [ -n "$DB_CON" ]; then
    echo -e "${YELLOW}Using connection string from DB_CON${NC}"
    CONNECTION_STRING="$DB_CON"
else
    echo -e "${YELLOW}Using individual connection parameters${NC}"
    CONNECTION_STRING="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
fi

echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}   Kasir API - Database Migration${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""

# Function to run migration
run_migration() {
    local migration_file=$1
    echo -e "${YELLOW}Running migration: ${migration_file}${NC}"
    
    if command -v psql &> /dev/null; then
        PGPASSWORD=$DB_PASSWORD psql "$CONNECTION_STRING" -f "$migration_file"
    else
        # Try using docker
        echo -e "${YELLOW}psql not found, trying docker...${NC}"
        docker exec -i kasir-db psql -U postgres -d kasir < "$migration_file"
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Migration completed successfully${NC}"
    else
        echo -e "${RED}✗ Migration failed${NC}"
        exit 1
    fi
}

# Check command
case "${1:-up}" in
    up)
        echo -e "${GREEN}Running UP migrations...${NC}"
        for f in migrations/*.up.sql; do
            if [ -f "$f" ]; then
                run_migration "$f"
            fi
        done
        ;;
    down)
        echo -e "${RED}Running DOWN migrations (ROLLBACK)...${NC}"
        # Run down migrations in reverse order
        for f in $(ls -r migrations/*.down.sql 2>/dev/null); do
            if [ -f "$f" ]; then
                run_migration "$f"
            fi
        done
        ;;
    init)
        echo -e "${GREEN}Running init.sql...${NC}"
        run_migration "migrations/init.sql"
        ;;
    *)
        echo "Usage: $0 {up|down|init}"
        echo "  up   - Run all UP migrations"
        echo "  down - Run all DOWN migrations (rollback)"
        echo "  init - Run the init.sql script"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}   Migration completed!${NC}"
echo -e "${GREEN}======================================${NC}"
