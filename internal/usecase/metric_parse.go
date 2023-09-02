package usecase

import (
	"fmt"
	"github.com/Kotletta-TT/MonoGo/internal/entity"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

func validateMetric(partsMetric []string) (string, string, any, error) {
	typeMetric := partsMetric[1]
	if len(partsMetric[2]) < 1 {
		return "", "", "", NoNameMetric{}
	}
	nameMetric := partsMetric[2]
	switch strings.ToLower(typeMetric) {
	case GAUGE:
		valueMetric, err := strconv.ParseFloat(partsMetric[3], 64)
		if err != nil {
			return "", "", "", IncorrectValueMetrics{Type: typeMetric, Value: partsMetric[3], Err: err.Error()}
		}
		return typeMetric, nameMetric, valueMetric, nil
	case COUNTER:
		valueMetric, err := strconv.ParseInt(partsMetric[3], 10, 64)
		if err != nil {
			return "", "", "", IncorrectValueMetrics{Type: typeMetric, Value: partsMetric[3], Err: err.Error()}
		}
		return typeMetric, nameMetric, valueMetric, nil
	default:
		return "", "", "", IncorrectTypeMetrics{Type: typeMetric}
	}
}

func parseUrl(url *url.URL) (typeMetric, nameMetric string, valueMetric any, err error) {
	trimUrl := strings.Trim(url.Path, "/")
	partsUrl := strings.Split(trimUrl, "/")
	if len(partsUrl) != 4 {
		return "", "", "", fmt.Errorf("parsing err: expected len == 4 got: %d", len(partsUrl))
	}
	return validateMetric(partsUrl)
}

func Parse(req *http.Request) (*entity.Metric, error) {
	typeMetric, nameMetric, valueMetric, err := parseUrl(req.URL)
	if err != nil {
		return nil, err
	}
	fmt.Println(typeMetric, nameMetric, valueMetric)
	return &entity.Metric{}, nil
}
