# Context Keeper CLI — Agent Instructions

Commands:
- context init
- context save --key "k" --value "v" [--channel "c"] [--priority "high|normal|low"] [--json]
- context get --key "k" [--json]
- context list [--channel "c"] [--limit N] [--json]
- context delete --key "k" [--json]

Notes:
- Use --json for machine-readable output.
- Data is stored at ./.context-keeper/context.db within the project root.
- Writes are limited to the .context-keeper directory.
