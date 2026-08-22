package syslogrx

func (r *Receiver) Recent(n int) []*Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n <= 0 || n > len(r.ring) {
		n = len(r.ring)
	}
	// Return isolated copies so callers cannot mutate the in-ring
	// originals (e.g. redacting RawBytes for a preview). Sharing the
	// ring's *Message pointers let preview edits corrupt the ring and
	// any later snapshot/search view of the same entries.
	out := make([]*Message, n)
	for i, m := range r.ring[len(r.ring)-n:] {
		out[i] = cloneMsg(m)
	}
	return out
}
