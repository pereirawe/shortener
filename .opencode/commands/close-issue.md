---
description: Close a tracked issue locally and remotely
agent: issue-auditor
subtask: true
---
Use known issue id `$ARGUMENTS`.

Read:
- @.config/opencode/known_issues.md

Run the close script:
`.config/opencode/scripts/close_issue.sh $ARGUMENTS`

Then verify `.config/opencode/known_issues.md` and normalize formatting if needed.
