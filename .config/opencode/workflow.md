## Development Workflow

### 1. Select Issue
- Developer selects issue from `known_issues.md`

Status rules:
- `backlog`: captured but not ready
- `ready`: refined and approved
- `open`: selected locally and waiting for remote issue creation
- `in-progress`: remote issue created and work started
- `resolved`: completed or closed

Allowed transitions:
- `backlog -> ready`
- `ready -> backlog`
- `backlog -> open`
- `ready -> open`
- `open -> in-progress`
- `open -> resolved`
- `in-progress -> resolved`
- `resolved -> open` only by explicit manual reopen

### 2. Create Remote Issue
Run:
```
.config/opencode/scripts/create_issue.sh <local_issue_id>
```

If the item is still `backlog` or `ready`, promote it first:
```
.config/opencode/scripts/promote.sh <local_issue_id>
```

Rules enforced by scripts:
- `promote.sh`: only `backlog` or `ready` can move to `open`
- `create_issue.sh`: only `open` can move to `in-progress`
- `close_issue.sh`: only `open` or `in-progress` can move to `resolved`

### 3. Branch Naming
- Pattern: `issue-<id>-<slug>`

### 4. Development Rules
- Tests are mandatory
- API changes must update Bruno docs
- Follow conventions

### 5. Pre-commit
- Runs tests
- Warns if `known_issues` not updated

### 6. Pull/Merge Request
- Must include:
  - Tests passing
  - Issue reference
  - Updated docs

### 7. Merge
- After merge:
  - Mark the entry as `resolved` in `known_issues.md`
  - Keep tracked in remote (GitHub/GitLab)
