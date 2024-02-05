// Package sender implements the Sender interface.
package sender

import (
	"crypto/tls"
	"fmt"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"io"
	"log"
	"math"
	"net/url"
	"time"

	"github.com/Kotletta-TT/MonoGo/cmd/agent/config"
	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"github.com/Kotletta-TT/MonoGo/internal/common"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
	JSON    = "json"
	TEXT    = "text"
)

type metricsStore interface {
	GetMetrics() map[string]*entity.Value
	GetMetricsSlice() []*common.Metrics
}

type Sender interface {
	Send()
}

type ResultWork struct {
	StatusCode int
	Body       []byte
	Err        error
}

type HTTPSender struct {
	repo   metricsStore
	client *resty.Client
	cfg    *config.Config
}

type TextPlainSender HTTPSender
type JSONSender HTTPSender

// NewRestyClient creates a new Resty client.
//
// It does not take any parameters.
// It returns a pointer to a resty.Client object.
func NewRestyClient(cfg *config.Config) *resty.Client {
	client := resty.New()
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)
	if cfg.SSL {
		cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			logger.Errorf("Failed to load client certificate: %s", err)
		}
		client.SetCertificates(cert)
		client.SetRootCertificate(cfg.CaPath)
	}
	return client
}

// NewHTTPSender returns a Sender based on the specified configuration.
//
// The function takes a metricsStore object and a *config.Config object as parameters.
// It switches on the SendType field of the config.Config object and returns a Sender based on the value.
// If the SendType is JSON, it returns a JSONSender object.
// If the SendType is TEXT, it returns a TextPlainSender object.
// If the SendType is neither JSON nor TEXT, it panics with the message "Send type unknown".
func NewHTTPSender(repo metricsStore, cfg *config.Config) Sender {
	switch cfg.SendType {
	case JSON:
		return &JSONSender{repo: repo, client: NewRestyClient(cfg), cfg: cfg}
	case TEXT:
		return &TextPlainSender{repo: repo, client: NewRestyClient(cfg), cfg: cfg}
	default:
		panic("Send type unknown")
	}
}

// compileURL compiles the URL for sending a metric.
//
// It takes a pointer to a TextPlainSender object (h) and a pointer to a Metrics object (metric) as parameters.
// It returns a string representing the compiled URL.
func (h *TextPlainSender) compileURL(metric *common.Metrics) string {
	compileURL := url.URL{Host: h.cfg.ServerHost, Scheme: "http"}
	switch metric.MType {
	case GAUGE:
		return compileURL.JoinPath("update", GAUGE, metric.ID, fmt.Sprintf("%f", *metric.Value)).String()
	case COUNTER:
		return compileURL.JoinPath("update", COUNTER, metric.ID, fmt.Sprintf("%d", *metric.Delta)).String()
	default:
		panic("Metric type unknown")
	}
}

// sendWorker sends metrics to a specified URL.
//
// The function takes in two parameters: jobs and results.
// - jobs is a channel that receives pointers to metrics of type *common.Metrics.
// - results is a channel that sends pointers to ResultWork structs.
//
// The function does the following:
// - For each metric received from the jobs channel, it compiles a URL using the compileURL method of the TextPlainSender struct.
// - It then logs the send URL using the log.Printf method.
// - Next, it sends an HTTP POST request to the send URL using the R method of the client struct.
// - The response and error from the request are stored in the resp and err variables respectively.
// - Finally, it sends a pointer to a ResultWork struct to the results channel, containing the status code, body, and error.
func (h *TextPlainSender) sendWorker(jobs <-chan *common.Metrics, results chan<- *ResultWork) {
	for metric := range jobs {
		sendURL := h.compileURL(metric)
		log.Printf("Send URL: %s", sendURL)
		resp, err := h.client.R().Post(sendURL)
		results <- &ResultWork{StatusCode: resp.StatusCode(), Body: resp.Body(), Err: err}
	}
}

