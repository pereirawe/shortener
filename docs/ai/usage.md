## AI Workflow Usage Guide

This document explains how to use the AI-assisted workflow for issues, improvements, and development.

---

## Core Concepts

- `known_issues.md`: Active tracked problems/tasks
- `improvements.md`: Backlog of ideas and enhancements
- Remote issues: GitHub/GitLab source of truth for execution

---

## Typical Flows

### 1. Fixing an Existing Issue

1. Pick an issue from `known_issues.md`
2. Create remote issue + branch:
```
./scripts/create_issue.sh 3
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
./scripts/close_issue.sh 3
```

---

### 2. Creating a New Feature (Improvement)

1. Add idea to `improvements.md`

Example:
```
### 6. Add analytics endpoint
- Status: backlog
- Description: Track click stats per URL
- Impact: Product insight
- Suggested implementation: Add new route + DB query
```

2. Promote to issue:
```
make promote id=6
```

3. Create remote issue:
```
./scripts/create_issue.sh 6
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

### Scripts

```
./scripts/create_issue.sh <id|title body>
./scripts/promote.sh <id>
./scripts/close_issue.sh <id>
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
- Remote: #123
- Description: Only checks prefix
- Impact: Invalid URLs accepted
- Suggested fix: Use net/url.Parse
```

---

## Example Improvement

```
### 6. Add analytics endpoint
- Status: backlog
- Description: Track clicks
- Impact: Product insights
- Suggested implementation: New route + DB
```

---

This workflow ensures traceability, consistency, and high engineering quality.
