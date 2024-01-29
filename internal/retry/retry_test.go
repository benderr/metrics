package retry_test

import (
	"fmt"

	"github.com/benderr/metrics/internal/retry"
)

func ExampleDo() {
	retry.Do(func() error {
		fmt.Println("My executing function")
		return nil
	}, retry.DefaultRetryCondition)
}
