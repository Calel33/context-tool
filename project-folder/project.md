### Phase 1: Discovery
This first phase is about defining the project's core vision, its users, and the value it provides. Based on your request, I've drafted the initial deliverables for this phase.
1. Problem Statement
AI agents operating in a local development environment lack persistent, structured memory. They lose valuable context between sessions or when their context window is reset, forcing them to re-learn project details, decisions, and history. This leads to redundant analysis, inconsistent behavior, and inefficient use of time and resources.
2. Value Proposition
The CLI tool will act as a persistent, queryable memory for an AI agent, enabling it to maintain context, recall past interactions, and make more informed decisions across multiple sessions and tasks. By offering a simple, command-driven interface for context management, it dramatically improves the agent's efficiency and continuity on long-term projects.
3. Personas (Target Audience)
Primary Persona: The AI Developer Agent
Description: An automated or semi-automated AI entity tasked with software development. It can read/write files, execute shell commands, and process CLI output.
Needs:
A simple, reliable command set to save and retrieve information.
Structured output (like JSON) for easy parsing.
A way to organize information by topic or task (channels/sessions).
A clear "how-to" guide to understand and use the tool's commands.
Pain Points: Forgetting previous steps, losing context after a restart, and re-analyzing the same files repeatedly.
Secondary Persona: The Human Developer
Description: The developer who directs the AI agent and occasionally inspects the agent's memory.
Needs:
A human-readable way to see what the agent has stored.
Simple commands to manage sessions or review stored context.
Confidence that the agent's memory is persistent and reliable.
Pain Points: Managing the agent's "state" manually, repeating instructions, and not knowing what the agent "knows."
Project Requirement Document (PRD): Context Keeper CLI
1. Project Title: Context Keeper CLI
2. Objective:
To create a standalone Go-based CLI application that provides a simple, persistent, and structured memory for an AI developer agent. The tool will enable the agent to save, retrieve, and organize context within a specific project directory, improving its long-term effectiveness and continuity.
3. Target Audience:
Primary: AI Developer Agents that can execute shell commands and parse structured output (JSON).
Secondary: Human developers overseeing the AI agent, who may need to inspect or manage the stored context.
4. Scope & Core Features:
This initial version will focus exclusively on the foundational context management features.
Feature: Project Initialization & Session Management
Description: The first time the tool is run within a directory, it will set up a dedicated folder (.context-keeper/) to store its database and instruction file. Each project directory is treated as a unique session.
Commands:
context init: Manually sets up the storage in the current directory.
User Stories:
As an AI agent, I want to run a setup command in a new project directory so that I have a persistent memory store ready to use.
As a human developer, I want the tool's data to be stored locally within my project so that it's self-contained and easy to manage with version control.
Feature: Context Storage
Description: The core functionality for the agent to save, retrieve, and delete key-value pairs of information. Context is organized into "channels" for better topic separation.
Commands:
context save --key "..." --value "..." [--channel "..."] [--priority "..."]: Saves a piece of information.
context get --key "..." [--json]: Retrieves a specific piece of information.
context list [--channel "..."] [--limit N] [--json]: Retrieves a list of stored items, with optional filters.
context delete --key "...": Removes a piece of information.
User Stories:
As an AI agent, I want to save a file summary with a specific key so that I don't have to read the file again later.
As an AI agent, I want to retrieve a decision I made earlier by its key so that I can apply it consistently.
As an AI agent, I want to list all context items in the 'testing' channel to review my progress on tests.
Feature: Agent Instruction File
Description: Upon initialization, the tool will generate a markdown file (AGENT_INSTRUCTIONS.md) within the .context-keeper directory. This file will contain a clear, concise guide for the AI agent on how to use the CLI commands.
User Stories:
As a human developer, I want an automatically generated instruction file so I can easily direct my AI agent to learn and use the tool.
5. Key Requirements:
Technology: Must be built in Go using the Bubble Tea and Bubbles libraries.
Data Storage: A single SQLite database file (context.db) will be stored in a .context-keeper folder in the project's root directory.
Output Format: Commands that retrieve data (get, list) must support a --json flag for easy parsing by the AI agent. The default output should be human-readable.
Interaction Model: The tool is non-interactive (no TUI). It executes a command, prints the output, and exits. The TUI components from Bubbles will be used for styling the output (e.g., tables for list).
6. Out of Scope (for this version):
The following advanced features from the original mcp-memory-keeper will not be included in this core build:
File Caching (and automatic change detection)
Checkpoints
Smart Compaction
Full-text or Semantic Search
Export/Import functionality
Knowledge Graph and Visualization
Multi-Agent System
Session Branching/Merging
Journaling and Timelines
Core Features List
Here is a summary of the primary commands we will build:
Command	Description
context init	Initializes the .context-keeper storage in the current project directory.
context save --key "k" --value "v"	Saves a key-value pair. Supports optional --channel and --priority.
context get --key "k"	Retrieves the value for a given key. Supports --json for structured output.
context list	Lists all saved context items. Supports filtering and --json.
context delete --key "k"	Deletes a context item by its key.
This App Flow Document will serve as the blueprint for the CLI's behavior.
App Flow Document: Context Keeper CLI
1. General Principles:
Execution Model: Each command is executed from the shell, performs its action, prints output to stdout (or stderr for errors), and exits.
Error Handling: If a command fails (e.g., database not found, invalid key), it will print a styled error message to stderr and exit with a non-zero status code.
Output Formatting: All data-retrieving commands (get, list) will default to a human-readable format and offer a --json flag for machine-readable output.
2. Command Flows
Flow: context init
Purpose: Sets up the necessary storage and instruction files in a new project directory.
Logic:
The command is run in the root of the user's project.
It checks for the existence of a .context-keeper directory.
If it doesn't exist, it creates the directory.
It creates the SQLite database file: .context-keeper/context.db.
It generates the AGENT_INSTRUCTIONS.md file inside .context-keeper/.
It prints a success message to stdout.
Human-Readable Output:
code
Code
✅ Context Keeper initialized successfully.
- Storage created at: ./.context-keeper/context.db
- Instructions for your AI agent are in: ./.context-keeper/AGENT_INSTRUCTIONS.md
JSON Output (--json):
code
JSON
{
  "success": true,
  "message": "Context Keeper initialized successfully.",
  "storagePath": "./.context-keeper/context.db",
  "instructionsPath": "./.context-keeper/AGENT_INSTRUCTIONS.md"
}
User Scenarios:
A developer starts a new project, runs context init, and then instructs their AI agent to read the AGENT_INSTRUCTIONS.md file to learn how to use the tool.
An AI agent, upon entering a new project directory, runs context init to set up its own memory system before starting its analysis.
Flow: context save
Purpose: Saves or updates a key-value pair in the context database.
Logic:
The command is run with --key and --value flags.
It parses optional --channel and --priority flags (defaulting to 'general' and 'normal').
It connects to .context-keeper/context.db.
It executes an INSERT OR REPLACE SQL statement to save the data.
It prints a confirmation to stdout.
Human-Readable Output:
code
Code
✅ Context saved for key: 'api_endpoint'
JSON Output (--json):
code
JSON
{
  "success": true,
  "key": "api_endpoint",
  "action": "saved"
}
User Stories:
As an AI agent, I want to save the project's main entry point by running context save --key "main_file" --value "src/index.js" so I can find it easily later.
As an AI agent, after summarizing a complex file, I want to store the summary with a high priority using context save --key "file_summary_auth" --value "..." --priority "high".
Flow: context list
Purpose: Retrieves a list of saved context items.
Logic:
The command is run.
It parses optional --channel, --limit, and --json flags.
It connects to the database and executes a SELECT query with the appropriate filters.
It formats the results as either a styled table (human-readable) or a JSON array.
Human-Readable Output (Styled with Bubbles):
code
Code
CHANNEL   PRIORITY   KEY
────────────────────────────────────────
 general   high       main_file
 general   normal     project_description
 testing   normal     test_setup_notes
