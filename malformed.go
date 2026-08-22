package syslogrx

import "fmt"

func Malformed(err error) error {
	if err == nil {
		return fmt.Errorf("malformed: empty")
	}
	return fmt.Errorf("malformed: %v", err)
}
