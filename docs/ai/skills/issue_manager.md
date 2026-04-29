## Skill: Issue Manager

Goal:
Maintain `docs/ai/known_issues.md` as the single source of truth for technical risks.

What it does:
- Scans code for bugs, security issues, and design problems
- Adds new issues with context
- Updates existing issues (no duplicates)
- Removes resolved issues
- Links to file:line when possible

Issue format:
- Title
- Severity: high | medium | low
- Location: file:line
- Description
- Impact
- Suggested fix

Heuristics:
- Security: URL handling, external calls, randomness
- Concurrency: goroutines, shared state
- Reliability: timeouts, retries, error handling
- Architecture: separation of concerns

Rules:
- Prefer updating existing entries over creating similar ones
- Keep entries concise and actionable
- Do not speculate without code evidence
