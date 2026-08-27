package settings

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestZipFile_ContainsSourceContentUnderEntryName(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "source.db")
	want := []byte("pretend this is sqlite bytes")
	if err := os.WriteFile(srcPath, want, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	zipPath := filepath.Join(dir, "backup.zip")
	if err := zipFile(srcPath, zipPath, "nutrition-planner.db"); err != nil {
		t.Fatalf("zipFile: %v", err)
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("zip.OpenReader: %v", err)
	}
	defer r.Close()

	if len(r.File) != 1 {
		t.Fatalf("expected exactly 1 entry, got %d", len(r.File))
	}
	entry := r.File[0]
	if entry.Name != "nutrition-planner.db" {
		t.Errorf("entry name = %q, want %q", entry.Name, "nutrition-planner.db")
	}

	rc, err := entry.Open()
	if err != nil {
		t.Fatalf("entry.Open: %v", err)
	}
	defer rc.Close()

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("entry content = %q, want %q", got, want)
	}
}

func TestZipFile_MissingSourceErrors(t *testing.T) {
	dir := t.TempDir()
	err := zipFile(filepath.Join(dir, "does-not-exist.db"), filepath.Join(dir, "backup.zip"), "x.db")
	if err == nil {
		t.Fatal("expected an error for a missing source file")
	}
}
