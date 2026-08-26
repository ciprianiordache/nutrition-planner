package pdf_export

import (
	"os/exec"
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

// openInDefaultApp launches path in the OS's default application for its
// file type. Wails' own runtime.BrowserOpenURL refuses file:// URLs
// outright (a deliberate security restriction in wails v2.13.0, not a
// bug) - shelling out ourselves is the supported way to open a local file
// from a Wails app.
func openInDefaultApp(path string) error {
	return commandForOS(runtime.GOOS, path).Start()
}
