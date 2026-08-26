// Package pdfgen is a small, generic PDF generation toolkit built on top of
// maroto v2 (https://github.com/johnfercher/maroto). It exposes a handful of
// domain-agnostic primitives (Title, SectionHeader/MealHeader, Table,
// Footer...) driven entirely by a GeneratorConfig, so it can be dropped into
// any application (nutrition-planner, SADE, etc.) without pulling in
// domain-specific concepts. Each application is expected to write its own
// thin translation layer that maps its own models onto these primitives.
//
// The default theme (colors + Noto Sans font, embedded - no external files
// needed) mirrors the NutriPlan app UI: cream page background, white cards,
// a muted sage-green accent, and light uppercase table headers with thin
// row dividers instead of solid-fill tables.
//
// Page background: maroto's WithBackgroundImage places an image whose
// vertical extent does NOT reliably reach the full page height - the exact
// rule isn't documented, and empirically its cutoff point shifts depending
// on how much content is on the page (shorter pages get a shorter cream
// fill, leaving a plain-white gap below the last row instead of reaching
// the bottom margin). Rather than depend on that undocumented behavior, we
// don't use a background image at all: every row emitted below carries an
// explicit cream Cell.BackgroundColor (the same technique already used, and
// already proven correct, for the DayBanner and SupplementsLine badges).
// The one remaining gap this can't cover is maroto's own internal "fill to
// bottom of page" row it adds automatically when finalizing a page (that
// row is added outside any hook we can style) - FillToBottom (called from
// Save, right before Generate) minimizes that by pre-filling the last
// page's remaining space with our own cream rows first.
package generator

