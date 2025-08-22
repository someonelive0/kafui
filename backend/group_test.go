package backend

import (
	"fmt"
	"testing"
)

// 这个函数是kafka-go库有错误
func TestGetGroupDesc(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	desc, err := kafkatool.GetGroupDesc("testgroup")
	if err != nil {
		t.Fatal("GetGroupDesc failed ", err)
	}
	t.Logf("desc: %#v", desc)
}

func TestGetGroupOffset(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	group_offsets, err := kafkatool.GetGroupOffset("testgroup")
	if err != nil {
		t.Fatalf("SetKafkaPartitionOffset failed: %s", err)
	}

	fmt.Printf("GetGroupOffset success, 'testgroup' group_offsets: %v\n", group_offsets)
}

func TestSetGroupOffset(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	err = kafkatool.SetGroupOffset("testgroup", "dnsTopic", 0, 1473)
	if err != nil {
		t.Fatalf("SetKafkaPartitionOffset failed: %s", err)
	}

	fmt.Println("SetGroupOffset success")
}
