# Development Scripts

Local development setup for gofins.

## Database

Start PostgreSQL in Docker:
```bash
./run_db.sh
```

Uses `docker-compose.yml` with local data in `~/.fins/postgres/data`.

## Backend

Run Go API server:
```bash
./run_server.sh
```

Starts on `localhost:8080` with `--user` flag support for testing multi-user.

## Frontend

Build UI:
```bash
./build_ui.sh
```

Run UI dev server:
```bash
./run_ui.sh
```

Starts on `localhost:5173` with hot reload.
