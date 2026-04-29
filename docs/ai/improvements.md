## Improvements Backlog

This file tracks potential enhancements before they become formal issues.

Format:

### <id>. <title>
- Status: backlog | ready | promoted
- Description:
- Impact:
- Suggested implementation:

---

### 1. Introduce service layer
- Status: backlog
- Description: Separate business logic from handlers
- Impact: Improves maintainability and testability
- Suggested implementation: Create internal/service layer and move logic

### 2. Replace math/rand with crypto/rand
- Status: backlog
- Description: Improve randomness security
- Impact: Prevent predictable short codes
- Suggested implementation: Use crypto/rand

### 3. Add URL validation with net/url
- Status: backlog
- Description: Stronger validation
- Impact: Prevent malformed input
- Suggested implementation: Use url.Parse and validation rules

### 4. Add SSRF protection
- Status: backlog
- Description: Restrict internal network access
- Impact: Security hardening
- Suggested implementation: Block private IP ranges

### 5. Improve cache structure
- Status: backlog
- Description: Cache full response
- Impact: Better performance and consistency
- Suggested implementation: Serialize DTO in cache
