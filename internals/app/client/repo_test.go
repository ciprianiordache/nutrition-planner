package client_test

import (
	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepo_CRUD(t *testing.T) {
	db := testutil.NewDatabase(t)
	r := client.NewRepo(db)

	c := &client.Client{Name: "Ciprian", Weight: 82.5, Height: 180, Target: "maintenance"}
	id, err := r.Create(c)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	require.Equal(t, id, c.ID)
	require.False(t, c.CreatedAt.IsZero())
	require.False(t, c.UpdatedAt.IsZero())

	got, err := r.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, c.Name, got.Name)
	require.Equal(t, c.Weight, got.Weight)
	require.Equal(t, c.Height, got.Height)
	require.Equal(t, c.Target, got.Target)

	got.Name = "Updated"
	got.Weight = 80
	before := got.UpdatedAt
	time.Sleep(2 * time.Millisecond)
	require.NoError(t, r.Update(got))

	updated, err := r.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Name)
	require.Equal(t, 80.0, updated.Weight)
	require.GreaterOrEqual(t, updated.UpdatedAt, before)

	list, err := r.List(0, 20)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, id, list[0].ID)

	require.NoError(t, r.Delete(id))
	_, err = r.GetByID(id)
	require.ErrorIs(t, err, client.ErrNotFound)
}

func TestRepo_ListPagination(t *testing.T) {
	db := testutil.NewDatabase(t)
	r := client.NewRepo(db)

	for i := 0; i < 3; i++ {
		_, err := r.Create(&client.Client{Name: string(rune('A' + i)), Weight: 70, Height: 175})
		require.NoError(t, err)
	}

	got, err := r.List(1, 1)
	require.NoError(t, err)
	require.Len(t, got, 1)
}
