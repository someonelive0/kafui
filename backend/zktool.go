package backend

import (
	"encoding/json"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type ZkTool struct {
}

func (p *ZkTool) ListBrokers(hosts []string) ([]string, error) {
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	items, _, err := conn.Children("/brokers/ids")
	if err != nil {
		return nil, err
	}

	brokers := make([]string, 0)
	for _, item := range items {
		data, _, err := conn.Get("/brokers/ids/" + item)
		if err != nil {
			return nil, err
		}
		m := make(map[string]interface{})
		if err = json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		if endpoints, exist := m["endpoints"]; exist {
			// fmt.Printf("endpoints %#v", endpoints)
			if ends, ok := endpoints.([]interface{}); ok {
				for _, endpoint := range ends {
					if end, ok := endpoint.(string); ok {
						brokers = append(brokers, end)
					}
				}
			}
		}
	}

	return brokers, nil
}

func (p *ZkTool) ListTopics(hosts []string) ([]string, error) {
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	items, _, err := conn.Children("/brokers/topics")
	if err != nil {
		return nil, err
	}

	return items, nil
}
