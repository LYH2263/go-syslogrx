package syslogrx_test

import (
	"context"
	"errors"
	"testing"

	syslogrx "github.com/LYH2263/go-syslogrx"
)

func TestBug04_NilSinkNoDirty(t *testing.T) {
	r := syslogrx.NewReceiver(8)
	defer r.Close()
	// intentionally no SetSink
	var panicked bool
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		_, err := r.Handle(context.Background(), []byte("<34>Oct 11 22:14:15 mymachine su: su root failed"))
		if panicked {
			return
		}
		if err == nil || !errors.Is(err, syslogrx.ErrNoSink) {
			t.Fatalf("%v", err)
		}
	}()
	if panicked {
		t.Fatal("panic")
	}
	if r.RingLen() != 0 {
		t.Fatalf("dirty ring %d", r.RingLen())
	}
	r.SetSink(syslogrx.NewMemSink())
	if _, err := r.Handle(context.Background(), []byte("<34>Oct 11 22:14:15 mymachine su: su root failed")); err != nil {
		t.Fatal(err)
	}
}
