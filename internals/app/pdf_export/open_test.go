package pdf_export

import (
	"os/exec"
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

func TestOpenInDefaultApp_UsesExecCommandSeam(t *testing.T) {
	var gotName string
	var gotArgs []string

	original := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		gotName = name
		gotArgs = args
		// A real, harmless command that Start() can succeed against on
		// whichever platform this test runs on, so we're only exercising
		// "did we build the right command", not "can the OS actually
		// launch a PDF viewer" (which isn't something a unit test should
		// depend on).
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/c", "exit", "0")
		}
		return exec.Command("true")
	}
	defer func() { execCommand = original }()

	if err := openInDefaultApp("/tmp/plan.pdf"); err != nil {
		t.Fatalf("openInDefaultApp: %v", err)
	}
	if gotName == "" {
		t.Fatal("expected execCommand to have been called")
	}
	if len(gotArgs) == 0 || gotArgs[len(gotArgs)-1] != "/tmp/plan.pdf" {
		t.Errorf("expected the path to be the last argument, got %v", gotArgs)
	}
}
