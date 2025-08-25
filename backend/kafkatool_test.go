package backend

import "testing"

var config_fileame = "../kafui.toml"

func TestListBrokers(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	brokers, err := kafkatool.ListBrokers()
	if err != nil {
		t.Fatal("ListBrokers failed ", err)
	}
	t.Logf("brokers: %#v", brokers)
}

func TestListTopics(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	topics, err := kafkatool.ListTopics()
	if err != nil {
		t.Fatal("ListTopics failed ", err)
	}
	t.Logf("topics: %#v", topics)
}

func TestListGroups(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	groups, err := kafkatool.ListGroups()
	if err != nil {
		t.Fatal("ListGroups failed ", err)
	}
	t.Logf("topics: %#v", groups)
}

func TestGetTopicMeta(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	metas, err := kafkatool.GetTopicMeta("dbTopic")
	if err != nil {
		t.Fatal("GetTopicMeta failed ", err)
	}
	t.Logf("metas: %#v", metas)
}

func TestGetTopicConfig(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	configs, err := kafkatool.GetTopicConfig("k1AssetApp")
	if err != nil {
		t.Fatal("GetTopicConfig failed ", err)
	}
	t.Logf("configs: %#v", configs)
}

func TestGetBrokerConfig(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	configs, err := kafkatool.GetBrokerConfig("1")
	if err != nil {
		t.Fatal("GetBrokerConfig failed ", err)
	}
	t.Logf("configs: %#v", configs)
}

func TestGetClusterConfig(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	configs, err := kafkatool.GetClusterConfig("1")
	if err != nil {
		t.Fatal("GetClusterConfig failed ", err)
	}
	t.Logf("configs: %#v", configs)
}

// 该函数不准，通常返回-1，应该根据TopicPartitin获得Offset
func TestGetTopicOffset(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	desc, err := kafkatool.GetTopicOffset("httpTopic")
	if err != nil {
		t.Fatal("GetTopicOffset failed ", err)
	}
	t.Logf("offset: %#v", desc)
}

func TestGetTopicPartition(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	desc, err := kafkatool.GetTopicPartition("test1")
	if err != nil {
		t.Fatal("GetTopicPartition failed ", err)
	}
	t.Logf("partions: %#v", desc)
}

func TestGetTopicPartitionOffset(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	first, last, err := kafkatool.GetTopicPartitionOffset("httpTopic", 0)
	if err != nil {
		t.Fatal("GetTopicPartitionOffset failed ", err)
	}
	t.Logf("offset first, last: %d, %d", first, last)
}

func TestDeleteGroup(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	err = kafkatool.DeleteGroup("kafka2nats")
	if err != nil {
		t.Fatal("DeleteGroup failed ", err)
	}
}

func TestCreateTopic(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	err = kafkatool.CreateTopic("test-topic", 1, 1)
	if err != nil {
		t.Fatal("CreateTopic failed ", err)
	}
}

func TestDeleteTopic(t *testing.T) {
	myconfig, err := LoadConfig(config_fileame)
	if err != nil {
		t.Fatalf("LoadConfig [%s] failed: %s", config_fileame, err)
	}
	kafkatool := NewKafkaTool(&myconfig.Kafka)

	err = kafkatool.DeleteTopic("test-topic")
	if err != nil {
		t.Fatal("CreateTopic failed ", err)
	}
}
