package backend

import (
	"fmt"
	"testing"
)

// Group 没有config项，这个测试是会出错
func TestGetGroupConfig(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	configs, err := kafkatool.GetGroupConfig("testgroup")
	if err != nil {
		t.Fatal("GetGroupConfig failed ", err)
	}
	t.Logf("configs: %#v", configs)
}

// 这个函数是kafka-go库有错误, kafka-go v0.4.48 Now work fine
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
	fmt.Printf("desc: %s\n", desc)
	t.Logf("desc: %s", desc)
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
