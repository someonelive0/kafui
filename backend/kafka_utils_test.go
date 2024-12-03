package backend

import "testing"

func TestGetKafkaPartitionOffse(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}

	offsets, err := GetKafkaPartitionOffset(myconfig.Kafka.Brokers,
		myconfig.Kafka.User, myconfig.Kafka.Password, 0,
		"trs-stat-flow", "k1StatFlow", 0)
	if err != nil {
		t.Fatalf("SetKafkaPartitionOffset failed: %s", err)
	}
	t.Logf("offsets: %v", offsets)
}

func TestSetKafkaPartitionOffse(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}

	err = SetKafkaPartitionOffset(myconfig.Kafka.Brokers,
		myconfig.Kafka.User, myconfig.Kafka.Password, 0,
		"trs-stat-flow", "k1StatFlow", 0, 1)
	if err != nil {
		t.Fatalf("SetKafkaPartitionOffset failed: %s", err)
	}
}
