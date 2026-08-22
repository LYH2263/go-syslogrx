package syslogrx

import (
	"errors"
	"testing"
)

func TestParseLineMalformedSentinel(t *testing.T) {
	cases := []struct {
		name string
		line string
	}{
		{"empty", ""},
		{"no prefix", "not a syslog line"},
		{"no closing bracket", "<34 Jan 2 15:04:05 host app content"},
		{"bad priority", "<xx>1 - host - - - - msg"},
		{"too few fields", "<34>1 - host"},
		{"short body", "<13>abc"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := ParseLine(c.line)
			if err == nil {
				t.Fatalf("expected an error for line %q, got nil", c.line)
			}
			if !errors.Is(err, ErrMalformed) {
				t.Fatalf("expected errors.Is(err, ErrMalformed) for line %q, got %v", c.line, err)
			}
		})
	}
}

func TestParseLineMalformedChainsUnderlying(t *testing.T) {
	// A malformed line that fails RFC3164 with ErrInvalid should still match
	// ErrInvalid through the wrapped error chain, so existing callers keep working.
	_, err := ParseLine("no closing bracket here")
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected errors.Is(err, ErrInvalid), got %v", err)
	}
}

func TestParseLineValid(t *testing.T) {
	// Sanity: valid RFC5424 and RFC3164 lines still parse and never carry ErrMalformed.
	lines := []string{
		"<34>1 2003-10-11T22:14:15.003Z host app procid msgid content",
		"<34>Oct 11 22:14:15 host app content",
	}
	for _, line := range lines {
		m, err := ParseLine(line)
		if err != nil {
			t.Fatalf("ParseLine(%q): %v", line, err)
		}
		if m == nil {
			t.Fatalf("ParseLine(%q): nil message", line)
		}
	}
}
