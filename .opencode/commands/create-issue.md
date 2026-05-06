---
description: Create the remote issue and local branch for a tracked item
agent: issue-auditor
subtask: true
---
Use known issue id or title/body arguments: `$ARGUMENTS`

Read:
- @.config/opencode/known_issues.md

Run the issue creation script with the provided arguments using bash:
`.config/opencode/scripts/create_issue.sh $ARGUMENTS`

Then summarize:
- the remote issue created
- the local branch name
- any updates written back to `.config/opencode/known_issues.md`
