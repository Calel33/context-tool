# AGENTS.md · Project Rules for AI Agents

## Environment
- OS: Windows
- Shell: PowerShell (prefer `&` call operator for paths with spaces)
- Repo root: this file governs this project; defer to openspec/specs for canonical requirements

## Core Principles
- KISS and YAGNI: minimal, straightforward implementations only
- Single‑responsibility files (≤500 lines); vertical slice architecture
- No secrets in code, commits, or logs

## Domains
- CLI: Go 1.21, uses `internal/db` (SQLite via modernc.org/sqlite), `internal/cli`, `internal/tui`
- GUI: "bubble" library UI over the same repository; non‑blocking, minimal telemetry

## Security & Storage
- Parameterized SQL exclusively; never concatenate user input into queries
- All storage restricted to `.context-keeper/` under the selected project directory
- No network calls unless explicitly requested by specs

## Output Modes
- Human readable (styled) and `--json` machine output for data commands
- GUI must remain responsive; surface errors via modal/toast without crashing

## Testing
- Keep unit tests for repository CRUD; integration tests for CLI end‑to‑end
- For GUI, manual E2E validation per OpenSpec tasks before shipping

## Git & OpenSpec
- Branch from `main`; open PRs with concise, why‑focused messages and Co‑authored‑by when applicable
- OpenSpec flow: propose under `openspec/changes/<id>/`, validate, land specs under `openspec/specs/`, archive change under `openspec/changes/archive/`

## Style Tokens
- Reuse existing CLI design tokens (primary, success, error, muted) where applicable in GUI
- No hard‑coded style values outside token definitions

## Operational Rules
- Prefer repo utilities and existing patterns; do not introduce new deps without need
- Avoid breaking CLI contract or DB schema; additive changes only unless spec says otherwise
- Respect Windows/PowerShell behaviors (e.g., search path, CRLF)
