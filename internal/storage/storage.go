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

type InMemoryMetricRepository struct {
	Counters []MetricCounterInfo
	Gauges   []MetricGaugeInfo
}

func (m *InMemoryMetricRepository) UpdateCounter(counter MetricCounterInfo) error {
	metric, err := m.GetCounter(counter.Name)
	if err != nil {
		return err
	}
	if metric != nil {
		metric.Value += counter.Value
	} else {
		m.Counters = append(m.Counters, counter)
	}
	return nil
}

func (m *InMemoryMetricRepository) UpdateGauge(gauge MetricGaugeInfo) error {
	metric, err := m.GetGauge(gauge.Name)
	if err != nil {
		return err
	}
	if metric != nil {
		metric.Value = gauge.Value
	} else {
		m.Gauges = append(m.Gauges, gauge)
	}
	return nil
}

func (m *InMemoryMetricRepository) GetCounter(name string) (*MetricCounterInfo, error) {
	for i, metric := range m.Counters {
		if metric.Name == name {
			return &m.Counters[i], nil
		}
	}
	return nil, nil
}

func (m *InMemoryMetricRepository) GetGauge(name string) (*MetricGaugeInfo, error) {
	for i, metric := range m.Gauges {
		if metric.Name == name {
			return &m.Gauges[i], nil
		}
	}
	return nil, nil
}

func (m *InMemoryMetricRepository) GetCounters() ([]MetricCounterInfo, error) {
	return m.Counters, nil
}

func (m *InMemoryMetricRepository) GetGauges() ([]MetricGaugeInfo, error) {
	return m.Gauges, nil
}
