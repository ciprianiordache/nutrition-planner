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

func TestRepo_Duplicate_CopiesFullTreeWithFreshIDs(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})

	foodRepo := food.NewRepo(db)
	supRepo := supplement.NewRepo(db)
	mealRepo := meal.NewRepo(db, foodRepo, supRepo)
	dvRepo := day_variant.NewRepo(db, mealRepo)
	r := plan.NewRepo(db, dvRepo)

	sourceID, err := r.Create(&plan.Plan{
		ClientID:               clientID,
		Title:                  "Original",
		GeneralRecommendations: "Drink water",
		Footer:                 "Some footer",
	})
	require.NoError(t, err)
	dvID, err := dvRepo.Create(&day_variant.DayVariant{PlanID: sourceID, Name: "Training day", Order: 1})
	require.NoError(t, err)
	mealID, err := mealRepo.Create(&meal.Meal{DayVariantID: dvID, Name: "Breakfast", Note: "Early", Order: 1})
	require.NoError(t, err)
	_, err = foodRepo.Create(&food.Food{MealID: mealID, BaseFood: "Oats", BaseQuantity: 80, BaseUnit: "g", Order: 1})
	require.NoError(t, err)
	_, err = supRepo.Create(&supplement.Supplement{MealID: mealID, Supplement: "Omega 3", Quantity: 1, Order: 1})
	require.NoError(t, err)

	source, err := r.GetByIDFull(sourceID)
	require.NoError(t, err)

	newID, err := r.Duplicate(source, "Original (Copy)")
	require.NoError(t, err)
	require.NotEqual(t, sourceID, newID, "the duplicate must be a distinct plan, not the same row")

	got, err := r.GetByIDFull(newID)
	require.NoError(t, err)
	require.Equal(t, "Original (Copy)", got.Title)
	require.Equal(t, clientID, got.ClientID)
	require.Equal(t, "Drink water", got.GeneralRecommendations)
	require.Equal(t, "Some footer", got.Footer)

	require.Len(t, got.DayVariants, 1)
	require.NotEqual(t, dvID, got.DayVariants[0].ID, "the day variant must get a fresh ID too")
	require.Equal(t, "Training day", got.DayVariants[0].Name)

	require.Len(t, got.DayVariants[0].Meals, 1)
	require.NotEqual(t, mealID, got.DayVariants[0].Meals[0].ID)
	require.Equal(t, "Breakfast", got.DayVariants[0].Meals[0].Name)
	require.Equal(t, "Early", got.DayVariants[0].Meals[0].Note)

	require.Len(t, got.DayVariants[0].Meals[0].Foods, 1)
	require.Equal(t, "Oats", got.DayVariants[0].Meals[0].Foods[0].BaseFood)
	require.Equal(t, 80.0, got.DayVariants[0].Meals[0].Foods[0].BaseQuantity)

	require.Len(t, got.DayVariants[0].Meals[0].Supplements, 1)
	require.Equal(t, "Omega 3", got.DayVariants[0].Meals[0].Supplements[0].Supplement)

	// The original must be completely untouched by duplicating it.
	original, err := r.GetByIDFull(sourceID)
	require.NoError(t, err)
	require.Equal(t, "Original", original.Title)
	require.Len(t, original.DayVariants, 1)
	require.Equal(t, dvID, original.DayVariants[0].ID)
}

func TestRepo_Duplicate_WithNoDayVariants(t *testing.T) {
	db := testutil.NewDatabase(t)
	clientID := testutil.SeedClient(t, db, &client.Client{Name: "Client", Weight: 80, Height: 180})

	foodRepo := food.NewRepo(db)
	supRepo := supplement.NewRepo(db)
	mealRepo := meal.NewRepo(db, foodRepo, supRepo)
	dvRepo := day_variant.NewRepo(db, mealRepo)
	r := plan.NewRepo(db, dvRepo)

	sourceID, err := r.Create(&plan.Plan{ClientID: clientID, Title: "Empty plan"})
	require.NoError(t, err)

	source, err := r.GetByIDFull(sourceID)
	require.NoError(t, err)

	newID, err := r.Duplicate(source, "Empty plan (Copy)")
	require.NoError(t, err)

	got, err := r.GetByIDFull(newID)
	require.NoError(t, err)
	require.Equal(t, "Empty plan (Copy)", got.Title)
	require.Empty(t, got.DayVariants)
}
