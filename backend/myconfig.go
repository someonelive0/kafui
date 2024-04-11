package backend

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	PASSWORD_PREFIX = "BASE64$"
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
	Password      string   `toml:"password" json:"password"`
}

type ZkConfig struct {
	Hosts []string `toml:"hosts" json:"hosts"` // hosts = [ "localhost:2181" ]
}

func LoadConfig(filename string) (*Myconfig, error) {
	myconfig := &Myconfig{Filename: filename}
	if _, err := toml.DecodeFile(filename, myconfig); err != nil {
		return nil, err
	}

	// if password begin with "BASE64$...", then decode weith base64
	if len(myconfig.Kafka.Password) > len(PASSWORD_PREFIX) && strings.Index(myconfig.Kafka.Password, PASSWORD_PREFIX) == 0 {
		b, err := base64.StdEncoding.DecodeString(myconfig.Kafka.Password[7:])
		if err != nil {
			return nil, err
		}
		myconfig.Kafka.Password = string(b)
	}

	return myconfig, nil
}

// save myconfig to filename
func SaveConfig(myconfig *Myconfig, filename string) error {
	var tmpfile = filename + ".tmp"
	os.Remove(tmpfile)
	fp, err := os.OpenFile(tmpfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	fp.WriteString("# save by app, on " + time.Now().Format(time.RFC3339))
	fp.WriteString("\n\n\n")

	// encode password with base64
	myconfig.Kafka.Password = PASSWORD_PREFIX + base64.StdEncoding.EncodeToString([]byte(myconfig.Kafka.Password))

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

	if err = os.Rename(tmpfile, filename); err != nil {
		return err
	}
	return nil
}

func (p *Myconfig) Dump() []byte {
	b, _ := json.MarshalIndent(p, "", "  ")
	return b
}
