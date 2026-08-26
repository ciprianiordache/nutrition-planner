// Package logger builds the application's *slog.Logger from
// config.LoggerConfig - level, format (json/text), and where logs go
// (stdout, or a rotating file). It doesn't invent its own logging
// abstraction: log/slog is already the standard, extensible logging
// interface in Go, so New just wires a slog.Handler up according to
// config and hands back a plain *slog.Logger every other package can use
// directly (or via slog.SetDefault, for call sites that don't want the
// logger injected).
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"nutrition-planner/config"

	"gopkg.in/natefinch/lumberjack.v2"
)

// logFileName is the fixed name of the log file created under
// config.LoggerConfig.FilePath when Output is "file" - the config only
// gives us a directory, not a filename, so we own the naming here.
const logFileName = "app.log"

// New builds a *slog.Logger fully configured from cfg. It returns an error
// for any field it doesn't recognize (level/format/output) rather than
// silently falling back, since a typo in config.yaml should surface
// immediately at startup, not produce logs nobody can find.
//
// The returned io.Closer must be closed when the logger is no longer
// needed (typically on application shutdown, right next to closing the
// database). When Output is "file", the underlying *lumberjack.Logger
// keeps the log file open for the process's entire lifetime - on Windows
// specifically, an open file cannot be deleted or moved by anything else
// (unlike POSIX, which allows unlinking open files), so skipping this
// close can leave the log file locked, break log rotation cleanup, and
// fail temp-directory removal in tests. When Output is "stdout", the
// returned closer is a no-op - closing os.Stdout would be wrong regardless
// of whether the logger is done with it.
func New(cfg config.LoggerConfig) (*slog.Logger, io.Closer, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, nil, err
	}

	writer, closer, err := resolveWriter(cfg)
	if err != nil {
		return nil, nil, err
	}

	handler, err := newHandler(cfg.Format, writer, level)
	if err != nil {
		return nil, nil, err
	}

	return slog.New(handler), closer, nil
}

func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("logger: unsupported level %q (expected debug, info, warn, or error)", level)
	}
}

// nopCloser is used for outputs (stdout) that must NOT be closed by us.
type nopCloser struct{}

func (nopCloser) Close() error { return nil }

// resolveWriter turns cfg.Output into the io.Writer logs get written to,
// plus the io.Closer the caller is responsible for closing on shutdown.
// "file" wires up automatic rotation via lumberjack, using cfg's
// MaxSize/MaxAge/MaxBackups/Compress as-is - those fields exist on
// LoggerConfig specifically to be passed straight through here. "both"
// writes to stdout AND the rotating file at the same time - useful since
// stdout is only visible during development (a packaged/installed app has
// no attached console at all), while the file is what you can actually go
// back and read later.
func resolveWriter(cfg config.LoggerConfig) (io.Writer, io.Closer, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Output)) {
	case "stdout", "":
		return os.Stdout, nopCloser{}, nil
	case "file":
		lj, err := newRotatingFile(cfg)
		if err != nil {
			return nil, nil, err
		}
		return lj, lj, nil // *lumberjack.Logger implements io.Closer itself
	case "both":
		lj, err := newRotatingFile(cfg)
		if err != nil {
			return nil, nil, err
		}
		return io.MultiWriter(os.Stdout, lj), lj, nil
	default:
		return nil, nil, fmt.Errorf("logger: unsupported output %q (expected \"stdout\", \"file\", or \"both\")", cfg.Output)
	}
}

// newRotatingFile builds the *lumberjack.Logger shared by the "file" and
// "both" output modes, creating cfg.FilePath if it doesn't exist yet.
func newRotatingFile(cfg config.LoggerConfig) (*lumberjack.Logger, error) {
	if cfg.FilePath == "" {
		return nil, fmt.Errorf("logger: output is %q but file_path is empty", cfg.Output)
	}
	if err := os.MkdirAll(cfg.FilePath, 0o755); err != nil {
		return nil, fmt.Errorf("logger: create log directory %q: %w", cfg.FilePath, err)
	}
	return &lumberjack.Logger{
		Filename:   filepath.Join(cfg.FilePath, logFileName),
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
	}, nil
}

func newHandler(format string, w io.Writer, level slog.Level) (slog.Handler, error) {
	opts := &slog.HandlerOptions{Level: level}
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json", "":
		return slog.NewJSONHandler(w, opts), nil
	case "text":
		return slog.NewTextHandler(w, opts), nil
	default:
		return nil, fmt.Errorf("logger: unsupported format %q (expected \"json\" or \"text\")", format)
	}
}
