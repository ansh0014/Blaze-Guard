package shared

import (
	"context"
	"time"
)

func DoWithRetry(ctx context.Context, attempts int, baseDelay time.Duration, fn func(context.Context) error) error {
	if attempts < 1 {
		attempts = 1
	}
	if baseDelay <= 0 {
		baseDelay = 300 * time.Millisecond
	}

	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := fn(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(baseDelay * time.Duration(1<<i)):
		}
	}
	return lastErr
}
