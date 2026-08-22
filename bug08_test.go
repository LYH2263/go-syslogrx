package syslogrx_test

import (
	"context"
	"testing"
	"time"

	syslogrx "github.com/LYH2263/go-syslogrx"
)

func TestBug08_WriteContextHonorsCancel(t *testing.T) {
	s := syslogrx.NewMemSink()
	m := &syslogrx.Message{RawBytes: []byte("x")}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	if err := syslogrx.WriteContext(ctx, s, m, time.Second); err == nil {
		t.Fatal("want ctx err")
	}
}
