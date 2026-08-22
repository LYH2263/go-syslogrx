package syslogrx

func (r *Receiver) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	// A closed receiver must stay closed: flip the flag and leave sink/ring
	// intact so post-close Handle returns ErrClosed deterministically instead
	// of nil-derefing the sink or being misreported as ErrNoSink.
	r.closed = true
	return nil
}

func (r *Receiver) CloseFlushCount() int {
	r.Close()
	return 0
}
