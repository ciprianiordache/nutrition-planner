package testutil

import (
	"testing"

	"nutrition-planner/config"
	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/day_variant"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/meal"
	"nutrition-planner/internals/app/plan"
	"nutrition-planner/internals/app/supplement"
	"nutrition-planner/internals/database"

	"github.com/ciprianiordache/schema-builder"
	_ "modernc.org/sqlite"
)

func NewDatabase(t *testing.T) *database.Database {
	t.Helper()

	cfg := config.DatabaseConfig{
		Driver: "sqlite",
		Name:   t.TempDir() + "/nutrition_test",
	}
	db := database.New(cfg, "")
	if err := db.Connect(); err != nil {
		t.Fatalf("connect test database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	builder := schema.New(db, schema.SQLite{})
	models := []any{
		client.Client{},
		plan.Plan{},
		day_variant.DayVariant{},
		meal.Meal{},
		food.Food{},
		supplement.Supplement{},
	}
	if err := schema.ValidateModels(models...); err != nil {
		t.Fatalf("validate application models: %v", err)
	}
	if err := builder.CreateSchema(models...); err != nil {
		t.Fatalf("create test schema: %v", err)
	}

	return db
}

func SeedClient(t *testing.T, db *database.Database, c *client.Client) string {
	t.Helper()
	id, err := client.NewRepo(db).Create(c)
	if err != nil {
		t.Fatalf("seed client: %v", err)
	}
	return id
}

func SeedPlan(t *testing.T, db *database.Database, p *plan.Plan) string {
	t.Helper()
	id, err := plan.NewRepo(db, nil).Create(p)
	if err != nil {
		t.Fatalf("seed plan: %v", err)
	}
	return id
}

func SeedDayVariant(t *testing.T, db *database.Database, dv *day_variant.DayVariant) string {
	t.Helper()
	id, err := day_variant.NewRepo(db, nil).Create(dv)
	if err != nil {
		t.Fatalf("seed day variant: %v", err)
	}
	return id
}

func SeedMeal(t *testing.T, db *database.Database, m *meal.Meal) string {
	t.Helper()
	fr := food.NewRepo(db)
	sr := supplement.NewRepo(db)
	id, err := meal.NewRepo(db, fr, sr).Create(m)
	if err != nil {
		t.Fatalf("seed meal: %v", err)
	}
	return id
}
