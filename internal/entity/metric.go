package entity

import "time"

type Metric struct {
	Name    string // TODO При ненужности данного поля - удалить
	Gauge   []Gauge
	Counter []Counter
}

type Gauge struct {
	Timestamp time.Time
	Value     float64
	Tags      map[string]string
}

type Counter struct {
	Timestamp time.Time
	Value     int64
	Tags      map[string]string
}
