package pdf_export

import (
	"fmt"
	"math"
	"strconv"

	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/plan"
	"nutrition-planner/internals/app/supplement"
	"nutrition-planner/internals/generator"
)

// renderPlan walks a fully-loaded plan (day variants -> meals -> foods /
// supplements) and a client and issues the corresponding pdfgen calls onto
// gen. This is the one place that knows about both the domain model and
// the generic pdfgen primitives - kept as a pure function (no I/O) so it's
// easy to unit test against fixture data.
func renderPlan(gen *generator.Generator, p plan.FullResponse, c client.Response) {
	gen.ClientHeader(c.Name, c.Target)
	gen.Title(p.Title)

	var kv [][2]string
	if c.Weight > 0 {
		kv = append(kv, [2]string{"Greutate", fmt.Sprintf("%s kg", formatQty(c.Weight))})
	}
	if c.Height > 0 {
		kv = append(kv, [2]string{"Înălțime", fmt.Sprintf("%s cm", formatQty(c.Height))})
	}
	if len(kv) > 0 {
		gen.KeyValueLine(kv)
	}
	gen.Spacer(3)

	for _, dv := range p.DayVariants {
		// Ține bannerul de zi lipit de prima lui masă - altfel bannerul
		// poate rămâne singur la finalul unei pagini, cu prima masă
		// începând abia pe pagina următoare.
		if len(dv.Meals) > 0 {
			first := dv.Meals[0]
			gen.EnsureFits(generator.DayBannerRowHeight + generator.BlockHeight(first.Note != "", len(first.Foods), len(first.Supplements)))
		} else {
			gen.EnsureFits(generator.DayBannerRowHeight)
		}
		gen.DayBanner(dv.Name)
		gen.Spacer(2)

		for i, m := range dv.Meals {
			// Ține fiecare masă (titlu + tabel + suplimente) ca un bloc unic -
			// altfel maroto poate rupe tabelul la mijloc când nu mai are loc.
			gen.EnsureFits(generator.BlockHeight(m.Note != "", len(m.Foods), len(m.Supplements)))

			gen.MealHeader(m.Name, m.Note)
			gen.Table(generator.TableData{
				Headers:   []string{"Aliment", "Cantitate", "Alternativă", "Cantitate"},
				Rows:      foodRows(m.Foods),
				ColWidths: []int{4, 2, 4, 2},
			})
			gen.SupplementsLine("Suplimente", supplementItems(m.Supplements))

			if i < len(dv.Meals)-1 {
				gen.Spacer(3)
				gen.Divider()
				gen.Spacer(3)
			}
		}
		gen.Spacer(6)
	}

	if p.GeneralRecommendations != "" {
		gen.Footer(p.GeneralRecommendations)
	}
	if p.Footer != "" {
		gen.Footer(p.Footer)
	}
}

func foodRows(foods []food.Response) [][]string {
	rows := make([][]string, len(foods))
	for i, f := range foods {
		alt, altQty := "", ""
		if f.AlternativeFood != nil {
			alt = *f.AlternativeFood
		}
		if f.AlternativeQuantity != nil {
			unit := ""
			if f.AlternativeUnit != nil {
				unit = *f.AlternativeUnit
			}
			altQty = fmt.Sprintf("%s %s", formatQty(*f.AlternativeQuantity), unit)
		}
		rows[i] = []string{
			f.BaseFood,
			fmt.Sprintf("%s %s", formatQty(f.BaseQuantity), f.BaseUnit),
			alt,
			altQty,
		}
	}
	return rows
}

func supplementItems(supplements []supplement.Response) []string {
	if len(supplements) == 0 {
		return nil
	}
	items := make([]string, len(supplements))
	for i, s := range supplements {
		items[i] = fmt.Sprintf("%s × %s", s.Supplement, formatQty(s.Quantity))
	}
	return items
}

// formatQty renders a float without a trailing ".0" for whole numbers
// (e.g. "60" instead of "60.000000"), while still showing one decimal for
// genuinely fractional quantities (e.g. "1.5").
func formatQty(f float64) string {
	if f == math.Trunc(f) {
		return strconv.FormatFloat(f, 'f', 0, 64)
	}
	return strconv.FormatFloat(f, 'f', 1, 64)
}
