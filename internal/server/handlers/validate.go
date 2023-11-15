package handlers

import (
	"errors"
	"strconv"

	"github.com/benderr/metrics/internal/server/repository"
)

func ParseCounter(memType, name, value string) (*repository.Metrics, error) {
	var metricInfo = repository.Metrics{}
	if memType == "counter" {

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

func ParseGauge(memType, name, value string) (*repository.Metrics, error) {
	var metricInfo = repository.Metrics{}
	if memType == "gauge" {

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
