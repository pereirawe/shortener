SHELL := /bin/bash

.PHONY: scan-issues review help sync-issues close-issue promote

help:
	@echo "Targets:"
	@echo "  make scan-issues  - run local scan + prompt assistant command"
	@echo "  make review       - show git diff + prompt assistant command"

scan-issues:
	@echo "[make] scan-issues"
	@chmod +x scripts/scan_issues.sh || true
	@./scripts/scan_issues.sh
	@echo "\nNext step: run command in assistant -> /scan-issues"

review:
	@echo "[make] review"
	@git status --porcelain
	@echo "\n--- DIFF (staged + unstaged) ---\n"
	@git diff
	@echo "\nNext step: run command in assistant -> /review-branch"

sync-issues:
	@echo "[make] sync-issues (basic)"
	@echo "This should reconcile local known_issues with remote issues (manual step for now)"

close-issue:
	@./scripts/close_issue.sh $(id)

promote:
	@./scripts/promote.sh $(id)