import (
	"fmt"
	"nutrition-planner/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	mrtconfig "github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/fontrepository"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// Row heights (roughly mm) used by the primitives below - ClientHeader,
// DayBanner, MealHeader, Table, SupplementsLine. Exposed as constants (not
// buried as literals inside each method) so BlockHeight below can compute,
// ahead of time, exactly how tall an upcoming "keep together" block will be
// - the two must never drift apart, or EnsureFits would be guessing wrong.
const (
	DayBannerRowHeight   = 10
	MealTitleRowHeight   = 8
	MealNoteRowHeight    = 5
	TableHeaderRowHeight = 7
	TableRowHeight       = 9
	SupplementsRowHeight = 7
)

// BlockHeight computes the total row height a MealHeader + Table +
// SupplementsLine call sequence will occupy, given whether a note is
// present and how many table/supplement rows there will be. Pair with
// EnsureFits right before emitting such a block, so it starts on a fresh
// page instead of being split across a page boundary mid-table.
func BlockHeight(hasNote bool, tableRows, supplementItems int) float64 {
	h := float64(MealTitleRowHeight)
	if hasNote {
		h += MealNoteRowHeight
	}
	h += TableHeaderRowHeight + TableRowHeight*float64(tableRows)
	if supplementItems > 0 {
		h += SupplementsRowHeight
	}
	return h
}

// textInset is a small left/top padding (mm) applied to text elements that
// otherwise sit flush against the edge of their own background-colored
// row - gives the text a bit of breathing room instead of touching the
// cream fill's edge directly.
const textInset = 3.0

// Generator wraps a maroto document and exposes generic, theme-aware
// building blocks for building a PDF page by page.
type Generator struct {
	cfg   config.GeneratorConfig
	theme Theme
	m     core.Maroto
}

// bg is a shorthand used by every primitive below to paint its row(s) with
// the page's cream background - see the package doc comment for why this
// replaces a page-level background image.
func (g *Generator) bg() *props.Cell {
	return &props.Cell{BackgroundColor: g.theme.PageBackground}
}

// NewGenerator builds a Generator from a GeneratorConfig. Missing fields in
// cfg are filled in with defaults via GeneratorConfig.WithDefaults. If
// FontRegularTTF/FontBoldTTF are left empty, the embedded Noto Sans font is
// used - no external font files are required to get a working PDF with
// Romanian diacritics.
func NewGenerator(cfg config.GeneratorConfig) (*Generator, error) {
	theme, err := NewTheme(cfg)
	if err != nil {
		return nil, err
	}

	pSize, err := resolvePageSize(cfg.PageSize)
	if err != nil {
		return nil, err
	}

	regularBytes, err := loadFontBytes(cfg.FontRegularTTF, defaultFontRegular)
	if err != nil {
		return nil, fmt.Errorf("pdfgen: loading regular font: %w", err)
	}
	boldBytes, err := loadFontBytes(cfg.FontBoldTTF, defaultFontBold)
	if err != nil {
		return nil, fmt.Errorf("pdfgen: loading bold font: %w", err)
	}

	// maroto v2.4.1 no longer accepts []*entity.CustomFont literals -
	// entity.CustomFont became an interface (GetFamily/GetStyle/GetFile/
	// GetBytes), and its only constructors live behind pkg/fontrepository.
	// AddUTF8FontFromBytes is what we want here since our fonts are
	// go:embed'd, not read from a path.
	customFonts, err := fontrepository.New().
		AddUTF8FontFromBytes(theme.FontFamily, fontstyle.Normal, regularBytes).
		AddUTF8FontFromBytes(theme.FontFamily, fontstyle.Bold, boldBytes).
		Load()
	if err != nil {
		return nil, fmt.Errorf("pdfgen: loading custom fonts: %w", err)
	}

	builder := mrtconfig.NewBuilder().
		WithPageSize(pSize).
		WithTopMargin(float64(cfg.MarginTop)).
		WithBottomMargin(float64(cfg.MarginBottom)).
		WithLeftMargin(float64(cfg.MarginLeft)).
		WithRightMargin(float64(cfg.MarginRight)).
		WithCustomFonts(customFonts)

	m := maroto.New(builder.Build())

	return &Generator{cfg: cfg, theme: theme, m: m}, nil
}

// loadFontBytes reads a TTF from path if given, otherwise returns fallback
// (the embedded default font).
func loadFontBytes(path string, fallback []byte) ([]byte, error) {
	if path == "" {
		return fallback, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %w", path, err)
	}
	return b, nil
}

func resolvePageSize(size string) (pagesize.Type, error) {
	switch strings.ToLower(strings.TrimSpace(size)) {
	case "a1":
		return pagesize.A1, nil
	case "a2":
		return pagesize.A2, nil
	case "a3":
		return pagesize.A3, nil
	case "a4", "":
		return pagesize.A4, nil
	case "a5":
		return pagesize.A5, nil
	case "a6":
		return pagesize.A6, nil
	case "letter":
		return pagesize.Letter, nil
	case "legal":
		return pagesize.Legal, nil
	case "tabloid":
		return pagesize.Tabloid, nil
	default:
		return "", fmt.Errorf("pdfgen: unsupported page size %q", size)
	}
}

// ClientHeader renders a small uppercase eyebrow line above the title, e.g.
// "ANDREEA POPESCU · SLĂBIRE". Pass any number of parts; empty parts are
// skipped. If no parts are given, nothing is rendered.
func (g *Generator) ClientHeader(parts ...string) {
	nonEmpty := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			nonEmpty = append(nonEmpty, strings.ToUpper(p))
		}
	}
	if len(nonEmpty) == 0 {
		return
	}
	g.m.AddRow(7, text.NewCol(12, strings.Join(nonEmpty, "  ·  "), props.Text{
		Size:   9,
		Style:  fontstyle.Bold,
		Family: g.theme.FontFamily,
		Color:  g.theme.Primary,
		Left:   textInset,
		Top:    1,
	})).WithStyle(g.bg())
}

// Title renders the main document title (e.g. "Plan alimentar — August").
func (g *Generator) Title(t string) {
	g.m.AddRow(13, text.NewCol(12, t, props.Text{
		Size:   18,
		Style:  fontstyle.Bold,
		Family: g.theme.FontFamily,
		Color:  g.theme.Text,
		Left:   textInset,
		Top:    1,
	})).WithStyle(g.bg())
}

// KeyValueLine renders a single row of "label: value" pairs spread evenly
// across the page width (e.g. Greutate / Target / Inaltime).
func (g *Generator) KeyValueLine(pairs [][2]string) {
	if len(pairs) == 0 {
		return
	}
	width := 12 / len(pairs)
	if width == 0 {
		width = 1
	}

	cols := make([]core.Col, 0, len(pairs))
	for i, p := range pairs {
		left := 0.0
		if i == 0 {
			left = textInset
		}
		cols = append(cols, text.NewCol(width, fmt.Sprintf("%s: %s", p[0], p[1]), props.Text{
			Size:   10,
			Family: g.theme.FontFamily,
			Color:  g.theme.Muted,
			Left:   left,
			Top:    1,
		}))
	}
	g.m.AddRow(7, cols...).WithStyle(g.bg())
}

