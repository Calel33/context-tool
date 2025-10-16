# Design: Context Keeper GUI

## Architecture
- Reuse Go module `context-keeper` with a new `cmd/context-gui` entry point.
- UI layer (bubble) interacts with `internal/db` repository for CRUD.
- A small adapter package `internal/gui` encapsulates view state, events, and mapping between UI widgets and repository models.

## Data Flow
1) App boot resolves project path (default: current dir); loads `.context-keeper/context.db`.
2) Repository initializes schema if needed (idempotent).
3) UI shows list view with filters and an editor panel.
4) Actions dispatch to repo (SaveItem, DeleteItem); list refreshes.

## Views
- List: table with columns [Key, Channel, Priority, Updated]
- Filters: channel dropdown, priority dropdown, text search on key
- Editor: key (disabled when editing), value (multiline), channel, priority, Save button
- Delete: button with confirm dialog

## Integration
- No schema change; respects parameterized SQL and storage path constraints.
- Shared style tokens: reuse lipgloss color tokens where applicable for consistency.

## Telemetry/Logging
- Minimal; surface errors in a non-blocking modal/toast and stderr logs in dev builds.
