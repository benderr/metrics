package metrics

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
)

func sendMetrics(m *Metrics, client *resty.Client) {
	for name, value := range m.Counters {
		go func(name string, value int) {
			url := fmt.Sprintf("/%v/%v/%v/%v", "update", "counter", name, value)
			client.R().Post(url)
		}(name, value)
	}

	for name, value := range m.Gauges {
		go func(name string, value float64) {
			url := fmt.Sprintf("/%v/%v/%v/%v", "update", "gauge", name, value)
			//fmt.Println("try send gauge", url)
			response, err := client.R().Post(url)
			if err != nil {
				//fmt.Println("error counter", err)
			} else {
				fmt.Println("finish gauge", response.StatusCode())
			}
		}(name, value)
	}
}

func updateReport(m *Metrics) {
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

func StartAgent(poolInterval int, reportInterval int, server string) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		pollTicker := time.NewTicker(time.Second * time.Duration(poolInterval))
		reportTicker := time.NewTicker(time.Second * time.Duration(reportInterval))

		defer pollTicker.Stop()
		defer reportTicker.Stop()

		report := &Metrics{
			Gauges:   make(map[string]float64),
			Counters: make(map[string]int),
		}

		client := resty.New().SetBaseURL(server)

		for {
			select {
			case <-pollTicker.C:
				updateReport(report)
			case <-reportTicker.C:
				sendMetrics(report, client)
			}
		}
	}()

	return done
}
