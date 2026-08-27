package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"nutrition-planner/config"
	"nutrition-planner/internals/app/client"
	"nutrition-planner/internals/app/day_variant"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/meal"
	"nutrition-planner/internals/app/pdf_export"
	"nutrition-planner/internals/app/plan"
	"nutrition-planner/internals/app/settings"
	"nutrition-planner/internals/app/supplement"
	"nutrition-planner/internals/database"
	"nutrition-planner/internals/logger"

	"github.com/ciprianiordache/schema-builder"
	_ "modernc.org/sqlite"
)

// App struct — links the Wails context to the database connection and
// exposes the handlers for each domain so they can be bound separately
// to the frontend. Each handler becomes its own namespace in
// frontend/wailsjs/go (e.g., client.Handler, plan.Handler, etc.).
type App struct {
	ctx    context.Context
	db     *database.Database
	log    *slog.Logger
	logger io.Closer
	// debugMode mirrors cfg.App.Debug - read by main.go when building the
	// Wails options, to decide whether to auto-open the DevTools inspector.
	debugMode bool

	Client     *client.Handler
	Plan       *plan.Handler
	DayVariant *day_variant.Handler
	Meal       *meal.Handler
	Food       *food.Handler
	Supplement *supplement.Handler
	Export     *pdf_export.Handler
	Settings   *settings.Handler
}

// NewApp creates a new App application struct. It loads the configuration,
// starts the structured logger, opens the database connection, ensures the
// schema exists, and then builds the repo -> service -> handler chain for
// each domain.
func NewApp() *App {
	// Bootstrap: until cfg.Logger is loaded we can't build the structured
	// logger, so a config.Load() failure still uses the stdlib "log"
	// package - this is the only place in the app where that's the case.
	cfg, err := config.Load("config.yaml", ".env")
	if err != nil {
		log.Fatalf("❌ Unable to load config: %v", err)
	}

	// App.Debug forces verbose logging regardless of what logger.level
	// says in config.yaml - handy when tracking down a problem without
	// having to go edit the config file first, then remembering to
	// revert it afterward.
	if cfg.App.Debug {
		cfg.Logger.Level = "debug"
	}

	lg, logCloser, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("❌ Unable to initialize logger: %v", err)
	}
	// Every log line also carries which environment produced it - useful
	// once you're comparing log files from more than one install (e.g. a
	// "testing" copy you experiment in vs. the "production" one with real
	// client data).
	lg = lg.With("environment", cfg.App.Environment)
	slog.SetDefault(lg)

	// Outside of "production", isolate the database and export folder by
	// environment so testing/experimenting never touches (or mixes with)
	// real client data: nutrition-planner.db becomes
	// nutrition-planner-testing.db, "NutriPlan Exports" becomes
	// "NutriPlan Exports (testing)", etc. An empty Environment behaves
	// like "production" (no suffix) - the field is optional in config.yaml.
	envSuffix := ""
	if cfg.App.Environment != "" && cfg.App.Environment != "production" {
		envSuffix = fmt.Sprintf(" (%s)", cfg.App.Environment)
		cfg.Database.Name = fmt.Sprintf("%s-%s", cfg.Database.Name, cfg.App.Environment)
	}

	db := database.New(cfg.Database, cfg.App.DataDirectory)
	if err := db.Connect(); err != nil {
		lg.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}

	if err := ensureSchema(db); err != nil {
		lg.Error("unable to create schema", "error", err)
		os.Exit(1)
	}

	// "Child" repos are built first — "parent" repos
	// need them to populate relationships upon reading (e.g., meal.Repo
	// needs food.Repo and supplement.Repo to populate Meal.Foods /
	// Meal.Supplements).
	foodRepo := food.NewRepo(db)
	supplementRepo := supplement.NewRepo(db)
	mealRepo := meal.NewRepo(db, foodRepo, supplementRepo)
	dayVariantRepo := day_variant.NewRepo(db, mealRepo)
	planRepo := plan.NewRepo(db, dayVariantRepo)
	clientRepo := client.NewRepo(db)
	settingsRepo := settings.NewRepo(db)

	foodSvc := food.NewService(foodRepo, lg)
	supplementSvc := supplement.NewService(supplementRepo, lg)
	mealSvc := meal.NewService(mealRepo, lg)
	dayVariantSvc := day_variant.NewService(dayVariantRepo, lg)
	planSvc := plan.NewService(planRepo, lg)
	clientSvc := client.NewService(clientRepo, lg)
	settingsSvc := settings.NewService(settingsRepo, lg)

	// Exporting depends on the plan + client (it requires the plan's
	// full tree and the client's name for folders/headers), not on
	// the database repository itself—its repository only manages
	// the per-client export directory on disk.
	//
	// BaseDirectory is optional in config.yaml (export.base_directory);
	// if missing, it defaults to "<user home>/Desktop/NutriPlan Exports"
	// so that the application installation works on another computer
	// or account without relying on a hardcoded username.
	exportBaseDir := cfg.Export.BaseDirectory
	if exportBaseDir == "" {
		var err error
		exportBaseDir, err = pdf_export.DefaultBaseDirectory()
		if err != nil {
			lg.Error("unable to resolve default export directory", "error", err)
			os.Exit(1)
		}
	}
	exportBaseDir += envSuffix
	lg.Info("export configured", "base_directory", exportBaseDir)
	exportRepo := pdf_export.NewRepo(exportBaseDir)
	exportSvc := pdf_export.NewService(planSvc, clientSvc, settingsSvc, exportRepo, cfg.Generator, cfg.App.TempDirectory, lg)

	return &App{
		db:         db,
		log:        lg,
		logger:     logCloser,
		debugMode:  cfg.App.Debug,
		Client:     client.NewHandler(clientSvc),
		Plan:       plan.NewHandler(planSvc),
		DayVariant: day_variant.NewHandler(dayVariantSvc),
		Meal:       meal.NewHandler(mealSvc),
		Food:       food.NewHandler(foodSvc),
		Supplement: supplement.NewHandler(supplementSvc),
		Export:     pdf_export.NewHandler(exportSvc),
		Settings:   settings.NewHandler(settingsSvc, db, lg),
	}
}

// ensureSchema creates the missing tables in the order dictated by foreign keys:
// clients -> plans -> day_variants -> meals -> foods/supplements. settings
// has no foreign keys (a standalone singleton row), so its position in the
// list doesn't matter.
// It is idempotent (CREATE TABLE IF NOT EXISTS), so it is safe to run
// every time the application starts.
func ensureSchema(db *database.Database) error {
	builder := schema.New(db, schema.SQLite{})
	return builder.CreateSchema(
		client.Client{},
		plan.Plan{},
		day_variant.DayVariant{},
		meal.Meal{},
		food.Food{},
		supplement.Supplement{},
		settings.Settings{},
	)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods - Settings.Handler needs it too, for
// its native folder-picker dialog (Wails only exposes that on the Go
// side, and only once a window context exists).
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	settings.WithContext(a.Settings, ctx)
}

// shutdown is called when the app terminates — close the
// database connection and flush/close the log file (if any).
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.log.Error("error closing database", "error", err)
		}
	}
	if a.logger != nil {
		if err := a.logger.Close(); err != nil {
			a.log.Error("error closing logger", "error", err)
		}
	}
}