// Send sends the Text/Plain metrics.
//
// No parameters.
// No return value.
func (h *TextPlainSender) Send() {
	log.Println("Start Text/Plain send metrics")
	metrics := h.repo.GetMetricsSlice()
	jobs := make(chan *common.Metrics, len(metrics))
	results := make(chan *ResultWork, len(metrics))
	for i := 0; i < h.cfg.RateLimit; i++ {
		go h.sendWorker(jobs, results)
	}
	for _, metric := range metrics {
		jobs <- metric
	}
	close(jobs)
	for i := 0; i < len(metrics); i++ {
		result := <-results
		if result.Err != nil {
			log.Printf("Error send metrics: %s\n", result.Err.Error())
			continue
		}
		if result.StatusCode != 200 {
			log.Printf("Error send metrics: %d\n", result.StatusCode)
			continue
		}
	}
}

// JSONMetricFabric constructs a JSON metric object based on the provided name and value.
//
// Parameters:
// - name: a string representing the ID of the metric.
// - value: a pointer to an entity.Value object containing the metric value and kind.
//
// Returns:
// - m: a pointer to a common.Metrics object representing the constructed JSON metric.
func JSONMetricFabric(name string, value *entity.Value) *common.Metrics {
	m := new(common.Metrics)
	m.ID = name
	switch value.Kind {
	case entity.KindGauge:
		m.MType = GAUGE
		val := math.Float64frombits(value.Metric)
		m.Value = &val
		return m
	case entity.KindCounter:
		m.MType = COUNTER
		delta := int64(value.Metric)
		m.Delta = &delta
		return m
	default:
		panic("Metric type unknown")
	}
}

// prepareBody prepares the body of the JSONSender function.
//
// The function takes in a variable number of *common.Metrics arguments.
// It returns a byte slice, a string, and an error.
func (j *JSONSender) prepareBody(m ...*common.Metrics) ([]byte, string, error) {
	var mBytes []byte
	var sign string
	var err error
	switch len(m) {
	case 0:
		mBytes = []byte{}
	case 1:
		mBytes, err = easyjson.Marshal(m[0])
	default:
		mArray := common.SliceMetrics(m)
		mBytes, err = easyjson.Marshal(mArray)
	}
	if err != nil {
		return nil, "", err
	}
	if j.cfg.HashKey != "" {
		sign, err = common.CreateHMACSignature(j.cfg.HashKey, mBytes)
		if err != nil {
			return nil, "", err
		}
	}
	if j.cfg.Compress == "gzip" {
		compressBytes, err := common.GzipCompress(mBytes)
		if err != nil {
			return nil, "", err
		}
		mBytes = compressBytes
	}
	return mBytes, sign, nil
}

// receiveResponse is a function that handles the response received from the server.
//
// It takes in two parameters:
//   - resp: a pointer to a resty.Response object that represents the server response.
//   - err: an error object that represents any error that occurred during the request.
//     If err is not nil and not io.EOF, it logs the error message and returns.
//
// The function checks if the err is not nil and not io.EOF. If it is, it logs the error message along with the
// response status code, response body, and the error message itself. Then it returns.
//
// If the function has a non-empty HashKey, it verifies the HMAC signature of the response body using the
// VerifyHMACSignature function from the common package. If there is an error during the verification, it logs
// the error message along with the response status code, response body, and the error message itself. Then it returns.
//
// If the response status code is not equal to 200, it logs the error message along with the response status code
// and response body.
func (j *JSONSender) receiveResponse(resp *resty.Response, err error) {
	if err != nil && err != io.EOF {
		log.Printf("error: Code: %d, Body: %s err: %s\n", resp.StatusCode(), resp.String(), err.Error())
		return
	}
	if j.cfg.HashKey != "" {
		if err := common.VerifyHMACSignature(resp.Header().Get("HashSHA256"), j.cfg.HashKey, resp.Body()); err != nil {
			log.Printf("error: Code: %d, Body: %s err: %s\n", resp.StatusCode(), resp.String(), err.Error())
			return
		}
	}
	if resp.StatusCode() != 200 {
		log.Printf("error: Code: %d, Body: %s\n", resp.StatusCode(), resp.String())
	}
}

