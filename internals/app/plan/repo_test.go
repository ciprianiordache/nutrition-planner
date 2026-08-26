package plan_test

import (
	"testing"

	"nutrition-planner/internals/app/plan"

	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/day_variant"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/meal"
	"nutrition-planner/internals/app/supplement"
	"nutrition-planner/internals/testutil"

	"github.com/stretchr/testify/require"
)

func TestRepo_CRUDAndListByClient(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})
	r := plan.NewRepo(db, nil)

	p := &plan.Plan{ClientID: clientID, Title: "Plan 1", GeneralRecommendations: "Water", Footer: "Footer"}
	id, err := r.Create(p)
	require.NoError(t, err)
	require.Equal(t, id, p.ID)
	require.False(t, p.CreatedAt.IsZero())

	got, err := r.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, "Plan 1", got.Title)
	require.Equal(t, clientID, got.ClientID)

	list, err := r.ListByClientID(clientID)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, id, list[0].ID)

	got.Title = "Updated"
	require.NoError(t, r.Update(got))
	updated, err := r.GetByID(id)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Title)

	require.NoError(t, r.Delete(id))
	_, err = r.GetByID(id)
	require.ErrorIs(t, err, plan.ErrNotFound)
}

func TestRepo_GetByIDFull(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})

	foodRepo := food.NewRepo(db)
	supRepo := supplement.NewRepo(db)
	mealRepo := meal.NewRepo(db, foodRepo, supRepo)
	dvRepo := day_variant.NewRepo(db, mealRepo)
	r := plan.NewRepo(db, dvRepo)

	planID, err := r.Create(&plan.Plan{ClientID: clientID, Title: "Full plan"})
	require.NoError(t, err)
	dvID, err := dvRepo.Create(&day_variant.DayVariant{PlanID: planID, Name: "Work", Order: 1})
	require.NoError(t, err)
	mealID, err := mealRepo.Create(&meal.Meal{DayVariantID: dvID, Name: "Breakfast", Order: 1})
	require.NoError(t, err)
	_, err = foodRepo.Create(&food.Food{MealID: mealID, BaseFood: "Oats", BaseQuantity: 80, BaseUnit: "g", Order: 1})
	require.NoError(t, err)
	_, err = supRepo.Create(&supplement.Supplement{MealID: mealID, Supplement: "Omega 3", Quantity: 1, Order: 1})
	require.NoError(t, err)

	got, err := r.GetByIDFull(planID)
	require.NoError(t, err)
	require.Len(t, got.DayVariants, 1)
	require.Len(t, got.DayVariants[0].Meals, 1)
	require.Len(t, got.DayVariants[0].Meals[0].Foods, 1)
	require.Len(t, got.DayVariants[0].Meals[0].Supplements, 1)
	require.Equal(t, "Oats", got.DayVariants[0].Meals[0].Foods[0].BaseFood)
}