JSON Output (--json):
code
JSON
[
  {
    "key": "main_file",
    "value": "src/index.js",
    "channel": "general",
    "priority": "high"
  },
  {
    "key": "project_description",
    "value": "A web application for...",
    "channel": "general",
    "priority": "normal"
  }
]
User Scenarios:
An AI agent runs context list --channel "database" --json to get all stored information related to the database schema and connection strings.
A developer runs context list to get a quick overview of what the agent has saved recently.
Flow: context get
Purpose: Retrieves the value for a single, specific key.
Logic:
The command is run with a required --key flag.
It connects to the database and executes a SELECT query for that key.
If found, it prints the value. If not, it prints an error.
Human-Readable Output:
code
Code
Value for key 'main_file':
src/index.js
JSON Output (--json):
code
JSON
{
  "key": "main_file",
  "value": "src/index.js",
  "channel": "general",
  "priority": "high"
}
User Scenarios:
Before editing a file, an AI agent runs context get --key "main_file" --json to retrieve the path to the main application entry point it had saved earlier.
Flow: context delete
Purpose: Removes a context item from the database.
Logic:
The command is run with a required --key flag.
It connects to the database and executes a DELETE statement for that key.
It prints a confirmation message.
Human-Readable Output:
code
Code
🗑️ Context for key 'old_api_key' deleted.
JSON Output (--json):
code
JSON
{
  "success": true,
  "key": "old_api_key",
  "action": "deleted"
}
User Scenarios:
After refactoring a module, an AI agent runs context delete --key "summary_module_v1" to remove the outdated summary.


