package components

/*
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// JSON response class from loki.
type LokiResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream struct {
				ClusterLogLevel  string `json:"cluster_log_level"`
				K8SNamespaceName string `json:"k8s_namespace_name"`
				OpenshiftCluster string `json:"openshift_cluster"`
			} `json:"stream"`
			Values [][]string `json:"values"`
		} `json:"result"`
		Stats struct {
			Summary struct {
				BytesProcessedPerSecond int     `json:"bytesProcessedPerSecond"`
				LinesProcessedPerSecond int     `json:"linesProcessedPerSecond"`
				TotalBytesProcessed     int     `json:"totalBytesProcessed"`
				TotalLinesProcessed     int     `json:"totalLinesProcessed"`
				ExecTime                float64 `json:"execTime"`
			} `json:"summary"`
			Store struct {
				TotalChunksRef        int     `json:"totalChunksRef"`
				TotalChunksDownloaded int     `json:"totalChunksDownloaded"`
				ChunksDownloadTime    float64 `json:"chunksDownloadTime"`
				HeadChunkBytes        int     `json:"headChunkBytes"`
				HeadChunkLines        int     `json:"headChunkLines"`
				DecompressedBytes     int     `json:"decompressedBytes"`
				DecompressedLines     int     `json:"decompressedLines"`
				CompressedBytes       int     `json:"compressedBytes"`
				TotalDuplicates       int     `json:"totalDuplicates"`
			} `json:"store"`
			Ingester struct {
				TotalReached       int `json:"totalReached"`
				TotalChunksMatched int `json:"totalChunksMatched"`
				TotalBatches       int `json:"totalBatches"`
				TotalLinesSent     int `json:"totalLinesSent"`
				HeadChunkBytes     int `json:"headChunkBytes"`
				HeadChunkLines     int `json:"headChunkLines"`
				DecompressedBytes  int `json:"decompressedBytes"`
				DecompressedLines  int `json:"decompressedLines"`
				CompressedBytes    int `json:"compressedBytes"`
				TotalDuplicates    int `json:"totalDuplicates"`
			} `json:"ingester"`
		} `json:"stats"`
	} `json:"data"`
}

// Fetch Job class for each data worker.
type FetchJob struct {
	start int
	end   int
	done  bool
	EOF   bool
}

// Main Datasource class.
type LokiDataSource struct {
	config     LogConfig
	logChan    chan Log
	jobChan    chan FetchJob
	numWorkers int
}

// Main goroutiene for fetching logs from loki.
// Spawns multiple workers for concurrent API calls
// Queues jobs for each worker
func (ds *LokiDataSource) Run() {
	for i := 0; i < ds.numWorkers; i++ {
		go ds.Worker()
	}

	ds.QueueJobs()
}

// QueueJobs queues jobs for each worker
func (ds *LokiDataSource) QueueJobs() {
	for i := ds.config.ParserConfig.Start; i <= ds.config.ParserConfig.End; i += 10 {
		//log.Printf("Queueing job for %d\t%d\t%d\n", i, ds.config.ParserConfig.Start, ds.config.ParserConfig.End)

		if i+10 > ds.config.ParserConfig.End {
			ds.jobChan <- FetchJob{i, ds.config.ParserConfig.End, false, false}
		} else {
			ds.jobChan <- FetchJob{i, i + 10, false, false}
		}
	}

	for i := 0; i < ds.numWorkers; i++ {
		if i == ds.numWorkers-1 {
			ds.jobChan <- FetchJob{0, 0, true, true}
		} else {
			ds.jobChan <- FetchJob{0, 0, true, false}
		}
	}
}

// Parses Json response from loki and forwards log to logChan
func (ds *LokiDataSource) ParseResponse(body []byte) {
	var response LokiResponse
	err := json.Unmarshal(body, &response)

	if err != nil {
		log.Print("JSON Error:", err)
		log.Print("JSON Response:", string(body))
		return
	}

	for _, result := range response.Data.Result {
		for _, values := range result.Values {
			ds.logChan <- Log{
				Val: []byte(values[1]),
				EOF: false,
			}
		}
	}
}

// Worker gorooutine for fetching logs from loki
// Receives jobs from jobChan, fetches logs from loki and parses response
func (ds *LokiDataSource) Worker() {
	for job := range ds.jobChan {
		if job.done {
			if job.EOF {
				ds.logChan <- Log{EOF: true}
			}

			break
		}

		err := ds.GetLogs(job.start, job.end)

		if err != nil {
			log.Println(err)
		}
	}
}

// Query Range API call
func (ds *LokiDataSource) GetLogs(start int, end int) error {
	url := "https://loki-frontend-opf-observatorium.apps.smaug.na.operate-first.cloud/loki/api/v1/query_range?"
	url += "query=" + ds.config.ParserConfig.Query
	url += "&start=" + strconv.Itoa(start)
	url += "&end=" + strconv.Itoa(end)
	url += "&limit=" + strconv.Itoa(ds.config.ParserConfig.Limit)

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("X-Scope-OrgID", ds.config.ParserConfig.XScopeOrgID)
	req.Header.Add("Authorization", ds.config.ParserConfig.AuthToken)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	ds.ParseResponse(body)
	return nil
}

func (ds *LokiDataSource) GetLogChan() chan Log {
	return ds.logChan
}

// Creates new LokiDataSource

func MakeLokiDataSource(config LogConfig) *LokiDataSource {
	lokiDataSource := LokiDataSource{
		config:     config,
		logChan:    make(chan Log),
		jobChan:    make(chan FetchJob),
		numWorkers: 8,
	}

	return &lokiDataSource
}
*/
