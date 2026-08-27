// Package osutil holds small OS-integration helpers shared across domains
// (pdf_export needs to open a generated PDF or its folder; settings needs
// to open a database backup's folder) - extracted here once the same need
// showed up in a second place, rather than duplicating it.
package osutil

import (
	"os/exec"
	"path/filepath"
	"runtime"
)

// execCommand is a seam for tests - swap it out to observe which
// command/args would have been run without actually launching anything.
var execCommand = exec.Command

// commandForOS builds the OS-specific command that opens path with
// whatever application the system has registered as the default handler
// for its file type (mirroring what double-clicking the file in
// Explorer/Finder/a file manager would do).
func commandForOS(goos, path string) *exec.Cmd {
	switch goos {
	case "windows":
		// rundll32's FileProtocolHandler is the standard way to do this on
		// Windows without the quoting footguns of `cmd /c start` (where
		// the first quoted argument is treated as a window title, not the
		// path, if the path itself needs quoting).
		return execCommand("rundll32", "url.dll,FileProtocolHandler", path)
	case "darwin":
		return execCommand("open", path)
	default: // linux and other unix-likes
		return execCommand("xdg-open", path)
	}
}

// OpenInDefaultApp launches path (a file OR a directory - all three OS
// commands above handle both the same way) in the OS's default
// application. Wails' own runtime.BrowserOpenURL refuses file:// URLs
// outright (a deliberate security restriction in wails v2.13.0, not a
// bug) - shelling out ourselves is the supported way to open a local path
// from a Wails app.
func OpenInDefaultApp(path string) error {
	return commandForOS(runtime.GOOS, path).Start()
}

// OpenContainingFolder opens the OS file manager at the directory
// containing filePath - the path manipulation happens here, in Go, rather
// than in the frontend, since splitting a Windows path
// ("C:\Users\...\plan.pdf") correctly in JS is fragile; filepath.Dir
// already knows how to do it properly for whatever OS this actually runs
// on.
func OpenContainingFolder(filePath string) error {
	return OpenInDefaultApp(filepath.Dir(filePath))
}
