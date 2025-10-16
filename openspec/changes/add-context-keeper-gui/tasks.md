# Tasks: Context Keeper GUI

1. Scaffold GUI entry point
   - Create `cmd/context-gui/main.go` with application bootstrap (project path resolution, repo init).
   - Wire bubble app root and basic window with list placeholder.

2. Repository reuse
   - Import and reuse `internal/db` repository for CRUD.
   - Add small helpers to adapt repo models to UI rows.

3. List view + filters
   - Implement table showing key, channel, priority, updated_at.
   - Add channel/priority filters and key search; update query accordingly.

4. Editor panel
   - Implement create/update form: key (text), value (multiline), channel (select), priority (select), Save.
   - Disable key editing in update mode; support insert-or-replace semantics.

5. Delete flow
   - Add Delete action with confirmation dialog; refresh list after success.

6. Project selection
   - Add File > Open Project to choose a folder; resolve `.context-keeper/context.db` under it.

7. Styling
   - Apply design tokens analogous to CLI (primary/success/error/muted) where supported by bubble components.

8. Packaging
   - Provide build script/notes for Windows binary `context-gui.exe`.

9. Validation
   - Manual test: full CRUD via GUI against a temp project folder.
   - Ensure no changes to CLI behavior; run existing unit/integration tests.
