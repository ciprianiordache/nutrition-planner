package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nutrition-planner/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/natefinch/lumberjack.v2"
)

// --- parseLevel ---------------------------------------------------------

func TestParseLevel_Valid(t *testing.T) {
	tests := []struct {
		in   string
		want slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"", slog.LevelInfo}, // empty defaults to info
		{"warn", slog.LevelWarn},
		{"warning", slog.LevelWarn},
		{"error", slog.LevelError},
		{"  Error  ", slog.LevelError},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := parseLevel(tt.in)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseLevel_Invalid(t *testing.T) {
	_, err := parseLevel("verbose")
	assert.Error(t, err)
}

// --- resolveWriter --------------------------------------------------------

func TestResolveWriter_Stdout(t *testing.T) {
	for _, output := range []string{"stdout", "", "STDOUT"} {
		t.Run(output, func(t *testing.T) {
			w, closer, err := resolveWriter(config.LoggerConfig{Output: output})
			require.NoError(t, err)
			assert.Equal(t, os.Stdout, w)
			// stdout must never be closed by us - resolveWriter must hand
			// back a no-op closer, not os.Stdout itself as an io.Closer.
			assert.NoError(t, closer.Close())
			_, statErr := os.Stdout.Stat()
			assert.NoError(t, statErr, "os.Stdout should still be usable after closer.Close()")
		})
	}
}

func TestResolveWriter_File(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "logs") // doesn't exist yet - must be created

	w, closer, err := resolveWriter(config.LoggerConfig{
		Output:     "file",
		FilePath:   dir,
		MaxSize:    50,
		MaxAge:     30,
		MaxBackups: 10,
		Compress:   true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = closer.Close() })

	lj, ok := w.(*lumberjack.Logger)
	require.True(t, ok, "expected a *lumberjack.Logger, got %T", w)
	assert.Equal(t, filepath.Join(dir, logFileName), lj.Filename)
	assert.Equal(t, 50, lj.MaxSize)
	assert.Equal(t, 30, lj.MaxAge)
	assert.Equal(t, 10, lj.MaxBackups)
	assert.True(t, lj.Compress)
	// The writer and closer must be the same underlying *lumberjack.Logger,
	// since that's what actually holds the open file handle.
	assert.Same(t, w, closer)

	info, statErr := os.Stat(dir)
	require.NoError(t, statErr, "resolveWriter should have created the missing directory")
	assert.True(t, info.IsDir())
}

func TestResolveWriter_FileWithoutPath(t *testing.T) {
	_, _, err := resolveWriter(config.LoggerConfig{Output: "file", FilePath: ""})
	assert.Error(t, err)
}

func TestResolveWriter_Both(t *testing.T) {
	dir := t.TempDir()

	w, closer, err := resolveWriter(config.LoggerConfig{
		Output:   "both",
		FilePath: dir,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = closer.Close() })

	lj, ok := closer.(*lumberjack.Logger)
	require.True(t, ok, "expected the closer to be a *lumberjack.Logger, got %T", closer)
	assert.Equal(t, filepath.Join(dir, logFileName), lj.Filename)

	// w must fan out to both stdout and the file - verify by writing
	// through it directly and checking the file received the bytes (stdout
	// itself isn't practical to assert on here; TestNew_Both below checks
	// the full stdout+file behavior end to end via New()).
	_, err = w.Write([]byte("hello both\n"))
	require.NoError(t, err)
	require.NoError(t, lj.Close())

	data, err := os.ReadFile(filepath.Join(dir, logFileName))
	require.NoError(t, err)
	assert.Contains(t, string(data), "hello both")
}

func TestResolveWriter_BothWithoutPath(t *testing.T) {
	_, _, err := resolveWriter(config.LoggerConfig{Output: "both", FilePath: ""})
	assert.Error(t, err)
}

func TestResolveWriter_InvalidOutput(t *testing.T) {
	_, _, err := resolveWriter(config.LoggerConfig{Output: "syslog"})
	assert.Error(t, err)
}

// --- newHandler -----------------------------------------------------------

func TestNewHandler_JSON(t *testing.T) {
	h, err := newHandler("json", io.Discard, slog.LevelInfo)
	require.NoError(t, err)
	_, ok := h.(*slog.JSONHandler)
	assert.True(t, ok, "expected *slog.JSONHandler, got %T", h)
}

func TestNewHandler_Text(t *testing.T) {
	h, err := newHandler("text", io.Discard, slog.LevelInfo)
	require.NoError(t, err)
	_, ok := h.(*slog.TextHandler)
	assert.True(t, ok, "expected *slog.TextHandler, got %T", h)
}

