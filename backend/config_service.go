package backend

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ConfigService struct {
	Appctx   *context.Context
	Filename string
	Myconfig *MyConfig
}

func (p *ConfigService) LoadConfig() error {
	myconfig, err := LoadConfig(p.Filename)
	if err != nil {
		runtime.LogInfof(*p.Appctx, "LoadConfig '%s' error: %v", p.Filename, err)
		return err
	}
	p.Myconfig = myconfig
	// runtime.LogInfof(*p.Appctx, "LoadConfig success: %s", p.Myconfig.Dump())

	return nil
}

func (p *ConfigService) GetKafkaConfigs() ([]KafkaConfig, error) {
	err := p.LoadConfig()
	if err != nil {
		return nil, err
	}

	return p.Myconfig.KafkaConfigs, nil
}

func (p *ConfigService) GetKafkaConfig(name string) (*KafkaConfig, error) {
	for _, v := range p.Myconfig.KafkaConfigs {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("没有找到配置项 '%s'", name)
}

func (p *ConfigService) AddKafkaConfig(name, templateName string) (*KafkaConfig, error) {
	newKafkaConfig := NewKafkaConfig()
	for _, v := range p.Myconfig.KafkaConfigs {
		if v.Name == name {
			return nil, fmt.Errorf("配置项 '%s' 已经存在", name)
		}
		if v.Name == templateName {
			*newKafkaConfig = v
		}
	}

	// new a ConnConfig and return it
	newKafkaConfig.Name = name
	p.Myconfig.KafkaConfigs = append(p.Myconfig.KafkaConfigs, *newKafkaConfig)
	SaveConfig(p.Myconfig, p.Filename)

	return newKafkaConfig, nil
}

func (p *ConfigService) UpdateKafkaConfig(kafkaConfig *KafkaConfig) error {
	for i := range p.Myconfig.KafkaConfigs {
		if p.Myconfig.KafkaConfigs[i].Name == kafkaConfig.Name {
			p.Myconfig.KafkaConfigs[i] = *kafkaConfig
			SaveConfig(p.Myconfig, p.Filename)
			return nil
		}
	}
	return fmt.Errorf("没有找到配置项 '%s'", kafkaConfig.Name)
}

func (p *ConfigService) DeleteKafkaConfig(name string) error {
	for i := range p.Myconfig.KafkaConfigs {
		if p.Myconfig.KafkaConfigs[i].Name == name {
			p.Myconfig.KafkaConfigs = append(p.Myconfig.KafkaConfigs[:i], p.Myconfig.KafkaConfigs[i+1:]...)
			SaveConfig(p.Myconfig, p.Filename)
			return nil
		}
	}
	return fmt.Errorf("没有找到配置项 '%s'", name)
}

// test conn config wether can connect database
// param connConfig maybe not stored in config file, it can be temp var
// return db version by map
func (p *ConfigService) TestKafkaConfig(kafkaConfig *KafkaConfig) (*Broker, error) {
	for _, v := range p.Myconfig.KafkaConfigs {
		if v.Name == kafkaConfig.Name {

			brokers, err := TestKafa(kafkaConfig)
			// fmt.Printf("brokers: %v\n", brokers)
			return brokers, err
		}
	}

	return nil, fmt.Errorf("没有找到配置项 '%s'", kafkaConfig.Name)
}
