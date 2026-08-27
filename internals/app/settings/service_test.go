package settings_test

import (
	"io"
	"log/slog"
	"testing"

	"nutrition-planner/internals/app/settings"
	"nutrition-planner/internals/testutil"

	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestService_Get_ReturnsDefaults(t *testing.T) {
	db := testutil.NewDatabase(t)
	svc := settings.NewService(settings.NewRepo(db), testLogger())

	resp, err := svc.Get()
	require.NoError(t, err)
	require.Empty(t, resp.ColorPrimary)
	require.Empty(t, resp.ExportBaseDirectory)
}

func TestService_Update_OnlyChangesProvidedFields(t *testing.T) {
	db := testutil.NewDatabase(t)
	svc := settings.NewService(settings.NewRepo(db), testLogger())

	primary := "1B5E3C"
	_, err := svc.Update(settings.UpdateRequest{ColorPrimary: &primary})
	require.NoError(t, err)

	// A second update touching a different field must not clobber the
	// first one - this is the whole point of UpdateRequest using pointers.
	exportDir := "/exports"
	resp, err := svc.Update(settings.UpdateRequest{ExportBaseDirectory: &exportDir})
	require.NoError(t, err)

	require.Equal(t, "1B5E3C", resp.ColorPrimary)
	require.Equal(t, "/exports", resp.ExportBaseDirectory)
}

func TestService_Update_PersistsAcrossGetCalls(t *testing.T) {
	db := testutil.NewDatabase(t)
	svc := settings.NewService(settings.NewRepo(db), testLogger())

	theme := "dark"
	_, err := svc.Update(settings.UpdateRequest{UITheme: &theme})
	require.NoError(t, err)

	resp, err := svc.Get()
	require.NoError(t, err)
	require.Equal(t, "dark", resp.UITheme)
}