### Design System Guidelines: Context Keeper CLI
1. Philosophy:
Clarity: Output must be immediately understandable. Use color and weight to guide the user's eye to the most important information.
Consistency: The same types of information (keys, values, success messages) should always be styled the same way.
Minimalism: Use a limited color palette. Avoid visual clutter. The interface should feel clean and fast.
2. Color Palette (Tokens):
We will use a simple, high-contrast palette that works well on both light and dark terminal themes.
Primary / Highlight: A vibrant color for user input, keys, and selected items.
Token: lipgloss.Color("212") (a bright magenta/pink)
Success: A clear green for confirmation messages.
Token: lipgloss.Color("42") (a vibrant green)
Error: A distinct red for error messages.
Token: lipgloss.Color("196") (a bright red)
Subtle / Muted: A soft gray for secondary information like help text, borders, and separators.
Token: lipgloss.Color("240") (a medium gray)
Normal Text: The terminal's default foreground color.
3. Typography (Font Tokens):
Font styles will be used to create a visual hierarchy.
Headers (e.g., table headers): Bold.
Keys / Identifiers: Bold and styled with the Primary / Highlight color.
Values: Normal weight, default text color.
Help Text: Normal weight, styled with the Subtle / Muted color.
4. Layout & Spacing Rules:
Global Margin: All output will have a left margin of 2 spaces to provide breathing room from the edge of the terminal.
Padding: Table cells and other bordered elements will have a horizontal padding of 1 space on each side.
Separators: When listing multiple items on a single line (like in help text), a dot (•) styled with the Subtle / Muted color will be used as a separator.
Component-Specific Styles
Here’s how these rules will apply to the specific outputs we defined in the App Flow:
Success Message (context save, context delete)
Structure: [Icon] [Message]
Styling:
Icon: ✅ or 🗑️
Message Text: Styled with the Success color.
Example: ✅ Context saved for key: 'api_endpoint'
Error Message
Structure: [Icon] [Message]
Styling:
Icon: ❌
Message Text: Styled with the Error color.
Example: ❌ Error: Key 'api_endpoint' not found.
Table Output (context list)
Headers: Bold, with a subtle underline or bottom border.
Cell (Key): Styled with Primary / Highlight color.
Cell (Other): Default text color.
Borders: Styled with Subtle / Muted color.
Key-Value Output (context get)
Structure:
code
Code
Value for key '[KEY]':
[VALUE]
Styling:
The label Value for key: will be normal text.
The [KEY] will be styled with the Primary / Highlight color.
The [VALUE] will be displayed as-is with default text color.

