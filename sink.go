package syslogrx

import (
	"fmt"
	"sync"
)

type Sink interface {
	Write(m *Message) error
	Flush() []*Message
	Clear()
	Close() error
}

type MemSink struct {
	mu      sync.Mutex
	items   []*Message
	fail    bool
	closed  bool
	flushed int
}

func NewMemSink() *MemSink {
	return &MemSink{}
}

func (s *MemSink) SetFail(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fail = v
}

func (s *MemSink) Write(m *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("sink closed")
	}
	if s.fail {
		return fmt.Errorf("sink write failed")
	}
	cp := *m
	cp.RawBytes = append([]byte(nil), m.RawBytes...)
	s.items = append(s.items, &cp)
	return nil
}

func (s *MemSink) Flush() []*Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*Message, len(s.items))
	copy(out, s.items)
	s.flushed += len(out)
	return out
}

func (s *MemSink) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = nil
}

func (s *MemSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

func (s *MemSink) Pending() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items)
}
