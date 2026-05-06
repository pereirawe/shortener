#!/usr/bin/env bash
set -euo pipefail

echo "[pre-commit] Running checks..."

# Run tests
if command -v go >/dev/null 2>&1; then
  go test ./... || { echo "Tests failed"; exit 1; }
fi

# Ensure known_issues was considered
if git diff --cached --name-only | grep -q "known_issues.md"; then
  echo "[pre-commit] known_issues updated"
else
  echo "[pre-commit] WARNING: known_issues not updated"
fi

echo "[pre-commit] OK"
