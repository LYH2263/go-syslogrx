package syslogrx

func (r *Receiver) Recent(n int) []*Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n <= 0 || n > len(r.ring) {
		n = len(r.ring)
	}
	return r.ring[len(r.ring)-n:]
}
