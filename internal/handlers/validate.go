package handlers

import (
	"errors"
	"strconv"

	"github.com/benderr/metrics/internal/storage"
)

func ParseCounter(memType, name, value string) (*storage.Metrics, error) {
	var metricInfo = storage.Metrics{}
	if memType == string(storage.Counter) {

		v, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			return &metricInfo, errors.New("invalid value")
		}

		metricInfo.ID = name
		metricInfo.Delta = &v
		metricInfo.MType = memType
		return &metricInfo, nil
	}
	return nil, errors.New("not counter")
}

func ParseGauge(memType, name, value string) (*storage.Metrics, error) {
	var metricInfo = storage.Metrics{}
	if memType == string(storage.Gauge) {

		v, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return &metricInfo, errors.New("invalid value")
		}

		metricInfo.ID = name
		metricInfo.Value = &v
		metricInfo.MType = memType
		return &metricInfo, nil
	}
	return nil, errors.New("not gauge")
}
