// cmd/seed populates the configured database (config.yaml/.env, same as
// the Wails application) with 4 clients, each with 2 or 3 plans, and each
// plan receives day variants, meals, foods, and supplements
// — so you can visually verify the entire UI flow.
//
// Run from the project root:
//
//	go run ./cmd/seed
//
// It is idempotent: if the database already contains clients, nothing is inserted
// (to avoid duplicating data on every run). To force re-seeding:
//
//	go run ./cmd/seed -force
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"nutrition-planner/config"
	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/day_variant"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/meal"
	"nutrition-planner/internals/app/plan"
	"nutrition-planner/internals/app/supplement"
	"nutrition-planner/internals/database"

	schema "github.com/ciprianiordache/schema-builder"
	_ "modernc.org/sqlite"
)

type seedClient struct {
	Name   string
	Weight float64
	Height float64
	Target string
	// PlanTitles — the number of titles determines how many plans the client receives (2 or 3).
	PlanTitles []string
}

type seedFood struct {
	BaseFood            string
	BaseQuantity        float64
	BaseUnit            string
	AlternativeFood     *string
	AlternativeQuantity *float64
	AlternativeUnit     *string
}

type seedSupplement struct {
	Supplement string
	Quantity   float64
}

type seedMeal struct {
	Name        string
	Note        string
	Order       int
	Foods       []seedFood
	Supplements []seedSupplement
}

type seedDayVariant struct {
	Name  string
	Order int
	Meals []seedMeal
}

func main() {
	force := flag.Bool("force", false, "seed even if the database already contains clients")
	flag.Parse()

	cfg, err := config.Load("config.yaml", ".env")
	if err != nil {
		log.Fatalf("❌ Unable to load config: %v", err)
	}

	// MUST match app.go's environment isolation exactly, or seeded data
	// silently lands in a different .db file than the one the running app
	// reads from - that's why this lives in a single shared method
	// (config.Config.EnvironmentIsolation) instead of being duplicated here.
	dbCfg, _ := cfg.EnvironmentIsolation()

	db := database.New(dbCfg, cfg.App.DataDirectory)
	if err := db.Connect(); err != nil {
		log.Fatalf("❌ Unable to connect to database: %v", err)
	}
	defer db.Close()

	builder := schema.New(db, schema.SQLite{})
	if err := builder.CreateSchema(
		client.Client{},
		plan.Plan{},
		day_variant.DayVariant{},
		meal.Meal{},
		food.Food{},
		supplement.Supplement{},
	); err != nil {
		log.Fatalf("❌ Unable to create schema: %v", err)
	}

	// Same wiring as app.go — "child" repositories are initialized before
	// their "parent" repositories.
	foodRepo := food.NewRepo(db)
	supplementRepo := supplement.NewRepo(db)
	mealRepo := meal.NewRepo(db, foodRepo, supplementRepo)
	dayVariantRepo := day_variant.NewRepo(db, mealRepo)
	planRepo := plan.NewRepo(db, dayVariantRepo)
	clientRepo := client.NewRepo(db)

	existing, err := clientRepo.List(0, 1)
	if err != nil {
		log.Fatalf("❌ Unable to check existing clients: %v", err)
	}
	if len(existing) > 0 && !*force {
		fmt.Println("⚠️  The database already contains clients — nothing will be added.")
		fmt.Println("    Run with -force if you want to seed the database anyway: go run ./cmd/seed -force")
		return
	}

	seedClients := []seedClient{
		{
			Name: "Andreea Popescu", Weight: 62, Height: 168, Target: "Weight loss",
			PlanTitles: []string{"Meal Plan — January", "Meal Plan — February"},
		},
		{
			Name: "Mihai Ionescu", Weight: 88, Height: 182, Target: "Muscle gain",
			PlanTitles: []string{"Meal Plan — Bulk Phase 1", "Meal Plan — Bulk Phase 2", "Meal Plan — Maintenance"},
		},
		{
			Name: "Elena Radu", Weight: 71, Height: 165, Target: "Maintenance",
			PlanTitles: []string{"Meal Plan — Summer", "Meal Plan — Autumn"},
		},
		{
			Name: "Vlad Marinescu", Weight: 95, Height: 190, Target: "Body recomposition",
			PlanTitles: []string{"Meal Plan — Phase 1", "Meal Plan — Phase 2", "Meal Plan — Phase 3"},
		},
	}

	const generalRecommendations = "3 liters of water daily, avoid fried foods, get at least 7 hours of sleep."
	const footer = "Discount Code: VITABOLIC, MihailVTB code - 10% off + free shipping"

	for _, sc := range seedClients {
		c := &client.Client{Name: sc.Name, Weight: sc.Weight, Height: sc.Height, Target: sc.Target}
		clientID, err := clientRepo.Create(c)
		if err != nil {
			log.Fatalf("❌ create client %q: %v", sc.Name, err)
		}
		fmt.Printf("✔ Client created: %s (%s) — %d plans\n", sc.Name, clientID, len(sc.PlanTitles))

		for _, title := range sc.PlanTitles {
			p := &plan.Plan{
				ClientID:               clientID,
				Title:                  title,
				GeneralRecommendations: generalRecommendations,
				Footer:                 footer,
			}
			planID, err := planRepo.Create(p)
			if err != nil {
				log.Fatalf("❌ create plan %q for %s: %v", title, sc.Name, err)
			}
			fmt.Printf("  ↳ Plan created: %s (%s)\n", title, planID)

			if err := seedPlanContent(planID, dayVariantRepo, mealRepo, foodRepo, supplementRepo); err != nil {
				log.Fatalf("❌ seed content for plan %q: %v", title, err)
			}
			fmt.Printf("    ↳ Content added: %d day variants\n", len(demoDayVariants()))
		}
	}

	fmt.Println("✅ Seed complete: 4 clients with complete plans (days, meals, foods, supplements).")
}

