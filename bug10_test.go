package syslogrx_test

import (
	"context"
	"testing"

	syslogrx "github.com/LYH2263/go-syslogrx"
)

func TestBug10_CloseFlushesSinkFirst(t *testing.T) {
	r := syslogrx.NewReceiver(8)
	r.SetSink(syslogrx.NewMemSink())
	if _, err := r.Handle(context.Background(), []byte("<34>Oct 11 22:14:15 mymachine su: su root failed")); err != nil {
		t.Fatal(err)
	}
	if r.CloseFlushCount() == 0 {
		t.Fatal("want flush")
	}
}
