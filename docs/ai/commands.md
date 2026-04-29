## Commands

These commands define how the assistant operates over this project context.

### /scan-issues
Deep analysis of the codebase.

Responsibilities:
- Detect new issues (security, concurrency, architecture)
- Update `docs/ai/known_issues.md`
- Normalize statuses (open/in-progress/resolved)
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
- Sync `known_issues.md`

---

### /plan-feature
Feature planning with risk awareness.

Responsibilities:
- Break down feature into steps
- Identify impacted layers
- Register risks in `known_issues.md`
- Add feature candidates into `improvements.md`

---

### /sync-improvements
Sync improvements backlog.

Responsibilities:
- Convert items from `improvements.md` into structured entries
- Ensure format consistency
- Prepare items for issue creation

---

### /promote <id>
Promote backlog item into a tracked issue.

Responsibilities:
- Move item from `improvements.md` → `known_issues.md`
- Assign Status: open
- Prepare for remote issue creation

Usage:
```
/promote 2
```

Flow:
1. Promote improvement → known_issues
2. Create remote issue:
   ./scripts/create_issue.sh 2
3. Start development on generated branch