### Tech Stack Document
This document specifies the core technologies for building the CLI application.
Language: Go
Reasoning: It's fast, produces a single cross-platform binary, has a strong standard library, and is the language in which Bubble Tea and Bubbles are written.
CLI & TUI Framework: charmbracelet/bubbletea
Reasoning: While our app is non-interactive, Bubble Tea provides a robust structure for managing I/O (like database calls) and rendering the final, styled output in a clean, declarative way.
UI Components & Styling: charmbracelet/bubbles & charmbracelet/lipgloss
Reasoning: We will use the bubbles/table component for the context list command's output. lipgloss will be used to implement the Design System we defined in Phase 4 for all styled text (colors, borders, etc.).
Database: SQLite (via mattn/go-sqlite3 driver)
Reasoning: A single-file, serverless database is perfect for this use case. It's portable, fast, and reliable. The data for each project will live in one .context-keeper/context.db file.
Command-Line Argument Parsing: Go Standard Library (flag package)
Reasoning: The command structure (--key, --value, --json) is simple enough that the built-in flag package is sufficient. It avoids adding external dependencies for a task the standard library handles well.
CLI Structure & Output Handling (Frontend Guidelines)
This section describes the Go project's folder structure and the application's internal flow.
1. Proposed Project Structure:
code
Code
context-keeper/
├── cmd/
│   └── context/
│       └── main.go         # Main entry point, command & flag parsing
├── internal/
│   ├── cli/                # Command handlers (e.g., handleSave, handleList)
│   ├── db/                 # Database repository (all SQL queries)
│   └── tui/                # Bubble Tea models for rendering output
└── AGENT_INSTRUCTIONS.md.template # Template for the generated instruction file
2. Main Entry Point (main.go):
This file will be responsible for parsing the subcommand (e.g., save, list) and its flags.
Based on the subcommand, it will call the appropriate handler function from the internal/cli package.
3. Output Rendering (internal/tui):
Each command that produces output will have a simple Bubble Tea model.
For example, list will have a model that contains a bubbles/table component. The handler will fetch the data from the database and pass it to this model.
The tea.NewProgram(...).Run() function will be called to render the final view and then exit immediately. This keeps the rendering logic clean and separate from the data-fetching logic.
Data & Logic Layer (Backend Structure)
This section defines the database schema and how the application will interact with it.
1. Database Schema:
The context.db file will contain a single table named context_items.
Table: context_items
key (TEXT, PRIMARY KEY): The unique identifier for the memory (e.g., "file_summary_auth").
value (TEXT NOT NULL): The stored information, which can be a long, multiline string.
channel (TEXT NOT NULL, DEFAULT 'general'): The organizational category for the memory.
priority (TEXT NOT NULL, DEFAULT 'normal'): The priority level ('high', 'normal', 'low').
created_at (TIMESTAMP DEFAULT CURRENT_TIMESTAMP): The timestamp when the memory was first saved.
updated_at (TIMESTAMP DEFAULT CURRENT_TIMESTAMP): The timestamp for the last modification.
2. Data Access Layer (internal/db/repository.go):
We will implement the Repository Pattern. A Repository struct will be created to handle all interactions with the SQLite database.
This keeps all SQL queries in one place, making the code cleaner and easier to maintain.
The repository will have methods like:
SaveItem(key, value, channel, priority)
GetItemByKey(key)
ListItems(channel, limit)
DeleteItem(key)

