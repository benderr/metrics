package storage

type MemType string

const (
	Gauge   MemType = "gauge"
	Counter MemType = "counter"
)

type MetricCounterInfo struct {
	Name  string
	Value int64
}

type MetricGaugeInfo struct {
	Name  string
	Value float64
}

type MemStorage struct {
	Counters []MetricCounterInfo
	Gauges   []MetricGaugeInfo
}

func (m *MemStorage) UpdateCounter(counter MetricCounterInfo) {
	if metric, ok := m.GetCounter(counter.Name); ok {
		metric.Value += counter.Value
	} else {
		m.Counters = append(m.Counters, counter)
	}
}

func (m *MemStorage) GetCounter(name string) (*MetricCounterInfo, bool) {
	for i, metric := range m.Counters {
		if metric.Name == name {
			return &m.Counters[i], true
		}
	}
	return nil, false
}

func (m *MemStorage) UpdateGauge(gauge MetricGaugeInfo) {
	if metric, ok := m.GetGauge(gauge.Name); ok {
		metric.Value = gauge.Value
	} else {
		m.Gauges = append(m.Gauges, gauge)
	}
}

func (m *MemStorage) GetGauge(name string) (*MetricGaugeInfo, bool) {
	for i, metric := range m.Gauges {
		if metric.Name == name {
			return &m.Gauges[i], true
		}
	}
	return nil, false
}

type MemoryRepository interface {
	UpdateCounter(counter MetricCounterInfo)
	UpdateGauge(gauge MetricGaugeInfo)
}