// seedPlanContent adds the same two demo day variants
// (Training Day / Rest Day) to a given plan, each with meals,
// foods, and supplements — enough to see the entire UI flow.
func seedPlanContent(
	planID string,
	dayVariantRepo day_variant.Repo,
	mealRepo meal.Repo,
	foodRepo food.Repo,
	supplementRepo supplement.Repo,
) error {
	for _, d := range demoDayVariants() {
		dv := &day_variant.DayVariant{PlanID: planID, Name: d.Name, Order: d.Order}
		dvID, err := dayVariantRepo.Create(dv)
		if err != nil {
			return fmt.Errorf("create day variant %q: %w", d.Name, err)
		}

		for _, sm := range d.Meals {
			m := &meal.Meal{DayVariantID: dvID, Name: sm.Name, Note: sm.Note, Order: sm.Order}
			mealID, err := mealRepo.Create(m)
			if err != nil {
				return fmt.Errorf("create meal %q: %w", sm.Name, err)
			}

			for i, sf := range sm.Foods {
				f := &food.Food{
					MealID:       mealID,
					BaseFood:     sf.BaseFood,
					BaseQuantity: sf.BaseQuantity,
					BaseUnit:     sf.BaseUnit,
					Order:        i,
				}
				if sf.AlternativeFood != nil {
					f.AlternativeFood = sql.NullString{String: *sf.AlternativeFood, Valid: true}
				}
				if sf.AlternativeQuantity != nil {
					f.AlternativeQuantity = sql.NullFloat64{Float64: *sf.AlternativeQuantity, Valid: true}
				}
				if sf.AlternativeUnit != nil {
					f.AlternativeUnit = sql.NullString{String: *sf.AlternativeUnit, Valid: true}
				}
				if _, err := foodRepo.Create(f); err != nil {
					return fmt.Errorf("create food %q: %w", sf.BaseFood, err)
				}
			}

			for i, ss := range sm.Supplements {
				s := &supplement.Supplement{MealID: mealID, Supplement: ss.Supplement, Quantity: ss.Quantity, Order: i}
				if _, err := supplementRepo.Create(s); err != nil {
					return fmt.Errorf("create supplement %q: %w", ss.Supplement, err)
				}
			}
		}
	}
	return nil
}

func demoDayVariants() []seedDayVariant {
	return []seedDayVariant{
		{
			Name: "Training Day", Order: 0,
			Meals: []seedMeal{
				{
					Name: "Breakfast", Note: "Within one hour of waking up", Order: 0,
					Foods: []seedFood{
						{BaseFood: "Oats", BaseQuantity: 60, BaseUnit: "g", AlternativeFood: strPtr("Brown rice"), AlternativeQuantity: floatPtr(50), AlternativeUnit: strPtr("g")},
						{BaseFood: "Eggs", BaseQuantity: 3, BaseUnit: "pcs", AlternativeFood: strPtr("Cottage cheese"), AlternativeQuantity: floatPtr(150), AlternativeUnit: strPtr("g")},
						{BaseFood: "Blueberries", BaseQuantity: 80, BaseUnit: "g"},
					},
					Supplements: []seedSupplement{{Supplement: "Multivitamin", Quantity: 1}},
				},
				{
					Name: "Lunch", Note: "3–4 hours after the previous meal", Order: 1,
					Foods: []seedFood{
						{BaseFood: "Chicken breast", BaseQuantity: 180, BaseUnit: "g", AlternativeFood: strPtr("Salmon"), AlternativeQuantity: floatPtr(150), AlternativeUnit: strPtr("g")},
						{BaseFood: "Basmati rice", BaseQuantity: 80, BaseUnit: "g", AlternativeFood: strPtr("Sweet potato"), AlternativeQuantity: floatPtr(250), AlternativeUnit: strPtr("g")},
						{BaseFood: "Green beans", BaseQuantity: 200, BaseUnit: "g"},
					},
					Supplements: []seedSupplement{{Supplement: "Omega-3", Quantity: 2}},
				},
				{
					Name: "Dinner", Note: "2 hours before bedtime", Order: 2,
					Foods: []seedFood{
						{BaseFood: "Turkey", BaseQuantity: 150, BaseUnit: "g", AlternativeFood: strPtr("Tuna"), AlternativeQuantity: floatPtr(120), AlternativeUnit: strPtr("g")},
						{BaseFood: "Mixed salad", BaseQuantity: 200, BaseUnit: "g"},
					},
					Supplements: []seedSupplement{{Supplement: "Magnesium", Quantity: 1}},
				},
			},
		},
		{
			Name: "Rest Day", Order: 1,
			Meals: []seedMeal{
				{
					Name: "Breakfast", Order: 0,
					Foods: []seedFood{
						{BaseFood: "Greek yogurt", BaseQuantity: 200, BaseUnit: "g"},
						{BaseFood: "Walnuts", BaseQuantity: 30, BaseUnit: "g"},
					},
				},
				{
					Name: "Lunch", Order: 1,
					Foods: []seedFood{
						{BaseFood: "Salmon", BaseQuantity: 150, BaseUnit: "g"},
						{BaseFood: "Quinoa", BaseQuantity: 100, BaseUnit: "g"},
					},
					Supplements: []seedSupplement{{Supplement: "Vitamin D", Quantity: 1}},
				},
			},
		},
	}
}

func strPtr(s string) *string     { return &s }
func floatPtr(f float64) *float64 { return &f }
