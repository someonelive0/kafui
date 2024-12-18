package main

import (
	"context"
	"fmt"
	"kafui/backend"
)

// App struct
type App struct {
	ctx       context.Context
	myconfig  *backend.Myconfig
	kafkatool *backend.KafkaTool
	zktool    *backend.ZkTool
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
	if a.kafkatool != nil {
		a.kafkatool.Close()
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}

func (a *App) GetMyconfig() *backend.Myconfig {
	myconfig, err := backend.LoadConfig(backend.DEFAULT_CONFIG_FILE)
	if err != nil {
		return nil
	}

	// reset kafkatool with new myconfig
	a.myconfig = myconfig
	a.kafkatool.Init(&myconfig.Kafka)

	return a.myconfig
}

func (a *App) SetMyconfig(myconfig *backend.Myconfig) error {
	// log.Infof("SetMyconfig %#v", *myconfig)
	if len(myconfig.Kafka.Password) == 0 { // 如果Password为空，自动保存password的原有值
		myconfig.Kafka.Password = a.myconfig.Kafka.Password
	}
	return backend.SaveConfig(myconfig, a.myconfig.Filename)
}

func (a *App) TestKafka(kafkaConfig *backend.KafkaConfig) (*backend.Broker, error) {
	return backend.TestKafa(kafkaConfig)
}
