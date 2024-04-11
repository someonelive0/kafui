package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var topic = flag.String("topic", "my-topic", "kafka topic")
var address = flag.String("address", "127.0.0.1:9092", "")
var user = flag.String("user", "", "")
var password = flag.String("password", "", "")

func init() {
	flag.Parse()
	log.Println("topic: ", *topic)
	log.Println("address: ", *address)
	log.Println("user: ", *user)
}

func main() {
	getTopics()
	getTopicMeta()
	getGroups()
}

func getTopics() {
	// 指定要连接的topic和partition
	// partition := 0
	mechanism := plain.Mechanism{
		Username: *user,
		Password: *password,
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	// to connect to the kafka leader via an existing non-leader connection rather than using DialLeader
	conn, err := dialer.DialContext(context.Background(), "tcp", *address)
	// conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("controller: %#v\n", controller)

	var connLeader *kafka.Conn
	connLeader, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("leader: %#v\n", connLeader.Broker())
	defer connLeader.Close()

	apiversions, err := conn.ApiVersions()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("apiversions: %#v\n", apiversions)

	brokers, err := conn.Brokers()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("brokers: %#v\n", brokers)

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}
	// 遍历所有分区取topic
	for _, p := range partitions {
		// fmt.Println("partitions: ", p)
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		fmt.Println(k)
	}
}

func getTopicMeta() {
	mechanism := plain.Mechanism{
		Username: *user,
		Password: *password,
	}

	// Transport 负责管理连接池和其他资源,
	// 通常最好的使用方式是创建后在应用程序中共享使用它们。
	sharedTransport := &kafka.Transport{
		SASL:        mechanism,
		DialTimeout: 10 * time.Second,
		IdleTimeout: 600 * time.Second,
	}

	client := &kafka.Client{
		Addr:      kafka.TCP(*address),
		Timeout:   10 * time.Second,
		Transport: sharedTransport,
	}

	metareq := &kafka.MetadataRequest{
		Addr:   client.Addr,
		Topics: []string{"k1AssetApp"},
	}
	metaresp, err := client.Metadata(context.Background(), metareq)
	if err != nil {
		panic(err.Error())
	}
	b, _ := json.MarshalIndent(metaresp, "", " ")
	fmt.Printf("meta: %s\n", b)

	sharedTransport.CloseIdleConnections()
}

func getGroups() {
	// 指定要连接的topic和partition
	// partition := 0
	mechanism := plain.Mechanism{
		Username: *user,
		Password: *password,
	}

	// Transport 负责管理连接池和其他资源,
	// 通常最好的使用方式是创建后在应用程序中共享使用它们。
	sharedTransport := &kafka.Transport{
		SASL: mechanism,
	}

	client := &kafka.Client{
		Addr:      kafka.TCP(*address),
		Timeout:   10 * time.Second,
		Transport: sharedTransport,
	}

	groupsreq := &kafka.ListGroupsRequest{Addr: client.Addr}
	groupsrep, err := client.ListGroups(context.Background(), groupsreq)
	if err != nil {
		panic(err.Error())
	}
	b, _ := json.MarshalIndent(groupsrep, "", " ")
	fmt.Printf("groups: %s\n", b)

	groupid := "trs-dns"
	groupdescreq := &kafka.DescribeGroupsRequest{
		Addr:     client.Addr,
		GroupIDs: []string{groupid},
	}
	groupdescresp, err := client.DescribeGroups(context.Background(), groupdescreq)
	if err != nil {
		panic(err.Error())
	}
	b, _ = json.MarshalIndent(groupdescresp, "", " ")
	fmt.Printf("group descript %s: %s\n", groupid, b)
}
