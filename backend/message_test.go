package backend

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestReadMsgs2Ch(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	ch := make(chan *Message, 10)
	go func() {
		defer close(ch)
		defer cancel()

		err := kafkatool.ReadMsgs2Ch(ctx, "test", 0, -1, ch)
		if err != nil {
			fmt.Println("GetMessages", err)
		}
	}()

	for msg := range ch {
		t.Logf("msg of %#v %s %d, %s = %s", msg.Time, msg.Topic, msg.Partition, string(msg.Key), string(msg.Value))
	}
}

func TestReadMsgs(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	msgs, err := kafkatool.ReadMsgs("test", -1, 3)
	if err != nil {
		t.Fatalf("ReadMsgs failed: %s", err)
	}

	for _, msg := range msgs {
		t.Logf("msg of %#v %s %d, %s = %s", msg.Time, msg.Topic, msg.Partition, string(msg.Key), string(msg.Value))
	}
}

func TestWriteMsg(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	t0 := time.Now().Format(time.RFC3339)
	err = kafkatool.WriteMsg("test", "key_"+t0, "value_"+t0)
	if err != nil {
		t.Fatalf("WriteMsg failed: %s", err)
	}
}
