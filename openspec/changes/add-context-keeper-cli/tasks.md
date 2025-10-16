## 1. Implementation
- [ ] 1.1 Initialize Go module and project structure (cmd/context, internal/{cli,db,tui})
- [ ] 1.2 Add dependencies: bubbletea, bubbles, lipgloss, go-sqlite3
- [ ] 1.3 Implement database repository and schema creation (context_items table)
- [ ] 1.4 Implement CLI handlers: init, save, get, list, delete
- [ ] 1.5 Implement TUI renderers for list (table) and messages
- [ ] 1.6 Wire subcommands and flags in main.go; support --json where applicable
- [ ] 1.7 Generate AGENT_INSTRUCTIONS.md from template on init
- [ ] 1.8 Add unit tests for repository; integration tests for CLI handlers
- [ ] 1.9 Add --help text for root and subcommands
- [ ] 1.10 Build and manual test on sample project

## 2. Validation
- [ ] 2.1 Static checks: build succeeds; no unused deps
- [ ] 2.2 Security checks: parameterized SQL only; writes constrained to .context-keeper
- [ ] 2.3 Verify human-readable and --json outputs match spec examples
- [ ] 2.4 Confirm non-interactive execution model

## 3. Completion
- [ ] 3.1 Update this checklist to all checked when done
