package pdf_export

import (
	"nutrition-planner/config"
	"nutrition-planner/internals/app/settings"

	"github.com/ciprianiordache/pdfgen"
)

// PDFConfig maps the app's generator settings onto pdfgen. The app keeps
// its own GeneratorConfig (YAML/env keys, Settings-screen overrides);
// pdfgen only sees the resolved values. ColorCardBackground has no
// counterpart: no primitive fills a card. OutputDirectory is not pdfgen's
// concern either - the caller chooses the path it saves to.
func PDFConfig(cfg config.GeneratorConfig) pdfgen.Config {
	// pdfgen reads a zero margin as "use the default"; the app has always
	// meant zero, so pass a hair above it.
	mm := func(v int) float64 {
		if v == 0 {
			return 0.01
		}
		return float64(v)
	}
	return pdfgen.Config{
		PageSize:     cfg.PageSize,
		MarginTop:    mm(cfg.MarginTop),
		MarginBottom: mm(cfg.MarginBottom),
		MarginLeft:   mm(cfg.MarginLeft),
		MarginRight:  mm(cfg.MarginRight),
		Colors: pdfgen.Colors{
			PageBackground: cfg.ColorPageBackground,
			CardBorder:     cfg.ColorCardBorder,
			Primary:        cfg.ColorPrimary,
			PrimaryLight:   cfg.ColorPrimaryLight,
			OnPrimary:      cfg.ColorOnPrimary,
			Text:           cfg.ColorText,
			Muted:          cfg.ColorMuted,
			TableHeader:    cfg.ColorTableHeader,
			Divider:        cfg.ColorDivider,
		},
		FontFamily:      cfg.FontFamily,
		FontRegularPath: cfg.FontRegularTTF,
		FontBoldPath:    cfg.FontBoldTTF,
	}
}

// applySettingsTheme returns a copy of base with any non-empty color field
// from s overlaid on top - settings.Service.Get() always returns every
// field (empty string for anything the user never customized in the
// Settings screen), so "non-empty" is exactly "the user set this".
func applySettingsTheme(base config.GeneratorConfig, s settings.Response) config.GeneratorConfig {
	cfg := base
	overlay := func(dst *string, src string) {
		if src != "" {
			*dst = src
		}
	}
	overlay(&cfg.ColorPageBackground, s.ColorPageBackground)
	overlay(&cfg.ColorCardBackground, s.ColorCardBackground)
	overlay(&cfg.ColorCardBorder, s.ColorCardBorder)
	overlay(&cfg.ColorPrimary, s.ColorPrimary)
	overlay(&cfg.ColorPrimaryLight, s.ColorPrimaryLight)
	overlay(&cfg.ColorOnPrimary, s.ColorOnPrimary)
	overlay(&cfg.ColorText, s.ColorText)
	overlay(&cfg.ColorMuted, s.ColorMuted)
	overlay(&cfg.ColorTableHeader, s.ColorTableHeader)
	overlay(&cfg.ColorDivider, s.ColorDivider)
	return cfg
}
