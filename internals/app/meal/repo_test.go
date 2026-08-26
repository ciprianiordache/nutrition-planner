package meal_test

import (
	"testing"

	"nutrition-planner/internals/app/meal"

	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/day_variant"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/plan"
	"nutrition-planner/internals/app/supplement"
	"nutrition-planner/internals/testutil"

	"github.com/stretchr/testify/require"
)

func TestRepo_CRUDAndListByDayVariant(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	planID := testutil.SeedPlan(t, db, &plan.Plan{ClientID: clientID, Title: "Plan"})
	dvID := testutil.SeedDayVariant(t, db, &day_variant.DayVariant{PlanID: planID, Name: "Work"})

	fr := food.NewRepo(db)
	sr := supplement.NewRepo(db)
	r := meal.NewRepo(db, fr, sr)

	id, err := r.Create(&meal.Meal{DayVariantID: dvID, Name: "Lunch", Note: "Note", Order: 2})
	require.NoError(t, err)

	got, err := r.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, "Lunch", got.Name)

	_, err = r.Create(&meal.Meal{DayVariantID: dvID, Name: "Breakfast", Order: 1})
	require.NoError(t, err)
	list, err := r.ListByDayVariantID(dvID)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, "Breakfast", list[0].Name)

	got.Name = "Updated"
	require.NoError(t, r.Update(got))
	require.NoError(t, r.Delete(id))
	_, err = r.GetByID(id)
	require.ErrorIs(t, err, meal.ErrNotFound)
}

func TestRepo_ListByDayVariantIDFullAndCascadeDelete(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	planID := testutil.SeedPlan(t, db, &plan.Plan{ClientID: clientID, Title: "Plan"})
	dvID := testutil.SeedDayVariant(t, db, &day_variant.DayVariant{PlanID: planID, Name: "Work"})

	fr := food.NewRepo(db)
	sr := supplement.NewRepo(db)
	r := meal.NewRepo(db, fr, sr)

	mealID, err := r.Create(&meal.Meal{DayVariantID: dvID, Name: "Dinner"})
	require.NoError(t, err)
	_, err = fr.Create(&food.Food{MealID: mealID, BaseFood: "Chicken", BaseQuantity: 200, BaseUnit: "g"})
	require.NoError(t, err)
	_, err = sr.Create(&supplement.Supplement{MealID: mealID, Supplement: "Omega 3", Quantity: 1})
	require.NoError(t, err)

	full, err := r.ListByDayVariantIDFull(dvID)
	require.NoError(t, err)
	require.Len(t, full, 1)
	require.Len(t, full[0].Foods, 1)
	require.Len(t, full[0].Supplements, 1)

	require.NoError(t, r.Delete(mealID))
	foods, err := fr.ListByMealID(mealID)
	require.NoError(t, err)
	require.Empty(t, foods)
	supps, err := sr.ListByMealID(mealID)
	require.NoError(t, err)
	require.Empty(t, supps)
}

func TestRepo_DeleteByDayVariantID(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	planID := testutil.SeedPlan(t, db, &plan.Plan{ClientID: clientID, Title: "Plan"})
	dvID := testutil.SeedDayVariant(t, db, &day_variant.DayVariant{PlanID: planID, Name: "Work"})

	fr := food.NewRepo(db)
	sr := supplement.NewRepo(db)
	r := meal.NewRepo(db, fr, sr)
	_, err := r.Create(&meal.Meal{DayVariantID: dvID, Name: "Meal"})
	require.NoError(t, err)

	require.NoError(t, r.DeleteByDayVariantID(dvID))
	got, err := r.ListByDayVariantID(dvID)
	require.NoError(t, err)
	require.Empty(t, got)
}
