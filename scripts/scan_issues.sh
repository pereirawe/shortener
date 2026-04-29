#!/usr/bin/env bash
set -euo pipefail

echo "[scan-issues] Running static scan..."

# Basic heuristics using ripgrep if available
if command -v rg >/dev/null 2>&1; then
  echo "[scan-issues] Searching risky patterns..."
  rg -n "math/rand|rand\.Intn|go\s+|http\.Get|http\.DefaultClient|context\.Background\(|TODO|FIXME" || true
else
  echo "[scan-issues] ripgrep not found, skipping pattern scan"
fi

echo "[scan-issues] Done. Now run /scan-issues in the assistant to update docs/ai/known_issues.md"
