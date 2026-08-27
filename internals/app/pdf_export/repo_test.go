package pdf_export_test

import (
	"os"
	"path/filepath"
	"testing"

	"nutrition-planner/internals/app/pdf_export"
)

func TestRepo_ClientExportDir_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	repo := pdf_export.NewRepo(base)

	dir, err := repo.ClientExportDir("", "Andreea Popescu")
	if err != nil {
		t.Fatalf("ClientExportDir: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", dir)
	}

	want := filepath.Join(base, "Andreea Popescu")
	if dir != want {
		t.Errorf("dir = %q, want %q", dir, want)
	}
}

func TestRepo_ClientExportDir_SanitizesUnsafeCharacters(t *testing.T) {
	base := t.TempDir()
	repo := pdf_export.NewRepo(base)

	dir, err := repo.ClientExportDir("", `Client/With:Bad*Chars`)
	if err != nil {
		t.Fatalf("ClientExportDir: %v", err)
	}

	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("expected sanitized directory to exist: %v", err)
	}
	// The exact sanitized form is naming.go's business (covered in
	// naming_test.go) - here we only care that a *valid, existing* path
	// under base came back, regardless of the original unsafe input.
	if filepath.Dir(dir) != base {
		t.Errorf("expected directory to live directly under base %q, got %q", base, dir)
	}
}

func TestRepo_ClientExportDir_IdempotentForSameClient(t *testing.T) {
	base := t.TempDir()
	repo := pdf_export.NewRepo(base)

	dir1, err := repo.ClientExportDir("", "Mihai Ionescu")
	if err != nil {
		t.Fatalf("ClientExportDir (1st call): %v", err)
	}
	dir2, err := repo.ClientExportDir("", "Mihai Ionescu")
	if err != nil {
		t.Fatalf("ClientExportDir (2nd call): %v", err)
	}

	if dir1 != dir2 {
		t.Errorf("expected the same directory across calls, got %q and %q", dir1, dir2)
	}
}

func TestRepo_ClientExportDir_DifferentClientsGetDifferentDirs(t *testing.T) {
	base := t.TempDir()
	repo := pdf_export.NewRepo(base)

	dirA, err := repo.ClientExportDir("", "Client A")
	if err != nil {
		t.Fatalf("ClientExportDir: %v", err)
	}
	dirB, err := repo.ClientExportDir("", "Client B")
	if err != nil {
		t.Fatalf("ClientExportDir: %v", err)
	}

	if dirA == dirB {
		t.Errorf("expected different clients to get different directories, both got %q", dirA)
	}
}

func TestRepo_ClientExportDir_BaseOverrideWins(t *testing.T) {
	configuredBase := t.TempDir()
	overrideBase := t.TempDir()
	repo := pdf_export.NewRepo(configuredBase)

	dir, err := repo.ClientExportDir(overrideBase, "Andreea Popescu")
	if err != nil {
		t.Fatalf("ClientExportDir: %v", err)
	}

	want := filepath.Join(overrideBase, "Andreea Popescu")
	if dir != want {
		t.Errorf("dir = %q, want %q (override should win over the configured base)", dir, want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected the override directory to have been created: %v", err)
	}
}

func TestRepo_ClientExportDir_EmptyOverrideUsesConfiguredBase(t *testing.T) {
	configuredBase := t.TempDir()
	repo := pdf_export.NewRepo(configuredBase)

	dir, err := repo.ClientExportDir("", "Andreea Popescu")
	if err != nil {
		t.Fatalf("ClientExportDir: %v", err)
	}

	want := filepath.Join(configuredBase, "Andreea Popescu")
	if dir != want {
		t.Errorf("dir = %q, want %q", dir, want)
	}
}

func TestDefaultBaseDirectory(t *testing.T) {
	dir, err := pdf_export.DefaultBaseDirectory()
	if err != nil {
		t.Fatalf("DefaultBaseDirectory: %v", err)
	}
	if dir == "" {
		t.Fatal("expected a non-empty directory")
	}
	if filepath.Base(dir) != "NutriPlan Exports" {
		t.Errorf("expected the last path segment to be %q, got %q", "NutriPlan Exports", dir)
	}
	if filepath.Base(filepath.Dir(dir)) != "Desktop" {
		t.Errorf("expected the export folder to live directly under Desktop, got %q", dir)
	}
	if !filepath.IsAbs(dir) {
		t.Errorf("expected an absolute path, got %q", dir)
	}
}
