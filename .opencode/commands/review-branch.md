---
description: Review the current branch like a pull request
agent: branch-reviewer
subtask: true
---
Review the current branch changes and report findings first.

Use this context:
- @.config/opencode/known_issues.md
- @.config/opencode/conventions.md
- @.config/opencode/pr_template.md

Repository state:
!`git status --short`

Recent commits:
!`git log --oneline --decorate -10`

Diff to review:
!`git diff`

If the code review reveals new tracked issues or resolved items, update `.config/opencode/known_issues.md`.
