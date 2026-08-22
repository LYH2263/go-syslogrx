package syslogrx_test

import (
	"os"
	"path/filepath"
	"testing"

	syslogrx "github.com/LYH2263/go-syslogrx"
)

func TestBug09_RotateClosesOldFile(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.log")
	p2 := filepath.Join(dir, "b.log")
	s, err := syslogrx.OpenFileSink(p1)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Write(&syslogrx.Message{RawBytes: []byte("one")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Rotate(p2); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(p1); err != nil {
		t.Fatalf("locked %v", err)
	}
}
