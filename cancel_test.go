package syslogrx

import (
	"context"
	"errors"
	"testing"
	"time"
)

// A canceled upstream must fail the call AND leave the ring untouched.
// Regression for ghost-log bug: a canceled ctx used to fall through to
// sink.Write, leaving phantom entries that surfaced on reconnect.
func TestHandleCanceledContextDoesNotRing(t *testing.T) {
	r := NewReceiver(8)
	r.SetSink(NewMemSink())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := r.Handle(ctx, []byte("<13>Oct 11 22:14:15 host app: hello")); err == nil {
		t.Fatal("canceled ctx: expected error, got nil")
	} else if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled ctx: expected context.Canceled, got %v", err)
	}
	if got := r.RingLen(); got != 0 {
		t.Fatalf("canceled ctx: ring must stay empty, got %d", got)
	}
}

// ServeContext must honor an upstream cancel: stop draining and return the
// cause without processing the in-flight packet into the ring.
func TestServeContextCanceledMidStream(t *testing.T) {
	r := NewReceiver(8)
	r.SetSink(NewMemSink())

	ctx, cancel := context.WithCancel(context.Background())
	pkts := make(chan []byte)

	go func() {
		<-ctx.Done() // wait for cancel to take effect
		pkts <- []byte("<13>Oct 11 22:14:15 host app: ghost")
		close(pkts)
	}()
	cancel()

	err := r.ServeContext(ctx, pkts)
	if err == nil {
		t.Fatal("canceled upstream: expected error from ServeContext, got nil")
	}
	if got := r.RingLen(); got != 0 {
		t.Fatalf("canceled upstream: ring must stay empty, got %d", got)
	}
}

// Happy path still works after the fix: a live ctx lands exactly one entry.
func TestServeContextLivePacket(t *testing.T) {
	r := NewReceiver(8)
	r.SetSink(NewMemSink())

	pkts := make(chan []byte, 1)
	pkts <- []byte("<13>Oct 11 22:14:15 host app: ok")
	close(pkts)

	if err := r.ServeContext(context.Background(), pkts); err != nil {
		t.Fatalf("live upstream: unexpected error %v", err)
	}
	if got := r.RingLen(); got != 1 {
		t.Fatalf("live upstream: want ring len 1, got %d", got)
	}
	// unblock nothing; just keep linter happy about time import in CI envs
	_ = time.Now
}
