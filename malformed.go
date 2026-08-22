package syslogrx

import "fmt"

// Malformed wraps a parse miss as ErrMalformed for callers using errors.Is.
func Malformed(err error) error {
	if err == nil {
		return fmt.Errorf("%w: empty", ErrMalformed)
	}
	return fmt.Errorf("%w: %v", ErrMalformed, err)
}