// DayBanner renders a solid, primary-colored full-width banner - used for a
// day-variant/section switch (e.g. "Zi de antrenament"), mirroring the
// active-tab look from the app UI.
func (g *Generator) DayBanner(label string) {
	g.m.AddRow(DayBannerRowHeight, text.NewCol(12, label, props.Text{
		Size:   11,
		Style:  fontstyle.Bold,
		Family: g.theme.FontFamily,
		Color:  g.theme.OnPrimary,
		Align:  align.Left,
		Left:   4,
		Top:    2.5,
	})).WithStyle(&props.Cell{BackgroundColor: g.theme.Primary})
}

// MealHeader renders a card-like block header: a bold dark title (e.g.
// "Mic dejun") with an optional smaller muted note beneath it
// (e.g. "În maxim o oră după trezire").
func (g *Generator) MealHeader(title, note string) {
	g.m.AddRow(MealTitleRowHeight, text.NewCol(12, title, props.Text{
		Size:   13,
		Style:  fontstyle.Bold,
		Family: g.theme.FontFamily,
		Color:  g.theme.Text,
		Left:   textInset,
		Top:    1,
	})).WithStyle(g.bg())
	if note = strings.TrimSpace(note); note != "" {
		g.m.AddRow(MealNoteRowHeight, text.NewCol(12, note, props.Text{
			Size:   9,
			Family: g.theme.FontFamily,
			Color:  g.theme.Muted,
			Left:   textInset,
		})).WithStyle(g.bg())
	}
}

// TableData is a generic table model: a header row plus N data rows, all
// treated as plain strings. pdfgen has no notion of what the columns mean.
// Empty cell values render as "—", matching the app's own placeholder for
// "no alternative".
type TableData struct {
	Headers []string
	Rows    [][]string
	// ColWidths optionally overrides the column widths as a 12-column grid
	// (must sum to <= 12 and match len(Headers)). If nil, columns are spread
	// evenly.
	ColWidths []int
}

// Table renders an uppercase, muted header row with a thin bottom divider,
// followed by data rows separated by thin light dividers - no solid fill,
// matching the app's own table styling.
func (g *Generator) Table(data TableData) {
	if len(data.Headers) == 0 {
		return
	}

	widths := data.ColWidths
	if len(widths) != len(data.Headers) {
		widths = evenWidths(len(data.Headers))
	}

	headerStyle := &props.Cell{
		BorderColor:     g.theme.TableHeader,
		BorderType:      border.Bottom,
		BorderThickness: 0.3,
		BackgroundColor: g.theme.PageBackground,
	}
	headerCols := make([]core.Col, 0, len(data.Headers))
	for i, h := range data.Headers {
		left := 0.0
		if i == 0 {
			left = textInset
		}
		headerCols = append(headerCols, text.NewCol(widths[i], strings.ToUpper(h), props.Text{
			Size:   8,
			Style:  fontstyle.Bold,
			Family: g.theme.FontFamily,
			Color:  g.theme.TableHeader,
			Top:    1,
			Bottom: 2,
			Left:   left,
		}))
	}
	g.m.AddRow(TableHeaderRowHeight, headerCols...).WithStyle(headerStyle)

	rowStyle := &props.Cell{
		BorderColor:     g.theme.Divider,
		BorderType:      border.Bottom,
		BorderThickness: 0.2,
		BackgroundColor: g.theme.PageBackground,
	}
	for _, r := range data.Rows {
		rowCols := make([]core.Col, 0, len(r))
		for i, cell := range r {
			w := 1
			if i < len(widths) {
				w = widths[i]
			}
			if strings.TrimSpace(cell) == "" {
				cell = "—"
			}
			left := 0.0
			if i == 0 {
				left = textInset
			}
			rowCols = append(rowCols, text.NewCol(w, cell, props.Text{
				Size:   9,
				Family: g.theme.FontFamily,
				Color:  g.theme.Text,
				Top:    2,
				Bottom: 2,
				Left:   left,
			}))
		}
		g.m.AddRow(TableRowHeight, rowCols...).WithStyle(rowStyle)
	}
}

func evenWidths(n int) []int {
	if n <= 0 {
		return nil
	}
	base := 12 / n
	rem := 12 % n
	widths := make([]int, n)
	for i := range widths {
		widths[i] = base
		if i < rem {
			widths[i]++
		}
	}
	return widths
}

