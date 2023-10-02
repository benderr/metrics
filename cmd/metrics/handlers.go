package metrics

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func SendMetrics(m MetricsReadWrite, client *resty.Client) {
	for name, value := range m.GetCounters() {
		go func(name string, value int) {
			url := fmt.Sprintf("/%v/%v/%v/%v", "update", "counter", name, value)
			fmt.Println("try send counter", url)
			response, err := client.R().Post(url)
			if err != nil {
				fmt.Println("error counter", err)
			} else {
				fmt.Println("finish counter", response.StatusCode())
			}
		}(name, value)
	}

	for name, value := range m.GetGauges() {
		go func(name string, value float64) {
			url := fmt.Sprintf("/%v/%v/%v/%v", "update", "gauge", name, value)
			fmt.Println("try send gauge", url)
			response, err := client.R().Post(url)
			if err != nil {
				fmt.Println("error counter", err)
			} else {
				fmt.Println("finish gauge", response.StatusCode())
			}
		}(name, value)
	}
}
