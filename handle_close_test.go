package syslogrx

import (
	"context"
	"errors"
	"testing"
)

// validRFC3164 is a well-formed syslog line accepted by ParseLine.
const validRFC3164 = "<34>Jan  2 15:04:05 myhost app: boom"

// TestHandleAfterCloseStableErrClosed pins the contract that, once a
// Receiver is closed, every subsequent delivery (the "tail packet" that
// races in after shutdown) returns ErrClosed — deterministically and
// without panicking.
//
// Three regressions this guards against:
//
//  1. panic: Close used to nil r.sink, so Handle's r.sink.Write nil-derefed.
//  2. wrong error class: the post-close error was whatever sink.Write did
//     (or a panic), never ErrClosed, so upstream redelivery could not
//     terminate on errors.Is(err, ErrClosed).
//  3. ErrNoSink masquerading as closed: because Close cleared the sink,
//     HandleAck reported ErrNoSink for a *closed* receiver, so the
//     redelivery loop could not close out on the closed state and retried
//     forever.
func TestHandleAfterCloseStableErrClosed(t *testing.T) {
	r := NewReceiver(8)
	r.SetSink(NewMemSink())

	if err := r.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// The sink and ring must survive Close: closing is a state flag, not a
	// field wipe. A nil sink after Close is exactly the bug that let Handle
	// nil-deref and let HandleAck misreport ErrNoSink.
	r.mu.Lock()
	if r.sink == nil {
		t.Fatal("Close nilled sink; closed flag must not destroy the sink")
	}
	if r.ring == nil {
		t.Fatal("Close nilled ring; closed flag must not destroy the ring")
	}
	r.mu.Unlock()

	// Handle must not panic and must return ErrClosed for the tail packet.
	msg, err := r.Handle(context.Background(), []byte(validRFC3164))
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("Handle after close: want ErrClosed, got %v (msg=%v)", err, msg)
	}
	if msg != nil {
		t.Fatalf("Handle after close returned non-nil message: %v", msg)
	}

	// HandleAck must also report ErrClosed (not ErrNoSink) for a closed recv.
	amsg, aerr := r.HandleAck(context.Background(), []byte(validRFC3164))
	if !errors.Is(aerr, ErrClosed) {
		t.Fatalf("HandleAck after close: want ErrClosed, got %v (msg=%v)", aerr, amsg)
	}
	if amsg != nil {
		t.Fatalf("HandleAck after close returned non-nil message: %v", amsg)
	}

	// Closing again is idempotent and must keep the closed invariant.
	if err := r.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
	if _, err := r.Handle(context.Background(), []byte(validRFC3164)); !errors.Is(err, ErrClosed) {
		t.Fatalf("Handle after second close: want ErrClosed, got %v", err)
	}

	// The ring must not have grown: a closed receiver must not accept tail
	// packets into the ring buffer.
	if got := r.RingLen(); got != 0 {
		t.Fatalf("RingLen after closed Handle: want 0, got %d", got)
	}
}

// TestHandleRedeliveryConverges models the upstream redelivery loop: it
// retries until the error class indicates a terminal closed state. With the
// bug, Handle either panicked or returned a non-ErrClosed error (or
// ErrNoSink), so this loop never converged and the queue churned forever.
func TestHandleRedeliveryConverges(t *testing.T) {
	r := NewReceiver(4)
	r.SetSink(NewMemSink())
	r.Close()

	const maxAttempts = 5
	var last error
	converged := false
	for i := 0; i < maxAttempts; i++ {
		_, err := r.Handle(context.Background(), []byte(validRFC3164))
		last = err
		if errors.Is(err, ErrClosed) {
			converged = true
			break
		}
	}
	if !converged {
		t.Fatalf("redelivery did not converge on ErrClosed within %d attempts; last=%v", maxAttempts, last)
	}
}
