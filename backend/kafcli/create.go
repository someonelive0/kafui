package main

import (
	"fmt"
	"strconv"

	"github.com/desertbit/grumble"
)

func init() {
	createCommand := &grumble.Command{
		Name:     "create",
		Help:     "create topic",
		LongHelp: "create topic with number of partitions and replicas",
	}
	App.AddCommand(createCommand)

	// create topic
	createTopicCommand := &grumble.Command{
		Name:     "topic",
		Help:     "create topic",
		LongHelp: "create topic",
		Args: func(a *grumble.Args) {
			a.String("topic", "topic name", grumble.Default(""))
			a.String("partitions", "number of partitions", grumble.Default("1"))
			a.String("replicas", "number of replicas", grumble.Default("1"))
		},
		Run: func(c *grumble.Context) error {
			topic := c.Args.String("topic")
			if len(topic) == 0 {
				return fmt.Errorf("topic name is empty")
			}
			partitions, err := strconv.Atoi(c.Args.String("partitions"))
			if err != nil {
				return fmt.Errorf("partitions is not number")
			}
			replicas, err := strconv.Atoi(c.Args.String("replicas"))
			if err != nil {
				return fmt.Errorf("replicas is not number")
			}
			c.App.Printf("create topic [%s], partitions %d, replicas %d\n", topic, partitions, replicas)

			err = kafkatool.CreateTopic(topic, partitions, replicas)
			if err != nil {
				// fmt.Println("CreateTopic failed: ", err)
				return err
			}
			c.App.Printf("create topic [%s] success!\n", topic)

			return nil
		},
	}
	createCommand.AddCommand(createTopicCommand)
}
