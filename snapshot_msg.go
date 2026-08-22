package syslogrx

func (r *Receiver) SnapshotRecent() []*Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*Message, len(r.ring))
	// Clone each entry; copying the *Message pointers verbatim would
	// alias the ring's messages, so a caller redacting RawBytes on the
	// snapshot would also mutate the in-ring original.
	for i, m := range r.ring {
		out[i] = cloneMsg(m)
	}
	return out
}
