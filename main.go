package main

import (
	"embed"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"kafui/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	// backend.Chdir2PrgPath()
	prgdir, err := backend.GetPrgDir()
	if err != nil {
		fmt.Println("GetPrgDir faild: ", err)
	} else {
		fmt.Println("program dir: ", prgdir)
	}
	if _, err := backend.InitLog("kafui.log", true); err != nil {
		fmt.Println("InitLog failed ", err)
	}
}

func main() {
	// Load config
	myconfig, err := backend.LoadConfig(backend.DEFAULT_CONFIG_FILE)
	if err != nil {
		log.Fatalf("LoadConfig failed: %s", err)
	}
	log.Infof("myconfig %s, %v", myconfig.Kafka[0].Name, myconfig.Kafka[0].Brokers)

	// Create an instance of the app structure
	app := NewApp()
	app.myconfig = myconfig
	zktool := &backend.ZkTool{}
	kafkatool := backend.NewKafkaTool(&myconfig.Kafka[0])
	app.kafkatool = kafkatool
	app.zktool = zktool

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "kafui - Kafka Client GUI",
		Width:  1024,
		Height: 768,
		// MinWidth:          720,
		// MinHeight:         570,
		// MaxWidth:          1280,
		// MaxHeight:         740,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
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
