## Technical Decisions

- Go stdlib for HTTP (simplicity)
- GORM for persistence (fast setup)
- Redis as cache with TTL (2 weeks)
- Random short codes (7 chars)
- Allow custom codes with validation

Tradeoffs:
- AutoMigrate vs controlled migrations
- Simplicity vs scalability
