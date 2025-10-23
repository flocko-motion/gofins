#!/bin/bash

# GoFins UI Build Script
# Builds the React app for production and outputs to bin/

echo "Building GoFins UI for production..."
echo ""

# Get the script directory and navigate to project root
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Navigate to gofins-ui directory
cd "$PROJECT_ROOT/gofins-ui" || exit 1

# Build the project
npm run build

# Copy built files to development/bin directory
echo "Copying built files to bin/..."
mkdir -p "$SCRIPT_DIR/bin"
cp -r dist/* "$SCRIPT_DIR/bin/"

echo ""
echo "Build complete! Files are in the bin/ directory."
echo "You can serve them with any static file server."
echo ""
echo "Example:"
echo "  cd bin && python3 -m http.server 8080"
echo "  # or"
echo "  cd bin && npx serve ."
