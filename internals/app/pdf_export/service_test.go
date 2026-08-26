package pdf_export_test

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"nutrition-planner/config"
	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/day_variant"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/meal"
	"nutrition-planner/internals/app/pdf_export"
	"nutrition-planner/internals/app/plan"
)

// testLogger returns a *slog.Logger that discards everything - the tests
// assert on returned values/errors, not on log output, so a real logger
// would just be test noise.
func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// --- fakes ---------------------------------------------------------------

type fakePlanService struct {
	full    plan.FullResponse
	fullErr error
}

func (f *fakePlanService) Create(plan.CreateRequest) (plan.Response, error) {
	return plan.Response{}, nil
}
func (f *fakePlanService) GetByIDFull(id string) (plan.FullResponse, error) {
	if f.fullErr != nil {
		return plan.FullResponse{}, f.fullErr
	}
	return f.full, nil
}
func (f *fakePlanService) ListByClientID(string) ([]plan.Response, error) { return nil, nil }
func (f *fakePlanService) Update(plan.UpdateRequest) (plan.Response, error) {
	return plan.Response{}, nil
}
func (f *fakePlanService) Delete(string) error { return nil }

type fakeClientService struct {
	resp    client.Response
	respErr error
}

func (f *fakeClientService) Create(client.CreateRequest) (client.Response, error) {
	return client.Response{}, nil
}
func (f *fakeClientService) GetByID(id string) (client.Response, error) {
	if f.respErr != nil {
		return client.Response{}, f.respErr
	}
	return f.resp, nil
}
func (f *fakeClientService) List(int, int) ([]client.Response, error) { return nil, nil }
func (f *fakeClientService) Update(client.UpdateRequest) (client.Response, error) {
	return client.Response{}, nil
}
func (f *fakeClientService) Delete(string) error { return nil }

type fakeRepo struct {
	dir    string
	dirErr error
}

func (f *fakeRepo) ClientExportDir(string) (string, error) {
	if f.dirErr != nil {
		return "", f.dirErr
	}
	return f.dir, nil
}

// --- helpers ---------------------------------------------------------

func testGeneratorConfig() config.GeneratorConfig {
	return config.GeneratorConfig{
		PageSize:            "A4",
		MarginTop:           15,
		MarginBottom:        15,
		MarginLeft:          14,
		MarginRight:         14,
		ColorPageBackground: "F7F5EF",
		ColorCardBackground: "FFFFFF",
		ColorCardBorder:     "E7E4D9",
		ColorPrimary:        "2F5233",
		ColorPrimaryLight:   "E3EDE1",
		ColorOnPrimary:      "FFFFFF",
		ColorText:           "2B2B28",
		ColorMuted:          "6B6B66",
		ColorTableHeader:    "8A8A82",
		ColorDivider:        "E5E2D8",
		FontFamily:          "NotoSans",
	}
}

func minimalFullPlan(clientID string) plan.FullResponse {
	return plan.FullResponse{
		Response: plan.Response{
			ID:       "plan-1",
			ClientID: clientID,
			Title:    "Plan alimentar — Test",
		},
		DayVariants: []day_variant.Response{
			{
				Name: "Zi de antrenament",
				Meals: []meal.Response{
					{
						Name:  "Mic dejun",
						Foods: []food.Response{{BaseFood: "Ovăz", BaseQuantity: 60, BaseUnit: "g"}},
					},
				},
			},
		},
	}
}

func assertValidPDF(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", path, err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("output at %q is not a PDF", path)
	}
}

// --- ExportPlanPDF -------------------------------------------------------

func TestService_ExportPlanPDF_Success(t *testing.T) {
	exportDir := t.TempDir()
	tempDir := t.TempDir()

	svc := pdf_export.NewService(
		&fakePlanService{full: minimalFullPlan("client-1")},
		&fakeClientService{resp: client.Response{ID: "client-1", Name: "Andreea Popescu", Weight: 62, Height: 168}},
		&fakeRepo{dir: exportDir},
		testGeneratorConfig(),
		tempDir,
		testLogger(),
	)

	resp, err := svc.ExportPlanPDF("plan-1")
	if err != nil {
		t.Fatalf("ExportPlanPDF: %v", err)
	}
	if resp.FileName == "" || resp.Path == "" {
		t.Fatalf("expected non-empty FileName/Path, got %+v", resp)
	}
	if filepath.Dir(resp.Path) != exportDir {
		t.Errorf("expected file to be written under %q, got %q", exportDir, resp.Path)
	}
	assertValidPDF(t, resp.Path)
}

func TestService_ExportPlanPDF_StagingFileRemovedAfterMove(t *testing.T) {
	exportDir := t.TempDir()
	tempDir := t.TempDir()

	svc := pdf_export.NewService(
		&fakePlanService{full: minimalFullPlan("client-1")},
		&fakeClientService{resp: client.Response{ID: "client-1", Name: "Andreea Popescu"}},
		&fakeRepo{dir: exportDir},
		testGeneratorConfig(),
		tempDir,
		testLogger(),
	)

	resp, err := svc.ExportPlanPDF("plan-1")
	if err != nil {
		t.Fatalf("ExportPlanPDF: %v", err)
	}

	stagingPath := filepath.Join(tempDir, resp.FileName)
	if _, err := os.Stat(stagingPath); !os.IsNotExist(err) {
		t.Errorf("expected the staging file %q to be gone after a successful export (got err=%v)", stagingPath, err)
	}
}

