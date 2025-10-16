## Context
Provide a per-project, persistent memory for AI developer agents via a simple, non-interactive CLI. Storage is local to the project in `.context-keeper/` using a single SQLite DB; outputs are human-readable with optional JSON for machine parsing.

## Goals / Non-Goals
- Goals: Init storage; save/get/list/delete context; styled human-readable output; `--json` for data retrieval; safe parameterized SQL; filesystem writes confined to project.
- Non-Goals: Advanced features (file caching, checkpoints, compaction, search, import/export, knowledge graph, multi-agent, branching/merging, journaling).

## Decisions
- Language: Go (single binary, strong stdlib).
- CLI/TUI: charmbracelet/bubbletea (non-interactive render), bubbles (table), lipgloss (styling).
- DB: SQLite with `mattn/go-sqlite3`; single table `context_items(key TEXT PRIMARY KEY, value TEXT NOT NULL, channel TEXT DEFAULT 'general', priority TEXT DEFAULT 'normal', created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`.
- Arg parsing: Go `flag` package.
- Project structure: `cmd/context/main.go`; `internal/cli/*` handlers; `internal/db/repository.go`; `internal/tui/*`; template for `AGENT_INSTRUCTIONS.md`.
- Output: human-readable styled views; `--json` flag emits structured JSON for `get` and `list`.
- Security: prepared statements only; path writes limited to `.context-keeper/`.

## Risks / Trade-offs
- Plaintext SQLite (no app-level encryption) → Simplicity; relies on OS permissions/disk encryption.
- Using Bubble Tea without interactivity → Slight overhead for styling, cleaner separation of view logic.
- Single-table schema → Simpler to start; may evolve with future capabilities.

## Migration Plan
Initial release only; future proposals will add advanced features or schema evolutions with migration steps.

## Open Questions
- Do we require configurable storage location beyond project root? (Default is project-local only.)
