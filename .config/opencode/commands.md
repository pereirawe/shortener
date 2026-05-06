## Commands

These commands define how the assistant operates over this project context.

### /scan-issues
Deep analysis of the codebase.

Responsibilities:
- Detect new issues (security, concurrency, architecture)
- Update `.config/opencode/known_issues.md`
- Normalize `Status`, `Type`, `Severity`, and `Reported by`
- Avoid duplication

Triggers:
- After major changes
- Before PR/MR

---

### /review-branch
Full PR/MR-style review.

Responsibilities:
- Analyze `git diff`
- Detect bugs, regressions, missing tests
- Validate:
  - tests exist
  - Bruno docs updated (if API changed)
- Sync `.config/opencode/known_issues.md`

---

### /plan-feature
Feature planning with risk awareness.

Responsibilities:
- Break down feature into steps
- Identify impacted layers
- Register risks and planned work in `.config/opencode/known_issues.md`
- Add new entries with `Status: backlog` and the proper `Type`

---

### /promote <id>
Promote a tracked backlog item inside `known_issues.md`.

Responsibilities:
- Change `Status` from `backlog` or `ready` to `open`
- Reset `Remote` to `-` until the remote issue is created
- Prepare for remote issue creation

Usage:
```
/promote 2
```

Flow:
1. Promote tracked item inside `known_issues.md`
2. Create remote issue:
   `.config/opencode/scripts/create_issue.sh 2`
3. Start development on generated branch

Status rules:
- `backlog` and `ready` are planning states
- `open` is the pre-remote execution state
- `in-progress` requires a remote issue id
- `resolved` is the terminal state unless manually reopened
