package syslogrx

func (r *Receiver) SnapshotRecent() []*Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.ring
}
