package syslogrx

import "github.com/LYH2263/go-syslogrx/internal/clone"

// SnapshotRecent returns deep copies of ring messages for admin export.
func (r *Receiver) SnapshotRecent() []*Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*Message, len(r.ring))
	for i, m := range r.ring {
		cp := *m
		cp.RawBytes = clone.Bytes(m.RawBytes)
		out[i] = &cp
	}
	return out
}
