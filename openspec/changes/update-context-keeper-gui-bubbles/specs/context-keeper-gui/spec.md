## MODIFIED Requirements

### Requirement: List and Filter Items
The GUI SHALL display stored items using Bubbles components (table for list, textinput for search), with filters for channel, priority, and key text search.

#### Scenario: Filter by channel
- WHEN the user selects a channel filter
- THEN only items matching the channel are shown

#### Scenario: Text search on key
- WHEN the user types in the search box
- THEN the table updates to rows whose key contains the text (case-insensitive)

### Requirement: Create and Update Items
The GUI SHALL use Bubbles textinput and textarea for the editor form (key, value) and provide select-like controls for channel and priority.

#### Scenario: Create new item
- WHEN the user fills the form (key, value, channel, priority) and clicks Save
- THEN the item is stored and appears in the list

#### Scenario: Update item
- WHEN the user selects an existing row, edits the value or metadata, and clicks Save
- THEN the item is updated (key unchanged) and list refreshes

### Requirement: Delete Item
The GUI SHALL present a confirmation view/modal using Bubbles-compatible UI prior to deletion.

#### Scenario: Delete success
- WHEN the user confirms item deletion
- THEN the item is removed from the database and the list refreshes

### Requirement: Non-Blocking UI and Error Handling
The GUI SHALL surface errors via non-blocking toast/modal components compatible with Bubbles and SHALL not crash the app.

#### Scenario: Database error
- WHEN a DB operation fails
- THEN an error message is shown and the user can continue interacting or retry
