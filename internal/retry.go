package internal

import (
	"errors"
	"time"

	"github.com/ovh/go-ovh/ovh"
	log "github.com/sirupsen/logrus"
)

var sleepFn = time.Sleep

func isRetryable(err error) bool {
	var apiErr *ovh.APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code >= 500
	}
	return true
}

func retry[T any](attempts int, fn func() (T, error)) (T, error) {
	var zero T
	var err error
	for i := range attempts {
		var result T
		result, err = fn()
		if err == nil {
			return result, nil
		}
		if !isRetryable(err) {
			return zero, err
		}
		if i < attempts-1 {
			wait := time.Duration(1<<uint(i)) * time.Second
			log.Warnf("attempt %d/%d failed: %v, retrying in %s", i+1, attempts, err, wait)
			sleepFn(wait)
		}
	}
	return zero, err
}
