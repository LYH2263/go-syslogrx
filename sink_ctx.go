package syslogrx

import (
	"context"
	"time"
)

func WriteContext(ctx context.Context, s Sink, m *Message, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return s.Write(m)
	}
}
