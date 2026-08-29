package plan_test

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"nutrition-planner/internals/app/plan"

	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type fakeRepo struct {
	fullPlan *plan.Plan
	fullErr  error

	created    *plan.Plan
	createdErr error

	duplicateErr error

	// Captured arguments of the most recent Duplicate call, so tests can
	// assert on exactly what the service decided to pass through.
	gotTargetClientID string
	gotTitle          string
}

func (f *fakeRepo) Create(*plan.Plan) (string, error) { return "", nil }
func (f *fakeRepo) GetByID(string) (*plan.Plan, error) {
	if f.createdErr != nil {
		return nil, f.createdErr
	}
	return f.created, nil
}
func (f *fakeRepo) GetByIDFull(string) (*plan.Plan, error) {
	if f.fullErr != nil {
		return nil, f.fullErr
	}
	return f.fullPlan, nil
}
func (f *fakeRepo) ListByClientID(string) ([]plan.Plan, error) { return nil, nil }
func (f *fakeRepo) Update(*plan.Plan) error                    { return nil }
func (f *fakeRepo) Delete(string) error                        { return nil }
func (f *fakeRepo) Duplicate(source *plan.Plan, targetClientID, newTitle string) (string, error) {
	f.gotTargetClientID = targetClientID
	f.gotTitle = newTitle
	if f.duplicateErr != nil {
		return "", f.duplicateErr
	}
	return "new-plan-id", nil
}

func TestService_Duplicate_EmptyTargetClientFallsBackToSourceClient(t *testing.T) {
	repo := &fakeRepo{
		fullPlan: &plan.Plan{ID: "plan-1", ClientID: "client-source", Title: "Original"},
		created:  &plan.Plan{ID: "new-plan-id", ClientID: "client-source", Title: "Original (Copy)"},
	}
	svc := plan.NewService(repo, testLogger())

	resp, err := svc.Duplicate(plan.DuplicateRequest{SourcePlanID: "plan-1", Title: "Original (Copy)"})
	require.NoError(t, err)
	require.Equal(t, "client-source", repo.gotTargetClientID,
		"an empty TargetClientID must fall back to the source plan's own client")
	require.Equal(t, "client-source", resp.ClientID)
}

func TestService_Duplicate_ExplicitTargetClientIsUsed(t *testing.T) {
	repo := &fakeRepo{
		fullPlan: &plan.Plan{ID: "plan-1", ClientID: "client-source", Title: "Original"},
		created:  &plan.Plan{ID: "new-plan-id", ClientID: "client-other", Title: "Original (Copy)"},
	}
	svc := plan.NewService(repo, testLogger())

	resp, err := svc.Duplicate(plan.DuplicateRequest{
		SourcePlanID:   "plan-1",
		TargetClientID: "client-other",
		Title:          "Original (Copy)",
	})
	require.NoError(t, err)
	require.Equal(t, "client-other", repo.gotTargetClientID)
	require.Equal(t, "client-other", resp.ClientID)
}

func TestService_Duplicate_EmptySourcePlanID(t *testing.T) {
	svc := plan.NewService(&fakeRepo{}, testLogger())
	_, err := svc.Duplicate(plan.DuplicateRequest{Title: "x"})
	require.ErrorIs(t, err, plan.ErrInvalidInput)
}

func TestService_Duplicate_EmptyTitle(t *testing.T) {
	svc := plan.NewService(&fakeRepo{}, testLogger())
	_, err := svc.Duplicate(plan.DuplicateRequest{SourcePlanID: "plan-1"})
	require.ErrorIs(t, err, plan.ErrInvalidInput)
}

func TestService_Duplicate_SourceLookupFails(t *testing.T) {
	wantErr := errors.New("not found")
	svc := plan.NewService(&fakeRepo{fullErr: wantErr}, testLogger())

	_, err := svc.Duplicate(plan.DuplicateRequest{SourcePlanID: "plan-1", Title: "x"})
	require.ErrorIs(t, err, wantErr)
}

func TestService_Duplicate_RepoErrorIsWrapped(t *testing.T) {
	wantErr := errors.New("db error")
	repo := &fakeRepo{
		fullPlan:     &plan.Plan{ID: "plan-1", ClientID: "c1"},
		duplicateErr: wantErr,
	}
	svc := plan.NewService(repo, testLogger())

	_, err := svc.Duplicate(plan.DuplicateRequest{SourcePlanID: "plan-1", Title: "x"})
	require.ErrorIs(t, err, wantErr)
}
