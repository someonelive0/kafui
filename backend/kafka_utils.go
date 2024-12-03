package backend

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

// return map[int]int64 means partition and offset pair
func GetKafkaPartitionOffset(brokers []string,
	username, password string, dial_timeout int,
	group, topic string, backoffset int64) (map[int]int64, error) {

	cfg := sarama.NewConfig()
	// cfg.Consumer.Offsets.CommitInterval = 1 * time.Second
	// cfg.Version = sarama.V2_7_0_0 // sarama.V0_10_0_1

	if dial_timeout <= 0 {
		dial_timeout = 60
	}
	cfg.Net.DialTimeout = time.Duration(dial_timeout) * time.Second
	if len(username) > 0 {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.Mechanism = "PLAIN"
		cfg.Net.SASL.User = username
		cfg.Net.SASL.Password = password
	}

	//brokers := []string{addr}
	client, err := sarama.NewClient(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("client create error: %s", err)
	}
	defer client.Close()

	//有个要注意的地方，如果想获取某个partition的offset位置，
	//需要这个offsetManager的groupId和consumer的一致，否则拿到的offset是不正确的。
	offsetManager, err := sarama.NewOffsetManagerFromClient(group, client)
	if err != nil {
		return nil, fmt.Errorf("offsetManager create error: %s", err)
	}
	defer offsetManager.Close()

	var offsets = make(map[int]int64)
	for i := 0; i < 4096; i++ {
		partitionOffsetManager, err := offsetManager.ManagePartition(topic, int32(i))
		if err != nil {
			break
			//return nil, fmt.Errorf(partitionOffsetManager create for %s partition %d failed: %s", topic, i, err)
		}
		defer partitionOffsetManager.Close()

		// nextOffset == -1 means there is no offset anymore, tag some times is empty string
		nextOffset, _ := partitionOffsetManager.NextOffset() // return nextoffset and tag
		if nextOffset == -1 {
			break // -1 means there is no parition, and offsetManager.ManagePartition() won't return error
		}
		//fmt.Printf("nextOffset of partition %d: %d, %s\n", i, nextOffset, tag)

		// will cause hang here. why, maybe it will exit before reset offset. it's work
		if nextOffset > 0 && backoffset != 0 {
			n := backoffset
			if n < 0 || n > nextOffset {
				n = nextOffset
			}
			partitionOffsetManager.ResetOffset(nextOffset-n, "")
			//fmt.Println("ResetOffset: ", nextOffset-n)
			nextOffset = nextOffset - n
		}
		offsets[i] = nextOffset
	}

	return offsets, nil
}

// return map[int]int64 means partition and offset pair
func SetKafkaPartitionOffset(brokers []string,
	username, password string, dial_timeout int,
	group, topic string, partition int, nextoffset int64) error {

	cfg := sarama.NewConfig()
	// cfg.Consumer.Offsets.CommitInterval = 1 * time.Second
	// cfg.Version = sarama.V2_7_0_0 // sarama.V0_10_0_1

	if dial_timeout <= 0 {
		dial_timeout = 60
	}
	cfg.Net.DialTimeout = time.Duration(dial_timeout) * time.Second
	if len(username) > 0 {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.Mechanism = "PLAIN"
		cfg.Net.SASL.User = username
		cfg.Net.SASL.Password = password
	}

	//brokers := []string{addr}
	client, err := sarama.NewClient(brokers, cfg)
	if err != nil {
		return fmt.Errorf("client create error: %s", err)
	}
	defer client.Close()

	//有个要注意的地方，如果想获取某个partition的offset位置，
	//需要这个offsetManager的groupId和consumer的一致，否则拿到的offset是不正确的。
	offsetManager, err := sarama.NewOffsetManagerFromClient(group, client)
	if err != nil {
		return fmt.Errorf("offsetManager create error: %s", err)
	}
	defer offsetManager.Close()

	partitionOffsetManager, err := offsetManager.ManagePartition(topic, int32(partition))
	if err != nil {
		return fmt.Errorf("partitionOffsetManager create for %s partition %d failed: %s",
			topic, partition, err)
	}
	defer partitionOffsetManager.Close()

	// nextOffset == -1 means there is no offset anymore, tag some times is empty string
	nextOffset, _ := partitionOffsetManager.NextOffset() // return nextoffset and tag
	if nextOffset == -1 {
		return nil // -1 means there is no parition, and offsetManager.ManagePartition() won't return error
	}

	partitionOffsetManager.ResetOffset(nextoffset, "")

	return nil
}
