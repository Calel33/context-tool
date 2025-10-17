## Why
Standardize the GUI on Charmbracelet Bubbles components to ensure consistent UX, accessibility, and maintainable patterns across list, filters, editor, and dialogs.

## What Changes
- Adopt Bubbles components for all GUI primitives: table, text input, textarea, and modal/confirmation patterns.
- Align styling with existing CLI design tokens (primary/success/error/muted) applied through Lipgloss.
- Specify non-blocking UI behavior and error surfacing via Bubbles-friendly toasts/modals.

## Impact
- Affected specs: context-keeper-gui
- Affected code: cmd/context-gui/main.go and related UI helpers
- No DB/CLI behavior changes; GUI-only architecture and patterns update.
