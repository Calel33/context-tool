# Proposal: Context Keeper GUI

## Summary
Add a graphical user interface (GUI) for the existing Context Keeper, built with the selected "bubble" library. The GUI provides a simple window to view, create, edit, and delete context items backed by the same local SQLite storage used by the CLI.

## Goals
- Deliver a minimal, fast GUI for common operations: list, filter, add, edit, delete.
- Reuse existing repository logic and database schema without breaking changes.
- Keep styling consistent with our design tokens where applicable.

## Non-Goals
- No full-text/semantic search, journaling, or multi-agent visualization.
- No sync/cloud features; storage remains local `.context-keeper/context.db`.
- No change to CLI contract or database schema beyond safe additive UI needs.

## Users
- Human developers who prefer a visual interface to inspect and manage stored context.
- AI agents remain CLI-first; GUI is a complementary tool for humans.

## Scope
- Single-window GUI application with: item table, filters (channel, priority), form to create/update item, and delete action.
- Read/write the same DB file resolved from the current working directory (or chosen project folder via File > Open Project).

## Risks & Mitigations
- Risk: Divergent behaviors between CLI and GUI. Mitigation: call into shared `internal/db` repository and reuse validation.
- Risk: Platform specifics. Mitigation: begin with Windows support; keep abstractions simple to enable later cross-platform.