Security Guidelines Document: Context Keeper CLI
1. Guiding Principles
Principle of Least Privilege: The application will only request the permissions it absolutely needs, which is read/write access within the directory it is run from. It will not attempt to access other parts of the filesystem or network.
Secure by Default: All security practices will be built-in and enabled by default. There will be no "insecure" modes of operation.
Transparency: The tool will be transparent about where it stores data. There will be no hidden files or unexpected side effects.
2. Input Validation (Critical for AI Agent Interaction)
This is the most important area of security for this tool. All input provided by the AI agent (or a human user) via command-line flags must be rigorously validated and handled safely.
No Command Injection: The application will never construct and execute shell commands from user-provided input. All interactions with the system will be through the Go standard library's safe APIs.
SQL Injection Prevention: All database queries that include user-provided data (e.g., from --key, --value, --channel) must use parameterized queries. We will never use string concatenation to build SQL statements with user input. This is a non-negotiable rule.
Correct (Safe): db.Exec("SELECT * FROM items WHERE key = ?", keyFromFlag)
Incorrect (Unsafe): db.Exec("SELECT * FROM items WHERE key = '" + keyFromFlag + "'")
Path Traversal Prevention: The application will only create one directory: .context-keeper in the current working directory. We will ensure that no user-provided input can be used to alter the file paths and write files outside of this intended location.
3. Data Storage and Handling
Storage Location: All data will be stored exclusively within the .context-keeper directory in the project's root. The tool will not write files anywhere else.
Data Encryption: Data at rest (in the context.db file) will not be encrypted by the application. The data is stored in plain text. Security relies on the user's standard filesystem permissions and disk encryption (if enabled on their OS). This is a transparent design choice to keep the tool simple and the data accessible.
Permissions: The .context-keeper directory and context.db file will be created with standard user-level read/write permissions, inheriting them from the user who runs the command.
4. Filesystem Interaction
Initialization (context init):
The init command will check if the .context-keeper directory already exists and will not overwrite it or its contents.
It will only create files (context.db, AGENT_INSTRUCTIONS.md) inside the .context-keeper directory.
5. Dependency Management
We will use a minimal set of well-known and reputable Go modules (bubbletea, bubbles, lipgloss, go-sqlite3).
Dependencies will be regularly updated to patch any security vulnerabilities discovered by the community.
6. Compliance and Transparency
Licensing: The project will include a clear LICENSE file (e.g., MIT).
Documentation: The AGENT_INSTRUCTIONS.md file and the CLI's --help output will be explicit about how the tool works and where data is stored, ensuring there is no hidden behavior.
Development Phases Overview: Context Keeper CLI
This document outlines the recommended stages for building the CLI tool.
Stage 1: Project Setup & Foundation
Goal: Create the Go project structure and initialize dependencies.
Tasks:
Initialize a new Go module: go mod init context-keeper.
Set up the project directory structure as defined in Phase 5 (cmd/, internal/).
Add the required Go dependencies by running go get:
github.com/charmbracelet/bubbletea
github.com/charmbracelet/bubbles
github.com/charmbracelet/lipgloss
github.com/mattn/go-sqlite3
Create the main entry point in cmd/context/main.go with basic command-line parsing for subcommands (init, save, list, etc.).
Stage 2: Core Logic - Database Repository
Goal: Implement the data access layer to handle all database interactions.
Tasks:
Create the Repository struct in internal/db/repository.go.
Implement the NewRepository function, which initializes the database connection and runs the SQL CREATE TABLE statement defined in Phase 5.
Implement the core data methods:
SaveItem(key, value, channel, priority)
GetItemByKey(key)
ListItems(channel, limit)
DeleteItem(key)
Security: Ensure all SQL queries are parameterized to prevent SQL injection.
RAG Recommendation:
To find examples of using go-sqlite3 for CRUD operations, you can perform a RAG query on your Go documentation sources:
Query: "go-sqlite3 examples for INSERT OR REPLACE and SELECT with WHERE clause"
Stage 3: CLI Command Implementation
Goal: Wire up the CLI commands to the database repository.
Tasks:
In internal/cli/, create handler functions for each command: handleInit, handleSave, handleList, handleGet, handleDelete.
handleInit: Implement the logic to create the .context-keeper directory, initialize the database repository, and generate the AGENT_INSTRUCTIONS.md file.
handleSave, handleGet, handleDelete: Connect the command-line flags to the corresponding repository methods.
handleList: Fetch data using the repository and prepare it for rendering.
Stage 4: Output Rendering with Bubble Tea & Bubbles
Goal: Create the styled, human-readable output for each command.
Tasks:
In internal/tui/, create a simple Bubble Tea model for each command that produces output.
List Command:
Create a model that contains a bubbles/table.
Pass the data fetched in handleList to this model.
Configure the table's styles using lipgloss according to the Design System (Phase 4).
Other Commands: Create simple models that just render the final styled string (e.g., the success or error message).
In each handler function, call tea.NewProgram(...).Run() with the appropriate model to render the output.
RAG Recommendation:
To correctly style the table and other components, you can use RAG to query the provided charmbracelet documentation:
Query: "How to style a bubbles/table header and selected row using lipgloss"
Query: "Example of using Bubble Tea to print a single view and then exit"
Stage 5: Finalization & Testing
Goal: Polish the tool, add help text, and ensure it's robust.
Tasks:
Implement comprehensive --help text for the main command and all subcommands.
Write unit tests for the database repository methods in internal/db/.
Write integration tests for the CLI handlers in internal/cli/.
Manually test all commands and flags to ensure they behave as defined in the App Flow Document.
Build the final binary (go build ./cmd/context).