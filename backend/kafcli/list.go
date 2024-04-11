package main

import (
	"fmt"
	"os"

	"github.com/desertbit/grumble"
	"github.com/olekukonko/tablewriter"

	"kafui/backend"
)

func init() {
	listCommand := &grumble.Command{
		Name:     "list",
		Help:     "list tools",
		LongHelp: "list some objects",
		Aliases:  []string{"l"},
	}
	App.AddCommand(listCommand)

	listCommand.AddCommand(&grumble.Command{
		Name: "brokers",
		Help: "kafka brokers",
		Run: func(c *grumble.Context) error {

			brokers, err := kafkatool.ListBrokers()
			if err != nil {
				fmt.Println("ListBrokers", err)
				return err
			}
			fmt.Printf("Found brokers %d\n", len(brokers))

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(backend.BrokerHeader())
			for i := range brokers {
				table.Append(brokers[i].ToStrings())
			}
			table.Render()

			return nil
		},
	})

	listCommand.AddCommand(&grumble.Command{
		Name: "topics",
		Help: "kafka topics",
		Run: func(c *grumble.Context) error {

			topics, err := kafkatool.ListTopics()
			if err != nil {
				fmt.Println("ListTopics", err)
				return err
			}
			fmt.Printf("Found topics %d\n", len(topics))

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Topic"})
			for i := range topics {
				table.Append([]string{topics[i]})
			}
			table.Render()

			return nil
		},
	})

	listCommand.AddCommand(&grumble.Command{
		Name: "groups",
		Help: "kafka consumer grouos",
		Run: func(c *grumble.Context) error {

			groups, err := kafkatool.ListGroups()
			if err != nil {
				fmt.Println("ListGroups", err)
				return err
			}
			fmt.Printf("Found groups %d\n", len(groups))

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Group"})
			for i := range groups {
				table.Append([]string{groups[i]})
			}
			table.Render()

			return nil
		},
	})
}
