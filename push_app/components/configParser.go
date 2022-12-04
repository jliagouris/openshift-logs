package components

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type LogConfig struct {
	ParserConfig struct {
		Topic     string `yaml:"topic"`
		LogSchema []struct {
			Key      string `yaml:"key"`
			Nullable bool   `yaml:"nullable"`
		} `yaml:"logSchema"`
	} `yaml:"parserConfig"`
}

func (cfg *LogConfig) LoadConfig() error {
	if data, err := ioutil.ReadFile("config.yaml"); err != nil {
		return err
	} else {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return err
		}
	}

	return nil
}
