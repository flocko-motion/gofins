# Installation Guide

## Prerequisites

```bash
# Install Docker
sudo apt update
sudo apt install docker.io docker-compose

# Enable Apache modules
sudo a2enmod proxy proxy_http headers rewrite
sudo systemctl restart apache2
```

## First-Time Setup

```bash
# 1. Clone repo
cd /opt
sudo git clone https://github.com/flocko-motion/gofins.git
sudo chown -R $(whoami):$(whoami) gofins

# 2. Configure environment
cd /opt/gofins/deployment
cp .env.example .env
nano .env  # Set DB_PASSWORD

# 3. Create user(s)
sudo htpasswd -c .htpasswd yourusername
sudo htpasswd .htpasswd friend1  # Add more users

# 4. Deploy (handles Apache, systemd, Docker automatically)
sudo ./deploy.sh
```

That's it! The deploy script will:
- Set up Apache config symlink
- Set up systemd service
- Build Docker images
- Start all containers

## Updates

```bash
sudo /opt/gofins/deployment/deploy.sh
```

## Useful Commands

```bash
# Check status
sudo systemctl status gofins
docker ps

# View logs
docker logs -f gofins-api
docker logs -f gofins-ui
docker logs -f gofins-db
sudo tail -f /var/log/apache2/gofins-error.log

# Restart services
sudo systemctl restart gofins
docker-compose restart

# Stop everything
sudo systemctl stop gofins
```

## Ports

- **7701** - PostgreSQL (localhost only)
- **7702** - Go API (localhost only)
- **7703** - UI nginx (localhost only)
- **80** - Apache (public, proxies to containers)
