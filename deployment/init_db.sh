#!/bin/bash
# Initialize database schema for production deployment

set -e

cd "$(dirname "$0")"

echo "=== Initializing database schema ==="

# Check if database container is running
if ! docker ps --format '{{.Names}}' | grep -q "^gofins-db$"; then
    echo "Starting database container..."
    docker compose up -d db
    sleep 3
fi

echo "Waiting for database to be ready..."
timeout=30
counter=0
while ! docker exec gofins-db pg_isready -U gofins -d gofins >/dev/null 2>&1; do
    sleep 1
    counter=$((counter + 1))
    if [ $counter -ge $timeout ]; then
        echo "ERROR: Database failed to start within $timeout seconds"
        exit 1
    fi
done
echo "✓ Database is ready"

echo "Checking database schema..."
table_count=$(docker exec gofins-db psql -U gofins -d gofins -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';" | tr -d ' ')

if [ "$table_count" -eq "0" ]; then
    echo "No tables found, importing schema..."
    if [ -f "schema.sql" ]; then
        cat schema.sql | docker exec -i gofins-db psql -U gofins -d gofins
        echo "✓ Schema imported successfully"
    else
        echo "ERROR: schema.sql not found in deployment/"
        echo "Run: cd ../development && ./db_export_schema.sh"
        exit 1
    fi
else
    echo "✓ Database already has $table_count tables"
    echo ""
    read -p "Reimport schema? This will DROP all existing data! (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Dropping existing schema..."
        docker exec gofins-db psql -U gofins -d gofins -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
        echo "Importing schema..."
        cat schema.sql | docker exec -i gofins-db psql -U gofins -d gofins
        echo "✓ Schema reimported successfully"
    else
        echo "Skipping schema import"
    fi
fi

echo ""
echo "=== Database initialization complete ==="
