## Why
AI agents operating locally lack persistent, structured memory across sessions, causing repeated analysis and inconsistent behavior. A simple CLI that provides durable, queryable context within each project directory will improve continuity and efficiency.

## What Changes
- Add a new Go-based CLI "Context Keeper" that manages per-project memory in `.context-keeper/` using a single SQLite database (`context.db`).
- Provide subcommands: `context init`, `context save --key --value [--channel --priority]`, `context get --key [--json]`, `context list [--channel --limit --json]`, `context delete --key`.
- Generate `AGENT_INSTRUCTIONS.md` on init with concise usage guidance.
- Human-readable output by default; support `--json` on data-retrieval commands.
- Non-interactive execution; use Bubble Tea/Bubbles/Lipgloss only for styled output (e.g., tables for `list`).
- Security baselines: parameterized SQL; restrict filesystem writes to project-local `.context-keeper/`.

## Impact
- Affected specs: `context-keeper` (new capability with requirements for init, save, get, list, delete, and safety constraints)
- Affected code (planned):
  - `cmd/context/main.go` (entrypoint, flag parsing)
  - `internal/cli/*` (command handlers)
  - `internal/db/repository.go` (SQLite access, schema init)
  - `internal/tui/*` (styled output models)
  - `.context-keeper/AGENT_INSTRUCTIONS.md` (generated at runtime)

## Notes
- Out of scope for this change (can be future proposals): file caching, checkpoints, compaction, search, import/export, knowledge graph, multi-agent features, branching/merging, journaling.
