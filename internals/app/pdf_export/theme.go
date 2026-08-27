package pdf_export

import (
	"nutrition-planner/config"
	"nutrition-planner/internals/app/settings"
)

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
