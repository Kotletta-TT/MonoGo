package usecase

import (
	"fmt"
	"github.com/Kotletta-TT/MonoGo/internal/server/entity"
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"log"
	"net/http"
	"net/url"
)

var serverHost = "localhost:8080"

type Sender interface {
	Send()
}

type HTTPSender struct {
	repo   repository.Repository
	client *http.Client
}

func NewHTTPSender(repo repository.Repository) HTTPSender {
	return HTTPSender{repo: repo, client: &http.Client{}}
}

func (h *HTTPSender) compileURL(typeMetric, nameMetric, value string) string {
	compileURL := url.URL{Host: serverHost, Scheme: "http"}
	return compileURL.JoinPath("update", typeMetric, nameMetric, value).String()
}

func (h *HTTPSender) Send() {
	allMetrics := h.repo.GetAllMetrics()
	metricURL := ""
	for _, metric := range allMetrics {
		switch metric.GetMetricKind() {
		case entity.KindGauge:
			metricURL = h.compileURL("gauge", metric.Name, fmt.Sprintf("%f", metric.GetGaugeValue()))
		case entity.KindCounter:
			metricURL = h.compileURL("counter", metric.Name, fmt.Sprintf("%d", metric.GetCounterValue()))
		}
		resp, err := h.client.Post(metricURL, "text/plain", nil)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Println("Error resp body close")
			}
		}()
		if err != nil {
			panic(err)
		}
		log.Printf("Status code: %d Metric: %s\n", resp.StatusCode, metric.Name)
	}
}
