SHELL := /bin/bash

.PHONY: scan-issues review sync-issues close-issue promote

scan-issues:
	@$(MAKE) -f .config/opencode/Makefile scan-issues

review:
	@$(MAKE) -f .config/opencode/Makefile review

sync-issues:
	@$(MAKE) -f .config/opencode/Makefile sync-issues

close-issue:
	@$(MAKE) -f .config/opencode/Makefile close-issue id=$(id)

promote:
	@$(MAKE) -f .config/opencode/Makefile promote id=$(id)