func TestNewHandler_EmptyDefaultsToJSON(t *testing.T) {
	h, err := newHandler("", io.Discard, slog.LevelInfo)
	require.NoError(t, err)
	_, ok := h.(*slog.JSONHandler)
	assert.True(t, ok)
}

func TestNewHandler_Invalid(t *testing.T) {
	_, err := newHandler("xml", io.Discard, slog.LevelInfo)
	assert.Error(t, err)
}

// --- New (end to end) -------------------------------------------------

// newFileLogger is a small test helper: builds a New() logger writing to a
// file under t.TempDir() and registers the returned closer to run BEFORE
// TempDir's own cleanup (t.Cleanup runs LIFO, and TempDir() already
// registered its removal by the time we register this one) - on Windows,
// deleting a directory while lumberjack still holds its log file open
// fails outright, so every file-output test must close its logger before
// the temp dir gets removed.
func newFileLogger(t *testing.T, cfg config.LoggerConfig) *slog.Logger {
	t.Helper()
	lg, closer, err := New(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = closer.Close() })
	return lg
}

func TestNew_WritesJSONToFile(t *testing.T) {
	dir := t.TempDir()
	lg := newFileLogger(t, config.LoggerConfig{
		Level:    "info",
		Format:   "json",
		Output:   "file",
		FilePath: dir,
	})

	lg.Info("plan exported", "client", "Andreea Popescu", "planID", "plan-1")

	data, err := os.ReadFile(filepath.Join(dir, logFileName))
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var entry map[string]any
	// The file may contain a trailing newline after the single JSON line -
	// decode just the first line.
	line := strings.SplitN(string(data), "\n", 2)[0]
	require.NoError(t, json.Unmarshal([]byte(line), &entry))

	assert.Equal(t, "plan exported", entry["msg"])
	assert.Equal(t, "Andreea Popescu", entry["client"])
	assert.Equal(t, "plan-1", entry["planID"])
	assert.Equal(t, "INFO", entry["level"])
}

func TestNew_RespectsLevelFiltering(t *testing.T) {
	dir := t.TempDir()
	lg := newFileLogger(t, config.LoggerConfig{Level: "warn", Format: "json", Output: "file", FilePath: dir})

	lg.Info("should be filtered out")
	lg.Warn("should appear")

	data, err := os.ReadFile(filepath.Join(dir, logFileName))
	require.NoError(t, err)

	assert.NotContains(t, string(data), "should be filtered out")
	assert.Contains(t, string(data), "should appear")
}

func TestNew_TextFormat(t *testing.T) {
	dir := t.TempDir()
	lg := newFileLogger(t, config.LoggerConfig{Level: "info", Format: "text", Output: "file", FilePath: dir})

	lg.Info("hello from text handler")

	data, err := os.ReadFile(filepath.Join(dir, logFileName))
	require.NoError(t, err)
	assert.Contains(t, string(data), "msg=\"hello from text handler\"")
	// Text format should NOT be valid JSON - sanity check we didn't just
	// get a JSON handler regardless of the format field.
	var js map[string]any
	assert.Error(t, json.Unmarshal(data, &js))
}

func TestNew_WritesToStdout(t *testing.T) {
	// Redirect os.Stdout for the duration of this test only - resolveWriter
	// resolves os.Stdout at call time, so swapping the package var before
	// calling New() is enough to capture output without needing New() to
	// accept a writer directly.
	r, w, err := os.Pipe()
	require.NoError(t, err)
	original := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = original }()

	lg, closer, err := New(config.LoggerConfig{Level: "info", Format: "json", Output: "stdout"})
	require.NoError(t, err)
	defer func() { _ = closer.Close() }() // no-op for stdout, but exercised for completeness

	lg.Info("hello stdout")

	require.NoError(t, w.Close())
	os.Stdout = original

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "hello stdout")
}

func TestNew_InvalidConfigPropagatesError(t *testing.T) {
	_, _, err := New(config.LoggerConfig{Level: "not-a-level"})
	assert.Error(t, err)
}

func TestNew_Both(t *testing.T) {
	dir := t.TempDir()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	original := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = original }()

	lg, closer, err := New(config.LoggerConfig{Level: "info", Format: "json", Output: "both", FilePath: dir})
	require.NoError(t, err)

	lg.Info("hello both outputs")

	require.NoError(t, closer.Close()) // flush/close the file half
	require.NoError(t, w.Close())      // close the write end of the pipe so io.Copy below returns
	os.Stdout = original

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "hello both outputs")

	data, err := os.ReadFile(filepath.Join(dir, logFileName))
	require.NoError(t, err)
	assert.Contains(t, string(data), "hello both outputs")
}
