package syslogrx

import "sync"

type Receiver struct {
	mu       sync.Mutex
	closed   bool
	sink     Sink
	ring     []*Message
	capacity int
}

func NewReceiver(capacity int) *Receiver {
	if capacity <= 0 {
		capacity = 64
	}
	return &Receiver{capacity: capacity, ring: make([]*Message, 0, capacity)}
}

func (r *Receiver) SetSink(s Sink) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sink = s
}

func (r *Receiver) RingLen() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.ring)
}
