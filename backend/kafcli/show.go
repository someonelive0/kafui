package main

import (
	"fmt"
	"kafui/backend"
	"os"

	"github.com/desertbit/grumble"
	"github.com/olekukonko/tablewriter"
)

func init() {
	showCommand := &grumble.Command{
		Name:     "show",
		Help:     "show tools",
		LongHelp: "show config or meta of object",
	}
	App.AddCommand(showCommand)

	// show topic meta|config
	showTopicCommand := &grumble.Command{
		Name:     "topic",
		Help:     "show topic tools",
		LongHelp: "show config or meta of topic",
	}
	showCommand.AddCommand(showTopicCommand)

	showTopicCommand.AddCommand(&grumble.Command{
		Name: "meta",
		Help: "show meta of topic, eg: show topic meta {TOPIC_NAME}",
		Args: func(a *grumble.Args) {
			a.String("topic", "which topic to show", grumble.Default(""))
		},
		Run: func(c *grumble.Context) error {
			if len(c.Args.String("topic")) == 0 {
				return fmt.Errorf("topic name is empty")
			}
			c.App.Println("show topic meta:", c.Args.String("topic"))

			metas, err := kafkatool.GetTopicMeta(c.Args.String("topic"))
			if err != nil {
				fmt.Println("GetTopicMeta", err)
				return err
			}
			c.App.Printf("metas: %#v\n", metas)

			return nil
		},
	})

	showTopicCommand.AddCommand(&grumble.Command{
		Name: "config",
		Help: "show config of topic, eg: show topic config {TOPIC_NAME}",
		Args: func(a *grumble.Args) {
			a.String("topic", "which topic to show", grumble.Default(""))
		},
		Run: func(c *grumble.Context) error {
			if len(c.Args.String("topic")) == 0 {
				return fmt.Errorf("topic name is empty")
			}
			c.App.Println("show topic config:", c.Args.String("topic"))

			configs, err := kafkatool.GetTopicConfig(c.Args.String("topic"))
			if err != nil {
				fmt.Println("GetTopicConfig", err)
				return err
			}
			c.App.Printf("Found configs %d\n", len(configs))

			table := tablewriter.NewWriter(os.Stdout)
			table.Header(backend.ConfigHeader())
			for i := range configs {
				table.Append(configs[i].ToStrings())
			}
			table.Render()

			return nil
		},
	})

	showTopicCommand.AddCommand(&grumble.Command{
		Name: "partition",
		Help: "show partition of topic, eg: show topic partition {TOPIC_NAME}",
		Args: func(a *grumble.Args) {
			a.String("topic", "which topic to show", grumble.Default(""))
		},
		Run: func(c *grumble.Context) error {
			if len(c.Args.String("topic")) == 0 {
				return fmt.Errorf("topic name is empty")
			}
			c.App.Println("show topic patition:", c.Args.String("topic"))

			partitions, err := kafkatool.GetTopicPartition(c.Args.String("topic"))
			if err != nil {
				return err
			}
			c.App.Printf("Found partitions %d\n", len(partitions))

			table := tablewriter.NewWriter(os.Stdout)
			table.Header(backend.PartitionrHeader())
			for i := range partitions {
				table.Append(partitions[i].ToStrings())
			}
			table.Render()

			return nil
		},
	})

	// show group config
	showGroupCommand := &grumble.Command{
		Name:     "group",
		Help:     "show group tools",
		LongHelp: "show config of group",
	}
	showCommand.AddCommand(showGroupCommand)

	// showGroupCommand.AddCommand(&grumble.Command{
	// 	Name: "config",
	// 	Help: "show config of group, eg: show group config {TOPIC_NAME}",
	// 	Args: func(a *grumble.Args) {
	// 		a.String("group", "which group to show", grumble.Default(""))
	// 	},
	// 	Run: func(c *grumble.Context) error {
	// 		if len(c.Args.String("group")) == 0 {
	// 			return fmt.Errorf("group name is empty")
	// 		}
	// 		c.App.Println("show group config:", c.Args.String("group"))

	// 		configs, err := kafkatool.GetGroupConfig(c.Args.String("group"))
	// 		if err != nil {
	// 			fmt.Println("GetGroupConfig", err)
	// 			return err
	// 		}
	// 		c.App.Printf("configs: %#v\n", configs)

	// 		return nil
	// 	},
	// })

	showGroupCommand.AddCommand(&grumble.Command{
		Name: "offset",
		Help: "show offset of group, eg: show group offset {TOPIC_NAME}",
		Args: func(a *grumble.Args) {
			a.String("group", "which group to show", grumble.Default(""))
		},
		Run: func(c *grumble.Context) error {
			if len(c.Args.String("group")) == 0 {
				return fmt.Errorf("group name is empty")
			}
			c.App.Println("show group offset:", c.Args.String("group"))

			offsets, err := kafkatool.GetGroupOffset(c.Args.String("group"))
			if err != nil {
				fmt.Println("GetGroupOffset", err)
				return err
			}
			c.App.Printf("Found offsets %d\n", len(offsets))

			table := tablewriter.NewWriter(os.Stdout)
			table.Header(backend.GroupOffsetHeader())
			for i := range offsets {
				table.Append(offsets[i].ToStrings())
			}
			table.Render()

			return nil
		},
	})

	// show broker config
	showBrokerCommand := &grumble.Command{
		Name:     "broker",
		Help:     "show broker tools",
		LongHelp: "show config of broker",
	}
	showCommand.AddCommand(showBrokerCommand)

	showBrokerCommand.AddCommand(&grumble.Command{
		Name: "config",
		Help: "show config of broker id, id sometimes is 1, eg: show broker config {BROKER_ID}",
		Args: func(a *grumble.Args) {
			a.String("broker", "which broker to show", grumble.Default(""))
		},
		Run: func(c *grumble.Context) error {
			if len(c.Args.String("broker")) == 0 {
				return fmt.Errorf("broker id is empty")
			}
			c.App.Println("show broker config:", c.Args.String("broker"))

			configs, err := kafkatool.GetBrokerConfig(c.Args.String("broker"))
			if err != nil {
				fmt.Println("GetBrokerConfig", err)
				return err
			}
			c.App.Printf("Found configs %d\n", len(configs))

			table := tablewriter.NewWriter(os.Stdout)
			table.Header(backend.ConfigHeader())
			for i := range configs {
				table.Append(configs[i].ToStrings())
			}
			table.Render()

			return nil
		},
	})

	// show cluster config
	showClusterCommand := &grumble.Command{
		Name:     "cluster",
		Help:     "show cluster tools",
		LongHelp: "show config of cluster",
	}
	showCommand.AddCommand(showClusterCommand)

	showClusterCommand.AddCommand(&grumble.Command{
		Name: "config",
		Help: "show config of cluster id, id sometimes is 1, eg: show cluster config {CLUSTER_ID}",
		Args: func(a *grumble.Args) {
			a.String("cluster", "which cluster to show", grumble.Default(""))
		},
		Run: func(c *grumble.Context) error {
			if len(c.Args.String("cluster")) == 0 {
				return fmt.Errorf("cluster id is empty")
			}
			c.App.Println("show cluster config:", c.Args.String("cluster"))

			configs, err := kafkatool.GetClusterConfig(c.Args.String("cluster"))
			if err != nil {
				fmt.Println("GetClusterConfig", err)
				return err
			}
			c.App.Printf("Found configs %d\n", len(configs))

			table := tablewriter.NewWriter(os.Stdout)
			table.Header(backend.ConfigHeader())
			for i := range configs {
				table.Append(configs[i].ToStrings())
			}
			table.Render()

			return nil
		},
	})
}
