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

// base on kafui.toml
type Myconfig struct {
	Filename  string        `toml:"-" json:"-"`
	Title     string        `toml:"title" json:"title"`
	License   string        `toml:"license" json:"license"`
	Kafka     []KafkaConfig `toml:"kafka" json:"kafka"`
	Zookeeper ZkConfig      `toml:"zookeeper" json:"zookeeper"`
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
	// check filename is exists
	if _, err := os.Stat(filename); err != nil {
		fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}
		defer fp.Close()
		fp.WriteString(CONFIG_FILE_TEMPLATE)
	}

	myconfig := &Myconfig{Filename: filename}
	if _, err := toml.DecodeFile(filename, myconfig); err != nil {
		return nil, err
	}

	for _, kafka := range myconfig.Kafka {
		// if password begin with "BASE64$...", then decode weith base64
		if len(kafka.Password) > len(PASSWORD_PREFIX) && strings.Index(kafka.Password, PASSWORD_PREFIX) == 0 {
			b, err := base64.StdEncoding.DecodeString(kafka.Password[7:])
			if err != nil {
				return nil, err
			}
			kafka.Password = string(b)
		}
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

	for _, kafka := range myconfig.Kafka {
		// encode password with base64
		kafka.Password = PASSWORD_PREFIX + base64.StdEncoding.EncodeToString([]byte(kafka.Password))
	}

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

const (
	DEFAULT_CONFIG_FILE  = "kafui.toml"
	PASSWORD_PREFIX      = "BASE64$"
	CONFIG_FILE_TEMPLATE = `
# Kafui config file template


title = "Kafui"
license = "Copyright @ 2024"
	
	
[kafka]
	name = "localhost"
	brokers = [ "127.0.0.1:9092" ]
	# sasl mechanism should be empty or "SASL_PLAINTEXT",
	# if mechanism is "SASL_PLAINTEXT", then set user and password
	sasl_mechanism = ""
	user = ""
	password = ""
	
[zookeeper]
	hosts = [ "127.0.0.1:2181" ]
`
)
