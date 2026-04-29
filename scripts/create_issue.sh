#!/usr/bin/env bash
set -euo pipefail

# Usage:
# ./scripts/create_issue.sh "Issue Title" "Issue Body"
# ./scripts/create_issue.sh <local_issue_number>

INPUT=${1:-}
TITLE=""
BODY=""

# If only a number is provided, fetch from known_issues
if [[ "$INPUT" =~ ^[0-9]+$ ]]; then
  FILE="docs/ai/known_issues.md"
  if [[ ! -f "$FILE" ]]; then
    echo "known_issues.md not found"
    exit 1
  fi

  SECTION=$(awk "/^### $INPUT\./,/^### /" "$FILE" | sed '$d')

  TITLE=$(echo "$SECTION" | head -n1 | sed 's/^### [0-9]*\. //')
  BODY=$(echo "$SECTION" | tail -n +2)
else
  TITLE=${1:-}
  BODY=${2:-}
fi

if [[ -z "$TITLE" || -z "$BODY" ]]; then
  echo "Usage: create_issue.sh \"title\" \"body\""
  exit 1
fi

REMOTE_URL=$(git config --get remote.origin.url)

echo "[issue] Remote: $REMOTE_URL"

if [[ "$REMOTE_URL" == *"github.com"* ]]; then
  if ! command -v gh >/dev/null 2>&1; then
    echo "GitHub CLI (gh) not installed"
    exit 1
  fi
ISSUE_URL=$(gh issue create --title "$TITLE" --body "$BODY")
ISSUE_ID=$(basename "$ISSUE_URL")

elif [[ "$REMOTE_URL" == *"gitlab"* ]]; then
  if ! command -v glab >/dev/null 2>&1; then
    echo "GitLab CLI (glab) not installed"
    exit 1
  fi
  ISSUE_URL=$(glab issue create --title "$TITLE" --description "$BODY" --yes | grep -Eo 'https?://[^ ]+')
  ISSUE_ID=$(basename "$ISSUE_URL")

else
  echo "Unsupported remote"
  exit 1
fi

echo "[issue] Created: $ISSUE_URL"

# Update known_issues with Remote ID and status
FILE="docs/ai/known_issues.md"
if [[ -f "$FILE" && "$INPUT" =~ ^[0-9]+$ ]]; then
  awk -v id="$INPUT" -v rid="$ISSUE_ID" '
  BEGIN{found=0}
  /^### [0-9]+\./{
    if(found==1){found=0}
  }
  $0 ~ "^### "id"\."{found=1}
  {
    if(found==1 && $0 ~ /^- Status:/){print "- Status: in-progress"; next}
    if(found==1 && $0 ~ /^- Remote:/){print "- Remote: #"rid; next}
    print
  }' "$FILE" > "$FILE.tmp" && mv "$FILE.tmp" "$FILE"
fi

# Create branch
SLUG=$(echo "$TITLE" | tr '[:upper:]' '[:lower:]' | tr -cd 'a-z0-9 ' | tr ' ' '-')
BRANCH="issue-${ISSUE_ID}-${SLUG}"

git checkout -b "$BRANCH"

echo "[issue] Branch created: $BRANCH"

echo "$ISSUE_ID"
