package settings

import (
	"context"
	"log/slog"

	"nutrition-planner/internals/database"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Handler exposes settings functionality to the frontend via Wails
// binding, including a native "choose a folder" dialog for the export
// directory - Wails only exposes OpenDirectoryDialog on the Go side, so
// this is the one method here that needs a *context.Context (set once
// from App.startup, since it isn't available yet when the handler chain
// is constructed in NewApp).
//
// ctx is deliberately unexported and set via the package-level WithContext
// function below, NOT an exported SetContext method - Wails binds every
// exported method of whatever struct is passed to Bind, with no filtering
// by parameter type, so an exported method taking a context.Context would
// itself get exposed to the frontend as a callable (and broken, since
// context.Context can't be marshaled to/from JS) binding.
type Handler struct {
	svc Service
	db  *database.Database
	ctx context.Context
	log *slog.Logger
}

func NewHandler(svc Service, db *database.Database, log *slog.Logger) *Handler {
	return &Handler{svc: svc, db: db, log: log}
}

// WithContext must be called once from App.startup - every method that
// needs to show a native dialog is a no-op (returns an error) until then.
func WithContext(h *Handler, ctx context.Context) { h.ctx = ctx }

func (h *Handler) GetSettings() (Response, error) {
	return h.svc.Get()
}

func (h *Handler) UpdateSettings(req UpdateRequest) (Response, error) {
	return h.svc.Update(req)
}

// ChooseExportFolder opens the OS's native folder picker and returns the
// chosen absolute path, or "" (with a nil error) if the user cancelled -
// it does NOT save the choice; the frontend is expected to follow up with
// UpdateSettings once the user confirms the rest of the form.
func (h *Handler) ChooseExportFolder() (string, error) {
	if h.ctx == nil {
		return "", errCtxNotReady
	}
	path, err := runtime.OpenDirectoryDialog(h.ctx, runtime.OpenDialogOptions{
		Title: "Alege folderul pentru exporturi PDF",
	})
	if err != nil {
		h.log.Error("failed to open folder picker", "error", err)
		return "", err
	}
	return path, nil
}
