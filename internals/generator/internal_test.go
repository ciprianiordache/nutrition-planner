package generator

import (
	"testing"

	"nutrition-planner/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- hexToColor -------------------------------------------------------

func TestHexToColor_Valid(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		wantR   int
		wantG   int
		wantB   int
	}{
		{"lowercase no hash", "1b5e3c", 0x1b, 0x5e, 0x3c},
		{"uppercase with hash", "#1B5E3C", 0x1b, 0x5e, 0x3c},
		{"black", "000000", 0, 0, 0},
		{"white", "FFFFFF", 255, 255, 255},
		{"padded with whitespace", "  2F5233  ", 0x2f, 0x52, 0x33},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := hexToColor(tt.hex)
			require.NoError(t, err)
			require.NotNil(t, c)
			assert.Equal(t, tt.wantR, c.Red)
			assert.Equal(t, tt.wantG, c.Green)
			assert.Equal(t, tt.wantB, c.Blue)
		})
	}
}

func TestHexToColor_Invalid(t *testing.T) {
	tests := []struct {
		name string
		hex  string
	}{
		{"empty", ""},
		{"too short", "1B5E3"},
		{"too long", "1B5E3C1"},
		{"non-hex characters", "GGHHII"},
		{"just a hash", "#"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := hexToColor(tt.hex)
			assert.Error(t, err)
		})
	}
}

// --- NewTheme -----------------------------------------------------------

func validGeneratorConfig() config.GeneratorConfig {
	// Mirrors the defaults declared via struct tags in config.GeneratorConfig
	// (normally filled in by the app's config loader) - generator.NewTheme/
	// NewGenerator expect a fully-populated config, since defaulting now
	// happens at the config-loading layer, not inside this package.
	return config.GeneratorConfig{
		OutputDirectory:     "./output",
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

func TestNewTheme_Valid(t *testing.T) {
	cfg := validGeneratorConfig()

	theme, err := NewTheme(cfg)
	require.NoError(t, err)

	require.NotNil(t, theme.PageBackground)
	require.NotNil(t, theme.CardBackground)
	require.NotNil(t, theme.CardBorder)
	require.NotNil(t, theme.Primary)
	require.NotNil(t, theme.PrimaryLight)
	require.NotNil(t, theme.OnPrimary)
	require.NotNil(t, theme.Text)
	require.NotNil(t, theme.Muted)
	require.NotNil(t, theme.TableHeader)
	require.NotNil(t, theme.Divider)

	assert.Equal(t, "NotoSans", theme.FontFamily)
	// Spot-check one resolved color end to end.
	assert.Equal(t, 0x2f, theme.Primary.Red)
	assert.Equal(t, 0x52, theme.Primary.Green)
	assert.Equal(t, 0x33, theme.Primary.Blue)
}

func TestNewTheme_InvalidColor(t *testing.T) {
	cfg := validGeneratorConfig()
	cfg.ColorPrimary = "not-a-color"

	_, err := NewTheme(cfg)
	assert.Error(t, err)
}

// --- resolvePageSize ------------------------------------------------------

func TestResolvePageSize_Valid(t *testing.T) {
	// Every page size pdfgen claims to support should resolve without error,
	// case-insensitively and regardless of surrounding whitespace.
	sizes := []string{"a1", "A2", " a3 ", "A4", "a5", "A6", "Letter", "LEGAL", "tabloid", ""}

	for _, s := range sizes {
		t.Run(s, func(t *testing.T) {
			_, err := resolvePageSize(s)
			assert.NoError(t, err)
		})
	}
}

func TestResolvePageSize_Invalid(t *testing.T) {
	_, err := resolvePageSize("A0")
	assert.Error(t, err)
}

// --- evenWidths -------------------------------------------------------

func TestEvenWidths(t *testing.T) {
	tests := []struct {
		n    int
		want []int
	}{
		{0, nil},
		{-1, nil},
		{1, []int{12}},
		{2, []int{6, 6}},
		{3, []int{4, 4, 4}},
		{4, []int{3, 3, 3, 3}},
		// 12 doesn't divide evenly by 5 - the remainder should be spread
		// across the first columns so the total still sums to 12.
		{5, []int{3, 3, 2, 2, 2}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := evenWidths(tt.n)
			assert.Equal(t, tt.want, got)

			if tt.n > 0 {
				sum := 0
				for _, w := range got {
					sum += w
				}
				assert.Equal(t, 12, sum, "column widths must always sum to 12")
			}
		})
	}
}

// --- loadFontBytes ------------------------------------------------------

func TestLoadFontBytes_EmptyPathUsesFallback(t *testing.T) {
	fallback := []byte("fallback-bytes")

	got, err := loadFontBytes("", fallback)
	require.NoError(t, err)
	assert.Equal(t, fallback, got)
}

func TestLoadFontBytes_MissingFileErrors(t *testing.T) {
	_, err := loadFontBytes("/path/does/not/exist.ttf", []byte("fallback"))
	assert.Error(t, err)
}
