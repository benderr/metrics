package metrics

import (
	"math/rand"
	"runtime"
)

type Metrics struct {
	Counter map[string]int
	Gauges  map[string]float64
}

func (m *Metrics) Update() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	m.InscrementCounter("PollCount", 1)
	m.UpdateGauge("RandomValue", rand.Float64())
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

type IMetrics interface {
	Update()
	GetCounters() map[string]int
	GetGauges() map[string]float64
}
