package syslogrx

// SnapshotRecent returns an isolated snapshot of the ring.
//
// Each returned Message (and its RawBytes) is a deep copy, so mutating the
// return value never pollutes the in-ring records, and the snapshot stays
// independent of the live ring as later packets are handled. This decouples
// the snapshot from the caller's raw buffer: the next packet received into a
// reused raw buffer cannot dirty either the ring or a prior snapshot.
func (r *Receiver) SnapshotRecent() []*Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*Message, len(r.ring))
	for i, m := range r.ring {
		out[i] = cloneMsg(m)
	}
	return out
}
