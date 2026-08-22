package syslogrx

func (r *Receiver) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	r.closed = true
	var err error
	if r.sink != nil {
		err = r.sink.Close()
	}
	return err
}

func (r *Receiver) CloseFlushCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return 0
	}
	n := 0
	if r.sink != nil {
		flushed := r.sink.Flush()
		n = len(flushed)
		r.sink.Clear()
		_ = r.sink.Close()
	}
	r.ring = nil
	r.closed = true
	return n
}