// SupplementsLine renders a small uppercase muted "label" (e.g. "SUPLIMENTE")
// followed by a soft pill-style block listing items on a light primary
// background - mirroring the app's badge look (e.g. "Multivitamine × 1").
// If items is empty, nothing is rendered.
func (g *Generator) SupplementsLine(label string, items []string) {
	if len(items) == 0 {
		return
	}
	g.m.AddRow(SupplementsRowHeight,
		text.NewCol(3, strings.ToUpper(label), props.Text{
			Size:   8,
			Style:  fontstyle.Bold,
			Family: g.theme.FontFamily,
			Color:  g.theme.Muted,
			Left:   textInset,
			Top:    2,
		}).WithStyle(g.bg()),
		text.NewCol(9, strings.Join(items, "   ·   "), props.Text{
			Size:   9,
			Family: g.theme.FontFamily,
			Color:  g.theme.Primary,
			Left:   3,
			Top:    1.5,
			Bottom: 1.5,
		}).WithStyle(&props.Cell{BackgroundColor: g.theme.PrimaryLight}),
	)
}

// Note renders a small muted line of free-form text.
func (g *Generator) Note(t string) {
	if strings.TrimSpace(t) == "" {
		return
	}
	g.m.AddRow(5, text.NewCol(12, t, props.Text{
		Size:   9,
		Family: g.theme.FontFamily,
		Color:  g.theme.Muted,
		Left:   textInset,
	})).WithStyle(g.bg())
}

// Divider renders a thin full-width horizontal rule in the card-border
// color, used to visually separate blocks (e.g. between meals) the way a
// card boundary would in the app UI.
func (g *Generator) Divider() {
	g.m.AddRow(3, text.NewCol(12, "", props.Text{})).
		WithStyle(&props.Cell{
			BorderColor:     g.theme.CardBorder,
			BorderType:      border.Bottom,
			BorderThickness: 0.3,
			BackgroundColor: g.theme.PageBackground,
		})
}

// Footer renders a small muted closing line
// (e.g. supplement discount code, hydration/sleep recommendations).
func (g *Generator) Footer(t string) {
	if strings.TrimSpace(t) == "" {
		return
	}
	g.m.AddRow(7, text.NewCol(12, t, props.Text{
		Size:   8,
		Family: g.theme.FontFamily,
		Color:  g.theme.Muted,
		Left:   textInset,
		Top:    1,
	})).WithStyle(g.bg())
}

// Spacer adds vertical whitespace of the given height (in the same units as
// row heights, roughly mm). Carries the page's cream background, same as
// every other row - an unstyled spacer would show as a plain white gap.
func (g *Generator) Spacer(height float64) {
	g.m.AddRow(height).WithStyle(g.bg())
}

// EnsureFits forces a page break, if needed, before a "keep together" block
// of the given height (see BlockHeight) - so it never gets torn across a
// page boundary (e.g. a meal's title ending up on one page and its food
// table starting on the next). It works by nudging the cursor down with
// small cream-colored filler rows until maroto's own pagination
// (FitlnCurrentPage) kicks in and moves us to a fresh page; the trade-off
// is a little (cream, not white) blank space at the bottom of the previous
// page instead of a torn block. If height itself is taller than a full
// page, EnsureFits gives up after a bounded number of nudges and lets the
// block render (and split) as it normally would - there's no page it could
// ever fit on regardless.
func (g *Generator) EnsureFits(height float64) {
	const nudge = 4.0
	const maxNudges = 100
	for i := 0; i < maxNudges && !g.m.FitlnCurrentPage(height); i++ {
		g.Spacer(nudge)
	}
}

// FillToBottom pre-fills whatever space remains on the CURRENT page with
// our own cream-colored rows, stopping just short of overflowing onto a new
// page. Call this once, right before Save/Generate: maroto always adds one
// final, unstyled (plain white) row of its own to close out the last page
// - by minimizing how much space is left for that internal row to cover,
// this keeps that trailing white sliver down to a couple of millimeters
// instead of however much of the page happened to be empty.
func (g *Generator) FillToBottom() {
	const step = 2.0
	const maxSteps = 300
	for i := 0; i < maxSteps && g.m.FitlnCurrentPage(step); i++ {
		g.Spacer(step)
	}
}

// Save generates the PDF and writes it to <OutputDirectory>/<filename>,
// returning the full path written.
func (g *Generator) Save(filename string) (string, error) {
	g.FillToBottom()

	doc, err := g.m.Generate()
	if err != nil {
		return "", fmt.Errorf("pdfgen: generate: %w", err)
	}

	outPath := filepath.Join(g.cfg.OutputDirectory, filename)
	if err := doc.Save(outPath); err != nil {
		return "", fmt.Errorf("pdfgen: save: %w", err)
	}
	return outPath, nil
}
