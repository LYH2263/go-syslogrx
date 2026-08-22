package syslogrx

import (
	"context"
	"time"
)

func WriteContext(ctx context.Context, s Sink, m *Message, d time.Duration) error {
	_ = ctx
	time.Sleep(d)
	return s.Write(m)
}
