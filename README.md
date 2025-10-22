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

## Quick Start

**Development**:
```bash
cd development
./run_db.sh          # Start PostgreSQL
cd ../gofins
go run . server      # Start API on :8080
cd ../gofins-ui
npm run dev          # Start UI on :5173
```

**Deployment**:
```bash
cd deployment
cp .env.example .env  # Set DB_PASSWORD
./deploy.sh           # Build & start Docker containers
```

## Database
Use gofins binary to interact with the database. In gofins/ directory:

```bash
# View schema
go run . db schema

# Execute SQL
go run . db sql -q "SELECT COUNT(*) FROM symbols"
```

## License

WTFPL - Do What The Fuck You Want To
