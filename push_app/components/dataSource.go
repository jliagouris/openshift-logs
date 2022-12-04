package components

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"push_app/Response"
	"push_app/configs"
	"strconv"
	"time"
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
	Key    DataShareKey
	EOF    bool
}

func (ps PrometheusDataSource) Run() error {
	params := url.Values{}
	params.Add("query", ps.Conf.Query)
	queryUrl := "https://" + ps.Conf.Route + "/api/v1/query?" + params.Encode()
	var bearer = "Bearer " + ps.Conf.Token

	// Create a new request using http
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		fmt.Printf("Error occurred creating new request: %v\n", err)
		return err
	}

	interval, err := time.ParseDuration(ps.Conf.Interval)
	if err != nil {
		fmt.Printf("Time Duration parsing error: %v\n", err)
		return err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}

	var queryId uint32
	queryId = 0
	for true {
		fmt.Printf("-------------------------------------------------------------------------- Query %v\n", queryId)
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
			return err
		}
		respJson := Response.PrometheusResponse{}
		ps.ParseResponse(body, &respJson)
		for _, result := range respJson.Data.Result {
			var seqNum uint32
			seqNum = 0
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
				Key: DataShareKey{
					ClientId: ps.Conf.ClientId,
					QueryId:  queryId,
					SeqNum:   seqNum,
				},
			}
			seqNum++
		}
		ps.RawDataChan <- PrometheusMetric{EOF: true}
		resp.Body.Close()
		queryId++
		time.Sleep(interval)
	}
	return nil
}

func (ps PrometheusDataSource) ParseResponse(body []byte, promResp *Response.PrometheusResponse) {
	err := json.Unmarshal(body, promResp)
	if err != nil {
		log.Println("JSON Error in Parse Response:", err)
		return
	}
}

func (ps PrometheusDataSource) GetDataChan() chan PrometheusMetric {
	return ps.RawDataChan
}

func MakePrometheusDataSource(opConfig configs.OperatorConf) *PrometheusDataSource {
	ps := PrometheusDataSource{}
	ps.RawDataChan = make(chan PrometheusMetric, opConfig.ChanBufSize)
	promConf := &configs.PrometheusConf{}
	promConf.LoadConfig()
	promConf.ClientId = opConfig.ClientId
	promConf.Interval = opConfig.TimeInterval
	ps.Conf = promConf
	return &ps
}
