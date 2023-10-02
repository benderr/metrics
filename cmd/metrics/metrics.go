package metrics

import (
	"math/rand"
	"runtime"
)

type Metrics struct {
	Counter map[string]int
	Gauges  map[string]float64
}

func (m *Metrics) Write() {
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

func (m *Metrics) InscrementCounter(name string, value int) {
	if v, ok := m.Counter[name]; ok {
		m.Counter[name] = v + value
	} else {
		m.Counter[name] = v
	}
}

func (m *Metrics) UpdateGauge(name string, value float64) {
	m.Gauges[name] = value
}

func (m *Metrics) GetCounters() map[string]int {
	return m.Counter
}

func (m *Metrics) GetGauges() map[string]float64 {
	return m.Gauges
}

type MetricsReadWrite interface {
	Write()
	GetCounters() map[string]int
	GetGauges() map[string]float64
}
