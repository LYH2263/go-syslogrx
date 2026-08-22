package syslogrx_test

import (
	"errors"
	"testing"

	syslogrx "github.com/LYH2263/go-syslogrx"
)

func TestBug05_MalformedWrapped(t *testing.T) {
	_, err := syslogrx.ParseLine("not-a-syslog")
	if err == nil || !errors.Is(err, syslogrx.ErrMalformed) {
		t.Fatalf("%v", err)
	}
	err2 := syslogrx.Malformed(err)
	if !errors.Is(err2, syslogrx.ErrMalformed) {
		t.Fatalf("malformed wrap %v", err2)
	}
}
