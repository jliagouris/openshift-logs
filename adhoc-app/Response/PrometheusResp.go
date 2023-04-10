package Response

// JSON response class from Prometheus.
type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				MetricName string `json:"__name__"`
				Container  string `json:"container"`
				Endpoint   string `json:"endpoint"`
				Instance   string `json:"instance"`
				Pod        string `json:"pod"`
				Job        string `json:"job"`
				Service    string `json:"service"`
			} `json:"metric"`
			Values []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}
