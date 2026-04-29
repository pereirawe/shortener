## Architecture

Flow:
HTTP -> Handler -> (future Service layer) -> Store (Postgres) + Cache (Redis)

Current state:
- Handler contains business logic
- Cache used as read-through fallback
- SEO fetch done inline in request lifecycle

Gaps:
- Missing service layer
- No clear separation of concerns
- External calls without timeout
