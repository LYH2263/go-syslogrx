package syslogrx_test

import (
	"context"
	"testing"

	syslogrx "github.com/LYH2263/go-syslogrx"
)

func TestBug07_ServeContextHonorsCancel(t *testing.T) {
	r := syslogrx.NewReceiver(8)
	defer r.Close()
	r.SetSink(syslogrx.NewMemSink())
	before := r.RingLen()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ch := make(chan []byte, 1)
	ch <- []byte("<34>Oct 11 22:14:15 mymachine su: su root failed")
	close(ch)
	if err := r.ServeContext(ctx, ch); err == nil {
		t.Fatal("want cancel")
	}
	if r.RingLen() != before {
		t.Fatalf("%d->%d", before, r.RingLen())
	}
}
