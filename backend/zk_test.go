package backend

import (
	"flag"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-zookeeper/zk"
)

var address = flag.String("address", "127.0.0.1:2181", "")

// go test .\backend -v -run TestZk -args -address 127.0.0.1:2181
func TestZk(t *testing.T) {
	flag.Parse()
	log.Println("address: ", *address)

	// 创建zk连接地址
	hosts := []string{*address}
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	println("connected: ", conn.Server())

	zklist(conn, "/brokers/ids")
	zkget(conn, "/controller")
	zklist(conn, "/brokers/topics")
}

func zklist(conn *zk.Conn, path string) {
	items, stat, err := conn.Children(path)
	if err != nil {
		fmt.Printf("查询Children %s失败, err: %v\n", path, err)
		return
	}

	fmt.Printf("Children %s 的值为 %d, %#v, %#v\n", path, len(items), items, *stat)
}

func zkget(conn *zk.Conn, path string) {
	data, stat, err := conn.Get(path)
	if err != nil {
		fmt.Printf("查询%s失败, err: %v\n", path, err)
		return
	}

	fmt.Printf("%s 的值为 [%s], %#v\n", path, string(data), *stat)
}

// go test .\backend -v -run TestZktool -args -address 127.0.0.1:2181
func TestZktool(t *testing.T) {
	flag.Parse()
	log.Println("address: ", *address)
	hosts := []string{*address}

	zktool := &ZkTool{}
	brokers, err := zktool.ListBrokers(hosts)
	if err != nil {
		t.Fatalf("ListBrokers failed, %s", err)
		return
	}
	t.Logf("brokers: %#v", brokers)

	topics, err := zktool.ListTopics(hosts)
	if err != nil {
		t.Fatalf("ListTopics failed, %s", err)
		return
	}
	t.Logf("brokers: %#v", topics)
}
