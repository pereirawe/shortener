## Known Issues

### 1. Insecure random generation
- Status: open
- Remote: -

- Severity: high
- Location: internal/api/handler.go:219
- Description: Uses math/rand without seed for short code generation
- Impact: Predictable codes and possible collisions
- Suggested fix: Use crypto/rand

### 2. Missing timeout in SEO fetch
- Status: open
- Remote: -

- Severity: high
- Location: internal/api/handler.go:136
- Description: External request without timeout or context
- Impact: Request may hang indefinitely
- Suggested fix: Use http.Client with timeout and context

### 3. Weak URL validation
- Status: open
- Remote: -

- Severity: high
- Location: internal/api/handler.go:87
- Description: Only checks string prefix
- Impact: Invalid or malicious URLs accepted
- Suggested fix: Use net/url.Parse and validate host

### 4. Potential SSRF vulnerability
- Status: open
- Remote: -

- Severity: high
- Location: internal/api/seo.go
- Description: Fetching arbitrary URLs without restriction
- Impact: Internal network exposure
- Suggested fix: Block private IP ranges

### 5. No rate limiting
- Status: open
- Remote: -

- Severity: medium
- Location: API layer
- Description: Public endpoint without throttling
- Impact: Abuse and resource exhaustion
- Suggested fix: Add middleware (token bucket / Redis)

### 6. Cache inconsistency
- Status: open
- Remote: -

- Severity: medium
- Location: internal/api/handler.go:153
- Description: Only original URL cached
- Impact: Loss of metadata consistency
- Suggested fix: Cache full response object

### 7. Uncontrolled goroutines
- Status: open
- Remote: -

- Severity: medium
- Location: internal/api/handler.go:187
- Description: Fire-and-forget goroutines
- Impact: Resource leaks under load
- Suggested fix: Use worker queue or sync handling

### 8. Business logic in handler
- Status: open
- Remote: -

- Severity: medium
- Location: internal/api/handler.go
- Description: Handler responsible for domain logic
- Impact: Hard to test and maintain
- Suggested fix: Introduce service layer

### 9. Logging not structured
- Status: open
- Remote: -

- Severity: low
- Location: global
- Description: Uses log.Printf
- Impact: Hard observability
- Suggested fix: Use structured logger

### 10. AutoMigrate in production
- Status: open
- Remote: -

- Severity: medium
- Location: persistence layer
- Description: Schema managed automatically
- Impact: Risky migrations
- Suggested fix: Use versioned migrations

### 11. Missing config validation
- Status: open
- Remote: -

- Severity: medium
- Location: internal/config/config.go
- Description: Env vars not validated
- Impact: Runtime failures
- Suggested fix: Validate on startup

### 12. Silent cache failures
- Status: open
- Remote: -

- Severity: low
- Location: handler.go
- Description: Cache errors only logged
- Impact: Hidden degradation
- Suggested fix: Add metrics or retry

### 13. Duplicate constructor
- Status: open
- Remote: -

- Severity: low
- Location: handler.go
- Description: NewHandler and NewHandlerWithMocks identical
- Impact: Redundant code
- Suggested fix: Remove duplication

### 14. Missing edge-case tests
- Status: open
- Remote: -

- Severity: medium
- Location: tests
- Description: No tests for concurrency, cache failure, collisions
- Impact: Undetected bugs
- Suggested fix: Expand test suite

### 15. No observability/tracing
- Status: open
- Remote: -

- Severity: low
- Location: global
- Description: No tracing or correlation IDs
- Impact: Hard debugging in production
- Suggested fix: Add tracing/metrics
