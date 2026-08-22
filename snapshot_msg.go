package syslogrx

func (r *Receiver) SnapshotRecent() []*Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*Message, len(r.ring))
	for i, m := range r.ring {
		out[i] = m
	}
	return out
}
