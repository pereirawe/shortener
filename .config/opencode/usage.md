## AI Workflow Usage Guide

This document explains how to use the AI-assisted workflow for tracked work and development.

---

## Core Concepts

- `known_issues.md`: Single source of truth for bugs, features, docs, and chores
- Remote issues: GitHub/GitLab source of truth for execution

Every entry must include `Reported by`:
- use the user name when a person creates or reports the item
- use the model name when an AI creates or reports the item

---

## Typical Flows

## Status Lifecycle

- `backlog`: item captured, still being triaged or refined
- `ready`: item is clear and approved to be picked up
- `open`: item selected locally, waiting for remote issue creation
- `in-progress`: remote issue created and implementation started
- `resolved`: item completed or explicitly closed

Allowed transitions:

- `backlog -> ready`
- `ready -> backlog`
- `backlog -> open` via `promote.sh`
- `ready -> open` via `promote.sh`
- `open -> in-progress` via `create_issue.sh`
- `open -> resolved` via `close_issue.sh`
- `in-progress -> resolved` via `close_issue.sh`
- `resolved -> open` only by explicit manual reopen in `known_issues.md`

Script rules:

- `promote.sh` resets `Remote` to `-`
- `create_issue.sh` refuses items that are not `open` or already have a remote id
- `close_issue.sh` refuses items that are still `backlog`, `ready`, or already `resolved`

### 1. Fixing an Existing Issue

1. Pick an issue from `known_issues.md`
2. Create remote issue + branch:
```
.config/opencode/scripts/create_issue.sh 3
```
3. Implement changes
4. Commit with reference:
```
git commit -m "fix: improve cache (#123)"
```
5. Review:
```
make review
/review-branch
```
6. Close:
```
.config/opencode/scripts/close_issue.sh 3
```

---

### 2. Planning and Starting a New Feature

1. Add the item to `known_issues.md` with `Status: backlog` and `Type: feat`

Example:
```
### 6. Add analytics endpoint
- Status: backlog
- Type: feat
- Severity: medium
- Reported by: William Pereira
- Remote: -
- Description: Track click stats per URL
- Impact: Product insight
- Location: internal/api/handler.go
- Suggested implementation: Add new route + DB query
```

2. Promote the local item when it is ready to be worked on:
```
make promote id=6
```

3. Create remote issue:
```
.config/opencode/scripts/create_issue.sh 6
```

4. Develop + test + document

---

## Commands Summary

### Local Commands

```
make scan-issues
make review
make promote id=<n>
make close-issue id=<n>
```

Alternative explicit form:
```
make -f .config/opencode/Makefile scan-issues
make -f .config/opencode/Makefile review
make -f .config/opencode/Makefile promote id=<n>
make -f .config/opencode/Makefile close-issue id=<n>
```

### Scripts

```
.config/opencode/scripts/create_issue.sh <id|title body>
.config/opencode/scripts/promote.sh <id>
.config/opencode/scripts/close_issue.sh <id>
```

### AI Commands

```
/scan-issues
/review-branch
/plan-feature
/promote <id>
```

---

## Conventions

- Always reference issue in commits: `(#id)`
- Tests are mandatory
- API changes require Bruno update
- Keep `known_issues.md` updated

---

## Example Known Issue

```
### 3. Weak URL validation
- Status: in-progress
- Type: bug
- Severity: high
- Reported by: openai/gpt-5.4
- Remote: #123
- Location: internal/api/handler.go:82-89
- Description: Only checks prefix
- Impact: Invalid URLs accepted
- Suggested fix: Use net/url.Parse
```

---

This workflow ensures traceability, consistency, and high engineering quality.
