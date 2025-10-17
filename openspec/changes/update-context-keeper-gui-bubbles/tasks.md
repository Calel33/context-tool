## Main Goal
Adopt Charmbracelet Bubbles components across the GUI (table, inputs, textarea, confirmation) with Lipgloss tokens to deliver a consistent, non-blocking UX without changing CLI behavior.

## Subtask Goals
- 1.1 Confirm Bubbles usage and tokens: Ensure table/search/editor use Bubbles and align colors with CLI tokens.
- 1.2 Confirmation modal pattern: Provide clear, keyboard-accessible delete confirmation using Bubbles-compatible view.
- 1.3 Non-blocking toasts: Show success/error toasts without freezing UI; errors must not crash the app.
- 1.4 Keyboard/help affordances: Short, discoverable keybinds and inline help text consistent with Bubbles patterns.
- 1.5 Validation: All tests pass and manual CRUD flows succeed; no CLI contract changes.
## 1. Implementation
- [ ] 1.1 Confirm Bubbles usage for list (table), inputs, and textarea; integrate Lipgloss tokens
- [ ] 1.2 Add confirmation modal pattern for delete using Bubbles-compatible view
- [ ] 1.3 Ensure non-blocking updates with toast messages on save/delete/open errors
- [ ] 1.4 Validate keyboard affordances and help text per patterns
- [ ] 1.5 Run tests and manual validation; no CLI behavioral change
