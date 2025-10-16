## ADDED Requirements

### Requirement: GUI Launch and Project Resolution
The app SHALL open a GUI window and resolve the target project directory to locate `.context-keeper/context.db`.

#### Scenario: Launch in project root
- WHEN a user starts the GUI from a project root
- THEN it loads `.context-keeper/context.db` (creating schema if needed) and shows the list of items

#### Scenario: Open Project folder
- WHEN a user selects File > Open Project and chooses a folder
- THEN the app resolves the folder’s `.context-keeper/context.db` and reloads data

### Requirement: List and Filter Items
The GUI SHALL display stored items in a table with filters for channel, priority, and key text search.

#### Scenario: Filter by channel
- WHEN the user selects a channel filter
- THEN only items matching the channel are shown

#### Scenario: Text search on key
- WHEN the user types in the search box
- THEN the table updates to rows whose key contains the text (case-insensitive)

### Requirement: Create and Update Items
The GUI SHALL allow creating a new item and updating an existing item’s value, channel, and priority.

#### Scenario: Create new item
- WHEN the user fills the form (key, value, channel, priority) and clicks Save
- THEN the item is stored and appears in the list

#### Scenario: Update item
- WHEN the user selects an existing row, edits the value or metadata, and clicks Save
- THEN the item is updated (key unchanged) and list refreshes

### Requirement: Delete Item
The GUI SHALL allow deleting an item with confirmation.

#### Scenario: Delete success
- WHEN the user confirms item deletion
- THEN the item is removed from the database and the list refreshes

### Requirement: Non-Blocking UI and Error Handling
The GUI SHALL remain responsive; errors SHALL be surfaced in a modal/toast and not crash the app.

#### Scenario: Database error
- WHEN a DB operation fails
- THEN an error message is shown and the user can continue interacting or retry
