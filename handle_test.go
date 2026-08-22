package syslogrx

import (
	"context"
	"errors"
	"testing"
)

// validLine is a syslog line ParseLine accepts (RFC3164).
const validLine = "<13>Oct 11 22:14:15 mymachine app: hello"

func TestHandleNoSinkReturnsErrNoSink(t *testing.T) {
	r := NewReceiver(8)

	// No SetSink call. Previously this path panicked on r.sink.Write and
	// left a half-record in the ring; it must now return a checkable error
	// and leave the ring empty.
	m, err := r.Handle(context.Background(), []byte(validLine))
	if !errors.Is(err, ErrNoSink) {
		t.Fatalf("Handle without sink: err = %v, want errors.Is ErrNoSink", err)
	}
	if m != nil {
		t.Fatalf("Handle without sink: msg = %+v, want nil", m)
	}
	if got := r.RingLen(); got != 0 {
		t.Fatalf("ring polluted by no-sink path: RingLen = %d, want 0", got)
	}
	// Recent must be empty too — no dirty half-records for alerting to surface.
	if recs := r.Recent(10); len(recs) != 0 {
		t.Fatalf("Recent after no-sink path: %d records, want 0", len(recs))
	}
}

func TestHandleSinkWriteFailureDoesNotEnterRing(t *testing.T) {
	r := NewReceiver(8)
	r.SetSink(NewMemSink())
	r.sink.(*MemSink).SetFail(true)

	_, err := r.Handle(context.Background(), []byte(validLine))
	if err == nil {
		t.Fatal("Handle with failing sink: err = nil, want non-nil")
	}
	if got := r.RingLen(); got != 0 {
		t.Fatalf("ring polluted by failed write: RingLen = %d, want 0", got)
	}
}

func TestHandleSuccessEntersRingAndSink(t *testing.T) {
	r := NewReceiver(8)
	sink := NewMemSink()
	r.SetSink(sink)

	m, err := r.Handle(context.Background(), []byte(validLine))
	if err != nil {
		t.Fatalf("Handle: err = %v", err)
	}
	if m == nil || m.Tag != "app" {
		t.Fatalf("Handle: msg = %+v, want tag app", m)
	}
	if got := r.RingLen(); got != 1 {
		t.Fatalf("RingLen = %d, want 1", got)
	}
	if got := SinkPending(sink); got != 1 {
		t.Fatalf("sink pending = %d, want 1", got)
	}
}
