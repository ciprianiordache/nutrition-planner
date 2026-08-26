package day_variant_test

import (
	"testing"

	"nutrition-planner/internals/app/day_variant"

	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/meal"
	"nutrition-planner/internals/app/plan"
	"nutrition-planner/internals/app/supplement"
	"nutrition-planner/internals/testutil"

	"github.com/stretchr/testify/require"
)

func TestRepo_CRUDAndListByPlan(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	planID := testutil.SeedPlan(t, db, &plan.Plan{ClientID: clientID, Title: "Plan"})
	r := day_variant.NewRepo(db, nil)

	id, err := r.Create(&day_variant.DayVariant{PlanID: planID, Name: "Working", Order: 2})
	require.NoError(t, err)

	got, err := r.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, "Working", got.Name)

	_, err = r.Create(&day_variant.DayVariant{PlanID: planID, Name: "Rest", Order: 1})
	require.NoError(t, err)
	list, err := r.ListByPlanID(planID)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, "Rest", list[0].Name)

	got.Name = "Updated"
	require.NoError(t, r.Update(got))
	require.Equal(t, "Updated", mustGet(t, r, id).Name)

	// Regression test: Delete must physically delete the row.
	require.NoError(t, r.Delete(id))
	_, err = r.GetByID(id)
	require.ErrorIs(t, err, day_variant.ErrNotFound)
}

func TestRepo_ListByPlanIDFull(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	planID := testutil.SeedPlan(t, db, &plan.Plan{ClientID: clientID, Title: "Plan"})

	fr := food.NewRepo(db)
	sr := supplement.NewRepo(db)
	mr := meal.NewRepo(db, fr, sr)
	r := day_variant.NewRepo(db, mr)

	dvID, err := r.Create(&day_variant.DayVariant{PlanID: planID, Name: "Work", Order: 1})
	require.NoError(t, err)
	mealID, err := mr.Create(&meal.Meal{DayVariantID: dvID, Name: "Lunch", Order: 1})
	require.NoError(t, err)
	_, err = fr.Create(&food.Food{MealID: mealID, BaseFood: "Rice", BaseQuantity: 150, BaseUnit: "g"})
	require.NoError(t, err)
	_, err = sr.Create(&supplement.Supplement{MealID: mealID, Supplement: "Vitamin D", Quantity: 1})
	require.NoError(t, err)

	got, err := r.ListByPlanIDFull(planID)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Len(t, got[0].Meals, 1)
	require.Len(t, got[0].Meals[0].Foods, 1)
	require.Len(t, got[0].Meals[0].Supplements, 1)
}

func TestRepo_DeleteByPlanID(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	planID := testutil.SeedPlan(t, db, &plan.Plan{ClientID: clientID, Title: "Plan"})
	r := day_variant.NewRepo(db, nil)
	_, err := r.Create(&day_variant.DayVariant{PlanID: planID, Name: "Work"})
	require.NoError(t, err)

	require.NoError(t, r.DeleteByPlanID(planID))
	got, err := r.ListByPlanID(planID)
	require.NoError(t, err)
	require.Empty(t, got)
}

func mustGet(t *testing.T, r day_variant.Repo, id string) *day_variant.DayVariant {
	t.Helper()
	v, err := r.GetByID(id)
	require.NoError(t, err)
	return v
}
