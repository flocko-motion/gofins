# gofins

Stock screener with YoY analysis, multi-user support, and research notes.

## What it does

- Fetches stock data from FMP API (profiles, prices, quotes)
- Calculates year-over-year returns for custom time ranges
- Analyzes volatility (mean, stddev) across thousands of symbols
- Tracks personal ratings, favorites, and notes per user
- Multi-user isolation via Apache .htaccess auth or config file

## Directory Structure

```
gofins/              # Go backend (API server, CLI tools)
gofins-ui/           # React frontend (Vite + TailwindCSS)
development/         # Local dev scripts (run_db.sh, run_server.sh, etc.)
deployment/          # Production Docker setup (docker-compose, deploy.sh)
```

## Authentication

**Production**: Apache sets `X-Remote-User` header via .htaccess  
**Development**: Uses `~/.gofins/config.yaml` default user or `--user` flag  
**User ID**: Username hashed to stable UUID for database isolation

All user data (ratings, favorites, notes, analyses) scoped per user.

## MCP Integration with Claude

FINS can be connected to Claude via the Model Context Protocol (MCP) to expose financial data and user preferences directly to Claude for enhanced analysis and interaction.

### Available via MCP:

**User Data:**
- Favorites (favs)
- Personal ratings
- Notes

**Market Data:**
- Symbol information (company profiles)
- Price histories
- Real-time quotes
- Statistics (e.g. CAGR, variance)

**Authentication:**
- Header-based authentication via `X-Remote-User`
- Maintains user isolation and data privacy

This enables Claude to assist with portfolio analysis, stock research, and personalized financial insights while respecting user data boundaries.

## Quick Start

### Development

```bash
cd development
./run_db.sh          # Start PostgreSQL
cd ../gofins
go run . server      # Start API on :8080
cd ../gofins-ui
npm install          # First time only
npm run dev          # Start UI on :5173
```

**UI Development Notes:**
- UI runs on http://localhost:5173
- Connects to API at http://localhost:8080 by default
- See "UI API Configuration" below for connecting to remote backends

### Production Deployment

**Prerequisites**:
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

**First-Time Setup**:
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
sed -i 's/yourdomain.com/<youractualdomain.whatever>/g' apache-gofins.conf

# 4. Create user(s)
sudo htpasswd -c .htpasswd yourusername
sudo htpasswd .htpasswd friend1  # Add more users

# 5. Deploy (handles Apache, systemd, Docker automatically)
sudo bash deploy.sh
```

**Updates**:
```bash
sudo /opt/gofins/deployment/deploy.sh
```

## User Management

**Automatic User Creation**: Users are created automatically when they first access the application.

### How it works:
1. Apache's `.htaccess` authentication prompts for username/password
2. On successful authentication, Apache sets the `X-Remote-User` header with the username
3. The Go backend receives this header and automatically creates a user record in the database if it doesn't exist
4. **The first user to access the system becomes the admin** (has access to the Errors tab and other admin features)
5. All subsequent users are regular users with access to their own data only

### User Isolation:
- Each user's ratings, favorites, notes, and analyses are completely isolated
- Users cannot see or access other users' data
- The admin user can see system errors but not other users' personal data

### Adding New Users:
```bash
cd /opt/gofins/deployment
sudo htpasswd .htpasswd newusername
# User will be created in database on first login
```

### Removing Users:
```bash
sudo htpasswd -D .htpasswd username
# Optionally delete user data:
docker exec -it gofins-db psql -U gofins -d gofins -c "DELETE FROM users WHERE username = 'username';"
```

## UI API Configuration

The UI can connect to different API backends without rebuilding. Configure in order of priority:

### 1. Browser Console (Runtime)
Change API URL on-the-fly without rebuilding:

```javascript
// Switch to production backend
localStorage.setItem('apiUrl', 'https://omnitopos.net/gofins/api')

// Set HTTP Basic Auth credentials (if required)
localStorage.setItem('apiAuth', btoa('username:password'))

location.reload()

// Switch back to localhost
localStorage.removeItem('apiUrl')
localStorage.removeItem('apiAuth')
location.reload()
```

### 2. Environment Variable
Create `gofins-ui/.env.local`:
```bash
VITE_API_URL=https://your-domain.com/api
```

### 3. Default
- Development: `http://localhost:8080/api`
- Production: relative `/api` path

## Debugging

### FMP API Request Logging

Enable verbose logging to see timing for every FMP API request:

```bash
# Enable verbose logging (shows every request with timing)
curl -X POST http://localhost:7702/api/debug/fmp-verbose/true

# Disable verbose logging
curl -X POST http://localhost:7702/api/debug/fmp-verbose/false
```

When enabled, logs will show:
```
📊 v3/historical-price-full/AAPL → 0.523s (status 200)
📊 v3/historical-price-full/MSFT → 0.487s (status 200)
⚠️  SLOW request to v3/profile/GOOGL took 3.21s (status 200)
```

Useful for diagnosing:
- Slow API responses
- Network issues
- Rate limiting problems
- Performance bottlenecks

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

## Database & CLI Commands

### Development
Use gofins binary to interact with the database. In gofins/ directory:

```bash
# View schema
go run . db schema

# Execute SQL
go run . db sql -q "SELECT COUNT(*) FROM symbols"

# Update quotes
go run . update quotes

# Fetch profiles
go run . fetch profiles --limit 100
```

### Production
Use the wrapper script to run commands inside the Docker container:

```bash
cd /opt/gofins/deployment

# View schema
./gofins db schema

# Execute SQL
./gofins db sql -q "SELECT COUNT(*) FROM symbols"

# Update quotes
./gofins update quotes

# Fetch profiles
./gofins fetch profiles --limit 100

# View all available commands
./gofins --help
```

The wrapper script automatically executes commands inside the running `gofins-api` container.

## License

WTFPL - Do What The Fuck You Want To
