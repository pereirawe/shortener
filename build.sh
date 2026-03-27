#!/bin/bash

# Build script for shortener application

set -e

echo "Building shortener application..."

# Get the version from git or use a default
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
TIMESTAMP=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")

# Create output directory
mkdir -p bin

# Build the application
go build -o bin/shortener .

echo "✓ Build completed successfully"
echo "✓ Binary location: bin/shortener"
echo "✓ Version: ${VERSION}"
echo "✓ Commit: ${COMMIT}"
