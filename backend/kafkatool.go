package backend

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type KafkaTool struct {
	Appctx          *context.Context
	KafkaConfig     *KafkaConfig
	mechanism       sasl.Mechanism
	dialer          *kafka.Dialer
	sharedTransport *kafka.Transport

	leader Broker
}

func NewKafkaTool(KafkaConfig *KafkaConfig) *KafkaTool {
	kafkatool := &KafkaTool{
		KafkaConfig: KafkaConfig,
	}
	kafkatool.Init(KafkaConfig)

	return kafkatool
}

func (p *KafkaTool) Init(KafkaConfig *KafkaConfig) {
	p.KafkaConfig = KafkaConfig
	if p.Appctx != nil {
		runtime.LogInfof(*p.Appctx, "set kafka config SaslMechanism: %#v", p.KafkaConfig.SaslMechanism)
	}

	// init sasl mechanism
	if KafkaConfig.SaslMechanism == "SASL_PLAINTEXT" {
		p.mechanism = &plain.Mechanism{
			Username: KafkaConfig.User,
			Password: KafkaConfig.Password,
		}
	} else {
		p.mechanism = nil
	}

	if p.sharedTransport != nil {
		p.sharedTransport.CloseIdleConnections()
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

func (p *KafkaTool) Close() {
	if p.sharedTransport != nil {
		p.sharedTransport.CloseIdleConnections()
	}
}

// 列出所有broker，同时保存leader broker
func (p *KafkaTool) ListBrokers() ([]Broker, error) {
	conn, err := p.dialer.DialContext(context.Background(), "tcp", p.KafkaConfig.Brokers[0])
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return nil, err
	}
	p.leader.Copy(&controller)
	runtime.LogInfof(*p.Appctx, "dial kafka leader success: %#v", p.leader)

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
	conn, err := p.dialer.DialContext(context.Background(), "tcp", p.KafkaConfig.Brokers[0])
	if err != nil {
		runtime.LogErrorf(*p.Appctx, "DialContext failed %s", err)
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

func (p *KafkaTool) GetTopicMeta(topic string) ([]string, error) {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.KafkaConfig.Brokers[0]),
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
	runtime.LogInfof(*p.Appctx, "GetTopicMeta '%s' failed: %#v", topic, metaresp)

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
		Addr:      kafka.TCP(p.KafkaConfig.Brokers[0]),
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

func (p *KafkaTool) SetTopicConfig(topic, configName, configValue string) error {
	return p.SetConfig("topic", topic, configName, configValue)
}

func (p *KafkaTool) SetBrokerConfig(topic, configName, configValue string) error {
	return p.SetConfig("broker", topic, configName, configValue)
}

func (p *KafkaTool) SetClusterConfig(topic, configName, configValue string) error {
	return p.SetConfig("cluster", topic, configName, configValue)
}

func (p *KafkaTool) SetConfig(resourceType, resourceName, configName, configValue string) error {
	client := &kafka.Client{
		Addr:      kafka.TCP(p.KafkaConfig.Brokers[0]),
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

	// 不要调用 AlterConfigs，这个需要把所有参数都作为参数输入，否则修改一个参数后其他参数又变成默认值了。
	// 需要调用 IncrementalAlterConfigs，意思增量修改，表示可以修改其中一个参数，其他保持不变。
	req := &kafka.IncrementalAlterConfigsRequest{
		Addr: client.Addr,
		Resources: []kafka.IncrementalAlterConfigsRequestResource{
			0: {
				ResourceType: rType,
				ResourceName: resourceName,
				Configs: []kafka.IncrementalAlterConfigsRequestConfig{
					0: {
						Name:            configName,
						Value:           configValue,
						ConfigOperation: kafka.ConfigOperationSet,
					},
				},
			},
		},
	}
	resp, err := client.IncrementalAlterConfigs(context.Background(), req)
	if err != nil {
		return err
	}

	if resp.Resources == nil || len(resp.Resources) == 0 {
		return fmt.Errorf("not found")
	}
	s := ""
	for i := range resp.Resources {
		if resp.Resources[i].Error != nil {
			s += fmt.Sprintf("type=%d, name '%s' setconfig failed: %v ",
				resp.Resources[i].ResourceType, resp.Resources[i].ResourceName, resp.Resources[i].Error)
		}
	}
	if len(s) > 0 {
		return fmt.Errorf("%s", s)
	}

	return nil
}

// static functions
func TestKafa(KafkaConfig *KafkaConfig) (*Broker, error) {
	var mechanism sasl.Mechanism = nil
	if KafkaConfig.SaslMechanism == "SASL_PLAINTEXT" {
		mechanism = &plain.Mechanism{
			Username: KafkaConfig.User,
			Password: KafkaConfig.Password,
		}
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	conn, err := dialer.DialContext(context.Background(), "tcp", KafkaConfig.Brokers[0])
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
