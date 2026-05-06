## Support Docs

The operational OpenCode entrypoints for this repository are:
- `AGENTS.md`
- `opencode.json`
- `.opencode/agents/`
- `.opencode/commands/`
- `.opencode/skills/`

This directory keeps supporting project workflow documentation and local helper scripts.

## Source Of Truth

- Tracked work register: `.config/opencode/known_issues.md`
- Operational scripts: `.config/opencode/scripts/`

## Local Helpers

Preferred commands:
```
make -f .config/opencode/Makefile scan-issues
make -f .config/opencode/Makefile review
make -f .config/opencode/Makefile promote id=<n>
make -f .config/opencode/Makefile close-issue id=<n>
```

Direct scripts:
```
.config/opencode/scripts/create_issue.sh <id|"title" "body">
.config/opencode/scripts/promote.sh <id>
.config/opencode/scripts/close_issue.sh <id>
```

## Notes

- Keep `.config/opencode/known_issues.md` current when fixing or discovering issues
- Tests are mandatory for code changes
- API changes must update `bruno/`
