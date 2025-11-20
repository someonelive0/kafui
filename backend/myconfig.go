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
type MyConfig struct {
	Filename     string        `toml:"-" json:"-"`
	Title        string        `toml:"title" json:"title"`
	License      string        `toml:"license" json:"license"`
	KafkaConfigs []KafkaConfig `toml:"kafka" json:"kafka"`
	Zookeeper    ZkConfig      `toml:"zookeeper" json:"zookeeper"`
}

type KafkaConfig struct {
	Name          string   `toml:"name" json:"name"`
	Brokers       []string `toml:"brokers" json:"brokers"`               // brokers = [ "localhost:9092" ]
	SaslMechanism string   `toml:"sasl_mechanism" json:"sasl_mechanism"` // "" or "SASL_PLAINTEXT"
	User          string   `toml:"user" json:"user"`
	Password      string   `toml:"password" json:"password"`
	Timeout       int      `toml:"timeout" json:"timeout"`
}

type ZkConfig struct {
	Hosts []string `toml:"hosts" json:"hosts"` // hosts = [ "localhost:2181" ]
}

func NewKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Name:          "default",
		Brokers:       []string{"127.0.0.1:9092"},
		SaslMechanism: "",
		User:          "",
		Password:      "",
		Timeout:       DEFAULT_TIMEOUT,
	}
}

func (p *KafkaConfig) Dump() []byte {
	b, _ := json.MarshalIndent(p, "", "  ")
	return b
}

func LoadConfig(filename string) (*MyConfig, error) {
	// check filename is exists
	if _, err := os.Stat(filename); err != nil {
		fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}
		defer fp.Close()
		fp.WriteString(CONFIG_FILE_TEMPLATE)
	}

	myconfig := &MyConfig{Filename: filename}
	if _, err := toml.DecodeFile(filename, myconfig); err != nil {
		return nil, err
	}

	// if password begin with "BASE64$...", then decode weith base64
	for i := range myconfig.KafkaConfigs {
		tmp := myconfig.KafkaConfigs[i].Password
		if len(tmp) > len(PASSWORD_PREFIX) &&
			strings.Index(tmp, PASSWORD_PREFIX) == 0 {
			b, err := base64.StdEncoding.DecodeString(tmp[len(PASSWORD_PREFIX):])
			if err != nil {
				return nil, err
			}
			myconfig.KafkaConfigs[i].Password = string(b)
		}

		if myconfig.KafkaConfigs[i].Timeout <= 0 {
			myconfig.KafkaConfigs[i].Timeout = DEFAULT_TIMEOUT
		}
	}

	return myconfig, nil
}

// save myconfig to filename
func SaveConfig(myconfig *MyConfig, filename string) error {
	var tmpconfig = *myconfig // deep copy myconfig, ATTENTION this is not a deep copy of slice inside struct
	tmpconfig.KafkaConfigs = make([]KafkaConfig, len(myconfig.KafkaConfigs))
	copy(tmpconfig.KafkaConfigs, myconfig.KafkaConfigs)

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
	for i := range tmpconfig.KafkaConfigs {
		tmp := tmpconfig.KafkaConfigs[i].Password
		if len(tmp) > len(PASSWORD_PREFIX) &&
			strings.Index(tmp, PASSWORD_PREFIX) == 0 {
			continue
		}
		tmpconfig.KafkaConfigs[i].Password = PASSWORD_PREFIX + base64.StdEncoding.EncodeToString([]byte(tmp))
	}

	buf := new(bytes.Buffer)
	if err = toml.NewEncoder(buf).Encode(tmpconfig); err != nil {
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

func (p *MyConfig) Dump() []byte {
	b, _ := json.MarshalIndent(p, "", "  ")
	return b
}

const (
	DEFAULT_CONFIG_FILE  = "kafui.toml"
	DEFAULT_TIMEOUT      = 10
	PASSWORD_PREFIX      = "BASE64$"
	CONFIG_FILE_TEMPLATE = `
# Kafui config file template


title = "Kafui"
license = "Copyright @ 2024"
	
	
[[kafka]]
	name = "localhost"
	brokers = [ "127.0.0.1:9092" ]
	# sasl mechanism should be empty or "SASL_PLAINTEXT",
	# if mechanism is "SASL_PLAINTEXT", then set user and password
	sasl_mechanism = ""
	user = ""
	password = ""
	# timeout for connect and read timeout in seconds, default 10 seconds
	timeout = 10
	
[zookeeper]
	hosts = [ "127.0.0.1:2181" ]
`
)
