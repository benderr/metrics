package report

import (
	"math/rand"
	"runtime"
)

type Report struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

type MetricItem struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func New() *Report {
	return &Report{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

func (m *Report) InscrementCounter(name string, value int64) {
	if v, ok := m.Counters[name]; ok {
		m.Counters[name] = v + value
	} else {
		m.Counters[name] = value
	}
}

func (m *Report) UpdateGauge(name string, value float64) {
	m.Gauges[name] = value
}

func (m *Report) UpdateReport() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.UpdateGauge("Alloc", float64(rtm.Alloc))
	m.UpdateGauge("BuckHashSys", float64(rtm.BuckHashSys))
	m.UpdateGauge("Frees", float64(rtm.Frees))
	m.UpdateGauge("GCCPUFraction", float64(rtm.GCCPUFraction))
	m.UpdateGauge("GCSys", float64(rtm.GCSys))
	m.UpdateGauge("HeapAlloc", float64(rtm.HeapAlloc))
	m.UpdateGauge("HeapIdle", float64(rtm.HeapIdle))
	m.UpdateGauge("HeapInuse", float64(rtm.HeapInuse))
	m.UpdateGauge("HeapObjects", float64(rtm.HeapObjects))
	m.UpdateGauge("HeapReleased", float64(rtm.HeapReleased))
	m.UpdateGauge("HeapSys", float64(rtm.HeapSys))
	m.UpdateGauge("LastGC", float64(rtm.LastGC))
	m.UpdateGauge("Lookups", float64(rtm.Lookups))
	m.UpdateGauge("MCacheInuse", float64(rtm.MCacheInuse))
	m.UpdateGauge("MCacheSys", float64(rtm.MCacheSys))
	m.UpdateGauge("MSpanInuse", float64(rtm.MSpanInuse))
	m.UpdateGauge("MSpanSys", float64(rtm.MSpanSys))
	m.UpdateGauge("Mallocs", float64(rtm.Mallocs))
	m.UpdateGauge("NextGC", float64(rtm.NextGC))
	m.UpdateGauge("NumForcedGC", float64(rtm.NumForcedGC))
	m.UpdateGauge("NumGC", float64(rtm.NumGC))
	m.UpdateGauge("OtherSys", float64(rtm.OtherSys))
	m.UpdateGauge("PauseTotalNs", float64(rtm.PauseTotalNs))
	m.UpdateGauge("StackInuse", float64(rtm.StackInuse))
	m.UpdateGauge("StackSys", float64(rtm.StackSys))
	m.UpdateGauge("Sys", float64(rtm.Sys))
	m.UpdateGauge("TotalAlloc", float64(rtm.TotalAlloc))
	m.UpdateGauge("RandomValue", rand.Float64())
	m.InscrementCounter("PollCount", 1)
}
