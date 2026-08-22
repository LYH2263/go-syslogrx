package syslogrx

import (
	"context"
	"time"
)

// WriteContext writes m to s after a delay d, but honors ctx cancellation so
// that once ctx is done the write path exits promptly instead of blocking for
// the remainder of an internal sleep.
func WriteContext(ctx context.Context, s Sink, m *Message, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-t.C:
		return s.Write(m)
	case <-ctx.Done():
		return ctx.Err()
	}
}
