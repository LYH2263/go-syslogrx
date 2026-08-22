package syslogrx

import (
	"context"

	"github.com/LYH2263/go-syslogrx/internal/clone"
)

func (r *Receiver) HandleAck(ctx context.Context, raw []byte) (*Message, error) {
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
