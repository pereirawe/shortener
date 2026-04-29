#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/promote.sh <id>

ID=${1:-}
if [[ -z "$ID" ]]; then
  echo "Usage: promote.sh <id>"
  exit 1
fi

IMP_FILE="docs/ai/improvements.md"
ISS_FILE="docs/ai/known_issues.md"

if [[ ! -f "$IMP_FILE" || ! -f "$ISS_FILE" ]]; then
  echo "Required files not found"
  exit 1
fi

# Extract improvement section
SECTION=$(awk "/^### $ID\./,/^### /" "$IMP_FILE" | sed '$d')

if [[ -z "$SECTION" ]]; then
  echo "Improvement $ID not found"
  exit 1
fi

TITLE=$(echo "$SECTION" | head -n1 | sed 's/^### [0-9]*\. //')

# Determine next issue number
NEXT_ID=$(grep -E "^### [0-9]+\." "$ISS_FILE" | tail -n1 | awk '{print $2}' | tr -d '.' )
NEXT_ID=$((NEXT_ID+1))

# Build new issue block
{
  echo ""
  echo "### $NEXT_ID. $TITLE"
  echo "- Status: open"
  echo "- Remote: -"
  echo "$(echo "$SECTION" | tail -n +2)"
} >> "$ISS_FILE"

# Mark improvement as promoted
awk -v id="$ID" '
BEGIN{found=0}
/^### [0-9]+\./{
  if(found==1){found=0}
}
$0 ~ "^### "id"\."{found=1}
{
  if(found==1 && $0 ~ /^- Status:/){print "- Status: promoted"; next}
  print
}' "$IMP_FILE" > "$IMP_FILE.tmp" && mv "$IMP_FILE.tmp" "$IMP_FILE"

echo "[promote] Improvement $ID promoted to issue $NEXT_ID"
