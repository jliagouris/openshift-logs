package main

import (
	"io/ioutil"

	"gopkg.in/yaml"
)

type Config struct {
	ParserConfig struct {
		AuthToken   string `yaml:"authToken"`
		XScopeOrgID string `yaml:"xScopeOrgID"`
		Query       string `yaml:"query"`
		Start       int    `yaml:"start"`
		End         int    `yaml:"end"`
		Limit       int    `yaml:"limit"`
		LogSchema   []struct {
			Key      string `yaml:"key"`
			Nullable bool   `yaml:"nullable"`
		} `yaml:"logSchema"`
	} `yaml:"parserConfig"`
}

func (cfg *Config) LoadConfig() error {
	if data, err := ioutil.ReadFile("config.yaml"); err != nil {
		return err
	} else {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return err
		}
	}

	return nil
}
