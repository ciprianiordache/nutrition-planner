// cmd/demo exercisează librăria pdfgen direct (fără DB, fără
// export.Service) cu date shape-uite ca ecranele reale ale aplicației - util
// ca smoke test rapid când modifici tema/layout-ul PDF-ului.
//
// Rulare - oricare din variantele de mai jos funcționează, indiferent din ce
// director îl lansezi (căile config.yaml/.env se calculează din locația
// acestui fișier sursă, NU din directorul curent al terminalului):
//
//	go run ./cmd/demo          // din rădăcina proiectului
//	cd cmd/demo && go run .    // din interiorul folderului
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"nutrition-planner/config"
	"nutrition-planner/internals/app/pdf_export"

	"github.com/ciprianiordache/pdfgen"
)

// Food/Meal/DayVariant are minimal stand-ins for the domain model NutriPlan
// already has. This file exists purely to exercise pdfgen end-to-end
// against data shaped like the app's own screens.
type Food struct {
	BaseName string
	BaseQty  string
	AltName  string
	AltQty   string
}

type Meal struct {
	Name        string
	Note        string
	Foods       []Food
	Supplements []string
}

type DayVariant struct {
	Label string
	Meals []Meal
}

func main() {
	// Calculăm rădăcina proiectului din locația fișierului sursă (fixă la
	// compilare), NU din directorul de lucru curent al terminalului -
	// "go run ./cmd/demo" păstrează directorul curent al shell-ului
	// (oricare ar fi el), deci căi relative de tipul "../../config.yaml"
	// s-ar rezolva greșit dacă nu rulezi neapărat din interiorul cmd/demo.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("nu pot determina locația fișierului sursă (runtime.Caller a eșuat)")
	}
	projectRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")

	configPath := filepath.Join(projectRoot, "config.yaml")
	envPath := filepath.Join(projectRoot, ".env")
	fmt.Println("Se încarcă config din:", configPath)

	cfg, err := config.Load(configPath, envPath)
	if err != nil {
		panic(err)
	}

	// doc.Save() scrie fișierul cu os.WriteFile, care NU creează directoare
	// lipsă - spre deosebire de fluxul real din aplicație (export.Repo),
	// care apelează os.MkdirAll pentru folderul per-client. Acest demo
	// ocolește export.Repo, deci ne asigurăm noi că directorul există.
	absOutputDir, err := filepath.Abs(cfg.Generator.OutputDirectory)
	if err != nil {
		log.Fatalf("resolve output directory: %v", err)
	}
	fmt.Println("Se scrie în:", absOutputDir)

	if err := os.MkdirAll(absOutputDir, 0o755); err != nil {
		log.Fatalf("create output directory %q: %v", absOutputDir, err)
	}

	doc, err := pdfgen.New(pdf_export.PDFConfig(cfg.Generator))
	if err != nil {
		log.Fatalf("pdfgen: %v", err)
	}

	doc.Eyebrow("Andreea Popescu", "Slăbire")
	doc.Title("Plan alimentar — August")
	doc.Spacer(3)

	dayVariants := []DayVariant{
		{
			Label: "Zi de antrenament",
			Meals: []Meal{
				{
					Name: "Mic dejun",
					Note: "În maxim o oră după trezire",
					Foods: []Food{
						{"Ovăz", "60 g", "Orez brun", "50 g"},
						{"Ouă", "3 buc", "Brânză de vaci", "150 g"},
						{"Afine", "80 g", "", ""},
					},
					Supplements: []string{"Multivitamine × 1"},
				},
				{
					Name: "Prânz",
					Note: "La 3-4 ore după masa anterioară",
					Foods: []Food{
						{"Piept de pui", "180 g", "Somon", "150 g"},
						{"Orez basmati", "80 g", "Cartof dulce", "250 g"},
						{"Fasole verde", "200 g", "", ""},
					},
				},
				{
					Name: "Gustare",
					Note: "La 2-3 ore după prânz",
					Foods: []Food{
						{"Iaurt grecesc", "200 g", "Skyr", "200 g"},
						{"Nuci", "20 g", "Migdale", "20 g"},
					},
					Supplements: []string{"Omega 3 × 1"},
				},
			},
		},
		{
			Label: "Zi de repaus",
			Meals: []Meal{
				{
					Name: "Mic dejun",
					Note: "În maxim o oră după trezire",
					Foods: []Food{
						{"Ovăz", "50 g", "Orez brun", "40 g"},
						{"Ouă", "2 buc", "Brânză de vaci", "120 g"},
					},
					Supplements: []string{"Multivitamine × 1"},
				},
				{
					Name: "Prânz",
					Note: "La 3-4 ore după masa anterioară",
					Foods: []Food{
						{"Piept de pui", "150 g", "Somon", "130 g"},
						{"Cartof dulce", "200 g", "", ""},
					},
				},
			},
		},
	}

	for _, dv := range dayVariants {
		if len(dv.Meals) > 0 {
			first := dv.Meals[0]
			doc.EnsureFits(pdfgen.BannerHeight + pdfgen.BlockHeight(first.Note != "", len(first.Foods), len(first.Supplements) > 0))
		} else {
			doc.EnsureFits(pdfgen.BannerHeight)
		}
		doc.Banner(dv.Label)
		doc.Spacer(2)

		for i, meal := range dv.Meals {
			doc.EnsureFits(pdfgen.BlockHeight(meal.Note != "", len(meal.Foods), len(meal.Supplements) > 0))

			doc.Heading(meal.Name, meal.Note)

			rows := make([][]string, 0, len(meal.Foods))
			for _, f := range meal.Foods {
				rows = append(rows, []string{f.BaseName, f.BaseQty, f.AltName, f.AltQty})
			}

			doc.Table(pdfgen.TableData{
				Headers:   []string{"Aliment", "Cantitate", "Alternativă", "Cantitate"},
				Rows:      rows,
				ColWidths: []int{4, 2, 4, 2},
			})

			doc.Badges("Suplimente", meal.Supplements)

			if i < len(dv.Meals)-1 {
				doc.Spacer(3)
				doc.Divider()
				doc.Spacer(3)
			}
		}
		doc.Spacer(6)
	}

	outPath := filepath.Join(absOutputDir, "plan_alimentar_demo.pdf")
	if err := doc.Save(outPath); err != nil {
		log.Fatalf("save: %v", err)
	}
	fmt.Println("PDF generat cu succes:", outPath)
}
