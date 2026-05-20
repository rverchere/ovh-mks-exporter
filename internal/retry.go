package internal

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func retry[T any](attempts int, fn func() (T, error)) (T, error) {
	var zero T
	var err error
	for i := range attempts {
		var result T
		result, err = fn()
		if err == nil {
			return result, nil
		}
		if i < attempts-1 {
			wait := time.Duration(1<<uint(i)) * time.Second
			log.Warnf("attempt %d/%d failed: %v, retrying in %s", i+1, attempts, err, wait)
			time.Sleep(wait)
		}
	}
	return zero, err
}
