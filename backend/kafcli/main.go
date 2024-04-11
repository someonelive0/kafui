package main

import (
	"kafui/backend"
	"os"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"github.com/labstack/gommon/log"
)

var App = grumble.New(&grumble.Config{
	Name:                  "kafcli",
	Description:           "An CLI of Kafka",
	HistoryFile:           "kafcli.hist",
	Prompt:                "kafcli » ",
	PromptColor:           color.New(color.FgGreen, color.Bold),
	HelpHeadlineColor:     color.New(color.FgGreen),
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,

	Flags: func(f *grumble.Flags) {
		// f.String("f", "config", "kafui.toml", "set config file")
		f.String("a", "address", "127.0.0.1:9092", "set kafka host")
		f.String("u", "user", "DEFAULT", "set kafka password of SASL PLAIN")
		f.String("p", "password", "", "set kafka password of SASL PLAIN")
		f.Bool("v", "verbose", false, "enable verbose mode")
	},
})

var kafkatool *backend.KafkaTool = nil

func init() {
	App.SetPrintASCIILogo(func(a *grumble.App) {
		a.Println("  Kafui CLI ")
		a.Println()
	})

	quitCommand := &grumble.Command{
		Name:     "quit",
		Help:     "quit this cli",
		LongHelp: "quit this cli",
		Aliases:  []string{"q"},
		Run: func(c *grumble.Context) error {
			os.Exit(0)
			return nil
		},
	}
	App.AddCommand(quitCommand)
}

func main() {
	myconfig, err := backend.LoadConfig("kafui.toml")
	if err != nil {
		log.Fatalf("LoadConfig failed: %s", err)
	}
	log.Infof("myconfig brokers %V", myconfig.Kafka.Brokers)
	kafkatool = backend.NewKafkaTool(&myconfig.Kafka)

	grumble.Main(App)
}
