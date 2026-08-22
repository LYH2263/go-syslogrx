package syslogrx_test

import (
	"bytes"
	"context"
	"testing"

	syslogrx "github.com/LYH2263/go-syslogrx"
)

func TestBug01_HandleIsolatesRaw(t *testing.T) {
	r := syslogrx.NewReceiver(8)
	defer r.Close()
	r.SetSink(syslogrx.NewMemSink())
	buf := []byte("<34>Oct 11 22:14:15 mymachine su: su root failed")
	if _, err := r.Handle(context.Background(), buf); err != nil {
		t.Fatal(err)
	}
	buf[0] = 'Z'
	snap := r.SnapshotRecent()
	if len(snap) == 0 || bytes.Equal(snap[0].RawBytes, buf) {
		t.Fatalf("handle leak %q", snap[0].RawBytes)
	}
	snap[0].RawBytes[0] = 'Q'
	if r.SnapshotRecent()[0].RawBytes[0] == 'Q' {
		t.Fatal("snapshot leak")
	}
}
