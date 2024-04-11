package main

import (
	"context"
	"fmt"
	"kafui/backend"
)

// App struct
type App struct {
	ctx      context.Context
	myconfig *backend.Myconfig
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}

func (a *App) GetMyconfig() *backend.Myconfig {
	return a.myconfig
}

func (a *App) SetMyconfig(myconfig *backend.Myconfig) error {
	return backend.SaveConfig(myconfig, a.myconfig.Filename)
}
