package syslogrx

import (
	"context"

	"github.com/LYH2263/go-syslogrx/internal/clone"
)

func cloneMsg(m *Message) *Message {
	if m == nil {
		return nil
	}
	out := *m
	out.RawBytes = clone.Bytes(m.RawBytes)
	return &out
}

func (r *Receiver) Handle(ctx context.Context, raw []byte) (*Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil, ErrClosed
	}
	if r.sink == nil {
		return nil, ErrNoSink
	}
	m, err := ParseLine(string(raw))
	if err != nil {
		return nil, err
	}
	m.RawBytes = clone.Bytes(raw)
	// Sink write must succeed before the message enters the ring; otherwise a
	// half-written/failed record pollutes Recent() and drives alert noise.
	if err := r.sink.Write(cloneMsg(m)); err != nil {
		return nil, err
	}
	cp := cloneMsg(m)
	r.ring = append(r.ring, cp)
	if len(r.ring) > r.capacity {
		r.ring = r.ring[len(r.ring)-r.capacity:]
	}
	return cloneMsg(m), nil
}
