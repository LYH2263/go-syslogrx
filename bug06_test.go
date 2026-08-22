package syslogrx_test

import (
	"context"
	"testing"

	syslogrx "github.com/LYH2263/go-syslogrx"
)

func TestBug06_SinkFailNoAck(t *testing.T) {
	r := syslogrx.NewReceiver(8)
	defer r.Close()
	s := syslogrx.NewMemSink()
	s.SetFail(true)
	r.SetSink(s)
	if _, err := r.HandleAck(context.Background(), []byte("<34>Oct 11 22:14:15 mymachine su: su root failed")); err == nil {
		t.Fatal("want fail")
	}
}
