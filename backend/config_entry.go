package backend

import (
	"strconv"

	"github.com/segmentio/kafka-go"
)

// Copy from segmentio/kafka-go kafka.DescribeConfigResponseConfigEntry
type ConfigEntry struct {
	ConfigName  string `json:"config_name"`
	ConfigValue string `json:"config_value"`
	ReadOnly    bool   `json:"readonly"`

	// Ignored if API version is greater than v0
	IsDefault bool `json:"is_default"`

	// Ignored if API version is less than v1
	ConfigSource int8 `json:"config_source"`

	IsSensitive bool `json:"is_sensitive"`

	// Ignored if API version is less than v1
	ConfigSynonyms []ConfigSynonym `json:"config_synonym"`

	// Ignored if API version is less than v3
	ConfigType int8 `json:"config_type"`

	// Ignored if API version is less than v3
	ConfigDocumentation string `json:"config_document"`
}

type ConfigSynonym struct {
	// Ignored if API version is less than v1
	ConfigName string `json:"config_name"`

	// Ignored if API version is less than v1
	ConfigValue string `json:"config_value"`

	// Ignored if API version is less than v1
	ConfigSource int8 `json:"config_source"`
}

func NewConfigFromSegmentio(configEntry *kafka.DescribeConfigResponseConfigEntry) *ConfigEntry {
	config := &ConfigEntry{
		ConfigName:          configEntry.ConfigName,
		ConfigValue:         configEntry.ConfigValue,
		ReadOnly:            configEntry.ReadOnly,
		IsDefault:           configEntry.IsDefault,
		ConfigSource:        configEntry.ConfigSource,
		IsSensitive:         configEntry.IsSensitive,
		ConfigSynonyms:      make([]ConfigSynonym, 0, len(configEntry.ConfigSynonyms)),
		ConfigType:          configEntry.ConfigType,
		ConfigDocumentation: configEntry.ConfigDocumentation,
	}

	for i := range configEntry.ConfigSynonyms {
		synonym := ConfigSynonym{
			ConfigName:   configEntry.ConfigSynonyms[i].ConfigName,
			ConfigValue:  configEntry.ConfigSynonyms[i].ConfigValue,
			ConfigSource: configEntry.ConfigSynonyms[i].ConfigSource,
		}
		config.ConfigSynonyms = append(config.ConfigSynonyms, synonym)
	}

	return config
}

func NewConfigArrayFromSegmentio(configEntrys []kafka.DescribeConfigResponseConfigEntry) []ConfigEntry {
	configs := make([]ConfigEntry, 0, len(configEntrys))
	for i := range configEntrys {
		config := NewConfigFromSegmentio(&configEntrys[i])
		configs = append(configs, *config)
	}
	return configs
}

func (p *ConfigEntry) ToStrings() []string {
	return []string{p.ConfigName, p.ConfigValue, strconv.FormatBool(p.ReadOnly),
		strconv.FormatBool(p.IsDefault), strconv.Itoa(int(p.ConfigSource)), strconv.FormatBool(p.IsSensitive),
		strconv.Itoa(int(p.ConfigType)), p.ConfigDocumentation}
}

func ConfigHeader() []string {
	return []string{"ConfigName", "ConfigValue", "ReadOnly", "IsDefault",
		"ConfigSource", "IsSensitive", "ConfigType", "ConfigDocumentation"}
}
