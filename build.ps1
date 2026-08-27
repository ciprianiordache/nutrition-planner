# build.ps1 — runs `wails build` and then copies the runtime config files
# (config.yaml, .env) next to the resulting executable in build/bin/.
#
# Why this is needed: wails build does NOT copy config.yaml/.env into
# build/bin/ on its own - they aren't frontend assets, they're runtime
# config the app reads with a relative path ("config.yaml", ".env") at
# startup. When you run the built .exe, its working directory is wherever
# it's launched from (often build/bin/ itself), so without this step the
# app can't find your real config and silently generates a brand-new one
# with default values (different theme, different logger output, etc.).
#
# Usage (from the project root):
#   .\build.ps1

$ErrorActionPreference = "Stop"

Write-Host "==> Running wails build..." -ForegroundColor Cyan
wails build
if ($LASTEXITCODE -ne 0) {
    Write-Host "wails build failed (exit code $LASTEXITCODE)" -ForegroundColor Red
    exit $LASTEXITCODE
}

$binDir = Join-Path $PSScriptRoot "build\bin"
if (-not (Test-Path $binDir)) {
    Write-Host "Expected output directory not found: $binDir" -ForegroundColor Red
    exit 1
}

Write-Host "==> Copying config.yaml and .env into $binDir" -ForegroundColor Cyan
Copy-Item -Path (Join-Path $PSScriptRoot "config.yaml") -Destination $binDir -Force
Copy-Item -Path (Join-Path $PSScriptRoot ".env") -Destination $binDir -Force

Write-Host "==> Done. Executable + config are both in: $binDir" -ForegroundColor Green
Write-Host "    Run the .exe directly from that folder to test with your real config."
