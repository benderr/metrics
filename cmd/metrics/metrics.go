package metrics

type Metrics struct {
	Counters map[string]int
	Gauges   map[string]float64
}

func (m *Metrics) InscrementCounter(name string, value int) {
	if v, ok := m.Counters[name]; ok {
		m.Counters[name] = v + value
	} else {
		m.Counters[name] = value
	}
}

func (m *Metrics) UpdateGauge(name string, value float64) {
	m.Gauges[name] = value
}
