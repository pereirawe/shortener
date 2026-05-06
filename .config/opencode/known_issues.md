## Known Issues

Single source of truth for tracked work in this project.

Format:

### <id>. <title>
- Status: backlog | ready | open | in-progress | resolved
- Type: bug | feat | doc | chore
- Severity: critical | high | medium | low
- Reported by: <user-name> | <model-name>
- Remote: - | #<remote-id>
- Location:
- Description:
- Impact:
- Suggested fix:

Status lifecycle:

- `backlog`: item captured but not yet refined or prioritized. Allowed next statuses: `ready`.
- `ready`: item is clear enough to be picked up. Allowed next statuses: `backlog`, `open`.
- `open`: item selected locally and awaiting remote issue creation. Allowed next statuses: `in-progress`, `resolved`.
- `in-progress`: remote issue exists and work has started. Allowed next statuses: `resolved`.
- `resolved`: completed or closed item. Reopen only by explicit manual change back to `open`.

Operational rules:

- `promote.sh` only allows `backlog -> open` and `ready -> open`.
- `create_issue.sh` only allows `open -> in-progress` and requires `Remote: -`.
- `close_issue.sh` only allows `open -> resolved` or `in-progress -> resolved`.
- `backlog <-> ready` and explicit reopen from `resolved` are manual editorial changes in `known_issues.md`.

### 1. Predictable short-code generation
- Status: in-progress
- Type: bug
- Severity: high
- Reported by: openai/gpt-5.4
- Remote: #2
- Location: internal/api/handler.go:9,213-221
- Description: Auto-generated short codes use `math/rand` via `rand.Intn`.
- Impact: Codes are predictable and collision behavior is weaker than intended.
- Suggested fix: Use `crypto/rand` for code generation.

### 2. Weak URL validation before persistence and SEO fetch
- Status: open
- Type: bug
- Severity: high
- Reported by: openai/gpt-5.4
- Remote: -
- Location: internal/api/handler.go:82-89
- Description: `original_url` validation only checks for `http://` or `https://` prefixes.
- Impact: Malformed URLs can pass validation and fail later during fetch or redirect.
- Suggested fix: Parse with `net/url` and require valid scheme and host.

### 3. SEO fetch allows arbitrary outbound requests
- Status: open
- Type: bug
- Severity: high
- Reported by: openai/gpt-5.4
- Remote: -
- Location: internal/api/handler.go:136; internal/api/seo.go:31-42
- Description: User input is fetched directly for SEO metadata without host or IP restrictions.
- Impact: The API can be used to probe internal services or private network ranges.
- Suggested fix: Resolve and deny private, loopback, and link-local destinations before fetch.

### 4. Short-code uniqueness check is not atomic
- Status: open
- Type: bug
- Severity: high
- Reported by: openai/gpt-5.4
- Remote: -
- Location: internal/api/handler.go:113-125,147-150,213-231; internal/model/url.go:15
- Description: The handler checks `ExistsByShortCode` before `Create`, while the DB uniqueness guarantee is enforced separately.
- Impact: Concurrent requests can still hit duplicate-key errors and return a generic `500` instead of retrying or returning `409`.
- Suggested fix: Handle unique-constraint errors explicitly and retry auto-generated codes.

### 5. Click increments run in fire-and-forget goroutines
- Status: open
- Type: bug
- Severity: medium
- Reported by: openai/gpt-5.4
- Remote: -
- Location: internal/api/handler.go:187,209; internal/service/postgres_service.go:63-67
- Description: Redirects spawn unbounded goroutines for click counting and ignore write failures.
- Impact: Click updates can be dropped and goroutine count can grow under load.
- Suggested fix: Increment synchronously or move counting to a bounded worker/queue.

### 6. Handler owns transport, validation, SEO, caching, and persistence flow
- Status: open
- Type: chore
- Severity: medium
- Reported by: openai/gpt-5.4
- Remote: -
- Location: internal/api/handler.go:74-231
- Description: Most request orchestration and business rules live in the HTTP handler.
- Impact: Changes are harder to test and maintain, and service boundaries remain unclear.
- Suggested fix: Move domain workflow into a service layer and keep handlers thin.

### 7. Automatic schema migration runs on every startup
- Status: open
- Type: chore
- Severity: medium
- Reported by: openai/gpt-5.4
- Remote: -
- Location: internal/service/postgres_service.go:20-40
- Description: `NewPostgresService` calls `db.AutoMigrate(&model.URL{})` during application boot.
- Impact: Production schema changes are coupled to deploys and are harder to review or roll back.
- Suggested fix: Use explicit versioned migrations outside application startup.

### 8. Server has no graceful shutdown path
- Status: open
- Type: chore
- Severity: medium
- Reported by: openai/gpt-5.4
- Remote: -
- Location: cmd/api/main.go:35-39
- Description: The API starts with `http.ListenAndServe` and does not handle shutdown signals.
- Impact: In-flight requests can be dropped during restarts or container termination.
- Suggested fix: Use `http.Server` with signal handling and `Shutdown`.

### 9. SEO fetch timeout is already configured
- Status: resolved
- Type: bug
- Severity: high
- Reported by: openai/gpt-5.4
- Remote: -
- Location: internal/api/seo.go:19-27
- Description: SEO requests previously lacked an explicit client timeout.
- Impact: Requests could hang indefinitely.
- Suggested fix: Keep `seoHTTPClient` timeout and redirect limits in place.

### 10. Build script target is already correct
- Status: resolved
- Type: chore
- Severity: medium
- Reported by: openai/gpt-5.4
- Remote: -
- Location: build.sh:18
- Description: The build script previously targeted the wrong package.
- Impact: Builds could fail or compile the wrong entrypoint.
- Suggested fix: Keep `go build -o bin/shortener ./cmd/api` as the build target.

### 11. PostgreSQL port docs are aligned
- Status: resolved
- Type: doc
- Severity: medium
- Reported by: openai/gpt-5.4
- Remote: -
- Location: .env.example; docker-compose.yml
- Description: The documented PostgreSQL port was inconsistent with local setup.
- Impact: Fresh setups could fail to connect.
- Suggested fix: Keep `POSTGRES_PORT=5433` aligned across docs and local setup.

### 12. Context env var documentation is aligned
- Status: resolved
- Type: doc
- Severity: medium
- Reported by: openai/gpt-5.4
- Remote: -
- Location: .config/opencode/context.md:236-252; internal/config/config.go:38-55
- Description: Context documentation previously referenced outdated environment variable names.
- Impact: Developers and assistants could use the wrong configuration.
- Suggested fix: Keep `context.md` aligned with `config.Load()`.

### 13. Cache only stores the original URL string
- Status: in-progress
- Type: chore
- Severity: low
- Reported by: openai/gpt-5.4
- Remote: #321
- Location: internal/api/handler.go:169-209; internal/service/redis_service.go:12-30
- Description: Redis entries cache only the original URL, so handlers must rebuild response state from mixed cache and database paths.
- Impact: Future cache-backed reads and metadata reuse are harder to extend consistently.
- Suggested fix: Cache a structured payload for URL lookups instead of a single string value.
