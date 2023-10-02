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

func (m *MemStorage) UpdateGauge(gauge MetricGaugeInfo) {
	if metric, ok := m.GetGauge(gauge.Name); ok {
		metric.Value = gauge.Value
	} else {
		m.Gauges = append(m.Gauges, gauge)
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

func (m *MemStorage) GetGauge(name string) (*MetricGaugeInfo, bool) {
	for i, metric := range m.Gauges {
		if metric.Name == name {
			return &m.Gauges[i], true
		}
	}
	return nil, false
}

func (m *MemStorage) GetCounters() ([]MetricCounterInfo, error) {
	return m.Counters, nil
}

func (m *MemStorage) GetGauges() ([]MetricGaugeInfo, error) {
	return m.Gauges, nil
}

type MemoryRepository interface {
	UpdateCounter(counter MetricCounterInfo)
	UpdateGauge(gauge MetricGaugeInfo)
	GetCounter(name string) (*MetricCounterInfo, bool)
	GetGauge(name string) (*MetricGaugeInfo, bool)
	GetCounters() ([]MetricCounterInfo, error)
	GetGauges() ([]MetricGaugeInfo, error)
}
