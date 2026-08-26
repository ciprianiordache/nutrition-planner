package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "nutrition-planner",
		Width:     1024,
		Height:    768,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		// OpenInspectorOnStartup mirrors config.yaml's app.debug. Note this
		// only has any effect in a debug build (`wails dev`, or
		// `wails build -debug`) - a regular `wails build` strips the
		// DevTools inspector out entirely regardless of this flag.
		Debug: options.Debug{
			OpenInspectorOnStartup: app.debugMode,
		},
		Bind: []interface{}{
			app.Client,
			app.Plan,
			app.DayVariant,
			app.Meal,
			app.Food,
			app.Supplement,
			app.Export,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
