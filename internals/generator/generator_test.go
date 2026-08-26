package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	"nutrition-planner/config"
	"nutrition-planner/internals/generator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validConfig mirrors the defaults declared via struct tags on
// config.GeneratorConfig (normally filled in by the app's config loader at
// startup). NewGenerator expects a fully-populated config - it no longer
// applies its own defaults - so every test builds one via this helper and
// only overrides what it's actually testing.
func validConfig(t *testing.T) config.GeneratorConfig {
	t.Helper()
	return config.GeneratorConfig{
		OutputDirectory:     t.TempDir(),
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

// assertValidPDF does a minimal sanity check that path points at a
// non-empty, well-formed PDF file, without depending on any external tool
// (poppler, pdftotext, etc.) that may not be installed wherever these tests
// run.
func assertValidPDF(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0), "generated PDF should not be empty")

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(data), 5)
	assert.Equal(t, "%PDF-", string(data[:5]), "file should start with the PDF magic bytes")
	assert.Contains(t, string(data[len(data)-1024:]), "%%EOF", "well-formed PDFs end with an %%EOF marker")
}

func TestNewGenerator_ValidConfig(t *testing.T) {
	gen, err := generator.NewGenerator(validConfig(t))
	require.NoError(t, err)
	require.NotNil(t, gen)
}

func TestNewGenerator_InvalidColor(t *testing.T) {
	cfg := validConfig(t)
	cfg.ColorPrimary = "not-a-hex-color"

	_, err := generator.NewGenerator(cfg)
	assert.Error(t, err)
}

func TestNewGenerator_InvalidPageSize(t *testing.T) {
	cfg := validConfig(t)
	cfg.PageSize = "A0"

	_, err := generator.NewGenerator(cfg)
	assert.Error(t, err)
}

func TestNewGenerator_MissingCustomFontFile(t *testing.T) {
	cfg := validConfig(t)
	cfg.FontRegularTTF = "/path/does/not/exist.ttf"

	_, err := generator.NewGenerator(cfg)
	assert.Error(t, err)
}

func TestNewGenerator_UsesEmbeddedFontByDefault(t *testing.T) {
	// FontRegularTTF/FontBoldTTF are left empty on purpose - this is the
	// zero-config path every real caller in the app takes today, and it
	// must keep working without any font files on disk.
	cfg := validConfig(t)
	cfg.FontRegularTTF = ""
	cfg.FontBoldTTF = ""

	gen, err := generator.NewGenerator(cfg)
	require.NoError(t, err)

	path, err := gen.Save("embedded-font.pdf")
	require.NoError(t, err)
	assertValidPDF(t, path)
}

func TestGenerator_Save_FullDocument(t *testing.T) {
	cfg := validConfig(t)
	gen, err := generator.NewGenerator(cfg)
	require.NoError(t, err)

	// Exercise every public building block in one realistic pass, including
	// diacritics (ăâîșț) and the "no alternative" case (empty alt columns,
	// which Table renders as "—") and the "no supplements" case (should be
	// a silent no-op, not an empty badge).
	gen.ClientHeader("Andreea Popescu", "Slăbire")
	gen.Title("Plan alimentar — August")
	gen.KeyValueLine([][2]string{
		{"Greutate", "62 kg"},
		{"Înălțime", "168 cm"},
	})
	gen.Spacer(3)

	gen.DayBanner("Zi de antrenament")
	gen.Spacer(2)

	gen.MealHeader("Mic dejun", "În maxim o oră după trezire")
	gen.Table(generator.TableData{
		Headers: []string{"Aliment", "Cantitate", "Alternativă", "Cantitate"},
		Rows: [][]string{
			{"Ovăz", "60 g", "Orez brun", "50 g"},
			{"Afine", "80 g", "", ""}, // no alternative -> should render as "—"
		},
		ColWidths: []int{4, 2, 4, 2},
	})
	gen.SupplementsLine("Suplimente", []string{"Multivitamine × 1"})

	gen.Spacer(3)
	gen.Divider()
	gen.Spacer(3)

	gen.MealHeader("Prânz", "")
	gen.Table(generator.TableData{
		Headers: []string{"Aliment", "Cantitate"},
		Rows: [][]string{
			{"Piept de pui", "180 g"},
		},
	})
	// No supplements for this meal - must not add an empty badge row.
	gen.SupplementsLine("Suplimente", nil)

	gen.Note("O notă liberă.")
	gen.Footer("Recomandare suplimente: VITABOLIC, cod MihailVTB")

	path, err := gen.Save("plan_test.pdf")
	require.NoError(t, err)
	assertValidPDF(t, path)
	assert.Equal(t, filepath.Join(cfg.OutputDirectory, "plan_test.pdf"), path)
}

