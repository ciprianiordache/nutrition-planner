# build.ps1 — runs `wails build` (optionally with -Installer, to also
# produce the NSIS Windows installer) and then copies the runtime config
# files (config.yaml, .env) next to the resulting executable in build/bin/.
#
# Why the config-copy step is needed: wails build does NOT copy
# config.yaml/.env into build/bin/ on its own - they aren't frontend
# assets, they're runtime config the app reads with a relative path
# ("config.yaml", ".env") at startup. When you run the built .exe, its
# working directory is wherever it's launched from (often build/bin/
# itself), so without this step the app can't find your real config and
# silently generates a brand-new one with default values (different
# theme, different logger output, etc.). The NSIS installer itself
# already bundles its own config.yaml directly (see
# build/windows/installer/project.nsi, which bundles config.installer.yaml
# renamed to config.yaml) - this copy step is specifically for testing the
# raw .exe in build/bin/ before you package it.
#
# Usage (from the project root):
#   .\build.ps1                 # plain build, exe only
#   .\build.ps1 -Installer      # also builds build\bin\nutrition-planner-amd64-installer.exe
#
# Building the installer requires NSIS (makensis) on PATH:
#   https://nsis.sourceforge.io/Download
#
# IMPORTANT before building a real, distributable installer: config.yaml
# currently has app.debug: true and app.environment: development, which
# are local-development settings (auto-opens DevTools, verbose logging) -
# not what you want every end user's install to run with. The installer
# itself already ships a separate, production-appropriate config
# (config.installer.yaml -> debug: false, environment: production), so
# this only matters if you're editing that file too.

param(
    [switch]$Installer
)

$ErrorActionPreference = "Stop"

if ($Installer) {
    Write-Host "==> Running wails build -nsis..." -ForegroundColor Cyan
    wails build -nsis
} else {
    Write-Host "==> Running wails build..." -ForegroundColor Cyan
    wails build
}
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
if ($Installer) {
    Write-Host "    Installer: $binDir\nutrition-planner-amd64-installer.exe" -ForegroundColor Green
}
