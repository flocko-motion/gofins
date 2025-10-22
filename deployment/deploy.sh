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
    echo "Copy deployment/.env.example to deployment/.env and set DB_PASSWORD"
    exit 1
fi
export $(cat deployment/.env | xargs)

echo "=== Setting up Apache config ==="
if [ ! -L /etc/apache2/sites-available/gofins.conf ]; then
    echo "Creating Apache config symlink..."
    sudo ln -s /opt/gofins/deployment/apache-gofins.conf /etc/apache2/sites-available/gofins.conf
    sudo a2ensite gofins
    echo "Reloading Apache..."
    sudo systemctl reload apache2
else
    echo "Apache config already linked"
fi

echo "=== Setting up systemd service ==="
if [ ! -L /etc/systemd/system/gofins.service ]; then
    echo "Creating systemd service symlink..."
    sudo ln -s /opt/gofins/deployment/gofins.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl enable gofins
else
    echo "Systemd service already linked"
fi

echo "=== Building Docker images ==="
cd deployment
docker compose build --build-arg GIT_HASH=$GIT_HASH

echo "=== Restarting gofins service ==="
sudo systemctl daemon-reload
sudo systemctl restart gofins

echo "=== Waiting for services to start ==="
sleep 5

echo "=== Deployment complete ==="
echo ""
echo "Systemd service status:"
sudo systemctl status gofins --no-pager --lines=5
echo ""
echo "Docker containers:"
docker compose ps
echo ""
echo "Container logs (last 10 lines):"
echo "--- gofins-api ---"
docker logs --tail 10 gofins-api
echo ""
echo "--- gofins-ui ---"
docker logs --tail 10 gofins-ui
echo ""
echo "Services running on:"
echo "  - PostgreSQL: localhost:7701"
echo "  - API:        localhost:7702"
echo "  - UI:         localhost:7703"
echo "  - Public:     https://$(hostname -f)/gofins (via Apache)"
