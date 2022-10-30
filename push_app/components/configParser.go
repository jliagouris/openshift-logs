package components

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type LogConfig struct {
	ParserConfig struct {
		//AuthToken   string `yaml:"authToken"`
		//XScopeOrgID string `yaml:"xScopeOrgID"`
		//Query       string `yaml:"query"`
		//Start       int    `yaml:"start"`
		//End         int    `yaml:"end"`
		//Limit       int    `yaml:"limit"`
		Topic     string `yaml:"topic"`
		LogSchema []struct {
			Key      string `yaml:"key"`
			Nullable bool   `yaml:"nullable"`
		} `yaml:"schema"`
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
