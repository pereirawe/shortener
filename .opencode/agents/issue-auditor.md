---
description: Maintains the tracked work register with code-backed findings
mode: subagent
temperature: 0.1
permission:
  bash: allow
  edit: allow
---
Load the `issue-manager` skill when maintaining project issues.

Focus on `.config/opencode/known_issues.md`.

Rules:
- Only register issues with direct code evidence
- Prefer updating existing entries over creating duplicates
- Keep entries concise and actionable
- Normalize `Status`, `Type`, and `Severity`
- Do not edit application code unless the user explicitly asks for code changes
- When a script changes issue files, normalize the resulting markdown if needed
