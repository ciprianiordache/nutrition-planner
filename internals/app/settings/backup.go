package settings

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"nutrition-planner/internals/osutil"
)

// BackupDatabase copies the live SQLite database file into a timestamped
// .zip inside a "backups" folder next to it, returning the resulting path.
// Zipping (rather than a raw .db copy) keeps a single, slightly smaller
// artifact per backup instead of a pile of loose .db files.
//
// This is a plain file copy, not a proper SQLite-aware backup (no
// VACUUM INTO, no WAL checkpoint) - fine for this app's write pattern
// (short-lived connections, no long-running transactions left open), but
// worth knowing if the database ever grows a different access pattern.
func (h *Handler) BackupDatabase() (string, error) {
	srcPath, err := h.db.SQLitePath()
	if err != nil {
		h.log.Error("backup failed: could not resolve database path", "error", err)
		return "", err
	}

	backupDir := filepath.Join(filepath.Dir(srcPath), "backups")
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		h.log.Error("backup failed: could not create backup directory", "backup_dir", backupDir, "error", err)
		return "", fmt.Errorf("create backup directory: %w", err)
	}

	zipName := fmt.Sprintf("nutrition-planner-backup-%s.zip", time.Now().Format("2006-01-02-150405"))
	zipPath := filepath.Join(backupDir, zipName)

	if err := zipFile(srcPath, zipPath, filepath.Base(srcPath)); err != nil {
		h.log.Error("backup failed", "src", srcPath, "dst", zipPath, "error", err)
		return "", fmt.Errorf("create backup archive: %w", err)
	}

	h.log.Info("database backed up", "path", zipPath)
	return zipPath, nil
}

// OpenBackupFolder opens the OS file manager at the directory containing
// filePath (typically the path BackupDatabase just returned).
func (h *Handler) OpenBackupFolder(filePath string) error {
	return osutil.OpenContainingFolder(filePath)
}

// zipFile writes a single-entry zip archive at zipPath containing
// srcPath's contents under entryName.
func zipFile(srcPath, zipPath, entryName string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("create archive file: %w", err)
	}
	defer dst.Close()

	zw := zip.NewWriter(dst)
	w, err := zw.Create(entryName)
	if err != nil {
		return fmt.Errorf("create archive entry: %w", err)
	}
	if _, err := io.Copy(w, src); err != nil {
		return fmt.Errorf("write archive entry: %w", err)
	}
	return zw.Close()
}
