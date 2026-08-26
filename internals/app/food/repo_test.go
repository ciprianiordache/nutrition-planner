package food_test

import (
	"database/sql"
	"testing"

	"nutrition-planner/internals/app/food"

	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/day_variant"
	"nutrition-planner/internals/app/meal"
	"nutrition-planner/internals/app/plan"
	"nutrition-planner/internals/app/supplement"
	"nutrition-planner/internals/testutil"

	"github.com/stretchr/testify/require"
)

func TestRepo_CRUDAndListByMeal(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	planID := testutil.SeedPlan(t, db, &plan.Plan{ClientID: clientID, Title: "Plan"})
	dvID := testutil.SeedDayVariant(t, db, &day_variant.DayVariant{PlanID: planID, Name: "Work"})
	fr := food.NewRepo(db)
	mr := meal.NewRepo(db, fr, supplement.NewRepo(db))
	mealID, err := mr.Create(&meal.Meal{DayVariantID: dvID, Name: "Lunch"})
	require.NoError(t, err)

	r := food.NewRepo(db)
	alt := "Green beans"
	q := 100.0
	id, err := r.Create(&food.Food{
		MealID: mealID, BaseFood: "Chicken", BaseQuantity: 200, BaseUnit: "g",
		AlternativeFood: sql.NullString{String: alt, Valid: true}, AlternativeQuantity: sql.NullFloat64{Float64: q, Valid: true}, Order: 2,
	})
	require.NoError(t, err)
	require.Equal(t, id, mustGet(t, r, id).ID)

	list, err := r.ListByMealID(mealID)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "Chicken", list[0].BaseFood)
	require.True(t, list[0].AlternativeFood.Valid)
	require.Equal(t, alt, list[0].AlternativeFood.String)

	got := list[0]
	got.BaseFood = "Turkey"
	require.NoError(t, r.Update(&got))
	require.Equal(t, "Turkey", mustGet(t, r, id).BaseFood)

	require.NoError(t, r.Delete(id))
	_, err = r.GetByID(id)
	require.ErrorIs(t, err, food.ErrNotFound)
}

func TestRepo_ListOrderingAndDeleteByMealID(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	planID := testutil.SeedPlan(t, db, &plan.Plan{ClientID: clientID, Title: "Plan"})
	dvID := testutil.SeedDayVariant(t, db, &day_variant.DayVariant{PlanID: planID, Name: "Work"})
	mr := meal.NewRepo(db, food.NewRepo(db), supplement.NewRepo(db))
	mealID, err := mr.Create(&meal.Meal{DayVariantID: dvID, Name: "Lunch"})
	require.NoError(t, err)

	r := food.NewRepo(db)
	_, err = r.Create(&food.Food{MealID: mealID, BaseFood: "Second", BaseQuantity: 1, BaseUnit: "g", Order: 2})
	require.NoError(t, err)
	_, err = r.Create(&food.Food{MealID: mealID, BaseFood: "First", BaseQuantity: 1, BaseUnit: "g", Order: 1})
	require.NoError(t, err)

	list, err := r.ListByMealID(mealID)
	require.NoError(t, err)
	require.Equal(t, "First", list[0].BaseFood)
	require.Equal(t, "Second", list[1].BaseFood)

	require.NoError(t, r.DeleteByMealID(mealID))
	list, err = r.ListByMealID(mealID)
	require.NoError(t, err)
	require.Empty(t, list)
}

func mustGet(t *testing.T, r food.Repo, id string) *food.Food {
	t.Helper()
	f, err := r.GetByID(id)
	require.NoError(t, err)
	return f
}
