package syslogrx

func (r *Receiver) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	r.closed = true
	if r.sink != nil {
		r.sink.Clear()
		return r.sink.Close()
	}
	return nil
}

func (r *Receiver) CloseFlushCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return 0
	}
	if r.sink != nil {
		r.sink.Clear()
		flushed := r.sink.Flush()
		_ = r.sink.Close()
		r.ring = nil
		r.closed = true
		return len(flushed)
	}
	r.closed = true
	return 0
}
