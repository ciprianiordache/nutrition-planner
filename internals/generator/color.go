package generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/johnfercher/maroto/v2/pkg/props"
)

// hexToColor parses a "RRGGBB" (with or without leading '#') hex string into
// a maroto props.Color.
func hexToColor(hex string) (*props.Color, error) {
	h := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(h) != 6 {
		return nil, fmt.Errorf("pdfgen: invalid hex color %q, expected 6 hex digits", hex)
	}

	r, err := strconv.ParseUint(h[0:2], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("pdfgen: invalid hex color %q: %w", hex, err)
	}
	g, err := strconv.ParseUint(h[2:4], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("pdfgen: invalid hex color %q: %w", hex, err)
	}
	b, err := strconv.ParseUint(h[4:6], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("pdfgen: invalid hex color %q: %w", hex, err)
	}

	return &props.Color{Red: int(r), Green: int(g), Blue: int(b)}, nil
}
