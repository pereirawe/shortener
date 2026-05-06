SHELL := /bin/bash

OPENCODE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ROOT_DIR := $(abspath $(OPENCODE_DIR)/../..)
.PHONY: scan-issues review help sync-issues close-issue promote

help:
	@echo "Targets:"
	@echo "  make scan-issues  - run local scan + prompt assistant command"
	@echo "  make review       - show git diff + prompt assistant command"
	@echo "  make promote id=n - move known_issues status to open"
	@echo "  make close-issue id=n - mark known_issues entry as resolved"

scan-issues:
	@echo "[make] scan-issues"
	@chmod +x $(OPENCODE_DIR)scripts/scan_issues.sh || true
	@cd $(ROOT_DIR) && $(OPENCODE_DIR)scripts/scan_issues.sh
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
	@cd $(ROOT_DIR) && $(OPENCODE_DIR)scripts/close_issue.sh $(id)

promote:
	@cd $(ROOT_DIR) && $(OPENCODE_DIR)scripts/promote.sh $(id)
