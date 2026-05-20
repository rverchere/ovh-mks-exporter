package internal

import (
	"errors"
	"testing"
	"time"

	"github.com/ovh/go-ovh/ovh"
)

func init() {
	sleepFn = func(time.Duration) {}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"generic error", errors.New("network timeout"), true},
		{"5xx api error", &ovh.APIError{Code: 500, Message: "server error"}, true},
		{"503 api error", &ovh.APIError{Code: 503, Message: "unavailable"}, true},
		{"404 api error", &ovh.APIError{Code: 404, Message: "not found"}, false},
		{"400 api error", &ovh.APIError{Code: 400, Message: "bad request"}, false},
		{"403 api error", &ovh.APIError{Code: 403, Message: "forbidden"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRetryable(tt.err); got != tt.want {
				t.Errorf("isRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetry_ImmediateSuccess(t *testing.T) {
	calls := 0
	got, err := retry(3, func() (string, error) {
		calls++
		return "ok", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != "ok" {
		t.Errorf("got %q, want %q", got, "ok")
	}
	if calls != 1 {
		t.Errorf("calls = %d, want 1", calls)
	}
}

func TestRetry_SuccessAfterTransientError(t *testing.T) {
	calls := 0
	got, err := retry(3, func() (int, error) {
		calls++
		if calls < 3 {
			return 0, errors.New("transient")
		}
		return 42, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
}

func TestRetry_StopsOnClientError(t *testing.T) {
	calls := 0
	_, err := retry(5, func() (string, error) {
		calls++
		return "", &ovh.APIError{Code: 404, Message: "not found"}
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if calls != 1 {
		t.Errorf("calls = %d, want 1", calls)
	}
}

func TestRetry_ExhaustsAttempts(t *testing.T) {
	calls := 0
	_, err := retry(3, func() (string, error) {
		calls++
		return "", errors.New("always fails")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
}
