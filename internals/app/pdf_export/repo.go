package pdf_export

import (
	"fmt"
	"os"
	"path/filepath"
)

// Repo owns the filesystem side of exporting - where PDFs land - kept
// separate from Service so the export directory strategy (one folder per
// client) can change without touching the plan -> PDF mapping logic.
type Repo interface {
	// ClientExportDir returns the directory PDFs for the given client
	// should be written to, creating it (and any missing parents) if it
	// doesn't exist yet. Safe to call repeatedly for the same client.
	//
	// baseOverride, when non-empty, is used as the root export folder
	// instead of the Repo's configured default - this is how a live
	// settings.ExportBaseDirectory takes effect without restarting the
	// app or rebuilding the Repo.
	ClientExportDir(baseOverride, clientName string) (string, error)
}

type repo struct {
	// baseDir is the root export folder (e.g. a "NutriPlan Exports" folder
	// on the user's Desktop, taken from config.ExportConfig.BaseDirectory
	// or computed by DefaultBaseDirectory) - ClientExportDir creates one
	// subdirectory per client underneath it.
	baseDir string
}

// NewRepo builds a Repo rooted at baseDir. baseDir is expected to already
// be an absolute, writable path - see DefaultBaseDirectory for the
// portable, per-machine default.
func NewRepo(baseDir string) Repo {
	return &repo{baseDir: baseDir}
}

// DefaultBaseDirectory resolves a portable default export location:
// "<current user's home directory>/Desktop/NutriPlan Exports". It relies
// on os.UserHomeDir(), which reads %USERPROFILE% on Windows and $HOME on
// macOS/Linux - so it resolves correctly on any machine/user account
// without ever hardcoding a username. Use this when
// config.ExportConfig.BaseDirectory is left empty.
func DefaultBaseDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home directory: %w", err)
	}
	return filepath.Join(home, "Desktop", "NutriPlan Exports"), nil
}

func (r *repo) ClientExportDir(baseOverride, clientName string) (string, error) {
	base := r.baseDir
	if baseOverride != "" {
		base = baseOverride
	}
	dir := filepath.Join(base, sanitizePathSegment(clientName, "client"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create export directory %q: %w", dir, err)
	}
	return dir, nil
}
