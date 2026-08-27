package osutil

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCommandForOS_Windows(t *testing.T) {
	cmd := commandForOS("windows", `C:\temp\plan.pdf`)
	if len(cmd.Args) != 3 {
		t.Fatalf("expected 3 args, got %v", cmd.Args)
	}
	if cmd.Args[0] != "rundll32" || cmd.Args[1] != "url.dll,FileProtocolHandler" || cmd.Args[2] != `C:\temp\plan.pdf` {
		t.Errorf("unexpected args: %v", cmd.Args)
	}
}

func TestCommandForOS_Darwin(t *testing.T) {
	cmd := commandForOS("darwin", "/tmp/plan.pdf")
	if len(cmd.Args) != 2 || cmd.Args[0] != "open" || cmd.Args[1] != "/tmp/plan.pdf" {
		t.Errorf("unexpected args: %v", cmd.Args)
	}
}

func TestCommandForOS_Linux(t *testing.T) {
	cmd := commandForOS("linux", "/tmp/plan.pdf")
	if len(cmd.Args) != 2 || cmd.Args[0] != "xdg-open" || cmd.Args[1] != "/tmp/plan.pdf" {
		t.Errorf("unexpected args: %v", cmd.Args)
	}
}

func TestCommandForOS_UnknownFallsBackToXdgOpen(t *testing.T) {
	cmd := commandForOS("plan9", "/tmp/plan.pdf")
	if len(cmd.Args) != 2 || cmd.Args[0] != "xdg-open" {
		t.Errorf("unexpected args for an unrecognized GOOS: %v", cmd.Args)
	}
}

// harmlessCommand returns a real, trivially-successful command for
// whichever OS the test suite actually runs on, so tests exercise "did we
// build/call the right thing" without depending on any real PDF
// viewer/file manager being installed.
func harmlessCommand() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", "exit", "0")
	}
	return exec.Command("true")
}

func TestOpenInDefaultApp_UsesExecCommandSeam(t *testing.T) {
	var gotName string
	var gotArgs []string

	original := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotName = name
		gotArgs = args
		return harmlessCommand()
	}
	defer func() { execCommand = original }()

	if err := OpenInDefaultApp("/tmp/plan.pdf"); err != nil {
		t.Fatalf("OpenInDefaultApp: %v", err)
	}
	if gotName == "" {
		t.Fatal("expected execCommand to have been called")
	}
	if len(gotArgs) == 0 || gotArgs[len(gotArgs)-1] != "/tmp/plan.pdf" {
		t.Errorf("expected the path to be the last argument, got %v", gotArgs)
	}
}

func TestOpenContainingFolder_OpensTheDirectoryNotTheFile(t *testing.T) {
	var gotPath string

	original := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotPath = args[len(args)-1]
		return harmlessCommand()
	}
	defer func() { execCommand = original }()

	filePath := filepath.Join("some", "nested", "dir", "plan.pdf")
	if err := OpenContainingFolder(filePath); err != nil {
		t.Fatalf("OpenContainingFolder: %v", err)
	}

	want := filepath.Join("some", "nested", "dir")
	if gotPath != want {
		t.Errorf("expected the containing directory %q to be opened, got %q", want, gotPath)
	}
}
