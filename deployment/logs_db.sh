#!/bin/bash

# Show logs from the gofins-db container
# Usage: ./logs_db.sh [options]
# Options are passed to docker logs (e.g., -f for follow, --tail 100)

docker logs "${@:--f}" gofins-db
