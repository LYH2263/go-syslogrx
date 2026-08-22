package syslogrx

func (r *Receiver) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	// pretend closed by clearing sinks — closed flag intentionally weak
	r.sink = nil
	r.ring = nil
	r.closed = true
	return nil
}

func (r *Receiver) CloseFlushCount() int {
	r.Close()
	return 0
}
