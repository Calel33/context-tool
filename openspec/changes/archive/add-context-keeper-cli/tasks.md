## 1. Implementation
- [x] 1.1 Initialize Go module and project structure (cmd/context, internal/{cli,db,tui})
- [x] 1.2 Add dependencies: bubbletea, bubbles, lipgloss, sqlite driver
- [x] 1.3 Implement database repository and schema creation (context_items table)
- [x] 1.4 Implement CLI handlers: init, save, get, list, delete
- [x] 1.5 Implement TUI renderers for list (table) and messages
- [x] 1.6 Wire subcommands and flags in main.go; support --json where applicable
- [x] 1.7 Generate AGENT_INSTRUCTIONS.md from template on init
- [ ] 1.8 Add unit tests for repository; integration tests for CLI handlers
- [x] 1.8 Add unit tests for repository; integration tests for CLI handlers
- [x] 1.9 Add --help text for root and subcommands
- [x] 1.10 Build and manual test on sample project

## 2. Validation
- [x] 2.1 Static checks: build succeeds; no unused deps
- [x] 2.2 Security checks: parameterized SQL only; writes constrained to .context-keeper
- [x] 2.3 Verify human-readable and --json outputs match spec examples
- [x] 2.4 Confirm non-interactive execution model

## 3. Completion
- [x] 3.1 Update this checklist to all checked when done
