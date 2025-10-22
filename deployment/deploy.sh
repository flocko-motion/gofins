#!/bin/bash
set -e

cd /opt/gofins

echo "=== Pulling latest code ==="
git pull

echo "=== Getting git hash ==="
export GIT_HASH=$(git rev-parse --short HEAD)
echo "Building with tag: $GIT_HASH"

echo "=== Loading environment ==="
if [ ! -f deployment/.env ]; then
    echo "ERROR: deployment/.env not found!"
    exit 1
fi
export $(cat deployment/.env | xargs)

echo "=== Building Docker images ==="
cd deployment
docker-compose build --build-arg GIT_HASH=$GIT_HASH

echo "=== Starting containers ==="
docker-compose up -d

echo "=== Deployment complete ==="
docker-compose ps
