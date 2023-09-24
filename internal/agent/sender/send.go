package sender

import (
	"github.com/Kotletta-TT/MonoGo/cmd/agent/config"
	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"github.com/Kotletta-TT/MonoGo/internal/shared"
	"github.com/go-resty/resty/v2"
	"io"
	"log"
	"math"
	"net/url"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
	JSON    = "json"
	TEXT    = "text"
)

type metricsStore interface {
	GetMetrics() map[string]*entity.Value
}

type Sender interface {
	Send()
}

type HTTPSender struct {
	repo       metricsStore
	client     *resty.Client
	serverAddr string
}

type TextPlainSender HTTPSender
type JSONSender HTTPSender

func NewHTTPSender(repo metricsStore, cfg *config.Config) Sender {
	switch cfg.SendType {
	case JSON:
		return &JSONSender{repo: repo, client: resty.New(), serverAddr: cfg.ServerHost}
	case TEXT:
		return &TextPlainSender{repo: repo, client: resty.New(), serverAddr: cfg.ServerHost}
	default:
		panic("Send type unknown")
	}
}

func (h *TextPlainSender) compileURL(nameMetric string, valueMetric *entity.Value) string {
	compileURL := url.URL{Host: h.serverAddr, Scheme: "http"}
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

func JSONMetricFabric(name string, value *entity.Value) *shared.Metrics {
	m := shared.NewMetrics()
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

func (j *JSONSender) Send() {
	log.Println("Start JSON send metrics")
	sendURL := url.URL{Host: j.serverAddr, Scheme: "http", Path: "/update/"}
	metrics := j.repo.GetMetrics()
	log.Printf("Send URL: %s\n", sendURL.String())
	for k, v := range metrics {
		m := JSONMetricFabric(k, v)
		mJSON, err := m.MarshalJSON()
		if err != nil {
			panic(err)
		}
		log.Printf("Send JSON: %s\n", mJSON)
		resp, err := j.client.R().SetHeader("Content-Type", "application/json").SetBody(mJSON).Post(sendURL.String())
		if err != nil && resp.StatusCode() != 200 && err != io.EOF {
			log.Printf("error: Code: %d, Body: %s err: %s\n", resp.StatusCode(), resp.String(), err.Error())
		}
	}
}
