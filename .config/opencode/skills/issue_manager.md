## Issue Manager Skill Notes

The operational OpenCode skill lives at:
`/.opencode/skills/issue-manager/SKILL.md`

This document remains as project documentation for the skill behavior.

Behavior summary:
- Maintain `.config/opencode/known_issues.md` as the source of truth for tracked work
- Scan code for bugs, security issues, concurrency problems, and design drift
- Update existing entries before creating new ones
- Remove or mark resolved issues when code evidence supports it
- Prefer file and line references when possible
- Keep `Reported by` aligned with the reporter: user name for humans, model name for AI
