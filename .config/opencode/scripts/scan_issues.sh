#!/usr/bin/env bash
set -euo pipefail

echo "[scan-issues] Running static scan..."

# Basic heuristics using ripgrep if available, grep otherwise
PATTERN="math/rand|rand\.Intn|go\s+|http\.Get|http\.DefaultClient|context\.Background\(|TODO|FIXME"
TARGETS=(./cmd ./internal ./.config/opencode/scripts ./build.sh)

echo "[scan-issues] Searching risky patterns..."

if command -v rg >/dev/null 2>&1; then
  rg -n "$PATTERN" "${TARGETS[@]}" || true
elif command -v grep >/dev/null 2>&1; then
  grep -RInE "$PATTERN" "${TARGETS[@]}" 2>/dev/null || true
else
  echo "[scan-issues] no text search tool available"
fi

echo "[scan-issues] Done. Now run /scan-issues in the assistant to update .config/opencode/known_issues.md"
