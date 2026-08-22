package syslogrx

import (
	"os"
	"path/filepath"
	"testing"
)

// newMessage builds a minimal Message suitable for FileSink.Write, which only
// cares about RawBytes.
func newMessage(raw string) *Message {
	return &Message{Raw: raw, RawBytes: []byte(raw)}
}

// TestFileSinkRotateReleasesOldHandle is the regression test for bug 9: after
// Rotate opens the new file, the previous file handle must be closed so the
// shipper can collect and delete the rotated file. On Windows (and any system
// that holds an inode while a handle is open) an unclosed handle keeps the old
// path busy, so os.RemoveAll fails.
func TestFileSinkRotateReleasesOldHandle(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.log")
	newPath := filepath.Join(dir, "new.log")

	s, err := OpenFileSink(oldPath)
	if err != nil {
		t.Fatalf("OpenFileSink: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	if err := s.Write(newMessage("first")); err != nil {
		t.Fatalf("write old: %v", err)
	}

	if err := s.Rotate(newPath); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	// New path is active: writes must land there.
	if err := s.Write(newMessage("second")); err != nil {
		t.Fatalf("write new: %v", err)
	}

	// The rotated (old) file must no longer be held by this process. If the
	// handle leaked, RemoveAll fails with a "file in use" style error on
	// Windows and is tolerated on POSIX — so also assert via OpenFile that we
	// can obtain a fresh exclusive write handle.
	if err := os.RemoveAll(oldPath); err != nil {
		t.Fatalf("RemoveAll old path after rotate: %v (handle not released)", err)
	}
	got, err := os.OpenFile(oldPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("reopen old path exclusively: %v (handle not released)", err)
	}
	got.Close()
	os.Remove(oldPath)

	// Sink must still be usable on the new path.
	if s.f == nil {
		t.Fatal("s.f nil after rotate")
	}
	if s.path != newPath {
		t.Fatalf("s.path = %q, want %q", s.path, newPath)
	}
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("new path missing: %v", err)
	}
}

// TestFileSinkRotatePreservesOldOnOpenFailure ensures that when opening the new
// file fails, the old handle is left intact and the sink stays usable.
func TestFileSinkRotatePreservesOldOnOpenFailure(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.log")
	// Pointing the new path at an existing directory forces OpenFile to fail
	// (Windows: access denied; POSIX: EISDIR) without disturbing the old handle.
	badNew := dir

	s, err := OpenFileSink(oldPath)
	if err != nil {
		t.Fatalf("OpenFileSink: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	if err := s.Write(newMessage("first")); err != nil {
		t.Fatalf("write old: %v", err)
	}

	if err := s.Rotate(badNew); err == nil {
		t.Fatal("Rotate want error, got nil")
	}

	// Old handle must still be live and writing must still go to oldPath.
	if s.f == nil {
		t.Fatal("s.f nil after failed rotate")
	}
	if s.path != oldPath {
		t.Fatalf("s.path = %q, want %q", s.path, oldPath)
	}
	if err := s.Write(newMessage("second")); err != nil {
		t.Fatalf("write after failed rotate: %v", err)
	}
}
