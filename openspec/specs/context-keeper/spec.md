## ADDED Requirements

### Requirement: Project Initialization & Session Management
The CLI SHALL initialize a project-local storage at `.context-keeper/` containing `context.db` and generate `AGENT_INSTRUCTIONS.md` on first run.

#### Scenario: Successful initialization
- **WHEN** a user runs `context init` in a project directory where `.context-keeper/` does not exist
- **THEN** the directory is created with `context.db` and `AGENT_INSTRUCTIONS.md`, and a success message is printed

#### Scenario: Idempotent initialization
- **WHEN** `context init` is run and `.context-keeper/` already exists
- **THEN** the command completes without destroying existing data and reports that storage already exists

### Requirement: Context Storage (Save)
The CLI SHALL save or update key-value pairs with optional `channel` and `priority` using parameterized SQL and store them in `context.db`.

#### Scenario: Save with defaults
- **WHEN** a user runs `context save --key k --value v`
- **THEN** the item is stored with `channel=general` and `priority=normal` and a confirmation is printed

#### Scenario: Save with channel and priority
- **WHEN** a user runs `context save --key k --value v --channel c --priority p`
- **THEN** the item is stored with the provided `channel` and `priority`

### Requirement: Retrieve Single Item (Get)
The CLI SHALL retrieve a single item by key and support `--json` output.

#### Scenario: Get human-readable
- **WHEN** a user runs `context get --key k`
- **THEN** the value is printed in a human-readable format

#### Scenario: Get JSON
- **WHEN** a user runs `context get --key k --json`
- **THEN** the output is a JSON object containing key, value, channel, and priority

### Requirement: List Items
The CLI SHALL list items with optional `--channel` and `--limit` filters and support human-readable table and `--json` output.

#### Scenario: List human-readable table
- **WHEN** a user runs `context list`
- **THEN** a styled table of items is printed with columns CHANNEL, PRIORITY, KEY

#### Scenario: List JSON
- **WHEN** a user runs `context list --json`
- **THEN** a JSON array of items is printed including key, value, channel, priority

### Requirement: Delete Item
The CLI SHALL delete an item by key.

#### Scenario: Delete success
- **WHEN** a user runs `context delete --key k`
- **THEN** the item is removed and a confirmation message is printed

### Requirement: Non-Interactive Styled Output
The CLI SHALL execute non-interactively and MAY use Bubble Tea/Bubbles/Lipgloss for formatting; data commands MUST support `--json`.

#### Scenario: Non-interactive run
- **WHEN** any command is executed
- **THEN** it performs its work, prints output (styled for human-readable mode), and exits

### Requirement: Security & Storage Constraints
All SQL operations MUST be parameterized; filesystem writes MUST be limited to `.context-keeper/` under the current project directory.

#### Scenario: Parameterized SQL only
- **WHEN** saving or querying data
- **THEN** the implementation uses prepared statements or placeholders rather than string concatenation

#### Scenario: Storage path restricted
- **WHEN** initializing or writing files
- **THEN** only files inside `.context-keeper/` are created or modified
