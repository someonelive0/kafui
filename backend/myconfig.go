package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// base on kafui.toml
type Myconfig struct {
	Filename  string      `toml:"-" json:"-"`
	Title     string      `toml:"title" json:"title"`
	License   string      `toml:"license" json:"license"`
	Kafka     KafkaConfig `toml:"kafka" json:"kafka"`
	Zookeeper ZkConfig    `toml:"zookeeper" json:"zookeeper"`
}

type KafkaConfig struct {
	Name          string   `toml:"name" json:"name"`
	Brokers       []string `toml:"brokers" json:"brokers"`               // brokers = [ "localhost:9092" ]
	SaslMechanism string   `toml:"sasl_mechanism" json:"sasl_mechanism"` // "" or "SASL_PLAINTEXT"
	User          string   `toml:"user" json:"user"`
	Password      string   `toml:"password" json:"-"`
}

type ZkConfig struct {
	Hosts []string `toml:"hosts" json:"hosts"` // hosts = [ "localhost:2181" ]
}

func LoadConfig(filename string) (*Myconfig, error) {
	myconfig := &Myconfig{Filename: filename}
	if _, err := toml.DecodeFile(filename, myconfig); err != nil {
		return nil, err
	}

	return myconfig, nil
}

func SaveConfig(myconfig *Myconfig) error {
	os.Remove(myconfig.Filename + ".tmp")
	fp, err := os.OpenFile(myconfig.Filename+".tmp", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	fp.WriteString("# save by app, on " + time.Now().Format(time.RFC3339))
	fp.WriteString("\n\n\n")

	buf := new(bytes.Buffer)
	if err = toml.NewEncoder(buf).Encode(myconfig); err != nil {
		return err
	}
	n, err := fp.Write(buf.Bytes())

	if err != nil {
		return err
	}
	if n != buf.Len() {
		return fmt.Errorf("write not enough bytes, %d < %d", n, buf.Len())
	}
	if err = fp.Close(); err != nil {
		return err
	}

	if err = os.Rename(myconfig.Filename+".tmp", myconfig.Filename); err != nil {
		return err
	}
	return nil
}

func (p *Myconfig) Dump() []byte {
	b, _ := json.MarshalIndent(p, "", "  ")
	return b
}
