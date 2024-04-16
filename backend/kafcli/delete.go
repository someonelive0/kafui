package main

import (
	"fmt"

	"github.com/desertbit/grumble"
)

func init() {
	deleteCommand := &grumble.Command{
		Name:     "delete",
		Help:     "create topic",
		LongHelp: "create topic with name",
	}
	App.AddCommand(deleteCommand)

	// delete topic
	deleteTopicCommand := &grumble.Command{
		Name:     "topic",
		Help:     "delete topic",
		LongHelp: "delete topic",
		Args: func(a *grumble.Args) {
			a.String("topic", "topic name", grumble.Default(""))
		},
		Run: func(c *grumble.Context) error {
			topic := c.Args.String("topic")
			if len(topic) == 0 {
				return fmt.Errorf("topic name is empty")
			}
			c.App.Printf("delete topic [%s]\n", topic)

			err := kafkatool.DeleteTopic(topic)
			if err != nil {
				// fmt.Println("DeleteTopic failed: ", err)
				return err
			}
			c.App.Printf("delete topic [%s] success!\n", topic)

			return nil
		},
	}
	deleteCommand.AddCommand(deleteTopicCommand)
}
