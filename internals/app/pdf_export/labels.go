package pdf_export

// pdfLabels holds every hardcoded string that ends up inside a generated
// PDF (table headers, "Weight"/"Height" line, "Supplements" badge label) -
// kept in one place so adding a language only means adding one map entry
// here, not hunting through render.go.
type pdfLabels struct {
	Weight      string
	Height      string
	Food        string
	Quantity    string
	Alternative string
	Supplements string
}

var pdfLabelsByLanguage = map[string]pdfLabels{
	"ro": {
		Weight:      "Greutate",
		Height:      "Înălțime",
		Food:        "Aliment",
		Quantity:    "Cantitate",
		Alternative: "Alternativă",
		Supplements: "Suplimente",
	},
	"en": {
		Weight:      "Weight",
		Height:      "Height",
		Food:        "Food",
		Quantity:    "Quantity",
		Alternative: "Alternative",
		Supplements: "Supplements",
	},
}

// labelsFor returns the label set for language, falling back to Romanian
// for an empty or unrecognized value - the app's default has always been
// Romanian, so an unset settings.Language must keep behaving exactly like
// it did before this field existed.
func labelsFor(language string) pdfLabels {
	if l, ok := pdfLabelsByLanguage[language]; ok {
		return l
	}
	return pdfLabelsByLanguage["ro"]
}
