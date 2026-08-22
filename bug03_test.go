package syslogrx_test

import (
	"context"
	"errors"
	"testing"

	syslogrx "github.com/LYH2263/go-syslogrx"
)

func TestBug03_HandleAfterClose(t *testing.T) {
	r := syslogrx.NewReceiver(8)
	r.SetSink(syslogrx.NewMemSink())
	_ = r.Close()
	var panicked bool
	var err error
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		_, err = r.Handle(context.Background(), []byte("<34>Oct 11 22:14:15 mymachine su: su root failed"))
	}()
	if panicked {
		t.Fatal("panic")
	}
	if !errors.Is(err, syslogrx.ErrClosed) {
		t.Fatalf("%v", err)
	}
}
