package settings_test

import (
	"testing"

	"nutrition-planner/internals/app/settings"
	"nutrition-planner/internals/testutil"

	"github.com/stretchr/testify/require"
)

func TestRepo_Get_CreatesRowOnFirstCall(t *testing.T) {
	db := testutil.NewDatabase(t)
	r := settings.NewRepo(db)

	s, err := r.Get()
	require.NoError(t, err)
	require.NotNil(t, s)
	// Everything should be empty on a freshly-created row - "empty" is
	// exactly what tells the rest of the app "fall back to config.yaml".
	require.Empty(t, s.ColorPrimary)
	require.Empty(t, s.ExportBaseDirectory)
	require.False(t, s.CreatedAt.IsZero())
}

func TestRepo_Get_IsIdempotent(t *testing.T) {
	db := testutil.NewDatabase(t)
	r := settings.NewRepo(db)

	first, err := r.Get()
	require.NoError(t, err)
	second, err := r.Get()
	require.NoError(t, err)

	require.Equal(t, first.ID, second.ID)
	// first.CreatedAt comes straight from Create() and still carries a
	// monotonic clock reading (time.Now() always does); second.CreatedAt
	// comes back from an actual SQLite round-trip, which strips it (any
	// value parsed/scanned from storage is monotonic-free by construction).
	// They represent the exact same instant, but require.Equal compares
	// the time.Time struct fields directly, and Go's docs are explicit
	// that a monotonic and a non-monotonic reading of the same instant are
	// != under that kind of structural comparison - only .Equal() (or
	// stripping it via .Round(0)) compares them the way "same point in
	// time" actually means. https://pkg.go.dev/time#hdr-Monotonic_Clocks
	require.True(t, first.CreatedAt.Equal(second.CreatedAt),
		"expected the same instant, got %v vs %v", first.CreatedAt, second.CreatedAt)
}

func TestRepo_Update_Persists(t *testing.T) {
	db := testutil.NewDatabase(t)
	r := settings.NewRepo(db)

	s, err := r.Get()
	require.NoError(t, err)

	s.ColorPrimary = "FF0000"
	s.ExportBaseDirectory = "/custom/export"
	require.NoError(t, r.Update(s))

	got, err := r.Get()
	require.NoError(t, err)
	require.Equal(t, "FF0000", got.ColorPrimary)
	require.Equal(t, "/custom/export", got.ExportBaseDirectory)
}

func TestRepo_Update_CannotCreateASecondRow(t *testing.T) {
	db := testutil.NewDatabase(t)
	r := settings.NewRepo(db)

	s, err := r.Get()
	require.NoError(t, err)

	// Even if a caller tampers with the ID before calling Update, the repo
	// must still target the one singleton row - there's no legitimate way
	// to end up with two settings rows.
	s.ID = "some-other-id"
	s.ColorPrimary = "00FF00"
	require.NoError(t, r.Update(s))

	got, err := r.Get()
	require.NoError(t, err)
	require.Equal(t, "00FF00", got.ColorPrimary)
}

func TestRepo_Get_KeepsFixedSingletonID_NotARandomUUID(t *testing.T) {
	// Settings.ID carries a "uuid" tag purely to satisfy schema-builder's
	// validation (a primary_key field must declare "auto" or "uuid") -
	// crud-depot's Create() only generates a random UUID when the field is
	// still zero-valued at insert time, and repo.Get() always pre-sets ID
	// to singletonID before calling Create, so the row's ID must come back
	// as that exact fixed string, never a generated UUID. This is exactly
	// the behavior the model.go comment promises - pin it down here so a
	// future crud-depot/schema-builder upgrade can't silently break it.
	db := testutil.NewDatabase(t)
	r := settings.NewRepo(db)

	s, err := r.Get()
	require.NoError(t, err)
	require.Equal(t, "app-settings", s.ID)
}
