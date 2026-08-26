// cmd/seed populează baza de date configurată (config.yaml/.env, la fel ca
// aplicația Wails) cu 4 clienți, fiecare cu 2 sau 3 planuri, iar fiecare
// plan primește la rândul lui variante de zi, mese, alimente și suplimente
// — ca să poți verifica vizual tot fluxul din UI.
//
// Rulare din rădăcina proiectului:
//
//	go run ./cmd/seed
//
// E idempotent: dacă baza de date are deja clienți, nu inserează nimic
// (ca să nu dublezi datele la fiecare rulare). Pentru a forța re-seedarea:
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
	// PlanTitles — cât de multe titluri, atâtea planuri primește clientul (2 sau 3).
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
	force := flag.Bool("force", false, "seedează chiar dacă baza de date are deja clienți")
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

	// Wiring identic cu app.go — repo-urile "copil" înainte de cele "părinte".
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
		fmt.Println("⚠️  Baza de date are deja clienți — nu adaug nimic.")
		fmt.Println("    Rulează cu -force dacă vrei să adaugi oricum: go run ./cmd/seed -force")
		return
	}

	seedClients := []seedClient{
		{
			Name: "Andreea Popescu", Weight: 62, Height: 168, Target: "Slăbire",
			PlanTitles: []string{"Plan alimentar — Ianuarie", "Plan alimentar — Februarie"},
		},
		{
			Name: "Mihai Ionescu", Weight: 88, Height: 182, Target: "Masă musculară",
			PlanTitles: []string{"Plan alimentar — Bulk faza 1", "Plan alimentar — Bulk faza 2", "Plan alimentar — Menținere"},
		},
		{
			Name: "Elena Radu", Weight: 71, Height: 165, Target: "Menținere",
			PlanTitles: []string{"Plan alimentar — Vară", "Plan alimentar — Toamnă"},
		},
		{
			Name: "Vlad Marinescu", Weight: 95, Height: 190, Target: "Recompoziție",
			PlanTitles: []string{"Plan alimentar — Faza 1", "Plan alimentar — Faza 2", "Plan alimentar — Faza 3"},
		},
	}

	const generalRecommendations = "3 litri de apă zilnic, evită prăjelile, minim 7 ore de somn."
	const footer = "Discount Code: VITABOLIC, cod MihailVTB - 10% reducere + transport gratuit"

	for _, sc := range seedClients {
		c := &client.Client{Name: sc.Name, Weight: sc.Weight, Height: sc.Height, Target: sc.Target}
		clientID, err := clientRepo.Create(c)
		if err != nil {
			log.Fatalf("❌ create client %q: %v", sc.Name, err)
		}
		fmt.Printf("✔ Client creat: %s (%s) — %d planuri\n", sc.Name, clientID, len(sc.PlanTitles))

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
			fmt.Printf("  ↳ Plan creat: %s (%s)\n", title, planID)

			if err := seedPlanContent(planID, dayVariantRepo, mealRepo, foodRepo, supplementRepo); err != nil {
				log.Fatalf("❌ seed content for plan %q: %v", title, err)
			}
			fmt.Printf("    ↳ Conținut adăugat: %d variante de zi\n", len(demoDayVariants()))
		}
	}

	fmt.Println("✅ Seed complet: 4 clienți, cu planuri complete (zile, mese, alimente, suplimente).")
}

// seedPlanContent adaugă, pentru un plan dat, aceleași două variante de zi
// demonstrative (Zi de antrenament / Zi de repaus), fiecare cu mese,
// alimente și suplimente — suficient ca să vezi tot fluxul din UI.
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
			Name: "Zi de antrenament", Order: 0,
			Meals: []seedMeal{
				{
					Name: "Mic dejun", Note: "În maxim o oră după trezire", Order: 0,
					Foods: []seedFood{
						{BaseFood: "Ovăz", BaseQuantity: 60, BaseUnit: "g", AlternativeFood: strPtr("Orez brun"), AlternativeQuantity: floatPtr(50), AlternativeUnit: strPtr("g")},
						{BaseFood: "Ouă", BaseQuantity: 3, BaseUnit: "buc", AlternativeFood: strPtr("Brânză de vaci"), AlternativeQuantity: floatPtr(150), AlternativeUnit: strPtr("g")},
						{BaseFood: "Afine", BaseQuantity: 80, BaseUnit: "g"},
					},
					Supplements: []seedSupplement{{Supplement: "Multivitamine", Quantity: 1}},
				},
				{
					Name: "Prânz", Note: "La 3–4 ore după masa anterioară", Order: 1,
					Foods: []seedFood{
						{BaseFood: "Piept de pui", BaseQuantity: 180, BaseUnit: "g", AlternativeFood: strPtr("Somon"), AlternativeQuantity: floatPtr(150), AlternativeUnit: strPtr("g")},
						{BaseFood: "Orez basmati", BaseQuantity: 80, BaseUnit: "g", AlternativeFood: strPtr("Cartof dulce"), AlternativeQuantity: floatPtr(250), AlternativeUnit: strPtr("g")},
						{BaseFood: "Fasole verde", BaseQuantity: 200, BaseUnit: "g"},
					},
					Supplements: []seedSupplement{{Supplement: "Omega-3", Quantity: 2}},
				},
				{
					Name: "Cină", Note: "Cu 2 ore înainte de somn", Order: 2,
					Foods: []seedFood{
						{BaseFood: "Curcan", BaseQuantity: 150, BaseUnit: "g", AlternativeFood: strPtr("Ton"), AlternativeQuantity: floatPtr(120), AlternativeUnit: strPtr("g")},
						{BaseFood: "Salată mixtă", BaseQuantity: 200, BaseUnit: "g"},
					},
					Supplements: []seedSupplement{{Supplement: "Magneziu", Quantity: 1}},
				},
			},
		},
		{
			Name: "Zi de repaus", Order: 1,
			Meals: []seedMeal{
				{
					Name: "Mic dejun", Order: 0,
					Foods: []seedFood{
						{BaseFood: "Iaurt grecesc", BaseQuantity: 200, BaseUnit: "g"},
						{BaseFood: "Nuci", BaseQuantity: 30, BaseUnit: "g"},
					},
				},
				{
					Name: "Prânz", Order: 1,
					Foods: []seedFood{
						{BaseFood: "Somon", BaseQuantity: 150, BaseUnit: "g"},
						{BaseFood: "Quinoa", BaseQuantity: 100, BaseUnit: "g"},
					},
					Supplements: []seedSupplement{{Supplement: "Vitamina D", Quantity: 1}},
				},
			},
		},
	}
}

func strPtr(s string) *string     { return &s }
func floatPtr(f float64) *float64 { return &f }
