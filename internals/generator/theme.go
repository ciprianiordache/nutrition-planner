package generator

import (
	"nutrition-planner/config"

	"github.com/johnfercher/maroto/v2/pkg/props"
)

// Theme holds the resolved (parsed) styling values derived from a
// GeneratorConfig. It's built once per Generator and reused for every
// component so config parsing/validation happens in a single place.
type Theme struct {
	PageBackground *props.Color
	CardBackground *props.Color
	CardBorder     *props.Color

	Primary      *props.Color
	PrimaryLight *props.Color
	OnPrimary    *props.Color

	Text        *props.Color
	Muted       *props.Color
	TableHeader *props.Color
	Divider     *props.Color

	FontFamily string
}

// NewTheme resolves a GeneratorConfig's hex colors into a Theme.
func NewTheme(cfg config.GeneratorConfig) (Theme, error) {
	colors := map[string]string{
		"page background": cfg.ColorPageBackground,
		"card background": cfg.ColorCardBackground,
		"card border":     cfg.ColorCardBorder,
		"primary":         cfg.ColorPrimary,
		"primary light":   cfg.ColorPrimaryLight,
		"on primary":      cfg.ColorOnPrimary,
		"text":            cfg.ColorText,
		"muted":           cfg.ColorMuted,
		"table header":    cfg.ColorTableHeader,
		"divider":         cfg.ColorDivider,
	}
	parsed := make(map[string]*props.Color, len(colors))
	for name, hex := range colors {
		c, err := hexToColor(hex)
		if err != nil {
			return Theme{}, err
		}
		parsed[name] = c
	}

	return Theme{
		PageBackground: parsed["page background"],
		CardBackground: parsed["card background"],
		CardBorder:     parsed["card border"],
		Primary:        parsed["primary"],
		PrimaryLight:   parsed["primary light"],
		OnPrimary:      parsed["on primary"],
		Text:           parsed["text"],
		Muted:          parsed["muted"],
		TableHeader:    parsed["table header"],
		Divider:        parsed["divider"],
		FontFamily:     cfg.FontFamily,
	}, nil
}
