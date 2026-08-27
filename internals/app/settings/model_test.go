package settings

import "testing"

func strp(s string) *string { return &s }

func TestApplyUpdate_NilFieldsLeaveCurrentValueUnchanged(t *testing.T) {
	s := Settings{ColorPrimary: "1B5E3C", UITheme: "dark"}

	applyUpdate(&s, UpdateRequest{}) // every field nil

	if s.ColorPrimary != "1B5E3C" {
		t.Errorf("ColorPrimary changed unexpectedly: %q", s.ColorPrimary)
	}
	if s.UITheme != "dark" {
		t.Errorf("UITheme changed unexpectedly: %q", s.UITheme)
	}
}

func TestApplyUpdate_SetFieldsOverwrite(t *testing.T) {
	s := Settings{ColorPrimary: "1B5E3C", ExportBaseDirectory: "/old/path"}

	applyUpdate(&s, UpdateRequest{
		ColorPrimary:        strp("FF0000"),
		ExportBaseDirectory: strp("/new/path"),
	})

	if s.ColorPrimary != "FF0000" {
		t.Errorf("ColorPrimary = %q, want %q", s.ColorPrimary, "FF0000")
	}
	if s.ExportBaseDirectory != "/new/path" {
		t.Errorf("ExportBaseDirectory = %q, want %q", s.ExportBaseDirectory, "/new/path")
	}
}

func TestApplyUpdate_EmptyStringIsAValidExplicitValue(t *testing.T) {
	// A pointer to "" is different from a nil pointer: the caller is
	// explicitly asking to CLEAR the override (fall back to config.yaml's
	// default again), not "leave it as-is".
	s := Settings{ColorPrimary: "1B5E3C"}

	applyUpdate(&s, UpdateRequest{ColorPrimary: strp("")})

	if s.ColorPrimary != "" {
		t.Errorf("expected ColorPrimary to be cleared to empty, got %q", s.ColorPrimary)
	}
}

func TestApplyUpdate_NeverTouchesID(t *testing.T) {
	s := Settings{ID: singletonID}

	applyUpdate(&s, UpdateRequest{ColorPrimary: strp("FF0000")})

	if s.ID != singletonID {
		t.Errorf("ID changed unexpectedly: %q", s.ID)
	}
}

func TestToResponse_MapsEveryField(t *testing.T) {
	s := Settings{
		TrainerName:         "Andreea Popescu",
		TrainerRole:         "Nutriționist",
		ExportBaseDirectory: "/export",
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
		UITheme:             "dark",
	}

	got := toResponse(s)

	want := Response{
		TrainerName:         "Andreea Popescu",
		TrainerRole:         "Nutriționist",
		ExportBaseDirectory: "/export",
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
		UITheme:             "dark",
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestApplyUpdate_TrainerNameAndRole(t *testing.T) {
	s := Settings{}

	applyUpdate(&s, UpdateRequest{
		TrainerName: strp("Andreea Popescu"),
		TrainerRole: strp("Nutriționist"),
	})

	if s.TrainerName != "Andreea Popescu" {
		t.Errorf("TrainerName = %q, want %q", s.TrainerName, "Andreea Popescu")
	}
	if s.TrainerRole != "Nutriționist" {
		t.Errorf("TrainerRole = %q, want %q", s.TrainerRole, "Nutriționist")
	}
}
