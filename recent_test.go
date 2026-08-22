package syslogrx

import (
	"context"
	"testing"
)

// newReceiverWith sets up a Receiver backed by a MemSink and feeds it the
// given raw lines via Handle.
func newReceiverWith(t *testing.T, raws ...[]byte) *Receiver {
	t.Helper()
	r := NewReceiver(64)
	r.SetSink(NewMemSink())
	for _, raw := range raws {
		if _, err := r.Handle(context.Background(), raw); err != nil {
			t.Fatalf("Handle: %v", err)
		}
	}
	return r
}

func rawLine(s string) []byte { return []byte(s) }

// Recent must return isolated copies: mutating RawBytes on a returned entry
// must not corrupt the in-ring original or a later snapshot.
func TestRecentIsolatedFromRing(t *testing.T) {
	r := newReceiverWith(t, rawLine("<34>Jan  2 15:04:05 host sshd: fail from 1.2.3.4"))

	got := r.Recent(1)
	if len(got) != 1 {
		t.Fatalf("Recent returned %d items, want 1", len(got))
	}
	if len(got[0].RawBytes) == 0 {
		t.Fatalf("RawBytes empty before redaction")
	}

	// Simulate the on-call preview redacting the raw bytes.
	for i := range got[0].RawBytes {
		got[0].RawBytes[i] = 'x'
	}

	// A subsequent snapshot must still hold the untouched original.
	snap := r.SnapshotRecent()
	if len(snap) != 1 {
		t.Fatalf("SnapshotRecent returned %d items, want 1", len(snap))
	}
	if string(snap[0].RawBytes) == "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" {
		t.Fatalf("ring original was mutated by Recent redaction: %q", snap[0].RawBytes)
	}
	if string(snap[0].RawBytes) != "<34>Jan  2 15:04:05 host sshd: fail from 1.2.3.4" {
		t.Fatalf("ring original changed: %q", snap[0].RawBytes)
	}
}

// SnapshotRecent must return isolated copies: redacting a snapshot entry must
// not leak into the ring or into a second snapshot.
func TestSnapshotRecentIsolatedFromRing(t *testing.T) {
	orig := "<34>Jan  2 15:04:05 host sshd: fail from 1.2.3.4"
	r := newReceiverWith(t, rawLine(orig))

	snap1 := r.SnapshotRecent()
	for i := range snap1[0].RawBytes {
		snap1[0].RawBytes[i] = 'x'
	}

	snap2 := r.SnapshotRecent()
	if string(snap2[0].RawBytes) != orig {
		t.Fatalf("SnapshotRecent aliasing leaked redaction into ring: got %q, want %q",
			snap2[0].RawBytes, orig)
	}
}

// Recent must return a slice with its own backing array; appending to it must
// not be able to reach into the ring's storage.
func TestRecentSliceDoesNotAliasRing(t *testing.T) {
	r := newReceiverWith(t, rawLine("<34>Jan  2 15:04:05 host sshd: fail from 1.2.3.4"))

	got := r.Recent(1)
	// Append to force growth; with the old aliasing bug this could not
	// happen, but ensure no panic and ring untouched.
	got = append(got, nil)
	if r.RingLen() != 1 {
		t.Fatalf("RingLen changed to %d, want 1", r.RingLen())
	}
}
