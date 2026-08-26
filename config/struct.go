package config

import "time"

// Config holds the entire configuration for the application, including app settings, database connection, logger settings, and PDF generator settings.
type Config struct {
	App       AppConfig       `yaml:"app"`
	Database  DatabaseConfig  `yaml:"database"`
	Logger    LoggerConfig    `yaml:"logger"`
	Generator GeneratorConfig `yaml:"generator"`
	Export    ExportConfig    `yaml:"export"`
}

// ExportConfig holds configuration for where generated plan PDFs get
// written on the user's machine - kept separate from GeneratorConfig
// because it's about file placement, not PDF rendering (colors/fonts/
// margins).
type ExportConfig struct {
	// BaseDirectory is the root folder under which one subfolder per
	// client gets created (e.g. "<BaseDirectory>/Andreea Popescu/plan.pdf").
	// Leave empty (the normal case) to auto-detect
	// "<user's home>/Desktop/NutriPlan Exports" at startup - this is what
	// makes exports work correctly on any machine/user account without
	// hardcoding a username anywhere. Only set this explicitly if a
	// specific install needs a different location.
	BaseDirectory string `yaml:"base_directory" env:"EXPORT_BASE_DIRECTORY" default:""`
}

// AppConfig holds the configuration for the application, including name, version, environment, debug mode, and directories for data and temporary files.
type AppConfig struct {
	Name          string `yaml:"name" env:"APP_NAME" default:"nutrition-planner"`          // Name of the application (e.g., "nutrition-planner") - default: "nutrition-planner"
	Version       string `yaml:"version" env:"APP_VERSION" default:"v0.1.0"`               // Version of the application - default: "v0.1.0"
	Environment   string `yaml:"environment" env:"APP_ENVIRONMENT" default:"development"`  // Environment in which the application is running - default: "development"
	Debug         bool   `yaml:"debug" env:"APP_DEBUG" default:"true"`                     // Debug mode flag - default: true
	DataDirectory string `yaml:"data_directory" env:"APP_DATA_DIRECTORY" default:"./data"` // Directory for storing data files - default: "./data"
	TempDirectory string `yaml:"temp_directory" env:"APP_TEMP_DIRECTORY" default:"./temp"` // Directory for storing temporary files - default: "./temp"
}

