package syslogrx_test

import (
	"context"
	"testing"

	syslogrx "github.com/LYH2263/go-syslogrx"
)

func TestBug02_RecentIsolated(t *testing.T) {
	r := syslogrx.NewReceiver(8)
	defer r.Close()
	r.SetSink(syslogrx.NewMemSink())
	raw := []byte("<34>Oct 11 22:14:15 mymachine su: su root failed")
	if _, err := r.Handle(context.Background(), raw); err != nil {
		t.Fatal(err)
	}
	got := r.Recent(1)
	got[0].RawBytes[0] ^= 0xff
	if r.Recent(1)[0].RawBytes[0] != raw[0] {
		t.Fatal("recent aliased")
	}
	snap := r.SnapshotRecent()
	snap[0].RawBytes[1] ^= 0xff
	if r.SnapshotRecent()[0].RawBytes[1] != raw[1] {
		t.Fatal("snapshot aliased")
	}
}
