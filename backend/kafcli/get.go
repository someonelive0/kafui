package main

import (
	"context"
	"fmt"
	"kafui/backend"
	"strconv"
	"time"

	"github.com/desertbit/grumble"
)

func init() {
	getCommand := &grumble.Command{
		Name:     "get",
		Help:     "get messages",
		LongHelp: "get messages from topic and partition",
	}
	App.AddCommand(getCommand)

	// get topic messages
	getTopicCommand := &grumble.Command{
		Name:     "topic",
		Help:     "get messages from topic",
		LongHelp: "get messages from topic",
		Args: func(a *grumble.Args) {
			a.String("topic", "which topic to get messages", grumble.Default(""))
			a.String("partition", "which partition to get messages", grumble.Default("-1"))
		},
		Run: func(c *grumble.Context) error {
			if len(c.Args.String("topic")) == 0 {
				return fmt.Errorf("topic name is empty")
			}
			partition, err := strconv.Atoi(c.Args.String("partition"))
			if err != nil {
				return fmt.Errorf("partition is not number")
			}
			c.App.Println("get topic message:", c.Args.String("topic"), partition)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			ch := make(chan *backend.Message, 10)

			go func() {
				defer close(ch)
				defer cancel()
				err := kafkatool.ReadMsgs2Ch(ctx, c.Args.String("topic"), partition, -1, ch)
				if err != nil {
					fmt.Println("GetMessages", err)
					return
				}
			}()

			count := 0
			for msg := range ch {
				c.App.Printf("%d > of %s %d : %s = %s\n", count,
					msg.Time, msg.Partition, string(msg.Key), string(msg.Value))
				count++
			}
			c.App.Printf("total get %d messages of topic %s\n", count, c.Args.String("topic"))

			return nil
		},
	}
	getCommand.AddCommand(getTopicCommand)

}
