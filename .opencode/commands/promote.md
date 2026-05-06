---
description: Promote a tracked backlog item inside known issues
agent: issue-auditor
subtask: true
---
Use known issue id `$ARGUMENTS`.

Read:
- @.config/opencode/known_issues.md

Run the promotion script:
`.config/opencode/scripts/promote.sh $ARGUMENTS`

Then review the resulting files and normalize formatting if needed.

Finally, explain the next step to create the remote issue using `.config/opencode/scripts/create_issue.sh $ARGUMENTS`.
