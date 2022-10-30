package components

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"push_app/Response"
	"push_app/configs"
	"strconv"
)

const (
	PROMETHEUS string = "Prometheus"
)

type PrometheusDataSource struct {
	Conf        *configs.PrometheusConf
	RawDataChan chan PrometheusMetric
}

type PrometheusMetric struct {
	Metric struct {
		MetricName string
		Container  string
		Endpoint   string
		Instance   string
		Pod        string
		Job        string
		Service    string
	}
	Values int
	EOF    bool
}

func (ps PrometheusDataSource) Run() error {
	url := "https://" + ps.Conf.Route + "/api/v1/query" + "?query=" + ps.Conf.Query
	var bearer = "Bearer " + ps.Conf.Token

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
		return err
	}

	respJson := Response.PrometheusResponse{}
	ps.ParseResponse(body, &respJson)
	for _, result := range respJson.Data.Result {
		valueStr := fmt.Sprintf("%v", result.Values[1])
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			log.Printf("Error turing value to int in prometheus datasource.Run()")
			return err
		}
		ps.RawDataChan <- PrometheusMetric{
			Metric: struct {
				MetricName string
				Container  string
				Endpoint   string
				Instance   string
				Pod        string
				Job        string
				Service    string
			}(result.Metric),
			Values: value,
		}
	}
	ps.RawDataChan <- PrometheusMetric{EOF: true}
	return nil
}

func (ps PrometheusDataSource) ParseResponse(body []byte, promResp *Response.PrometheusResponse) {
	err := json.Unmarshal(body, promResp)
	if err != nil {
		log.Println("JSON Error:", err)
		return
	}
	//log.Println("JSON Response:", string(body))
}

func (ps PrometheusDataSource) GetDataChan() chan PrometheusMetric {
	return ps.RawDataChan
}

/*
func MakeDataSource(opConfig configs.OperatorConf) DataSource {
	if strings.EqualFold(opConfig.DataSourceType, PROMETHEUS) {
		return MakePrometheusDataSource(opConfig)
	}
	log.Println("UNKNOWN DATA SOURCE!")
	return nil
}
*/

func MakePrometheusDataSource(opConfig configs.OperatorConf) *PrometheusDataSource {
	ps := PrometheusDataSource{}
	ps.RawDataChan = make(chan PrometheusMetric, opConfig.ChanBufSize)
	promConf := &configs.PrometheusConf{}
	promConf.LoadConfig()
	ps.Conf = promConf
	return &ps
}
