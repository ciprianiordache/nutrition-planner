package supplement_test

import (
	"testing"

	"nutrition-planner/internals/app/supplement"

	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/day_variant"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/meal"
	"nutrition-planner/internals/app/plan"
	"nutrition-planner/internals/testutil"

	"github.com/stretchr/testify/require"
)

func TestRepo_CRUDAndListByMeal(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	planID := testutil.SeedPlan(t, db, &plan.Plan{ClientID: clientID, Title: "Plan"})
	dvID := testutil.SeedDayVariant(t, db, &day_variant.DayVariant{PlanID: planID, Name: "Work"})
	mr := meal.NewRepo(db, food.NewRepo(db), supplement.NewRepo(db))
	mealID, err := mr.Create(&meal.Meal{DayVariantID: dvID, Name: "Dinner"})
	require.NoError(t, err)

	r := supplement.NewRepo(db)
	id, err := r.Create(&supplement.Supplement{MealID: mealID, Supplement: "Omega 3", Quantity: 1, Order: 2})
	require.NoError(t, err)

	got, err := r.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, "Omega 3", got.Supplement)

	got.Supplement = "Vitamin D"
	require.NoError(t, r.Update(got))
	require.Equal(t, "Vitamin D", mustGet(t, r, id).Supplement)

	list, err := r.ListByMealID(mealID)
	require.NoError(t, err)
	require.Len(t, list, 1)

	require.NoError(t, r.Delete(id))
	_, err = r.GetByID(id)
	require.ErrorIs(t, err, supplement.ErrNotFound)
}

func TestRepo_ListOrderingAndDeleteByMealID(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	planID := testutil.SeedPlan(t, db, &plan.Plan{ClientID: clientID, Title: "Plan"})
	dvID := testutil.SeedDayVariant(t, db, &day_variant.DayVariant{PlanID: planID, Name: "Work"})
	mr := meal.NewRepo(db, food.NewRepo(db), supplement.NewRepo(db))
	mealID, err := mr.Create(&meal.Meal{DayVariantID: dvID, Name: "Dinner"})
	require.NoError(t, err)

	r := supplement.NewRepo(db)
	_, err = r.Create(&supplement.Supplement{MealID: mealID, Supplement: "Second", Quantity: 1, Order: 2})
	require.NoError(t, err)
	_, err = r.Create(&supplement.Supplement{MealID: mealID, Supplement: "First", Quantity: 1, Order: 1})
	require.NoError(t, err)

	list, err := r.ListByMealID(mealID)
	require.NoError(t, err)
	require.Equal(t, "First", list[0].Supplement)
	require.Equal(t, "Second", list[1].Supplement)

	require.NoError(t, r.DeleteByMealID(mealID))
	list, err = r.ListByMealID(mealID)
	require.NoError(t, err)
	require.Empty(t, list)
}

func mustGet(t *testing.T, r supplement.Repo, id string) *supplement.Supplement {
	t.Helper()
	s, err := r.GetByID(id)
	require.NoError(t, err)
	return s
}
