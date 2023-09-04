package usecase

import (
	"fmt"
	"github.com/Kotletta-TT/MonoGo/internal/agent/infrastructure/repository"
	"github.com/go-resty/resty/v2"
	"net/url"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

type Sender interface {
	Send()
}

type HTTPSender struct {
	repo       repository.AgentRepository
	client     *resty.Client
	serverAddr string
}

func NewHTTPSender(repo repository.AgentRepository, serverHost string) HTTPSender {
	return HTTPSender{repo: repo, client: resty.New(), serverAddr: serverHost}
}

func (h *HTTPSender) compileURL(typeMetric, nameMetric, value string) string {
	compileURL := url.URL{Host: h.serverAddr, Scheme: "http"}
	return compileURL.JoinPath("update", typeMetric, nameMetric, value).String()
}

func (h *HTTPSender) Send() {
	var sendURL string
	metrics := h.repo.GetMetrics()
	for k, v := range metrics {
		switch v.(type) {
		case float64:
			sendURL = h.compileURL(GAUGE, k, fmt.Sprintf("%f", v))
		case int64:
			sendURL = h.compileURL(COUNTER, k, fmt.Sprintf("%d", v))
		}
		_, err := h.client.R().Post(sendURL)
		if err != nil {
			panic(err)
		}
	}
}
