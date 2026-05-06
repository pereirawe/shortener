---
description: Scan the project and sync known issues
agent: issue-auditor
subtask: true
---
Load the `issue-manager` skill if available.

Review the codebase and update `.config/opencode/known_issues.md`.

Use this context:
- @.config/opencode/known_issues.md
- @.config/opencode/context.md
- @.config/opencode/architecture.md

Static scan output:
!`.config/opencode/scripts/scan_issues.sh`

Rules:
- keep only code-backed issues
- remove duplicates
- mark resolved issues when the code already fixes them
- normalize `Status`, `Type`, `Severity`, and `Reported by`
- keep entries concise and actionable
