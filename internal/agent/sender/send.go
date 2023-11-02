package sender

import (
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

type HTTPSender struct {
	repo   metricsStore
	client *resty.Client
	cfg    *config.Config
}

type TextPlainSender HTTPSender
type JSONSender HTTPSender

func NewRestyClient() *resty.Client {

	client := resty.New()
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)
	return client
}

func NewHTTPSender(repo metricsStore, cfg *config.Config) Sender {
	switch cfg.SendType {
	case JSON:
		return &JSONSender{repo: repo, client: NewRestyClient(), cfg: cfg}
	case TEXT:
		return &TextPlainSender{repo: repo, client: NewRestyClient(), cfg: cfg}
	default:
		panic("Send type unknown")
	}
}

func (h *TextPlainSender) compileURL(nameMetric string, valueMetric *entity.Value) string {
	compileURL := url.URL{Host: h.cfg.ServerHost, Scheme: "http"}
	switch valueMetric.Kind {
	case entity.KindGauge:
		return compileURL.JoinPath("update", GAUGE, nameMetric, valueMetric.String()).String()
	case entity.KindCounter:
		return compileURL.JoinPath("update", COUNTER, nameMetric, valueMetric.String()).String()
	default:
		panic("Metric type unknown")
	}
}

func (h *TextPlainSender) Send() {
	log.Println("Start Text/Plain send metrics")
	var sendURL string
	metrics := h.repo.GetMetrics()
	for k, v := range metrics {
		sendURL = h.compileURL(k, v)
		log.Printf("Send URL: %s", sendURL)
		_, err := h.client.R().Post(sendURL)
		if err != nil {
			log.Println(sendURL)
			panic(err)
		}
	}
}

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

func (j *JSONSender) reciveResponse(resp *resty.Response, err error) {
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
		j.reciveResponse(req.Post(sendURL.String()))
	default:
		metrics := j.repo.GetMetrics()
		sendURL = url.URL{Host: j.cfg.ServerHost, Scheme: "http", Path: "/update/"}
		for k, v := range metrics {
			m := JSONMetricFabric(k, v)
			mJSON, sign, err := j.prepareBody(m)
			if err != nil {
				log.Printf("prepare body err: %s\n", err.Error())
				continue
			}
			if j.cfg.HashKey != "" && sign != "" {
				req.SetHeader("HashSHA256", sign)
			}
			req.SetBody(mJSON)
			j.reciveResponse(req.Post(sendURL.String()))
		}
	}
}
