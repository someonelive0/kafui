package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	log "github.com/sirupsen/logrus"
)

type KafkaTool struct {
	kafkaConfig     *KafkaConfig
	mechanism       sasl.Mechanism
	dialer          *kafka.Dialer
	sharedTransport *kafka.Transport

	leader Broker
}

func NewKafkaTool(kafkaConfig *KafkaConfig) *KafkaTool {
	kafkatool := &KafkaTool{
		kafkaConfig: kafkaConfig,
	}
	kafkatool.Init(kafkaConfig)

	return kafkatool
}

func (p *KafkaTool) Init(kafkaConfig *KafkaConfig) {
	p.kafkaConfig = kafkaConfig

	// init sasl mechanism
	if kafkaConfig.SaslMechanism == "SASL_PLAINTEXT" {
		p.mechanism = &plain.Mechanism{
			Username: kafkaConfig.User,
			Password: kafkaConfig.Password,
		}
	}

	p.dialer = &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: p.mechanism,
	}

	p.sharedTransport = &kafka.Transport{
		DialTimeout: 10 * time.Second,
		IdleTimeout: 600 * time.Second,
		SASL:        p.mechanism,
	}

}

// 列出所有broker，同时保存leader broker
func (p *KafkaTool) ListBrokers() ([]Broker, error) {
	conn, err := p.dialer.DialContext(context.Background(), "tcp", p.kafkaConfig.Brokers[0])
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return nil, err
	}
	p.leader.Copy(&controller)
	log.Infof("leader: %#v\n", p.leader)

	brokers, err := conn.Brokers()
	if err != nil {
		return nil, err
	}
	// fmt.Printf("brokers: %#v\n", brokers)

	mybrokers := make([]Broker, 0, len(brokers))
	for i := range brokers {
		broker := NewBrokerFromSegmentio(&brokers[i])
		mybrokers = append(mybrokers, *broker)
	}

	sort.SliceStable(mybrokers, func(i, j int) bool {
		return mybrokers[i].ID < mybrokers[j].ID
	})
	return mybrokers, nil
}

func (p *KafkaTool) ListTopics() ([]string, error) {
	conn, err := p.dialer.DialContext(context.Background(), "tcp", p.kafkaConfig.Brokers[0])
	if err != nil {
		log.Printf("DialContext failed %s", err)
		return nil, err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, err
	}

	m := map[string]struct{}{}
	// 遍历所有分区取topic
	for _, p := range partitions {
		// fmt.Println("partitions: ", p)
		m[p.Topic] = struct{}{}
	}
	topics := make([]string, 0, len(m))
	for k := range m {
		topics = append(topics, k)
	}

	sort.Strings(topics)
	return topics, nil
}

func (p *KafkaTool) ListGroups() ([]string, error) {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.kafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	groupsreq := &kafka.ListGroupsRequest{Addr: client.Addr}
	groupsrep, err := client.ListGroups(context.Background(), groupsreq)
	if err != nil {
		return nil, err
	}
	// b, _ := json.MarshalIndent(groupsrep, "", " ")
	// fmt.Printf("groups: %s\n", b)

	groups := make([]string, 0, len(groupsrep.Groups))
	for i := range groupsrep.Groups {
		groups = append(groups, groupsrep.Groups[i].GroupID)
	}

	sort.Strings(groups)
	return groups, nil
}

func (p *KafkaTool) GetTopicMeta(topic string) ([]string, error) {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.kafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	metareq := &kafka.MetadataRequest{
		Addr:   client.Addr,
		Topics: []string{topic},
	}
	metaresp, err := client.Metadata(context.Background(), metareq)
	if err != nil {
		return nil, err
	}
	b, _ := json.MarshalIndent(metaresp, "", " ")
	fmt.Printf("meta: %s\n", b)

	return nil, nil
}

func (p *KafkaTool) GetTopicConfig(topic string) ([]ConfigEntry, error) {
	return p.GetConfig("topic", topic)
}

func (p *KafkaTool) GetGroupConfig(group string) ([]ConfigEntry, error) {
	return p.GetConfig("group", group)
}

func (p *KafkaTool) GetBrokerConfig(brokerid string) ([]ConfigEntry, error) {
	return p.GetConfig("broker", brokerid)
}

func (p *KafkaTool) GetClusterConfig(clusterid string) ([]ConfigEntry, error) {
	return p.GetConfig("cluster", clusterid)
}

func (p *KafkaTool) GetConfig(resourceType, resourceName string) ([]ConfigEntry, error) {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.kafkaConfig.Brokers[0]),
		Transport: p.sharedTransport,
	}

	rType := kafka.ResourceTypeUnknown
	switch resourceType {
	case "topic":
		rType = kafka.ResourceTypeTopic
	case "broker":
		rType = kafka.ResourceTypeBroker // 输入参数是broker的ID号，例如 "1"
	case "cluster":
		rType = kafka.ResourceTypeCluster // 输入参数是clusterid，例如 "1"
	case "group":
		rType = kafka.ResourceTypeGroup // group 通常没有配置项
	}

	req := &kafka.DescribeConfigsRequest{
		Addr: client.Addr,
		Resources: []kafka.DescribeConfigRequestResource{
			0: {
				ResourceType: rType,
				ResourceName: resourceName,
			},
		},
	}
	resp, err := client.DescribeConfigs(context.Background(), req)
	if err != nil {
		return nil, err
	}

	if len(resp.Resources) == 0 {
		return nil, fmt.Errorf("not found")
	}
	if resp.Resources[0].Error != nil {
		return nil, resp.Resources[0].Error
	}

	configs := NewConfigArrayFromSegmentio(resp.Resources[0].ConfigEntries)
	return configs, nil
}

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

// static functions
func TestKafa(kafkaConfig *KafkaConfig) (*Broker, error) {
	var mechanism sasl.Mechanism = nil
	if kafkaConfig.SaslMechanism == "SASL_PLAINTEXT" {
		mechanism = &plain.Mechanism{
			Username: kafkaConfig.User,
			Password: kafkaConfig.Password,
		}
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	conn, err := dialer.DialContext(context.Background(), "tcp", kafkaConfig.Brokers[0])
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return nil, err
	}
	// log.Infof("leader: %#v\n", controller)
	leader := NewBrokerFromSegmentio(&controller)

	return leader, err
}