// sendWorker sends metrics to the specified URL.
//
// It takes in a channel of metrics to send, a channel to receive the results,
// and the URL to send the metrics to.
// The metrics are sent as JSON.
// If compression is enabled, the metrics are sent in gzip format.
// The function prepares the metric body, sets the necessary headers,
// and sends the metric using an HTTP POST request.
// The response status code, body, and any errors encountered are sent back
// through the results channel.
func (j *JSONSender) sendWorker(jobs <-chan *common.Metrics, results chan<- *ResultWork, url string) {
	for metric := range jobs {
		req := j.client.R()
		req.SetHeader("Content-Type", "application/json")
		if j.cfg.Compress == "gzip" {
			req.SetHeader("Content-Encoding", "gzip")
			req.SetHeader("Accept-Encoding", "gzip")
		}
		mJSON, sign, err := j.prepareBody(metric)
		if err != nil {
			log.Printf("prepare body err: %s\n", err.Error())
			continue
		}
		req.SetBody(mJSON)
		if j.cfg.HashKey != "" && sign != "" {
			req.SetHeader("HashSHA256", sign)
		}
		resp, err := req.Post(url)
		results <- &ResultWork{StatusCode: resp.StatusCode(), Body: resp.Body(), Err: err}
	}
}

// Send sends the JSON data to the specified URL.
//
// It sets the appropriate headers and prepares the JSON body based on the
// configuration. If compression is enabled, it sets the necessary headers for
// gzip compression. It then sends the request to the server and receives the
// response. If batch support is enabled, it prepares the metrics and URL for
// batch updates. Otherwise, it prepares the metrics and URL for individual
// updates. It then sends the metrics concurrently to the server using worker
// goroutines and waits for the results. If an error occurs during sending or
// if the response status code is not 200, it logs the error.
func (j *JSONSender) Send() {
	var sendURL url.URL
	log.Println("Start JSON send metrics")
	req := j.client.R()
	req.SetHeader("Content-Type", "application/json")
	if j.cfg.Compress == "gzip" {
		req.SetHeader("Content-Encoding", "gzip")
		req.SetHeader("Accept-Encoding", "gzip")
	}
	switch j.cfg.BatchSupport {
	case true:
		metrics := j.repo.GetMetricsSlice()
		sendURL = url.URL{Host: j.cfg.ServerHost, Scheme: "http", Path: "/updates/"}
		mJSON, sign, err := j.prepareBody(metrics...)
		if err != nil {
			log.Printf("prepare body err: %s\n", err.Error())
			return
		}
		if j.cfg.HashKey != "" && sign != "" {
			req.SetHeader("HashSHA256", sign)
		}
		req.SetBody(mJSON)
		j.receiveResponse(req.Post(sendURL.String()))
	default:
		metrics := j.repo.GetMetricsSlice()
		sendURL = url.URL{Host: j.cfg.ServerHost, Scheme: "http", Path: "/update/"}
		jobs := make(chan *common.Metrics, len(metrics))
		results := make(chan *ResultWork, len(metrics))
		for i := 0; i < j.cfg.RateLimit; i++ {
			go j.sendWorker(jobs, results, sendURL.String())
		}
		for _, metric := range metrics {
			jobs <- metric
		}
		close(jobs)
		for i := 0; i < len(metrics); i++ {
			result := <-results
			if result.Err != nil {
				log.Printf("Error send metrics: %s\n", result.Err.Error())
				continue
			}
			if result.StatusCode != 200 {
				log.Printf("Error send metrics: %d\n", result.StatusCode)
				continue
			}
		}
	}
}
