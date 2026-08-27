package pdf_export

import (
	"nutrition-planner/internals/osutil"
)

// Handler exposes export functionality to the frontend via Wails binding
// (App.Export in app.go becomes pdf_export.Handler in frontend/wailsjs/go).
type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

// ExportPlanPDF generates the PDF for the given plan and writes it to disk
// under a per-client folder, returning the resulting file name and path.
func (h *Handler) ExportPlanPDF(planID string) (Response, error) {
	return h.svc.ExportPlanPDF(planID)
}

// PreviewPlanPDF generates the PDF for the given plan into a scratch
// temp-directory location (not the client's export folder) and returns
// its file name and path, so the frontend can show it before the user
// decides to export for real.
func (h *Handler) PreviewPlanPDF(planID string) (Response, error) {
	return h.svc.PreviewPlanPDF(planID)
}

// OpenExportFolder opens the OS file manager at the directory containing
// filePath (typically Response.Path from a prior ExportPlanPDF call).
func (h *Handler) OpenExportFolder(filePath string) error {
	return osutil.OpenContainingFolder(filePath)
}
