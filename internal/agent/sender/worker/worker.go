package worker

import (
	"sync"

	"github.com/benderr/metrics/internal/agent/report"
)

type WorkerFunc func(*report.MetricItem) error

func Run(rateLimit int, metrics []report.MetricItem, fn WorkerFunc) []error {
	count := len(metrics)
	if count == 0 {
		return nil
	}

	jobs := make(chan *report.MetricItem, count)
	results := make(chan error, count)
	wg := &sync.WaitGroup{}

	//если указали 0 то запускаем один воркер
	limit := rateLimit
	if limit == 0 {
		limit = 1
	}

	wg.Add(limit)

	for i := 0; i < limit; i++ {
		go worker(i, wg, jobs, results, fn)
	}

	for _, m := range metrics {
		jobs <- &m
	}

	close(jobs)

	wg.Wait()
	close(results)

	allErrors := make([]error, 0)

	for err := range results {
		if err != nil {
			allErrors = append(allErrors, err)
		}
	}

	return allErrors
}

func worker(id int, wg *sync.WaitGroup, jobs <-chan *report.MetricItem, results chan<- error, fn WorkerFunc) {
	defer wg.Done()
	for m := range jobs {
		err := fn(m)
		results <- err
	}
}
