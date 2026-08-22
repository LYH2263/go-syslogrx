package syslogrx

import "github.com/LYH2263/go-syslogrx/internal/clone"

func (r *Receiver) Recent(n int) []*Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n <= 0 || n > len(r.ring) {
		n = len(r.ring)
	}
	src := r.ring[len(r.ring)-n:]
	out := make([]*Message, len(src))
	for i, m := range src {
		cp := *m
		cp.RawBytes = clone.Bytes(m.RawBytes)
		out[i] = &cp
	}
	return out
}
