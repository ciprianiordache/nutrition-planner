package pdf_export

import (
	"os"
	"path/filepath"
	"testing"

	"nutrition-planner/config"
	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/day_variant"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/meal"
	"nutrition-planner/internals/app/plan"
	"nutrition-planner/internals/app/supplement"

	"github.com/ciprianiordache/pdfgen"
)

// validGeneratorConfig mirrors the defaults declared via struct tags on
// config.GeneratorConfig (normally filled in by the app's config loader at
// startup).
func validGeneratorConfig(outputDir string) config.GeneratorConfig {
	return config.GeneratorConfig{
		OutputDirectory:     outputDir,
		PageSize:            "A4",
		MarginTop:           15,
		MarginBottom:        15,
		MarginLeft:          14,
		MarginRight:         14,
		ColorPageBackground: "F7F5EF",
		ColorCardBackground: "FFFFFF",
		ColorCardBorder:     "E7E4D9",
		ColorPrimary:        "2F5233",
		ColorPrimaryLight:   "E3EDE1",
		ColorOnPrimary:      "FFFFFF",
		ColorText:           "2B2B28",
		ColorMuted:          "6B6B66",
		ColorTableHeader:    "8A8A82",
		ColorDivider:        "E5E2D8",
		FontFamily:          "NotoSans",
	}
}

func strPtr(s string) *string   { return &s }
func f64Ptr(f float64) *float64 { return &f }

// fixturePlan builds a realistic FullResponse - two day variants, several
// meals, foods with and without an alternative, and a mix of meals with
// and without supplements - covering every branch renderPlan takes.
func fixturePlan() (plan.FullResponse, client.Response) {
	c := client.Response{
		ID:     "client-1",
		Name:   "Andreea Popescu",
		Weight: 62,
		Height: 168,
		Target: "Slăbire",
	}

	p := plan.FullResponse{
		Response: plan.Response{
			ID:                     "plan-1",
			ClientID:               c.ID,
			Title:                  "Plan alimentar — August",
			GeneralRecommendations: "3-5 l apă, 7-9h somn.",
			Footer:                 "Recomandare suplimente: VITABOLIC, cod MihailVTB",
		},
		DayVariants: []day_variant.Response{
			{
				ID:   "dv-1",
				Name: "Zi de antrenament",
				Meals: []meal.Response{
					{
						ID:   "meal-1",
						Name: "Mic dejun",
						Note: "În maxim o oră după trezire",
						Foods: []food.Response{
							{BaseFood: "Ovăz", BaseQuantity: 60, BaseUnit: "g", AlternativeFood: strPtr("Orez brun"), AlternativeQuantity: f64Ptr(50), AlternativeUnit: strPtr("g")},
							{BaseFood: "Afine", BaseQuantity: 80, BaseUnit: "g"}, // no alternative
						},
						Supplements: []supplement.Response{
							{Supplement: "Multivitamine", Quantity: 1},
						},
					},
					{
						ID:   "meal-2",
						Name: "Prânz",
						Note: "La 3-4 ore după masa anterioară",
						Foods: []food.Response{
							{BaseFood: "Piept de pui", BaseQuantity: 180, BaseUnit: "g"},
						},
						// no supplements
					},
				},
			},
			{
				ID:   "dv-2",
				Name: "Zi de repaus",
				Meals: []meal.Response{
					{
						ID:   "meal-3",
						Name: "Mic dejun",
						Foods: []food.Response{
							{BaseFood: "Ovăz", BaseQuantity: 50, BaseUnit: "g"},
						},
					},
				},
			},
		},
	}

	return p, c
}

func TestRenderPlan_ProducesValidPDF(t *testing.T) {
	dir := t.TempDir()
	doc, err := pdfgen.New(PDFConfig(validGeneratorConfig(dir)))
	if err != nil {
		t.Fatalf("pdfgen.New: %v", err)
	}

	p, c := fixturePlan()
	renderPlan(doc, p, c, "ro")

	path := filepath.Join(dir, "render_test.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("output is not a PDF (missing %%PDF- header)")
	}
}

func TestRenderPlan_EnglishLabelsProduceValidPDF(t *testing.T) {
	dir := t.TempDir()
	doc, err := pdfgen.New(PDFConfig(validGeneratorConfig(dir)))
	if err != nil {
		t.Fatalf("pdfgen.New: %v", err)
	}

	p, c := fixturePlan()
	renderPlan(doc, p, c, "en")

	path := filepath.Join(dir, "render_test_en.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("output is not a PDF (missing %%PDF- header)")
	}
}

func TestFoodRows(t *testing.T) {
	foods := []food.Response{
		{BaseFood: "Ovăz", BaseQuantity: 60, BaseUnit: "g", AlternativeFood: strPtr("Orez brun"), AlternativeQuantity: f64Ptr(50), AlternativeUnit: strPtr("g")},
		{BaseFood: "Afine", BaseQuantity: 80, BaseUnit: "g"},
	}

	rows := foodRows(foods)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}

	if rows[0][0] != "Ovăz" || rows[0][1] != "60 g" || rows[0][2] != "Orez brun" || rows[0][3] != "50 g" {
		t.Errorf("row 0 = %v", rows[0])
	}
	// No alternative -> both alternative columns must be empty strings
	// (Table itself is responsible for turning that into "—").
	if rows[1][0] != "Afine" || rows[1][1] != "80 g" || rows[1][2] != "" || rows[1][3] != "" {
		t.Errorf("row 1 = %v", rows[1])
	}
}

func TestSupplementItems(t *testing.T) {
	if got := supplementItems(nil); got != nil {
		t.Errorf("expected nil for no supplements, got %v", got)
	}

	items := supplementItems([]supplement.Response{
		{Supplement: "Multivitamine", Quantity: 1},
		{Supplement: "Omega 3", Quantity: 2},
	})
	want := []string{"Multivitamine × 1", "Omega 3 × 2"}
	if len(items) != len(want) {
		t.Fatalf("got %v, want %v", items, want)
	}
	for i := range want {
		if items[i] != want[i] {
			t.Errorf("items[%d] = %q, want %q", i, items[i], want[i])
		}
	}
}

func TestFormatQty(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{60, "60"},
		{0, "0"},
		{1.5, "1.5"},
		{200, "200"},
		{0.5, "0.5"},
	}
	for _, tt := range tests {
		if got := formatQty(tt.in); got != tt.want {
			t.Errorf("formatQty(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
