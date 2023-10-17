package validate

import (
	"errors"
	"strconv"

	"github.com/benderr/metrics/internal/storage"
)

func ParseCounter(memType, name, value string) (storage.MetricCounterInfo, error) {
	var metricInfo = storage.MetricCounterInfo{}
	if memType == string(storage.Counter) {

		v, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			return metricInfo, errors.New("invalid value")
		}

		metricInfo.Name = name
		metricInfo.Value = v
		return metricInfo, nil
	}
	return metricInfo, errors.New("not counter")
}

func ParseGauge(memType, name, value string) (storage.MetricGaugeInfo, error) {
	var metricInfo = storage.MetricGaugeInfo{}
	if memType == string(storage.Gauge) {

		v, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return metricInfo, errors.New("invalid value")
		}

		metricInfo.Name = name
		metricInfo.Value = v
		return metricInfo, nil
	}
	return metricInfo, errors.New("not gauge")
}
