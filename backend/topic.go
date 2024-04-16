package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

// kafka topic
// type Topic struct {
// }

// 该函数不准，通常返回-1，需要用 GetTopicPartitionOffset 函数才有正确返回
func (p *KafkaTool) GetTopicOffset(topic string) ([]string, error) {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.kafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	req := &kafka.ListOffsetsRequest{
		Addr:   client.Addr,
		Topics: make(map[string][]kafka.OffsetRequest),
	}

	offsetReqs := make([]kafka.OffsetRequest, 0)
	offsetReqs = append(offsetReqs, kafka.OffsetRequest{Partition: 0, Timestamp: time.Now().Unix() - 90000})
	req.Topics[topic] = offsetReqs
	resp, err := client.ListOffsets(context.Background(), req)
	if err != nil {
		return nil, err
	}
	b, _ := json.MarshalIndent(resp, "", " ")
	fmt.Printf("topic offset %s: %s\n", topic, b)

	return nil, nil
}

func (p *KafkaTool) GetTopicPartition(topic string) ([]Partition, error) {
	conn, err := p.dialer.DialContext(context.Background(), "tcp", p.kafkaConfig.Brokers[0])
	if err != nil {
		log.Printf("DialContext failed %s", err)
		return nil, err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		return nil, err
	}

	// 遍历所有分区取 partition
	mypartitions := make([]Partition, 0)
	for i := range partitions {
		part := Partition{
			Topic:           partitions[i].Topic,
			ID:              partitions[i].ID,
			Leader:          *NewBrokerFromSegmentio(&partitions[i].Leader),
			Replicas:        NewBrokerArrayFromSegmentio(partitions[i].Replicas),
			Isr:             NewBrokerArrayFromSegmentio(partitions[i].Isr),
			OfflineReplicas: NewBrokerArrayFromSegmentio(partitions[i].OfflineReplicas),
		}

		first, last, err := p.GetTopicPartitionOffset(topic, partitions[i].ID)
		if err != nil {
			return nil, err
		}
		part.FirstOffset = first
		part.LastOffset = last
		part.Number = last - first
		mypartitions = append(mypartitions, part)
	}

	sort.Sort(PartitionSlice(mypartitions))
	return mypartitions, nil
}

func (p *KafkaTool) GetTopicPartitionOffset(topic string, partition int) (int64, int64, error) {
	conn, err := p.dialer.DialLeader(context.Background(), "tcp", p.kafkaConfig.Brokers[0], topic, partition)
	if err != nil {
		// log.Printf("DialContext [%s] failed %s", p.mechanism.Username, err)
		return 0, 0, err
	}
	defer conn.Close()

	first, last, err := conn.ReadOffsets()
	return first, last, err
}

// Not work
// func (p *KafkaTool) CreateTopic0(topic string, partitions, replicas int) error {
// 	conn, err := p.dialer.Dial("tcp", p.kafkaConfig.Brokers[0])
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close()

// 	controller, err := conn.Controller()
// 	if err != nil {
// 		return err
// 	}
// 	var controllerConn *kafka.Conn
// 	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
// 	if err != nil {
// 		return err
// 	}
// 	defer controllerConn.Close()

// 	topicConfigs := []kafka.TopicConfig{
// 		{
// 			Topic:             topic,
// 			NumPartitions:     partitions,
// 			ReplicationFactor: replicas,
// 		},
// 	}

// 	err = controllerConn.CreateTopics(topicConfigs...)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (p *KafkaTool) CreateTopic(topic string, partitions, replicas int) error {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.kafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	req := &kafka.CreateTopicsRequest{
		Addr: client.Addr,
		Topics: []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     partitions,
				ReplicationFactor: replicas,
			},
		},
	}

	resp, err := client.CreateTopics(context.Background(), req)
	if err != nil {
		return err
	}
	if err, ok := resp.Errors[topic]; ok {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *KafkaTool) DeleteTopic(topic string) error {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.kafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	req := &kafka.DeleteTopicsRequest{
		Addr:   client.Addr,
		Topics: []string{topic},
	}

	resp, err := client.DeleteTopics(context.Background(), req)
	if err != nil {
		return err
	}
	if err, ok := resp.Errors[topic]; ok {
		if err != nil {
			return err
		}
	}

	return nil
}
