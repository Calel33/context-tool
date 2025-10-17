Param(
  [string]$Output = "context-gui.exe"
)

# Windows PowerShell build script for GUI binary
go build -o $Output ./cmd/context-gui
if ($LASTEXITCODE -ne 0) {
  Write-Error "Build failed"
  exit 1
}
Write-Host "Built $Output"
