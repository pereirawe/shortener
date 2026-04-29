#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/close_issue.sh <issue_id>

ID=${1:-}
if [[ -z "$ID" ]]; then
  echo "Usage: close_issue.sh <id>"
  exit 1
fi

REMOTE_URL=$(git config --get remote.origin.url)

if [[ "$REMOTE_URL" == *"github.com"* ]]; then
  gh issue close "$ID" || true
elif [[ "$REMOTE_URL" == *"gitlab"* ]]; then
  glab issue close "$ID" || true
fi

# Mark as resolved locally
FILE="docs/ai/known_issues.md"
if [[ -f "$FILE" ]]; then
  awk -v id="$ID" '
  BEGIN{found=0}
  /^### [0-9]+\./{
    if(found==1){found=0}
  }
  $0 ~ "^### "id"\."{found=1}
  {
    if(found==1 && $0 ~ /^- Status:/){print "- Status: resolved"; next}
    print
  }' "$FILE" > "$FILE.tmp" && mv "$FILE.tmp" "$FILE"
fi

echo "[issue] closed $ID"
