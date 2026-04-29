## Development Workflow

### 1. Select Issue
- Developer selects issue from `known_issues.md`

### 2. Create Remote Issue
Run:
```
./scripts/create_issue.sh "title" "description"
```

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
  - Remove issue from `known_issues.md`
  - Keep tracked in remote (GitHub/GitLab)
