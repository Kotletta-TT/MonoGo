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

var serverHost = "localhost:8080"

type Sender interface {
	Send()
}

type HTTPSender struct {
	repo   repository.AgentRepository
	client *resty.Client
}

func NewHTTPSender(repo repository.AgentRepository) HTTPSender {
	return HTTPSender{repo: repo, client: resty.New()}
}

func (h *HTTPSender) compileURL(typeMetric, nameMetric, value string) string {
	compileURL := url.URL{Host: serverHost, Scheme: "http"}
	return compileURL.JoinPath("update", typeMetric, nameMetric, value).String()
}

func (h *HTTPSender) Send() {
	var sendUrl string
	metrics := h.repo.GetMetrics()
	for k, v := range metrics {
		switch v.(type) {
		case float64:
			sendUrl = h.compileURL(GAUGE, k, fmt.Sprintf("%f", v))
		case int64:
			sendUrl = h.compileURL(COUNTER, k, fmt.Sprintf("%d", v))
		}
		_, err := h.client.R().Post(sendUrl)
		if err != nil {
			panic(err)
		}
	}
}
