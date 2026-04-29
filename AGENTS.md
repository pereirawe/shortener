## Overview

Go URL shortener API using:
- `net/http` (no framework)
- GORM (Postgres)
- Redis cache

Entrypoint:
```
cmd/api/main.go
```

Handlers contain business logic (no service layer yet).

---

## Run / Dev

Infra (required):
```
docker compose up -d
```

Run API:
```
go run ./cmd/api
```

Build:
```
go build -o bin/shortener ./cmd/api
```

Env:
- Copy `.env.example` → `.env`
- API depends on Redis + Postgres being up

---

## Tests

Run all:
```
go test ./...
```

Pre-commit runs tests automatically. Failing tests block commits.

---

## Key Architecture Notes

- Routing uses Go 1.22+ pattern syntax (`"POST /path"`)
- Cache is read-through (Redis → fallback Postgres)
- SEO metadata is fetched synchronously in request path
- Short code generation currently uses `math/rand`

Important:
- No timeouts on external HTTP calls (SEO)
- Handlers mix transport + domain logic

---

## AI Workflow (CRITICAL)

Project uses a structured workflow in `docs/ai/`:

- `known_issues.md` → active work
- `improvements.md` → backlog
- `commands.md` → AI behavior
- `usage.md` → how to operate

Never bypass this workflow when making non-trivial changes.

---

## Core Commands

Local:
```
make scan-issues
make review
make promote id=<n>
make close-issue id=<n>
```

Scripts:
```
./scripts/create_issue.sh <id|"title" "body">
./scripts/promote.sh <id>
./scripts/close_issue.sh <id>
```

---

## Required Dev Flow

1. Select issue from `known_issues.md`
2. Create remote issue + branch:
```
./scripts/create_issue.sh <id>
```
3. Implement + tests
4. Commit MUST include issue id:
```
(#123)
```
5. Review:
```
make review
```
6. Close:
```
./scripts/close_issue.sh <id>
```

Commits without `#id` will fail (`commit-msg` hook).

---

## Conventions (Non-Obvious)

- API changes → MUST update Bruno collection (`bruno/`)
- Tests are mandatory for changes
- Do not duplicate issues between `improvements` and `known_issues`
- `known_issues` must reflect current state (status + remote id)

---

## Gotchas

- App will fail silently if Redis/Postgres not running
- Cache only stores original URL (not full response)
- Goroutines used for click counting are fire-and-forget
- No migrations system (AutoMigrate used)

---

## When Editing Code

Prefer minimal changes.

If touching:
- handlers → consider impact on cache + SEO
- routes → update Bruno
- models → consider DB migration risk

Always consider updating:
```
docs/ai/known_issues.md
```
