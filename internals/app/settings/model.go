package settings

import "time"

// Settings is a singleton row - the app has exactly one settings record,
// always addressed by singletonID (see repo.go). Every field is optional:
// an empty string means "fall back to config.yaml's default", so adding a
// new overridable field never requires a migration of existing data.
type Settings struct {
	// ID carries the "uuid" tag purely to satisfy schema-builder's
	// validation, which requires every primary_key field to declare
	// "auto" or "uuid" (whether the DB or crud-depot generates the value)
	// - it does NOT mean this row gets a random ID. crud-depot's own
	// Create() only auto-generates a UUID when the field is still its zero
	// value at insert time; repo.go always sets ID to the fixed
	// singletonID *before* calling Create, so the "uuid" generation path
	// never actually triggers - the row's ID stays singletonID exactly as
	// intended. See repo_test.go for a regression test pinning this down.
	ID string `db:"id,primary_key,uuid"`

	// ExportBaseDirectory overrides config.yaml's export.base_directory.
	// Empty means "use the configured default (or the auto-detected
	// Desktop folder if that's empty too)".
	ExportBaseDirectory string `db:"export_base_directory"`

	// PDF theme - each mirrors a field on config.GeneratorConfig. Empty
	// means "use the config.yaml default for this color".
	ColorPageBackground string `db:"color_page_background"`
	ColorCardBackground string `db:"color_card_background"`
	ColorCardBorder     string `db:"color_card_border"`
	ColorPrimary        string `db:"color_primary"`
	ColorPrimaryLight   string `db:"color_primary_light"`
	ColorOnPrimary      string `db:"color_on_primary"`
	ColorText           string `db:"color_text"`
	ColorMuted          string `db:"color_muted"`
	ColorTableHeader    string `db:"color_table_header"`
	ColorDivider        string `db:"color_divider"`

	// UITheme is "light", "dark", or "system" (default). Not yet consumed
	// by the frontend - stored now so the settings screen can offer the
	// toggle without a later schema change.
	UITheme string `db:"ui_theme"`

	// Language is "ro" or "en" (empty defaults to "ro" in the frontend and
	// in pdf_export's label lookup). Drives both the UI's own text and the
	// labels baked into generated PDFs (table headers, "Greutate",
	// "Suplimente", etc.) - a single setting rather than two, since a
	// nutritionist demoing the app in English would also want an English
	// sample PDF to show.
	Language string `db:"language"`

	CreatedAt time.Time `db:"created_at,oncreate"`
	UpdatedAt time.Time `db:"updated_at,onwrite"`
}

func (Settings) TableName() string {
	return "settings"
}

// UpdateRequest — the data received from the frontend when saving settings.
// Every field is optional (a pointer): only fields the user actually
// touched in the form get sent, and omitted ones leave the stored value
// unchanged rather than being reset to empty.
type UpdateRequest struct {
	ExportBaseDirectory *string `json:"exportBaseDirectory,omitempty"`

	ColorPageBackground *string `json:"colorPageBackground,omitempty"`
	ColorCardBackground *string `json:"colorCardBackground,omitempty"`
	ColorCardBorder     *string `json:"colorCardBorder,omitempty"`
	ColorPrimary        *string `json:"colorPrimary,omitempty"`
	ColorPrimaryLight   *string `json:"colorPrimaryLight,omitempty"`
	ColorOnPrimary      *string `json:"colorOnPrimary,omitempty"`
	ColorText           *string `json:"colorText,omitempty"`
	ColorMuted          *string `json:"colorMuted,omitempty"`
	ColorTableHeader    *string `json:"colorTableHeader,omitempty"`
	ColorDivider        *string `json:"colorDivider,omitempty"`

	UITheme  *string `json:"uiTheme,omitempty"`
	Language *string `json:"language,omitempty"`
}

// Response — the data sent to the frontend.
type Response struct {
	ExportBaseDirectory string `json:"exportBaseDirectory"`

	ColorPageBackground string `json:"colorPageBackground"`
	ColorCardBackground string `json:"colorCardBackground"`
	ColorCardBorder     string `json:"colorCardBorder"`
	ColorPrimary        string `json:"colorPrimary"`
	ColorPrimaryLight   string `json:"colorPrimaryLight"`
	ColorOnPrimary      string `json:"colorOnPrimary"`
	ColorText           string `json:"colorText"`
	ColorMuted          string `json:"colorMuted"`
	ColorTableHeader    string `json:"colorTableHeader"`
	ColorDivider        string `json:"colorDivider"`

	UITheme  string `json:"uiTheme"`
	Language string `json:"language"`
}

func toResponse(s Settings) Response {
	return Response{
		ExportBaseDirectory: s.ExportBaseDirectory,
		ColorPageBackground: s.ColorPageBackground,
		ColorCardBackground: s.ColorCardBackground,
		ColorCardBorder:     s.ColorCardBorder,
		ColorPrimary:        s.ColorPrimary,
		ColorPrimaryLight:   s.ColorPrimaryLight,
		ColorOnPrimary:      s.ColorOnPrimary,
		ColorText:           s.ColorText,
		ColorMuted:          s.ColorMuted,
		ColorTableHeader:    s.ColorTableHeader,
		ColorDivider:        s.ColorDivider,
		UITheme:             s.UITheme,
		Language:            s.Language,
	}
}

// applyUpdate mutates s in place with every non-nil field from req -
// fields left nil in req keep their current stored value.
func applyUpdate(s *Settings, req UpdateRequest) {
	set := func(dst *string, src *string) {
		if src != nil {
			*dst = *src
		}
	}
	set(&s.ExportBaseDirectory, req.ExportBaseDirectory)
	set(&s.ColorPageBackground, req.ColorPageBackground)
	set(&s.ColorCardBackground, req.ColorCardBackground)
	set(&s.ColorCardBorder, req.ColorCardBorder)
	set(&s.ColorPrimary, req.ColorPrimary)
	set(&s.ColorPrimaryLight, req.ColorPrimaryLight)
	set(&s.ColorOnPrimary, req.ColorOnPrimary)
	set(&s.ColorText, req.ColorText)
	set(&s.ColorMuted, req.ColorMuted)
	set(&s.ColorTableHeader, req.ColorTableHeader)
	set(&s.ColorDivider, req.ColorDivider)
	set(&s.UITheme, req.UITheme)
	set(&s.Language, req.Language)
}
