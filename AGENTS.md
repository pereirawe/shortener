## Overview

Go URL shortener API using:
- `net/http`
- GORM + PostgreSQL
- Redis cache

Entrypoint:
```
cmd/api/main.go
```

## OpenCode Layout

Operational OpenCode files live in:
- `opencode.json`
- `.opencode/agents/`
- `.opencode/commands/`
- `.opencode/skills/`

Project documentation and workflow details stay in:
- `.config/opencode/`

Always treat `.config/opencode/known_issues.md` as the source of truth for tracked technical issues.

## Workflow

Use these project commands when relevant:
- `/scan-issues`
- `/review-branch`
- `/plan-feature <description>`
- `/promote <id>`
- `/create-issue <id>`
- `/close-issue <id>`

Local scripts live in `.config/opencode/scripts/`.

## Engineering Notes

- Prefer minimal code changes
- Tests are mandatory for code changes
- API changes must update `bruno/`
- Keep `.config/opencode/known_issues.md` current when fixing or discovering issues
