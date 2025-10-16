# Context Keeper CLI — Usage

## Install/Run
- Build: `go build ./cmd/context` (produces `context.exe`)
- Optional: move `context.exe` or add its folder to PATH to run from anywhere

### PowerShell invocation tips (Windows)
- If calling with a quoted full path, use the call operator `&`:
  - `& "C:\\Users\\user1\\Desktop\\context-tool\\context-keeper\\context.exe" init`
- From the binary folder: `./context.exe ...` (PowerShell searches current dir only with `./`)
- Avoid running the binary inside `./.context-keeper` (that directory is for data only).

## Commands
- `context init [--json]`
- `context save --key "k" --value "v" [--channel "c"] [--priority "high|normal|low"] [--json]`
- `context get --key "k" [--json]`
- `context list [--channel "c"] [--limit N] [--json]`
- `context delete --key "k" [--json]`

## Behavior
- Data stored in `./.context-keeper/context.db` in the current working directory
- `--json` outputs machine‑readable JSON; default is human‑readable styled output
- SQL uses parameterized statements; writes are limited to `.context-keeper/`

## Examples
```
context init
context save --key "project_description" --value "Go CLI persistent context" --priority high
context get --key "project_description"
context list --json
context delete --key "project_description" --json
```

## JSON Output Samples
- init
```
{"success":true,"message":"Context Keeper initialized successfully.","storagePath":".context-keeper/context.db","instructionsPath":".context-keeper/AGENT_INSTRUCTIONS.md"}
```
- list
```
[{"Key":"main_file","Value":"src/index.js","Channel":"general","Priority":"high"}]
```

## Notes
- Run the CLI from the directory whose context you want to manage
- Use channels to group related items (e.g., `--channel testing`)

## Add to PATH (PowerShell)
User PATH (no admin):
```
$toolDir = 'C:\\Users\\user1\\Desktop\\context-tool\\context-keeper'
$u = [Environment]::GetEnvironmentVariable('Path','User')
if ($u -notmatch [regex]::Escape($toolDir)) { setx Path "$u;$toolDir" }
```
System PATH (Admin PowerShell):
```
$toolDir = 'C:\\Users\\user1\\Desktop\\context-tool\\context-keeper'
$m = [Environment]::GetEnvironmentVariable('Path','Machine')
if ($m -notmatch [regex]::Escape($toolDir)) { setx /M Path "$m;$toolDir" }
```
Then open a NEW PowerShell and run `context.exe ...` from anywhere.