func TestGenerator_Table_EmptyHeadersIsNoop(t *testing.T) {
	gen, err := generator.NewGenerator(validConfig(t))
	require.NoError(t, err)

	// Should not panic and should not add any content that breaks
	// generation.
	gen.Table(generator.TableData{})

	_, err = gen.Save("empty-table.pdf")
	require.NoError(t, err)
}

func TestGenerator_ClientHeader_AllPartsEmptyIsNoop(t *testing.T) {
	gen, err := generator.NewGenerator(validConfig(t))
	require.NoError(t, err)

	gen.ClientHeader("", "  ", "")
	gen.Title("Should still render") // document must not be empty

	_, err = gen.Save("empty-header.pdf")
	require.NoError(t, err)
}

func TestGenerator_Save_CreatesFileAtExpectedPath(t *testing.T) {
	cfg := validConfig(t)
	gen, err := generator.NewGenerator(cfg)
	require.NoError(t, err)

	gen.Title("Minimal document")

	path, err := gen.Save("output.pdf")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(cfg.OutputDirectory, "output.pdf"), path)

	_, statErr := os.Stat(path)
	require.NoError(t, statErr)
}

func TestGenerator_Save_FailsWhenOutputDirectoryDoesNotExist(t *testing.T) {
	cfg := validConfig(t)
	// Point at a directory that was never created - Save() writes with
	// os.WriteFile, which does not create missing parent directories.
	cfg.OutputDirectory = filepath.Join(cfg.OutputDirectory, "does-not-exist")

	gen, err := generator.NewGenerator(cfg)
	require.NoError(t, err)

	gen.Title("Doesn't matter")

	_, err = gen.Save("output.pdf")
	assert.Error(t, err)
}

func TestGenerator_MultiPageDocument(t *testing.T) {
	// Enough content to force at least one automatic page break, mirroring
	// a real multi-day plan - regression guard for anything that only
	// works on a single page.
	cfg := validConfig(t)
	gen, err := generator.NewGenerator(cfg)
	require.NoError(t, err)

	gen.Title("Plan alimentar — testare pagini multiple")
	for day := 1; day <= 6; day++ {
		gen.DayBanner("Ziua de test")
		for meal := 1; meal <= 5; meal++ {
			gen.MealHeader("Masă", "Notă descriptivă pentru masă")
			gen.Table(generator.TableData{
				Headers: []string{"Aliment", "Cantitate", "Alternativă", "Cantitate"},
				Rows: [][]string{
					{"Ovăz", "60 g", "Orez brun", "50 g"},
					{"Ouă", "3 buc", "Brânză de vaci", "150 g"},
					{"Afine", "80 g", "", ""},
				},
			})
			gen.SupplementsLine("Suplimente", []string{"Multivitamine × 1"})
			gen.Spacer(3)
		}
	}

	path, err := gen.Save("multi-page.pdf")
	require.NoError(t, err)
	assertValidPDF(t, path)

	info, err := os.Stat(path)
	require.NoError(t, err)
	// A single-page A4 with this little content would be a few KB; a
	// multi-page document with 30 meal tables should be meaningfully
	// larger. Loose bound - just guards against silent truncation.
	assert.Greater(t, info.Size(), int64(5000))
}
