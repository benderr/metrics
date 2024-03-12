package retry_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/benderr/metrics/pkg/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleDo() {
	retry.Do(func() error {
		fmt.Println("My executing function")
		return nil
	}, retry.DefaultRetryCondition)
}

func TestDo_simple(t *testing.T) {
	flag := false
	retry.Do(func() error {
		flag = true
		return nil
	}, retry.DefaultRetryCondition)
	assert.Equal(t, true, flag)
}

func TestDo_second(t *testing.T) {
	retries := 0
	retry.Do(func() error {
		retries += 1
		if retries == 2 {
			return nil
		}
		return errors.New("failed")
	}, func(attempt int, err error) bool {
		return attempt <= 2
	})
	assert.Equal(t, 2, retries)
}

func TestDoWithValue(t *testing.T) {
	retries := 0
	v, err := retry.DoWithValue(func() (string, error) {
		retries += 1
		if retries == 2 {
			return "success", nil
		}
		return "error", errors.New("failed")
	}, func(attempt int, err error) bool {
		return attempt <= 2
	})

	require.NoError(t, err)
	assert.Equal(t, v, "success")
}