// DatabaseConfig holds the configuration for the database connection, including driver, host, port, user credentials, and connection pooling settings.
type DatabaseConfig struct {
	Driver          string        `yaml:"driver" env:"DATABASE_DRIVER" default:"sqlite"` // Database driver (e.g., "postgres", "mysql", "sqlite3", etc.) - default: "sqlite"
	Host            string        `yaml:"host" env:"DATABASE_HOST"`
	Port            int           `yaml:"port" env:"DATABASE_PORT"`
	User            string        `yaml:"user" env:"DATABASE_USER"`
	Password        string        `yaml:"password" env:"DATABASE_PASSWORD" secret:"true"`
	Name            string        `yaml:"name" env:"DATABASE_NAME" default:"nutrition-planner"`
	SSLMode         string        `yaml:"sslmode" env:"DATABASE_SSLMODE" default:"disable"`                   // SSL mode for the database connection (e.g., "disable", "require", "verify-ca", "verify-full") - default: "disable"
	Timezone        string        `yaml:"timezone" env:"DATABASE_TIMEZONE" default:"UTC"`                     // Timezone for the database connection (e.g., "UTC", "America/New_York", etc.) - default: "UTC"
	MaxOpenConns    int           `yaml:"max_open_conns" env:"DATABASE_MAX_OPEN_CONNS" default:"20"`          // Maximum number of open connections to the database - default: 20
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"DATABASE_MAX_IDLE_CONNS" default:"10"`          // Maximum number of idle connections in the pool - default: 10
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DATABASE_CONN_MAX_LIFETIME" default:"1h"`    // Maximum amount of time a connection may be reused - default: 1 hour
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env:"DATABASE_CONN_MAX_IDLE_TIME" default:"30m"` // Maximum amount of time a connection may be idle before being closed - default: 30 minutes
}

// LoggerConfig holds the configuration for the logger, including log level, format, output destination, and file rotation settings.
type LoggerConfig struct {
	Level      string `yaml:"level" env:"LOGGER_LEVLE" default:"info"`           // Log level (e.g., "debug", "info", "warn", "error") - default: "info"
	Format     string `yaml:"format" env:"LOGGER_FORMAT" default:"json"`         // Log format (e.g., "json", "text") - default: "json"
	Output     string `yaml:"output" env:"LOGGER_OUTPUT" default:"both"`         // Output destination (e.g., "stdout", "file") - default: "both"
	FilePath   string `yaml:"file_path" env:"LOGGER_FILE_PATH" default:"./logs"` // File path for log output - default: "./logs"
	MaxSize    int    `yaml:"max_size" env:"LOGGER_MAX_SIZE" default:"50"`       // Maximum size of log file in MB - default: 50
	MaxAge     int    `yaml:"max_age" env:"LOGGER_MAX_AGE" default:"30"`         // Maximum age of log file in days - default: 30
	MaxBackups int    `yaml:"max_backups" env:"LOGGER_MAX_BACKUPS" default:"10"`
	Compress   bool   `yaml:"compress" env:"LOGGER_COMPRESS" default:"true"`
}

// GeneratorConfig holds the configuration for the PDF generator.
// Every field is meant to be tweakable via YAML/env without recompiling,
// so this struct stays free of any domain-specific concept (no "meal",
// "invoice", etc. - just generic layout/theme knobs). Defaults below mirror
// the NutriPlan app's own look (cream page, white cards, sage green accent),
// but every value can be overridden per app via config.
type GeneratorConfig struct {
	OutputDirectory string `yaml:"output_directory" env:"GENERATOR_OUTPUT_DIRECTORY" default:"./output"`
	PageSize        string `yaml:"page_size" env:"GENERATOR_PAGE_SIZE" default:"A4"`
	MarginTop       int    `yaml:"margin_top" env:"GENERATOR_MARGIN_TOP" default:"6"`
	MarginBottom    int    `yaml:"margin_bottom" env:"GENERATOR_MARGIN_BOTTOM" default:"6"`
	MarginLeft      int    `yaml:"margin_left" env:"GENERATOR_MARGIN_LEFT" default:"5"`
	MarginRight     int    `yaml:"margin_right" env:"GENERATOR_MARGIN_RIGHT" default:"5"`

	// ColorScheme (hex without #) for the generated PDF - kept separate so
	// they can be adjusted via config without recompiling. Defaults match
	// the NutriPlan app UI: cream page background, white cards, sage-green
	// accent, muted gray table headers/dividers.
	ColorPageBackground string `yaml:"color_page_background" env:"GENERATOR_COLOR_PAGE_BACKGROUND" default:"F7F5EF"`
	ColorCardBackground string `yaml:"color_card_background" env:"GENERATOR_COLOR_CARD_BACKGROUND" default:"FFFFFF"`
	ColorCardBorder     string `yaml:"color_card_border" env:"GENERATOR_COLOR_CARD_BORDER" default:"E7E4D9"`

	ColorPrimary      string `yaml:"color_primary" env:"GENERATOR_COLOR_PRIMARY" default:"2F5233"`
	ColorPrimaryLight string `yaml:"color_primary_light" env:"GENERATOR_COLOR_PRIMARY_LIGHT" default:"E3EDE1"`
	// ColorOnPrimary is the text color used on top of a solid ColorPrimary
	// fill (e.g. a day-variant banner).
	ColorOnPrimary string `yaml:"color_on_primary" env:"GENERATOR_COLOR_ON_PRIMARY" default:"FFFFFF"`

	ColorText        string `yaml:"color_text" env:"GENERATOR_COLOR_TEXT" default:"2B2B28"`
	ColorMuted       string `yaml:"color_muted" env:"GENERATOR_COLOR_MUTED" default:"6B6B66"`
	ColorTableHeader string `yaml:"color_table_header" env:"GENERATOR_COLOR_TABLE_HEADER" default:"8A8A82"`
	ColorDivider     string `yaml:"color_divider" env:"GENERATOR_COLOR_DIVIDER" default:"E5E2D8"`

	FontFamily string `yaml:"font_family" env:"GENERATOR_FONT_FAMILY" default:"NotoSans"`
	// FontRegularTTF / FontBoldTTF optionally point to a custom TTF on disk.
	// Leave empty to use the embedded default (Noto Sans), which already
	// covers Romanian diacritics - no external font files required.
	FontRegularTTF string `yaml:"font_regular_ttf" env:"GENERATOR_FONT_REGULAR_TTF" default:""`
	FontBoldTTF    string `yaml:"font_bold_ttf" env:"GENERATOR_FONT_BOLD_TTF" default:""`
}
