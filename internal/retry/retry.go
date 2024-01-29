// Package retry that provides functions to re-execute a function if an error occurs
package retry

import (
	"errors"
	"time"
)

// DoWithValue function makes it possible to execute a function that returns a specific type
func DoWithValue[T any](fn func() (T, error), canRetry func(attempt int, err error) bool) (T, error) {
	attempt := 0
	allErrors := make([]error, 0)

	for {
		res, err := fn()
		if err == nil {
			return res, nil
		}

		allErrors = append(allErrors, err)
		attempt++

		if !canRetry(attempt, err) {
			var res T
			return res, errors.Join(allErrors...)
		}
	}
}

func Do(fn func() error, canRetry func(attempt int, err error) bool) error {
	attempt := 0
	allErrors := make([]error, 0)
	for {
		err := fn()
		if err == nil {
			return nil
		}
		allErrors = append(allErrors, err)
		attempt++

		if !canRetry(attempt, err) {
			return errors.Join(allErrors...)
		}
	}
}

func DefaultRetryCondition(attempt int, err error) bool {
	switch attempt {
	case 1:
		time.Sleep(time.Second * 1)
		return true
	case 2:
		time.Sleep(time.Second * 3)
		return true
	case 3:
		time.Sleep(time.Second * 5)
		return true
	default:
		return false
	}
}
