package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"kafui/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	// backend.Chdir2PrgPath() // will cause "AppData\Local\Temp\wails.json: The system cannot find the file specified"
	prgdir, err := backend.GetPrgDir()
	if err != nil {
		fmt.Println("GetPrgDir faild: ", err)
	} else {
		fmt.Println("program dir: ", prgdir)
	}
}

func main() {

	// Create an instance of the app structure
	app := NewApp()
	kafkatool := &app.kafkatool
	zktool := &app.zktool

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "kafui - Kafka Client GUI",
		Width:  1024,
		Height: 768,
		// MinWidth:          720,
		// MinHeight:         570,
		// MaxWidth:          1280,
		// MaxHeight:         740,
		DisableResize:      false,
		Fullscreen:         false,
		Frameless:          false,
		StartHidden:        false,
		HideWindowOnClose:  false,
		BackgroundColour:   &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:             assets,
		Logger:             logger.NewFileLogger("kafui.log"),
		LogLevel:           logger.INFO,
		LogLevelProduction: logger.INFO,
		OnStartup:          app.startup,
		OnDomReady:         app.domReady,
		OnShutdown:         app.shutdown,
		Bind: []interface{}{
			app,
			kafkatool,
			zktool,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
