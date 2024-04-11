package backend

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type ZKCli struct {
	Hosts   []string
	Timeout int // default time.Second*5

	conn *zk.Conn
}

func (p *ZKCli) Connect(hosts []string) error {
	if p.Timeout <= 0 {
		p.Timeout = 5 // default time.Second*5
	}

	conn, _, err := zk.Connect(hosts, time.Second*time.Duration(p.Timeout))
	if err != nil {
		return err
	}

	println("connected: ", conn.Server())
	p.Hosts = hosts
	p.conn = conn

	return nil
}

func (p *ZKCli) Close() error {
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	} else {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (p *ZKCli) List(path string) ([]string, error) {
	items, _, err := p.conn.Children(path)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (p *ZKCli) Get(path string) ([]byte, error) {
	data, _, err := p.conn.Get(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}
