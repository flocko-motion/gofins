#!/bin/bash
# Export database schema from local Docker to deployment directory

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOYMENT_DIR="$SCRIPT_DIR/../deployment"
SCHEMA_FILE="$DEPLOYMENT_DIR/schema.sql"

echo "=== Exporting database schema ==="

# Check if container exists and is running
if docker ps --format '{{.Names}}' | grep -q "^fins-postgres$"; then
    echo "✓ PostgreSQL container is running"
elif docker ps -a --format '{{.Names}}' | grep -q "^fins-postgres$"; then
    echo "Starting existing PostgreSQL container..."
    docker start fins-postgres
    sleep 2
else
    echo "ERROR: PostgreSQL container not found"
    echo "Start it with: ./run_db.sh start"
    exit 1
fi

# Export schema only (no data, no ownership)
echo "Exporting schema from local database..."
docker exec fins-postgres pg_dump -U fins -d fins --schema-only --no-owner --no-acl > "$SCHEMA_FILE"

# Replace local username 'fins' with production username 'gofins'
sed -i 's/\bfins\b/gofins/g' "$SCHEMA_FILE"

echo "✓ Schema exported to: $SCHEMA_FILE"
echo ""
echo "To import on the server:"
echo "  cat deployment/schema.sql | docker exec -i gofins-db psql -U gofins -d gofins"
