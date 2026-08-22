package syslogrx

import "fmt"

// Malformed wraps err as an ErrMalformed sentinel error so callers can
// identify syslog format failures via errors.Is(err, ErrMalformed) instead
// of string-matching the "malformed" text. The underlying error is wrapped
// with %w so errors.Is against it (e.g. ErrInvalid) still holds.
func Malformed(err error) error {
	if err == nil {
		return fmt.Errorf("%w: empty", ErrMalformed)
	}
	return fmt.Errorf("%w: %w", ErrMalformed, err)
}
