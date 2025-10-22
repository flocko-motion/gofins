# Installation Guide

## Prerequisites

```bash
# Install Docker (official method)
sudo apt update
sudo apt install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

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
sed -i "s/DB_PASSWORD=.*/DB_PASSWORD=$(openssl rand -base64 32)/" .env

# 3. Configure Apache
cp apache-gofins.conf.example apache-gofins.conf
nano apache-gofins.conf  # Edit ServerName and SSL certificate paths

# 4. Create user(s) (still in deployment/)
sudo htpasswd -c .htpasswd yourusername
sudo htpasswd .htpasswd friend1  # Add more users

# 5. Deploy (handles Apache, systemd, Docker automatically)
sudo bash deploy.sh
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
