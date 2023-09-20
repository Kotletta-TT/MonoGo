package sender

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"github.com/go-resty/resty/v2"
	"log"
	"net/url"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
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

func NewHTTPSender(repo metricsStore, serverHost string) HTTPSender {
	return HTTPSender{repo: repo, client: resty.New(), serverAddr: serverHost}
}

func (h *HTTPSender) compileURL(nameMetric string, valueMetric *entity.Value) string {
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

func (h *HTTPSender) Send() {
	var sendURL string
	metrics := h.repo.GetMetrics()
	for k, v := range metrics {
		sendURL = h.compileURL(k, v)
		_, err := h.client.R().Post(sendURL)
		if err != nil {
			log.Println(sendURL)
			panic(err)
		}
	}
}
