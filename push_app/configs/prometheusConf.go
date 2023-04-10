package configs

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type PrometheusConf struct {
	Query     string `yaml:"query"`
	NumWorker int    `yaml:"numWorker"`
	Route     string
	Token     string
	ClientId  string
	Interval  string
}

func (ps *PrometheusConf) LoadConfig() error {
	if data, err := ioutil.ReadFile("prom_config.yaml"); err != nil {
		fmt.Printf("Error reading prom config%v\n", err)
		return err
	} else {
		if err := yaml.Unmarshal(data, &ps); err != nil {
			return err
		}
	}
	fmt.Println(os.Environ())
	ps.Route = os.Getenv("PROM_URL")
	ps.Token = os.Getenv("PROM_TOKEN")
	fmt.Printf("Prom Route: %v\n", ps.Route)
	fmt.Printf("Prom Token: %v\n", ps.Token)
	return nil
}
