package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultAppDataDir returns a sensible, writable, per-user directory for
// this app's own data (database, logs, temp scratch files) - deliberately
// separate from wherever the program itself is installed.
//
// This matters for two real reasons, not just tidiness:
//  1. On Windows, a machine-wide install typically lands in
//     "C:\Program Files\...", which a standard (non-admin) user often
//     can't write to at runtime - writes there either fail outright or
//     get silently redirected by Windows' File System Virtualization into
//     a per-user VirtualStore folder, which is confusing and easy to miss.
//  2. Windows' file-system convention for installers is to blow away the
//     entire installation directory on uninstall. If the database/logs
//     lived inside it (as they would with a relative path like "./data"
//     when the app's working directory is the install directory), then
//     uninstalling the app would silently delete the user's client data
//     along with it.
//
// os.UserConfigDir() resolves to %AppData% on Windows, ~/Library/
// Application Support on macOS, and ~/.config (or $XDG_CONFIG_HOME) on
// Linux - all outside any program's own install location, and always
// writable by the current user without elevation.
func DefaultAppDataDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config directory: %w", err)
	}
	return filepath.Join(dir, "NutriPlan"), nil
}
