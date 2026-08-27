package pdf_export

import (
	"testing"

	"nutrition-planner/config"
	"nutrition-planner/internals/app/settings"
)

func TestApplySettingsTheme_AllEmptyKeepsBase(t *testing.T) {
	base := config.GeneratorConfig{ColorPrimary: "1B5E3C", ColorText: "222222"}

	got := applySettingsTheme(base, settings.Response{})

	if got != base {
		t.Errorf("expected an all-empty settings override to leave the config unchanged, got %+v", got)
	}
}

func TestApplySettingsTheme_OverridesOnlySetFields(t *testing.T) {
	base := config.GeneratorConfig{
		ColorPrimary: "1B5E3C",
		ColorText:    "222222",
		ColorMuted:   "6B6B6B",
	}

	got := applySettingsTheme(base, settings.Response{ColorPrimary: "FF0000"})

	if got.ColorPrimary != "FF0000" {
		t.Errorf("ColorPrimary = %q, want overridden %q", got.ColorPrimary, "FF0000")
	}
	if got.ColorText != "222222" {
		t.Errorf("ColorText = %q, want unchanged %q", got.ColorText, "222222")
	}
	if got.ColorMuted != "6B6B6B" {
		t.Errorf("ColorMuted = %q, want unchanged %q", got.ColorMuted, "6B6B6B")
	}
}

func TestApplySettingsTheme_OverridesEveryColorField(t *testing.T) {
	base := config.GeneratorConfig{}
	override := settings.Response{
		ColorPageBackground: "1",
		ColorCardBackground: "2",
		ColorCardBorder:     "3",
		ColorPrimary:        "4",
		ColorPrimaryLight:   "5",
		ColorOnPrimary:      "6",
		ColorText:           "7",
		ColorMuted:          "8",
		ColorTableHeader:    "9",
		ColorDivider:        "10",
	}

	got := applySettingsTheme(base, override)

	want := config.GeneratorConfig{
		ColorPageBackground: "1",
		ColorCardBackground: "2",
		ColorCardBorder:     "3",
		ColorPrimary:        "4",
		ColorPrimaryLight:   "5",
		ColorOnPrimary:      "6",
		ColorText:           "7",
		ColorMuted:          "8",
		ColorTableHeader:    "9",
		ColorDivider:        "10",
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestApplySettingsTheme_DoesNotTouchNonColorFields(t *testing.T) {
	base := config.GeneratorConfig{
		PageSize:        "A4",
		MarginTop:       15,
		OutputDirectory: "./output",
		FontFamily:      "NotoSans",
	}

	got := applySettingsTheme(base, settings.Response{ColorPrimary: "FF0000"})

	if got.PageSize != base.PageSize || got.MarginTop != base.MarginTop ||
		got.OutputDirectory != base.OutputDirectory || got.FontFamily != base.FontFamily {
		t.Errorf("expected non-color fields to stay untouched, got %+v", got)
	}
}
