package pdf_export

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"nutrition-planner/config"
	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/plan"
	"nutrition-planner/internals/generator"
)

// Service exports a plan (identified by ID) to a PDF file on disk, and can
// render a throwaway preview of the same plan without touching the
// client's real export folder.
type Service interface {
	// ExportPlanPDF renders planID's plan as a PDF and places it in the
	// client's export folder. The PDF is first rendered into a scratch
	// file under tempDir and only moved into place once fully written, so
	// a crash or write failure mid-render never leaves a partial file
	// where the client would see it.
	ExportPlanPDF(planID string) (Response, error)
	// PreviewPlanPDF renders planID's plan as a PDF into tempDir and
	// leaves it there - unlike ExportPlanPDF, nothing gets moved into the
	// client's export folder. Intended for a "Preview" action that lets
	// the user look at the PDF before deciding to export it for real.
	PreviewPlanPDF(planID string) (Response, error)
}

type service struct {
	planSvc   plan.Service
	clientSvc client.Service
	repo      Repo
	// genCfg is the base pdfgen configuration (theme, margins, page size)
	// loaded once from the app config. OutputDirectory is overridden per
	// render with a scratch directory (tempDir) or the client's export
	// directory, depending on the operation.
	genCfg config.GeneratorConfig
	// tempDir is where every PDF gets rendered first (config.AppConfig.
	// TempDirectory) - both the export staging file and preview files
	// live here.
	tempDir string
	log     *slog.Logger
}

func NewService(planSvc plan.Service, clientSvc client.Service, repo Repo, genCfg config.GeneratorConfig, tempDir string, log *slog.Logger) Service {
	return &service{planSvc: planSvc, clientSvc: clientSvc, repo: repo, genCfg: genCfg, tempDir: tempDir, log: log}
}

func (s *service) ExportPlanPDF(planID string) (Response, error) {
	if planID == "" {
		return Response{}, ErrInvalidInput
	}

	stagingPath, fileName, clientName, err := s.renderPDF(planID)
	if err != nil {
		return Response{}, err
	}

	dir, err := s.repo.ClientExportDir(clientName)
	if err != nil {
		s.log.Error("pdf export failed: could not resolve export directory", "plan_id", planID, "client", clientName, "error", err)
		return Response{}, fmt.Errorf("resolve export directory: %w", err)
	}

	finalPath := filepath.Join(dir, fileName)
	if err := moveFile(stagingPath, finalPath); err != nil {
		s.log.Error("pdf export failed: could not finalize staged file", "plan_id", planID, "error", err)
		return Response{}, fmt.Errorf("finalize export: %w", err)
	}

	s.log.Info("plan exported to pdf", "plan_id", planID, "client", clientName, "path", finalPath)
	return Response{FileName: fileName, Path: finalPath}, nil
}

func (s *service) PreviewPlanPDF(planID string) (Response, error) {
	if planID == "" {
		return Response{}, ErrInvalidInput
	}

	stagingPath, fileName, clientName, err := s.renderPDF(planID)
	if err != nil {
		return Response{}, err
	}

	// Best-effort: the preview PDF is already rendered and its path is
	// valid either way, so a failure to auto-open it (no PDF viewer
	// registered, sandboxed environment, etc.) doesn't fail the whole
	// operation - it's just logged, and the frontend still gets the path
	// back in case it wants to show it or let the user open it manually.
	if err := openInDefaultApp(stagingPath); err != nil {
		s.log.Warn("could not auto-open pdf preview", "plan_id", planID, "path", stagingPath, "error", err)
	}

	s.log.Info("plan preview generated", "plan_id", planID, "client", clientName, "path", stagingPath)
	return Response{FileName: fileName, Path: stagingPath}, nil
}

// renderPDF loads planID's full plan tree and client, and renders it as a
// PDF into a scratch file under s.tempDir. It's shared by ExportPlanPDF
// (which then moves the result into the client's export folder) and
// PreviewPlanPDF (which leaves it where it is) - every failure mode
// (missing plan/client, bad generator config, disk errors) only needs to
// be handled and logged in one place.
func (s *service) renderPDF(planID string) (stagingPath, fileName, clientName string, err error) {
	p, err := s.planSvc.GetByIDFull(planID)
	if err != nil {
		s.log.Error("pdf render failed: could not load plan", "plan_id", planID, "error", err)
		return "", "", "", fmt.Errorf("load plan: %w", err)
	}

	c, err := s.clientSvc.GetByID(p.ClientID)
	if err != nil {
		s.log.Error("pdf render failed: could not load client", "plan_id", planID, "client_id", p.ClientID, "error", err)
		return "", "", "", fmt.Errorf("load client: %w", err)
	}

	if err := os.MkdirAll(s.tempDir, 0o755); err != nil {
		s.log.Error("pdf render failed: could not create temp directory", "plan_id", planID, "temp_dir", s.tempDir, "error", err)
		return "", "", "", fmt.Errorf("create temp directory: %w", err)
	}

	cfg := s.genCfg
	cfg.OutputDirectory = s.tempDir

	gen, err := generator.NewGenerator(cfg)
	if err != nil {
		s.log.Error("pdf render failed: could not initialize pdf generator", "plan_id", planID, "error", err)
		return "", "", "", fmt.Errorf("init pdf generator: %w", err)
	}

	renderPlan(gen, p, c)

	fileName = buildFileName(p.Title)
	stagingPath, err = gen.Save(fileName)
	if err != nil {
		s.log.Error("pdf render failed: could not save pdf", "plan_id", planID, "file_name", fileName, "error", err)
		return "", "", "", fmt.Errorf("save pdf: %w", err)
	}

	return stagingPath, fileName, c.Name, nil
}

// buildFileName derives a filesystem-safe PDF name from the plan title,
// suffixed with the render date so re-exporting/re-previewing the same
// plan later in the day overwrites it (predictable) while a different day
// doesn't (avoids silently clobbering an older export people may still
// want).
func buildFileName(planTitle string) string {
	return fmt.Sprintf("%s_%s.pdf", sanitizePathSegment(planTitle, "plan"), time.Now().Format("2006-01-02"))
}

// moveFile moves src to dst, trying an atomic os.Rename first. Rename
// fails when src and dst live on different filesystems/volumes (e.g. the
// temp directory and the export directory are on different drives) - in
// that case we fall back to a plain copy-then-remove. A failure to remove
// the now-redundant staging file afterward isn't treated as an error: dst
// already has the correct, complete content, so the export itself
// succeeded - the leftover in tempDir is just cleanup cruft.
func moveFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read staged file: %w", err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return fmt.Errorf("write final file: %w", err)
	}
	_ = os.Remove(src)
	return nil
}
