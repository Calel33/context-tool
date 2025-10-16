# Project Context

## Purpose
Build and maintain a Context Keeper CLI that provides persistent, structured memory for AI agents within a project directory. It stores key/value context with channel and priority in a local SQLite database and supports both human‑readable and JSON output for automation.

## Tech Stack
- Language: Go 1.21
- Storage: SQLite (modernc.org/sqlite, pure‑Go)
- TUI styling: Bubble Tea, Bubbles, Lipgloss (non‑interactive use for formatting)
- Tests: Go testing (unit + integration)

## Project Conventions

### Code Style
- KISS and YAGNI; small, single‑responsibility files (≤500 lines)
- Parameterized SQL only; no string‑concatenated queries
- Minimal logs/output; support `--json` for machine use

### Architecture Patterns
- Vertical slice: `cmd/context` (entry), `internal/cli` (handlers), `internal/db` (repository), `internal/tui` (rendering)
- Filesystem writes restricted to `.context-keeper/` in the current project
- AGENT_INSTRUCTIONS generated from template on init

### Testing Strategy
- Unit tests for repository CRUD using temp DB paths
- Integration tests build the binary and run end‑to‑end (init/save/get/list/delete) with JSON and human outputs

### Git Workflow
- Feature branch → PR to `main`; include concise, why‑focused messages
- Use Co‑authored‑by when applicable; no secrets in commits
- OpenSpec changes land under `openspec/changes/` then archive to `openspec/changes/archive/`; canonical specs in `openspec/specs/`

## Domain Context
- CLI acts as local memory for agents; keys are unique; items have `channel` and `priority`
- Designed for Windows/PowerShell environment; executable may be called with PowerShell call operator `&` when path contains spaces

## Important Constraints
- Pure‑Go SQLite driver (no CGO); parameterized SQL everywhere
- Non‑interactive CLI; human mode may be styled but must not block
- Backwards compatible storage and commands once released

## External Dependencies
- None at runtime beyond the compiled binary; SQLite DB file lives at `.context-keeper/context.db`
