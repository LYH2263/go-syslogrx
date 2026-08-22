package syslogrx

import (
	"context"
	"errors"
	"testing"
)

// validLine is a well-formed RFC5424 syslog line that ParseLine accepts.
const validLine = "<34>1 2023-01-02T15:04:05Z myhost myapp 1234 - hello"

// TestHandleAck_SinkFailurePropagates guards the ack semantics: when the
// downstream sink cannot persist the message, HandleAck must NOT return a
// success Message to the upstream. Returning success here would let the
// upstream advance its cursor and silently drop the log on the floor, which
// reconciliation later finds missing on disk. The write error must surface.
func TestHandleAck_SinkFailurePropagates(t *testing.T) {
	ctx := context.Background()
	r := NewReceiver(8)
	s := NewMemSink()
	r.SetSink(s)

	before := r.RingLen()
	s.SetFail(true)

	m, err := r.HandleAck(ctx, []byte(validLine))
	if err == nil {
		t.Fatalf("HandleAck: expected error when sink write fails, got nil")
	}
	if m != nil {
		t.Fatalf("HandleAck: expected nil Message on sink failure, got %+v", m)
	}

	// The rolled-back entry must not linger in the recent ring: the upstream
	// was told the write failed, so the message must not appear as accepted.
	if got := r.RingLen(); got != before {
		t.Fatalf("ring length after failed ack = %d, want %d (entry must roll back)", got, before)
	}
}

// TestHandleAck_SinkSuccessAcknowledges is the positive control: a successful
// sink write yields an acknowledged Message and records it in the ring.
func TestHandleAck_SinkSuccessAcknowledges(t *testing.T) {
	ctx := context.Background()
	r := NewReceiver(8)
	s := NewMemSink()
	r.SetSink(s)

	before := r.RingLen()
	m, err := r.HandleAck(ctx, []byte(validLine))
	if err != nil {
		t.Fatalf("HandleAck: unexpected error on success: %v", err)
	}
	if m == nil {
		t.Fatalf("HandleAck: expected acknowledged Message, got nil")
	}
	if got := r.RingLen(); got != before+1 {
		t.Fatalf("ring length after ack = %d, want %d", got, before+1)
	}
	if got := s.Pending(); got != 1 {
		t.Fatalf("sink pending = %d, want 1", got)
	}
}

// TestHandleAck_ParseErrorUnchanged ensures input-validation errors still
// surface and never reach the sink.
func TestHandleAck_ParseErrorUnchanged(t *testing.T) {
	ctx := context.Background()
	r := NewReceiver(8)
	s := NewMemSink()
	r.SetSink(s)

	_, err := r.HandleAck(ctx, []byte("not-a-syslog-line"))
	if err == nil {
		t.Fatalf("HandleAck: expected parse error, got nil")
	}
	if got := s.Pending(); got != 0 {
		t.Fatalf("sink pending = %d, want 0 on parse error", got)
	}
}

// errSink is a Sink whose Write always fails, for error-chain assertions.
type errSink struct{ msg string }

func (s *errSink) Write(*Message) error { return errors.New(s.msg) }
func (s *errSink) Flush() []*Message     { return nil }
func (s *errSink) Clear()                {}
func (s *errSink) Close() error          { return nil }

// TestHandleAck_SinkErrorChains confirms the original sink error reaches the
// caller verbatim rather than being swallowed into a generic success.
func TestHandleAck_SinkErrorChains(t *testing.T) {
	ctx := context.Background()
	r := NewReceiver(8)
	r.SetSink(&errSink{msg: "disk full"})

	_, err := r.HandleAck(ctx, []byte(validLine))
	if err == nil || err.Error() != "disk full" {
		t.Fatalf("HandleAck: expected sink error %q, got %v", "disk full", err)
	}
}