func TestService_ExportPlanPDF_CreatesMissingTempDir(t *testing.T) {
	exportDir := t.TempDir()
	tempDir := filepath.Join(t.TempDir(), "does", "not", "exist", "yet")

	svc := pdf_export.NewService(
		&fakePlanService{full: minimalFullPlan("client-1")},
		&fakeClientService{resp: client.Response{ID: "client-1", Name: "Andreea Popescu"}},
		&fakeRepo{dir: exportDir},
		testGeneratorConfig(),
		tempDir,
		testLogger(),
	)

	if _, err := svc.ExportPlanPDF("plan-1"); err != nil {
		t.Fatalf("ExportPlanPDF: %v", err)
	}
}

func TestService_ExportPlanPDF_EmptyPlanID(t *testing.T) {
	svc := pdf_export.NewService(
		&fakePlanService{},
		&fakeClientService{},
		&fakeRepo{dir: t.TempDir()},
		testGeneratorConfig(),
		t.TempDir(),
		testLogger(),
	)

	_, err := svc.ExportPlanPDF("")
	if !errors.Is(err, pdf_export.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_ExportPlanPDF_PlanLookupFails(t *testing.T) {
	wantErr := errors.New("plan not found")
	svc := pdf_export.NewService(
		&fakePlanService{fullErr: wantErr},
		&fakeClientService{},
		&fakeRepo{dir: t.TempDir()},
		testGeneratorConfig(),
		t.TempDir(),
		testLogger(),
	)

	_, err := svc.ExportPlanPDF("missing-plan")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("expected wrapped %v, got %v", wantErr, err)
	}
}

func TestService_ExportPlanPDF_ClientLookupFails(t *testing.T) {
	wantErr := errors.New("client not found")
	svc := pdf_export.NewService(
		&fakePlanService{full: minimalFullPlan("missing-client")},
		&fakeClientService{respErr: wantErr},
		&fakeRepo{dir: t.TempDir()},
		testGeneratorConfig(),
		t.TempDir(),
		testLogger(),
	)

	_, err := svc.ExportPlanPDF("plan-1")
	if !errors.Is(err, wantErr) {
		t.Errorf("expected wrapped %v, got %v", wantErr, err)
	}
}

func TestService_ExportPlanPDF_RepoDirFails(t *testing.T) {
	wantErr := errors.New("disk full")
	svc := pdf_export.NewService(
		&fakePlanService{full: minimalFullPlan("client-1")},
		&fakeClientService{resp: client.Response{ID: "client-1", Name: "Andreea Popescu"}},
		&fakeRepo{dirErr: wantErr},
		testGeneratorConfig(),
		t.TempDir(),
		testLogger(),
	)

	_, err := svc.ExportPlanPDF("plan-1")
	if !errors.Is(err, wantErr) {
		t.Errorf("expected wrapped %v, got %v", wantErr, err)
	}
}

func TestService_ExportPlanPDF_InvalidGeneratorConfig(t *testing.T) {
	cfg := testGeneratorConfig()
	cfg.ColorPrimary = "not-a-color"

	svc := pdf_export.NewService(
		&fakePlanService{full: minimalFullPlan("client-1")},
		&fakeClientService{resp: client.Response{ID: "client-1", Name: "Andreea Popescu"}},
		&fakeRepo{dir: t.TempDir()},
		cfg,
		t.TempDir(),
		testLogger(),
	)

	_, err := svc.ExportPlanPDF("plan-1")
	if err == nil {
		t.Fatal("expected an error for an invalid generator config")
	}
}

// --- PreviewPlanPDF --------------------------------------------------

func TestService_PreviewPlanPDF_Success(t *testing.T) {
	exportDir := t.TempDir() // must stay untouched by a preview
	tempDir := t.TempDir()

	svc := pdf_export.NewService(
		&fakePlanService{full: minimalFullPlan("client-1")},
		&fakeClientService{resp: client.Response{ID: "client-1", Name: "Andreea Popescu"}},
		&fakeRepo{dir: exportDir},
		testGeneratorConfig(),
		tempDir,
		testLogger(),
	)

	resp, err := svc.PreviewPlanPDF("plan-1")
	if err != nil {
		t.Fatalf("PreviewPlanPDF: %v", err)
	}
	if filepath.Dir(resp.Path) != tempDir {
		t.Errorf("expected preview to be written under tempDir %q, got %q", tempDir, resp.Path)
	}
	assertValidPDF(t, resp.Path)

	entries, err := os.ReadDir(exportDir)
	if err != nil {
		t.Fatalf("ReadDir(exportDir): %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected the client's export directory to stay untouched by a preview, found %d entries", len(entries))
	}
}

func TestService_PreviewPlanPDF_EmptyPlanID(t *testing.T) {
	svc := pdf_export.NewService(
		&fakePlanService{},
		&fakeClientService{},
		&fakeRepo{dir: t.TempDir()},
		testGeneratorConfig(),
		t.TempDir(),
		testLogger(),
	)

	_, err := svc.PreviewPlanPDF("")
	if !errors.Is(err, pdf_export.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_PreviewPlanPDF_PlanLookupFails(t *testing.T) {
	wantErr := errors.New("plan not found")
	svc := pdf_export.NewService(
		&fakePlanService{fullErr: wantErr},
		&fakeClientService{},
		&fakeRepo{dir: t.TempDir()},
		testGeneratorConfig(),
		t.TempDir(),
		testLogger(),
	)

	_, err := svc.PreviewPlanPDF("plan-1")
	if !errors.Is(err, wantErr) {
		t.Errorf("expected wrapped %v, got %v", wantErr, err)
	}
}
